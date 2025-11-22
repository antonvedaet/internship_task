package service

import (
	"time"

	"antonvedaet/internship_task/internal/models"
	"antonvedaet/internship_task/internal/store"
)

type prService struct {
	db *store.DB
}

func NewPRService(db *store.DB) PRService {
	return &prService{db: db}
}

func (s *prService) CreatePR(prRequest *models.CreatePRRequest) (*models.PullRequest, error) {
	exists, err := s.db.PRExists(prRequest.PullRequestID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrPRExists
	}

	author, err := s.db.GetUser(prRequest.AuthorID)
	if err != nil {
		return nil, ErrNotFound
	}

	teamUsers, err := s.db.GetActiveTeamUsers(author.TeamName, prRequest.AuthorID)
	if err != nil {
		return nil, err
	}

	reviewers := s.selectRandomReviewers(teamUsers, 2)

	pr := &models.PullRequest{
		PullRequestID:     prRequest.PullRequestID,
		PullRequestName:   prRequest.PullRequestName,
		AuthorID:          prRequest.AuthorID,
		Status:            "OPEN",
		AssignedReviewers: reviewers,
		CreatedAt:         time.Now(),
	}

	if err := s.db.CreatePR(pr); err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *prService) MergePR(prID string) (*models.PullRequest, error) {
	pr, err := s.db.GetPR(prID)
	if err != nil {
		return nil, ErrNotFound
	}

	if pr.Status == "MERGED" {
		return pr, nil
	}

	pr.Status = "MERGED"
	now := time.Now()
	pr.MergedAt = &now

	if err := s.db.UpdatePR(pr); err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *prService) ReassignReviewer(prID, oldReviewerID string) (*models.PullRequest, string, error) {
	pr, err := s.db.GetPR(prID)
	if err != nil {
		return nil, "", ErrNotFound
	}

	if pr.Status == "MERGED" {
		return nil, "", ErrPRAlreadyMerged
	}

	if !contains(pr.AssignedReviewers, oldReviewerID) {
		return nil, "", ErrReviewerNotAssigned
	}

	oldReviewer, err := s.db.GetUser(oldReviewerID)
	if err != nil {
		return nil, "", ErrNotFound
	}

	teamUsers, err := s.db.GetActiveTeamUsers(oldReviewer.TeamName, oldReviewerID)
	if err != nil {
		return nil, "", err
	}

	var candidates []models.User
	for _, user := range teamUsers {
		if !contains(pr.AssignedReviewers, user.UserID) {
			candidates = append(candidates, user)
		}
	}

	if len(candidates) == 0 {
		return nil, "", ErrNoAvailableReviewers
	}

	newReviewer := candidates[randomInt(len(candidates))]

	for i, reviewer := range pr.AssignedReviewers {
		if reviewer == oldReviewerID {
			pr.AssignedReviewers[i] = newReviewer.UserID
			break
		}
	}

	if err := s.db.UpdatePR(pr); err != nil {
		return nil, "", err
	}

	return pr, newReviewer.UserID, nil
}

func (s *prService) selectRandomReviewers(users []models.User, max int) []string {
	if len(users) == 0 {
		return []string{}
	}

	shuffled := make([]models.User, len(users))
	copy(shuffled, users)
	for i := range shuffled {
		j := randomInt(len(shuffled))
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	count := min(len(shuffled), max)
	reviewers := make([]string, count)
	for i := 0; i < count; i++ {
		reviewers[i] = shuffled[i].UserID
	}

	return reviewers
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func randomInt(n int) int {
	return int(time.Now().UnixNano()) % n
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
