package service

import "errors"

var (
	ErrTeamExists           = errors.New("team already exists")
	ErrPRExists             = errors.New("PR already exists")
	ErrNotFound             = errors.New("not found")
	ErrPRAlreadyMerged      = errors.New("PR already merged")
	ErrReviewerNotAssigned  = errors.New("reviewer not assigned")
	ErrNoAvailableReviewers = errors.New("no available reviewers")
	ErrUserNotInTeam        = errors.New("user not in team")
)
