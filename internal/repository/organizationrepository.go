package repository

import (
	"fmt"
	"github.com/lib/pq"
	"net/http"
	"relay-backend/internal/enums"
	"relay-backend/internal/model"
	"relay-backend/internal/store"
	"relay-backend/internal/utils/exception"
	"strings"
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

	rows, err := or.store.Db.Query(
		"select id, first_name, last_name, patronymic, email from users join employees on users.id = employees.user_id where employees.organization_id = $1",
		orgId,
	)
	employeesList := make([]model.User, 0)

	for rows.Next() {
		empl := &model.User{}

		if err := rows.Scan(
			&empl.Id,
			&empl.FirstName,
			&empl.LastName,
			&empl.Patronymic,
			&empl.Email,
		); err != nil {
			return nil, err
		}

		employeesList = append(employeesList, *empl)
	}

	if err != nil {
		return nil, err
	}

	org.Employees = employeesList

	return org, nil
}
func (or *OrganizationRepository) GetList(userId int, page int, pageSize int) ([]model.Organization, error) {
	orgList := make([]model.Organization, 0)

	rows, err := or.store.Db.Query("select * from organizations where owner_id = $1 order by id desc limit $2 offset $3",
		userId,
		pageSize,
		(page-1)*pageSize,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		org := &model.Organization{}

		if err := rows.Scan(
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

		orgList = append(orgList, *org)
	}

	return orgList, nil
}

func (or *OrganizationRepository) Update(userId int, orgId int, organization *model.Organization) error {
	err := or.store.Db.QueryRow(`
			update organizations 
			set name = coalesce(nullif($1, ''), name), 
				description = coalesce(nullif($2, ''), description),
				address = coalesce(nullif($3, ''), address), 
				email = coalesce(nullif($4, ''), email) 
			where owner_id = $5 and id = $6
			returning *
		`,
		&organization.Name,
		&organization.Description,
		&organization.Address,
		&organization.Email,
		userId,
		orgId,
	).Scan(
		&organization.Id,
		&organization.OwnerId,
		&organization.Name,
		&organization.Description,
		&organization.Address,
		&organization.Email,
		&organization.CreationDate,
	)

	if err != nil {
		return err
	}

	return nil
}

func (or *OrganizationRepository) Delete(ownerId int, orgId int) error {
	res, err := or.store.Db.Exec(`delete from organizations where owner_id = $1 and id = $2`, ownerId, orgId)
	if err != nil {
		return exception.NewException(http.StatusInternalServerError, exception.Enum.InternalServerError)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return exception.NewException(http.StatusInternalServerError, exception.Enum.InternalServerError)
	}

	if count == 0 {
		return exception.NewException(http.StatusNotFound, exception.Enum.NotFound)
	}

	return nil
}

func (or *OrganizationRepository) AddEmployees(userId int, orgId int, employees []model.Employee) error {
	_, err := or.Find(userId, orgId)
	if err != nil {
		return exception.NewDetailsException(http.StatusNotFound, exception.Enum.OrganizationNotFound, map[string]interface{}{"id": orgId})
	}

	tx, err := or.store.Db.Begin()
	if err != nil {
		return err
	}

	for _, employee := range employees {
		if employee.Id == userId {
			continue
		}

		if employee.UserRole == enums.UserRoleEnum.OrganizationOwner {
			tx.Rollback()
			return exception.NewDetailsException(
				http.StatusBadRequest,
				exception.Enum.BadRequest,
				map[string]interface{}{
					"message": fmt.Sprintf("Can not set userRole OrganizationOwner for employees. Employee id: %d", employee.Id),
				},
			)
		}

		_, err := tx.Exec("insert into employees (organization_id, user_id, user_role) values ($1, $2, $3)", orgId, employee.Id, employee.UserRole)

		if err != nil {
			tx.Rollback()

			if err, ok := err.(*pq.Error); ok {
				return exception.NewException(http.StatusBadRequest, err.Detail)
			}

			if strings.Contains(err.Error(), "employees_user_id_fkey") {
				return exception.NewDetailsException(http.StatusBadRequest, exception.Enum.UserNotFound, map[string]interface{}{"id": employee.Id})
			}

			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (or *OrganizationRepository) DeleteAllEmployees(ownerId int, orgId int) error {
	_, err := or.Find(ownerId, orgId)
	if err != nil {
		return exception.NewException(http.StatusNotFound, exception.Enum.NotFound)
	}

	_, err = or.store.Db.Exec(`delete from employees where organization_id = $1`, orgId)
	if err != nil {
		return exception.NewException(http.StatusInternalServerError, exception.Enum.InternalServerError)
	}

	return nil
}

func (or *OrganizationRepository) GetUserRole(userId int, orgId int) enums.UserRole {
	var userRole int

	if _, err := or.Find(userId, orgId); err == nil {
		return enums.UserRoleEnum.OrganizationOwner
	}

	err := or.store.Db.QueryRow(
		`select user_role from employees where user_id = $1 and organization_id = $2`,
		userId,
		orgId,
	).Scan(
		userRole,
	)

	if err != nil {
		//TODO: consider about error. now returns 403 if org not found
		return enums.UserRoleEnum.None
	}

	return enums.UserRole(userRole)
}
