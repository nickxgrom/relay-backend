package enums

type UserRole byte

const (
	OrganizationOwner UserRole = 0
	OrganizationAdmin          = 1
	Operator                   = 2
)
