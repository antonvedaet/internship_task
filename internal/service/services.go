package service

import "antonvedaet/internship_task/internal/models"

type TeamService interface {
	CreateTeam(team *models.Team) error
	GetTeam(teamName string) (*models.Team, error)
}

type UserService interface {
	SetUserActive(userID string, isActive bool) (*models.User, error)
	GetUserReviewPRs(userID string) ([]models.PullRequest, error)
}

type PRService interface {
	CreatePR(prRequest *models.CreatePRRequest) (*models.PullRequest, error)
	MergePR(prID string) (*models.PullRequest, error)
	ReassignReviewer(prID, oldReviewerID string) (*models.PullRequest, string, error)
}
