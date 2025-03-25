package repository

import (
	"context"
	"database/sql"
	queryparams "montelukast/modules/admin/query_params"
	"montelukast/modules/user/entity"
	adminEntity "montelukast/modules/admin/entity"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
)

type AdminRepo interface {
	GetUserByEmail(c context.Context, email string) (entity.User, error)
	IsUserExistByEmail(c context.Context, email string) (bool, error)
	GetUsers(c context.Context, queryParams queryparams.QueryParams, totalItems int) ([]adminEntity.User, error)
	GetTotalUser(c context.Context, queryParams queryparams.QueryParams) (int, error)
}

type adminRepoImpl struct {
	db *sql.DB
}

func NewAdminRepository(dbConn *sql.DB) adminRepoImpl {
	return adminRepoImpl{
		db: dbConn,
	}
}
func (r adminRepoImpl) IsUserExistByEmail(c context.Context, email string) (bool, error) {
	checkQuery := `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL) `

	var exists bool
	err := r.db.QueryRowContext(c, checkQuery, email).Scan(&exists)
	if err != nil {
		return exists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return exists, nil
}

func (r adminRepoImpl) GetUserByEmail(c context.Context, email string) (admin entity.User, err error) {
	query := `SELECT id, password, role FROM users where email = $1 AND deleted_at IS NULL`
	err = r.db.QueryRowContext(c, query, email).Scan(&admin.ID, &admin.Password, &admin.Role)
	if err != nil {
		return admin, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	admin.Email = email
	return admin, nil
}



func (r adminRepoImpl) GetUsers(c context.Context, queryParams queryparams.QueryParams, totalItems int) ([]adminEntity.User, error) {
	users := []adminEntity.User{}

	query := `SELECT id, name, email, profile_photo, role
				FROM users u 
				WHERE deleted_at IS NULL AND role != 'pharmacist'`

	var params []any
	query += queryparams.AddQueryParams(&params, queryParams)

	rows, err := r.db.Query(query, params...)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	defer rows.Close()

	for rows.Next() {
		var user adminEntity.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.ProfilePhoto,
			&user.Role,
		)
		if err != nil {
			return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (r adminRepoImpl) GetTotalUser(c context.Context, queryParams queryparams.QueryParams) (int, error) {
	queryParams.Limit = 0
	queryParams.Page = 0
	queryParams.SortBy = ""
	queryParams.Order = ""

	query := `SELECT count(*)
				FROM users u 
				WHERE deleted_at IS NULL AND role != 'pharmacist'`

	var params []any
	query += queryparams.AddQueryParams(&params, queryParams)

	var totalItems int
	err := r.db.QueryRow(query, params...).Scan(&totalItems)
	if err != nil {
		return 0, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return totalItems, nil
}