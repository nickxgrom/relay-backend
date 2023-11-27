package enums

type UserRole byte

const (
	None              UserRole = 0
	OrganizationOwner          = 1
	OrganizationAdmin          = 2
	Operator                   = 3
	Any                        = 5
)
