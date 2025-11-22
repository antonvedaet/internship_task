package store

import (
	"antonvedaet/internship_task/internal/models"
	"fmt"
)

func (db *DB) CreateTeam(team *models.Team) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO teams (team_name) VALUES ($1) ON CONFLICT (team_name) DO NOTHING", team.TeamName)
	if err != nil {
		return err
	}

	for _, member := range team.Members {
		_, err = tx.Exec(`
            INSERT INTO users (user_id, username, team_name, is_active) 
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (user_id) DO UPDATE SET
                username = EXCLUDED.username,
                team_name = EXCLUDED.team_name,
                is_active = EXCLUDED.is_active
        `, member.UserID, member.Username, team.TeamName, member.IsActive)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (db *DB) GetTeam(teamName string) (*models.Team, error) {
	var team models.Team
	team.TeamName = teamName

	rows, err := db.Query(`
        SELECT user_id, username, is_active 
        FROM users 
        WHERE team_name = $1
    `, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var member models.TeamMember
		if err := rows.Scan(&member.UserID, &member.Username, &member.IsActive); err != nil {
			return nil, err
		}
		team.Members = append(team.Members, member)
	}

	if len(team.Members) == 0 {
		return nil, fmt.Errorf("team not found")
	}

	return &team, nil
}

// User
func (db *DB) GetUser(userID string) (*models.User, error) {
	var user models.User
	err := db.QueryRow(`
        SELECT user_id, username, team_name, is_active 
        FROM users 
        WHERE user_id = $1
    `, userID).Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *DB) UpdateUser(user *models.User) error {
	_, err := db.Exec(`
        UPDATE users 
        SET username = $1, team_name = $2, is_active = $3 
        WHERE user_id = $4
    `, user.Username, user.TeamName, user.IsActive, user.UserID)
	return err
}

func (db *DB) GetActiveTeamUsers(teamName, excludeUserID string) ([]models.User, error) {
	var users []models.User
	query := `
        SELECT user_id, username, team_name, is_active 
        FROM users 
        WHERE team_name = $1 AND is_active = true
    `
	args := []interface{}{teamName}

	if excludeUserID != "" {
		query += " AND user_id != $2"
		args = append(args, excludeUserID)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// PullRequest
func (db *DB) CreatePR(pr *models.PullRequest) error {
	_, err := db.Exec(`
        INSERT INTO pull_requests 
        (pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at) 
        VALUES ($1, $2, $3, $4, $5, $6)
    `, pr.PullRequestID, pr.PullRequestName, pr.AuthorID, pr.Status, pr.AssignedReviewers, pr.CreatedAt)
	return err
}

func (db *DB) GetPR(prID string) (*models.PullRequest, error) {
	var pr models.PullRequest
	err := db.QueryRow(`
        SELECT pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at, merged_at
        FROM pull_requests 
        WHERE pull_request_id = $1
    `, prID).Scan(
		&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status,
		&pr.AssignedReviewers, &pr.CreatedAt, &pr.MergedAt,
	)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

func (db *DB) UpdatePR(pr *models.PullRequest) error {
	_, err := db.Exec(`
        UPDATE pull_requests 
        SET status = $1, assigned_reviewers = $2, merged_at = $3 
        WHERE pull_request_id = $4
    `, pr.Status, pr.AssignedReviewers, pr.MergedAt, pr.PullRequestID)
	return err
}

func (db *DB) GetPRsByReviewer(userID string) ([]models.PullRequest, error) {
	var prs []models.PullRequest
	rows, err := db.Query(`
        SELECT pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at, merged_at
        FROM pull_requests 
        WHERE $1 = ANY(assigned_reviewers)
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pr models.PullRequest
		if err := rows.Scan(
			&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status,
			&pr.AssignedReviewers, &pr.CreatedAt, &pr.MergedAt,
		); err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}

	return prs, nil
}

func (db *DB) PRExists(prID string) (bool, error) {
	var exists bool
	err := db.QueryRow(`
        SELECT EXISTS(SELECT 1 FROM pull_requests WHERE pull_request_id = $1)
    `, prID).Scan(&exists)
	return exists, err
}
