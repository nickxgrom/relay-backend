package service

import (
	"relay-backend/internal/model"
	"relay-backend/internal/repository"
	"relay-backend/internal/store"
)

type OrganizationService struct {
	organizationRepository *repository.OrganizationRepository
}

func NewOrganizationService(s *store.Store) *OrganizationService {
	return &OrganizationService{
		organizationRepository: repository.NewOrganizationRepository(s),
	}
}

func (os *OrganizationService) CreateOrganization(organization *model.Organization) error {
	err := os.organizationRepository.Save(organization)
	if err != nil {
		return err
	}

	return nil
}

func (os *OrganizationService) GetOrganization(userId int, orgId int) (*model.Organization, error) {
	org, err := os.organizationRepository.Find(userId, orgId)
	if err != nil {
		return nil, err
	}

	return org, nil
}

func (os *OrganizationService) CreateOrganizationList() ([]*model.Organization, error) {
	return nil, nil
}

func (os *OrganizationService) UpdateOrganization() (*model.Organization, error) {
	return nil, nil
}

func (os *OrganizationService) DeleteOrganization() (*model.Organization, error) {
	return nil, nil
}
