package domain

type RoleType string

const (
	Admin   RoleType = "admin"
	User    RoleType = "user"
	Manager RoleType = "manager"
)

type UserModel struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Role     RoleType `json:"role"`
	Password string   `json:"password,omitempty"`
}
