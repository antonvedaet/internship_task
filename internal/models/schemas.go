package models

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type SetActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type SetActiveResponse struct {
	User *User `json:"user"`
}

type CreatePRRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type CreatePRResponse struct {
	PR *PullRequest `json:"pr"`
}

type MergePRRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

type MergePRResponse struct {
	PR *PullRequest `json:"pr"`
}

type ReassignRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_reviewer_id"`
}

type ReassignResponse struct {
	PR         *PullRequest `json:"pr"`
	ReplacedBy string       `json:"replaced_by"`
}

type UserReviewResponse struct {
	UserID       string             `json:"user_id"`
	PullRequests []PullRequestShort `json:"pull_requests"`
}

type TeamResponse struct {
	Team *Team `json:"team"`
}

type DeactivateTeamRequest struct {
	TeamName string `json:"team_name"`
}

type DeactivateTeamResponse struct {
	Message          string `json:"message"`
	DeactivatedCount int    `json:"deactivated_count"`
}
