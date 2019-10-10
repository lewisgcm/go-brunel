package security

// UserRole is the role of a user
type UserRole string

const (
	// UserRoleAdmin is the administrator role, managing the system/users/jobs
	UserRoleAdmin UserRole = "admin"

	// UserRoleReader is a read only role
	UserRoleReader UserRole = "reader"

	// UserRoleAnonymous is the role used for anonymous un-authenticated users
	UserRoleAnonymous UserRole = "anonymous"
)

// Identity models the basic identity of a user, username is unique and role holds their role
type Identity struct {
	Username string
	Role     UserRole
}
