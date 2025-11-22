package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"antonvedaet/internship_task/internal/models"
	"antonvedaet/internship_task/internal/store"
)

type Handlers struct {
	db *store.DB
}

func NewHandlers(db *store.DB) *Handlers {
	return &Handlers{db: db}
}

func (h *Handlers) AddTeam(w http.ResponseWriter, r *http.Request) {
	var team models.Team
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if team.TeamName == "" {
		h.sendErrorResponse(w, "INVALID_REQUEST", "team_name is required", http.StatusBadRequest)
		return
	}

	if err := h.db.CreateTeam(&team); err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			h.sendErrorResponse(w, "TEAM_EXISTS", "team_name already exists", http.StatusBadRequest)
		} else {
			h.sendError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.TeamResponse{Team: &team})
}

func (h *Handlers) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		h.sendErrorResponse(w, "INVALID_REQUEST", "team_name is required", http.StatusBadRequest)
		return
	}

	team, err := h.db.GetTeam(teamName)
	if err != nil {
		h.sendErrorResponse(w, "NOT_FOUND", "team not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(team)
}

func (h *Handlers) SetUserActive(w http.ResponseWriter, r *http.Request) {
	var req models.SetActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.db.GetUser(req.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			h.sendErrorResponse(w, "NOT_FOUND", "user not found", http.StatusNotFound)
		} else {
			h.sendError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	user.IsActive = req.IsActive
	if err := h.db.UpdateUser(user); err != nil {
		h.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.SetActiveResponse{User: user})
}

func (h *Handlers) CreatePR(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	exists, err := h.db.PRExists(req.PullRequestID)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if exists {
		h.sendErrorResponse(w, "PR_EXISTS", "PR id already exists", http.StatusConflict)
		return
	}

	// TODO: Implement business logic for assigning reviewers
	pr := &models.PullRequest{
		PullRequestID:     req.PullRequestID,
		PullRequestName:   req.PullRequestName,
		AuthorID:          req.AuthorID,
		Status:            "OPEN",
		AssignedReviewers: []string{},
	}

	if err := h.db.CreatePR(pr); err != nil {
		h.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.CreatePRResponse{PR: pr})
}

func (h *Handlers) MergePR(w http.ResponseWriter, r *http.Request) {
	var req models.MergePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	pr, err := h.db.GetPR(req.PullRequestID)
	if err != nil {
		if err == sql.ErrNoRows {
			h.sendErrorResponse(w, "NOT_FOUND", "PR not found", http.StatusNotFound)
		} else {
			h.sendError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// TODO: Implement merge logic
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.MergePRResponse{PR: pr})
}

func (h *Handlers) ReassignReviewer(w http.ResponseWriter, r *http.Request) {
	var req models.ReassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	pr, err := h.db.GetPR(req.PullRequestID)
	if err != nil {
		if err == sql.ErrNoRows {
			h.sendErrorResponse(w, "NOT_FOUND", "PR not found", http.StatusNotFound)
		} else {
			h.sendError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// TODO: Implement reassign logic
	response := models.ReassignResponse{
		PR:         pr,
		ReplacedBy: req.OldUserID, // Placeholder
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) GetUserReview(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		h.sendErrorResponse(w, "INVALID_REQUEST", "user_id is required", http.StatusBadRequest)
		return
	}

	prs, err := h.db.GetPRsByReviewer(userID)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shortPRs := make([]models.PullRequestShort, len(prs))
	for i, pr := range prs {
		shortPRs[i] = models.PullRequestShort{
			PullRequestID:   pr.PullRequestID,
			PullRequestName: pr.PullRequestName,
			AuthorID:        pr.AuthorID,
			Status:          pr.Status,
		}
	}

	response := models.UserReviewResponse{
		UserID:       userID,
		PullRequests: shortPRs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	if err := h.db.Ping(); err != nil {
		http.Error(w, "Database unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"healthy"}`))
}

func (h *Handlers) sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func (h *Handlers) sendErrorResponse(w http.ResponseWriter, code, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}{
			Code:    code,
			Message: message,
		},
	})
}
