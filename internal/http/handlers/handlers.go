package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"antonvedaet/internship_task/internal/models"
	"antonvedaet/internship_task/internal/service"
)

type Handlers struct {
	teamService service.TeamService
	userService service.UserService
	prService   service.PRService
}

func NewHandlers(teamService service.TeamService, userService service.UserService, prService service.PRService) *Handlers {
	return &Handlers{
		teamService: teamService,
		userService: userService,
		prService:   prService,
	}
}

func (h *Handlers) AddTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var team models.Team
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if team.TeamName == "" {
		h.sendErrorResponse(w, "INVALID_REQUEST", "team_name is required", http.StatusBadRequest)
		return
	}

	if err := h.teamService.CreateTeam(&team); err != nil {
		if strings.Contains(err.Error(), "unique constraint") || err == service.ErrTeamExists {
			h.sendErrorResponse(w, "TEAM_EXISTS", "team_name already exists", http.StatusBadRequest)
		} else {
			log.Printf("Error creating team: %v", err)
			h.sendError(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.TeamResponse{Team: &team})
}

func (h *Handlers) GetTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		h.sendErrorResponse(w, "INVALID_REQUEST", "team_name is required", http.StatusBadRequest)
		return
	}

	team, err := h.teamService.GetTeam(teamName)
	if err != nil {
		if err == service.ErrNotFound {
			h.sendErrorResponse(w, "NOT_FOUND", "team not found", http.StatusNotFound)
		} else {
			log.Printf("Error getting team: %v", err)
			h.sendError(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(team)
}

func (h *Handlers) SetUserActive(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.SetActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userService.SetUserActive(req.UserID, req.IsActive)
	if err != nil {
		if err == service.ErrNotFound {
			h.sendErrorResponse(w, "NOT_FOUND", "user not found", http.StatusNotFound)
		} else {
			log.Printf("Error setting user active: %v", err)
			h.sendError(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.SetActiveResponse{User: user})
}

func (h *Handlers) CreatePR(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.CreatePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	pr, err := h.prService.CreatePR(&req)
	if err != nil {
		switch err {
		case service.ErrPRExists:
			h.sendErrorResponse(w, "PR_EXISTS", "PR id already exists", http.StatusConflict)
		case service.ErrNotFound:
			h.sendErrorResponse(w, "NOT_FOUND", "author/team not found", http.StatusNotFound)
		default:
			log.Printf("Error creating PR: %v", err)
			h.sendError(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.CreatePRResponse{PR: pr})
}

func (h *Handlers) MergePR(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.MergePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	pr, err := h.prService.MergePR(req.PullRequestID)
	if err != nil {
		if err == service.ErrNotFound {
			h.sendErrorResponse(w, "NOT_FOUND", "PR not found", http.StatusNotFound)
		} else {
			log.Printf("Error merging PR: %v", err)
			h.sendError(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.MergePRResponse{PR: pr})
}

func (h *Handlers) ReassignReviewer(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.ReassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	pr, newReviewerID, err := h.prService.ReassignReviewer(req.PullRequestID, req.OldUserID)
	if err != nil {
		switch err {
		case service.ErrNotFound:
			h.sendErrorResponse(w, "NOT_FOUND", "PR or user not found", http.StatusNotFound)
		case service.ErrPRAlreadyMerged:
			h.sendErrorResponse(w, "PR_MERGED", "cannot reassign on merged PR", http.StatusConflict)
		case service.ErrReviewerNotAssigned:
			h.sendErrorResponse(w, "NOT_ASSIGNED", "reviewer is not assigned to this PR", http.StatusConflict)
		case service.ErrNoAvailableReviewers:
			h.sendErrorResponse(w, "NO_CANDIDATE", "no active replacement candidate in team", http.StatusConflict)
		default:
			log.Printf("Error reassigning reviewer: %v", err)
			h.sendError(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	response := models.ReassignResponse{
		PR:         pr,
		ReplacedBy: newReviewerID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) GetUserReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		h.sendErrorResponse(w, "INVALID_REQUEST", "user_id is required", http.StatusBadRequest)
		return
	}

	prs, err := h.userService.GetUserReviewPRs(userID)
	if err != nil {
		log.Printf("Error getting user review PRs: %v", err)
		h.sendError(w, "Internal server error", http.StatusInternalServerError)
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
	if r.Method != "GET" {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
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
