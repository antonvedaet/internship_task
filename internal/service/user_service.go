package service

import (
	"antonvedaet/internship_task/internal/models"
	"antonvedaet/internship_task/internal/store"
)

type userService struct {
	db *store.DB
}

func NewUserService(db *store.DB) UserService {
	return &userService{db: db}
}

func (s *userService) SetUserActive(userID string, isActive bool) (*models.User, error) {
	user, err := s.db.GetUser(userID)
	if err != nil {
		return nil, ErrNotFound
	}

	user.IsActive = isActive
	if err := s.db.UpdateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUserReviewPRs(userID string) ([]models.PullRequest, error) {
	return s.db.GetPRsByReviewer(userID)
}
