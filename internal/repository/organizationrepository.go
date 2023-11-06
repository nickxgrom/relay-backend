package repository

import (
	"relay-backend/internal/model"
	"relay-backend/internal/store"
)

type OrganizationRepository struct {
	store *store.Store
}

func NewOrganizationRepository(s *store.Store) *OrganizationRepository {
	return &OrganizationRepository{
		store: s,
	}
}

func (or *OrganizationRepository) Save(organization *model.Organization) error {
	if err := or.store.Db.QueryRow(
		"insert into organizations (owner_id, name, description, address, email, creation_date) values ((select id from users where id=$1), $2, $3, $4, $5, current_date) on conflict (owner_id, name) do nothing returning id, creation_date",
		&organization.OwnerId,
		&organization.Name,
		&organization.Description,
		&organization.Address,
		&organization.Email,
	).Scan(&organization.Id, &organization.CreationDate); err != nil {
		return err
	}

	return nil
}

// select users.* from users join employees on users.id = employees.user_id where employees.organization_id = $1
func (or *OrganizationRepository) Find(userId int, orgId int) (*model.Organization, error) {
	org := &model.Organization{}

	if err := or.store.Db.QueryRow("select * from organizations where id = $1 and owner_id = $2", orgId, userId).Scan(
		&org.Id,
		&org.OwnerId,
		&org.Name,
		&org.Description,
		&org.Address,
		&org.Email,
		&org.CreationDate,
	); err != nil {
		return nil, err
	}

	return org, nil
}
func (or *OrganizationRepository) GetList() ([]*model.Organization, error) {
	return nil, nil
}

func (or *OrganizationRepository) Update() error {
	return nil
}

func (or *OrganizationRepository) Delete() error {
	return nil
}
