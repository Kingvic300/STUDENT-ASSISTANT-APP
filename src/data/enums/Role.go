package enums

type Role string

const (
	Admin Role = "ADMIN"
	User  Role = "USER"
	Guest Role = "GUEST"
)

func (r Role) IsValid() bool {
	switch r {
	case Admin, User, Guest:
		return true
	default:
		return false
	}
}