package enums

type UserRole byte

type userRoleStruct struct {
	None              UserRole
	OrganizationOwner UserRole
	OrganizationAdmin UserRole
	Operator          UserRole
	Any               UserRole
}

var UserRoleEnum = userRoleStruct{
	None:              0,
	OrganizationOwner: 1,
	OrganizationAdmin: 2,
	Operator:          3,
	Any:               5,
}

type accessStruct struct {
	Owner         []UserRole
	OwnerAndAdmin []UserRole
	Any           []UserRole
}

var Access = accessStruct{
	Owner:         []UserRole{UserRoleEnum.OrganizationOwner},
	OwnerAndAdmin: []UserRole{UserRoleEnum.OrganizationOwner, UserRoleEnum.OrganizationAdmin},
	Any:           []UserRole{UserRoleEnum.Any},
}
