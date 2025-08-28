package types

import (
	"fmt"
	"strings"
)

// Role represents a user role in the system.
type Role int

const (
	// RoleUnknown is an unknown or unassigned role.
	RoleUnknown Role = iota
	// RoleAdmin has all permissions.
	RoleAdmin
	// RoleUser is a standard user with limited permissions.
	RoleUser
)

// String returns the string representation of a Role.
func (r Role) String() string {
	switch r {
	case RoleAdmin:
		return "admin"
	case RoleUser:
		return "user"
	default:
		return "unknown"
	}
}

// ParseRole converts a string to a Role.
func ParseRole(s string) (Role, error) {
	switch s {
	case "admin":
		return RoleAdmin, nil
	case "user":
		return RoleUser, nil
	default:
		return RoleUnknown, fmt.Errorf("unknown role: %s", s)
	}
}

// Roles is a slice of Role that implements xorm.Conversion to handle
// PostgreSQL's text[] type.
type Roles []Role

// FromDB is called by xorm to convert a database value to a Roles slice.
// It parses a PostgreSQL array string like "{admin,user}" into []Role.
func (r *Roles) FromDB(data []byte) error {
	// The data from postgres will be in the format "{admin,user}"
	s := string(data)
	if s == "{}" || s == "" {
		*r = []Role{}
		return nil
	}

	// Trim the curly braces and split by comma.
	// Note: This is a simple parser and will not handle roles with commas
	// or quotes in their names. For the current roles ("admin", "user"), this is safe.
	trimmed := strings.Trim(s, "{}")
	parts := strings.Split(trimmed, ",")
	roles := make([]Role, 0, len(parts))

	for _, part := range parts {
		role, err := ParseRole(part)
		if err != nil {
			return err
		}
		roles = append(roles, role)
	}

	*r = roles
	return nil
}

// ToDB is called by xorm to convert a Roles slice to a database value.
// It converts []Role into a PostgreSQL array string like "{admin,user}".
func (r Roles) ToDB() ([]byte, error) {
	if r == nil || len(r) == 0 {
		return []byte("{}"), nil
	}

	parts := make([]string, len(r))
	for i, role := range r {
		parts[i] = role.String()
	}

	// Format as a postgres array literal
	return []byte("{" + strings.Join(parts, ",") + "}"), nil
}

// HasRole checks if a specific role exists in the slice.
func (r Roles) HasRole(role Role) bool {
	for _, existingRole := range r {
		if existingRole == role {
			return true
		}
	}
	return false
}
