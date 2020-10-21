package restapi

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"reflect"
	"strconv"
	"time"
)

// ListUser function return user data from database
func (c *InitAPI) ListUser(ctx context.Context, req *GetUsers) (*GetUsers, error) {
	limit := 10

	if req.Limit != 0 {
		limit = int(req.Limit)
	}

	rows, err := c.Db.Query(`
		SELECT ID,
			username,
			email,
			status,
			role_id,
			created_at,
			updated_at
		FROM users LIMIT $1
	`, limit)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()
	var items []*User
	for rows.Next() {
		var item User
		var status string
		var updateTime sql.NullString
		err = rows.Scan(
			&item.ID,
			&item.Username,
			&item.Email,
			&status,
			&item.RoleID,
			&item.CreatedAt,
			&updateTime,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		item.UpdatedAt = updateTime.String
		items = append(items, &item)
	}

	if len(items) == 0 {
		return nil, errors.New("user-not-found")
	}

	return &GetUsers{
		List: items,
	}, nil
}

// CreateUser for adding new user
func (c *InitAPI) CreateUser(ctx context.Context, req *User, rolesID string) (*UserID, error) {
	var id string
	roles, err := c.GetRoles(rolesID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if roles != "ADMIN" {
		return nil, errors.New("invalid-roles")
	}

	status := strconv.Itoa(req.Status)
	err = c.Db.QueryRow(`INSERT INTO users (username, email, status, role_id, updated_at) VALUES  ($1, $2, $3, $4, $5) RETURNING id`,
		req.Username, req.Email, status, req.RoleID, time.Now().Format("RFC3339")).Scan(&id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &UserID{
		ID: id,
	}, nil

}

// UpdateUser for adding new user
func (c *InitAPI) UpdateUser(ctx context.Context, req *User, rolesID string) (*UserID, error) {
	var id string
	roles, err := c.GetRoles(rolesID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(reflect.TypeOf(rolesID))
	if roles != "ADMIN" {
		return nil, errors.New("invalid-roles")
	}
	status := strconv.Itoa(req.Status)
	log.Println(req.ID)
	err = c.Db.QueryRow(`UPDATE users
	SET username = $1,
	email = $2,
	status = $3,
	role_id = $4,
	updated_at = $5
	WHERE id = $6 RETURNING id`,
		req.Username, req.Email, status, req.RoleID, time.Now().Format("RFC3339"), req.ID).Scan(&id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &UserID{
		ID: id,
	}, nil
}

// DeleteUser for adding new user
func (c *InitAPI) DeleteUser(ctx context.Context, req *UserID, rolesID string) (*UserID, error) {
	var id string
	roles, err := c.GetRoles(rolesID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(reflect.TypeOf(rolesID))
	if roles != "ADMIN" {
		return nil, errors.New("invalid-roles")
	}
	log.Println(req.ID)
	err = c.Db.QueryRow(`DELETE FROM users
	WHERE id = $1 RETURNING id`,
		req.ID).Scan(&id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &UserID{
		ID: id,
	}, nil
}

// GetRoles to get user roles
func (c *InitAPI) GetRoles(id string) (string, error) {
	var roles string
	err := c.Db.QueryRow(`SELECT roles FROM roles WHERE id = $1`, id).Scan(&roles)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return roles, nil
}
