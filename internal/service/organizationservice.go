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

func (os *OrganizationService) GetOrganizationList(userId int, page int, pageSize int) ([]model.Organization, error) {
	orgList, err := os.organizationRepository.GetList(userId, page, pageSize)
	if err != nil {
		return nil, err
	}

	return orgList, nil
}

func (os *OrganizationService) UpdateOrganization(ownerId int, orgId int, organization *model.Organization) error {
	err := os.organizationRepository.Update(ownerId, orgId, organization)
	if err != nil {
		return err
	}

	return nil
}

func (os *OrganizationService) DeleteOrganization() (*model.Organization, error) {
	return nil, nil
}

func (os *OrganizationService) AddOrganizationEmployees(userId int, orgId int, employeeIds []int) error {
	err := os.organizationRepository.AddEmployees(userId, orgId, employeeIds)

	if err != nil {
		return err
	}

	return nil
}
