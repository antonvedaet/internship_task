package service

import (
	"antonvedaet/internship_task/internal/models"
	"antonvedaet/internship_task/internal/store"
)

type teamService struct {
	db *store.DB
}

func NewTeamService(db *store.DB) TeamService {
	return &teamService{db: db}
}

func (s *teamService) CreateTeam(team *models.Team) error {
	return s.db.CreateTeam(team)
}

func (s *teamService) GetTeam(teamName string) (*models.Team, error) {
	return s.db.GetTeam(teamName)
}

func (s *teamService) DeactivateTeamUsers(teamName string) (int, error) {
	return s.db.DeactivateTeamUsers(teamName)
}
