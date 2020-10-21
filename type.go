package restapi

import (
	"time"

	"github.com/jackc/pgx"
)

// GetUsers struct for query for getusers api
type GetUsers struct {
	Limit int32   `json:"limit"`
	List  []*User `json:"list"`
}

// User struct for table user
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Status    int       `json:"status"`
	RoleID    string    `json:"roleId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt string    `json:"updatedAt"`
}

// InitAPI init
type InitAPI struct {
	Db *pgx.ConnPool
}

// UserID get rolesid
type UserID struct {
	ID string `json:"id"`
}

// UserName get username
type UserName struct {
	Name string `json:"name"`
}
