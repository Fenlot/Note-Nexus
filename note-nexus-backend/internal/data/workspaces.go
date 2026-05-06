package data

import (
	"database/sql"
	"time"
)

type Workspace struct {
	ID               int64     `json:"id"`
	Name             string    `json:"name"`
	OwnerID          int64     `json:"ownerId"`
	SubscriptionTier string    `json:"subscriptionTier"`
	CreatedAt        time.Time `json:"createdAt"`
}

type WorkspaceMember struct {
	WorkspaceID int64     `json:"workspaceId"`
	UserID      int64     `json:"userId"`
	Role        string    `json:"role"`
	JoinedAt    time.Time `json:"joinedAt"`
}

type WorkspaceModel struct {
	DB *sql.DB
}

func (m WorkspaceModel) Insert(name string, ownerID int64) (*Workspace, error) {
	tx, err := m.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO workspaces (name, owner_id) 
		VALUES ($1, $2) 
		RETURNING id, subscription_tier, created_at`

	args := []interface{}{name, ownerID}

	workspace := &Workspace{
		Name:    name,
		OwnerID: ownerID,
	}

	err = tx.QueryRow(query, args...).Scan(&workspace.ID, &workspace.SubscriptionTier, &workspace.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Add owner as a member
	memberQuery := `
		INSERT INTO workspace_members (workspace_id, user_id, role)
		VALUES ($1, $2, 'owner')`

	_, err = tx.Exec(memberQuery, workspace.ID, ownerID)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return workspace, nil
}

func (m WorkspaceModel) GetUserWorkspaces(userID int64) ([]*Workspace, error) {
	query := `
		SELECT w.id, w.name, w.owner_id, w.subscription_tier, w.created_at
		FROM workspaces w
		INNER JOIN workspace_members wm ON wm.workspace_id = w.id
		WHERE wm.user_id = $1`

	rows, err := m.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workspaces []*Workspace
	for rows.Next() {
		var w Workspace
		err := rows.Scan(&w.ID, &w.Name, &w.OwnerID, &w.SubscriptionTier, &w.CreatedAt)
		if err != nil {
			return nil, err
		}
		workspaces = append(workspaces, &w)
	}

	return workspaces, nil
}

func (m WorkspaceModel) HasAccess(workspaceID, userID int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM workspace_members 
			WHERE workspace_id = $1 AND user_id = $2
		)`
	
	var exists bool
	err := m.DB.QueryRow(query, workspaceID, userID).Scan(&exists)
	return exists, err
}
