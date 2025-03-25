package repository

import (
	"context"
	"database/sql"
	"montelukast/modules/pharmacist/entity"
	queryparams "montelukast/modules/pharmacist/query_params"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/transaction"
	"os"
)

type PharmacistRepo interface {
	GetUserByEmail(c context.Context, email string) (entity.Pharmacist, error)
	IsUserExistByEmail(c context.Context, email string) (bool, error)
	AddPharmacist(c context.Context, pharmacists *entity.Pharmacist) error
	AddPharmacistDetail(c context.Context, pharmacists entity.Pharmacist) error
	IsPharmacistExistsByID(c context.Context, pharmacistID int) (bool, error)
	IsPharmacistExistsByEmail(c context.Context, email string) (bool, error)
	IsPharmacistExistsBySipaNumber(c context.Context, sipaNumber string) (bool, error)
	IsPharmacistExistsByPhoneNumber(c context.Context, phoneNumber string) (bool, error)
	IsPhoneNumberExists(c context.Context, pharmacistID int, phoneNumber string) (bool, error)
	IsPharmacyExistsByID(c context.Context, pharmacyID int) (bool, error)
	GetPharmacyNameByPharmacyID(c context.Context, pharmacyID int) (string, error)
	UpdatePharmacist(c context.Context, pharmacist entity.Pharmacist) error
	UpdatePharmacistPhoto(c context.Context, url string, pharmacistID int) error
	DeletePharmacist(c context.Context, pharmacyID int) error
	DeletePharmacistDetail(c context.Context, pharmacyID int) error
	GetPharmacistByID(c context.Context, id int) (*entity.Pharmacist, error)
	GetPharmacists(c context.Context, queryParams queryparams.QueryParams, totalItem int) ([]entity.Pharmacist, error)
	GetTotalPharmacist(c context.Context, queryParams queryparams.QueryParams) (int, error)
	GetPharmacyIDByPharmacistID(c context.Context, pharmacistID int) (*int, error)
	GetTotalPharmacistByPharmacyID(c context.Context, pharmacyID int) (int, error)
}

type pharmacistRepoImpl struct {
	db *sql.DB
}

func NewPharmacistsRepo(dbConn *sql.DB) pharmacistRepoImpl {
	return pharmacistRepoImpl{
		db: dbConn,
	}
}

func (r pharmacistRepoImpl) IsUserExistByEmail(c context.Context, email string) (bool, error) {
	checkQuery := `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL) `

	var exists bool
	err := r.db.QueryRowContext(c, checkQuery, email).Scan(&exists)
	if err != nil {
		return exists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return exists, nil
}

func (r pharmacistRepoImpl) GetUserByEmail(c context.Context, email string) (admin entity.Pharmacist, err error) {
	query := `SELECT id, password,role FROM users where email = $1 AND deleted_at IS NULL`
	err = r.db.QueryRowContext(c, query, email).Scan(&admin.ID, &admin.Password, &admin.Role)
	if err != nil {
		return admin, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	admin.Email = email
	return admin, nil
}

func (r pharmacistRepoImpl) IsPharmacistExistsByID(c context.Context, pharmacistID int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1 AND deleted_at IS NULL) FOR UPDATE `

	var exists bool
	err := r.db.QueryRow(query, pharmacistID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return exists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return exists, nil
}

func (r pharmacistRepoImpl) IsPharmacistExistsByEmail(c context.Context, email string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL) FOR UPDATE `

	var exists bool
	err := r.db.QueryRow(query, email).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return exists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return exists, nil
}

func (r pharmacistRepoImpl) IsPharmacistExistsBySipaNumber(c context.Context, sipaNumber string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM pharmacist_details WHERE sipa_number = $1 AND deleted_at IS NULL) FOR UPDATE `

	var exists bool
	err := r.db.QueryRow(query, sipaNumber).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return exists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return exists, nil
}

func (r pharmacistRepoImpl) IsPharmacistExistsByPhoneNumber(c context.Context, phoneNumber string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM pharmacist_details WHERE phone_number = $1 AND deleted_at IS NULL) FOR UPDATE `

	var exists bool
	err := r.db.QueryRow(query, phoneNumber).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return exists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return exists, nil
}

func (r pharmacistRepoImpl) IsPharmacyExistsByID(c context.Context, pharmacyID int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM pharmacies WHERE id = $1 AND deleted_at IS NULL) FOR UPDATE `

	var exists bool
	err := r.db.QueryRow(query, pharmacyID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return exists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return exists, nil
}

func (r pharmacistRepoImpl) IsPhoneNumberExists(c context.Context, pharmacistID int, phoneNumber string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM pharmacist_details WHERE pharmacist_id != $1 AND phone_number = $2 AND deleted_at IS NULL) FOR UPDATE `

	var exists bool
	err := r.db.QueryRow(query, pharmacistID, phoneNumber).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return exists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return exists, nil
}

func (r pharmacistRepoImpl) GetPharmacyNameByPharmacyID(c context.Context, pharmacyID int) (string, error) {
	tx := transaction.ExtractTx(c)

	query := `SELECT name FROM pharmacies WHERE id = $1 AND deleted_at IS NULL`

	var err error
	var pharmacyName string
	if tx != nil {
		err = tx.QueryRowContext(c, query, pharmacyID).Scan(&pharmacyName)
	} else {
		err = r.db.QueryRowContext(c, query, pharmacyID).Scan(&pharmacyName)
	}
	if err != nil {
		return "", apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return pharmacyName, nil
}

func (r pharmacistRepoImpl) AddPharmacist(c context.Context, pharmacists *entity.Pharmacist) error {
	tx := transaction.ExtractTx(c)
	profilePhoto := os.Getenv("DEFAULT_PROFILE")

	query := `INSERT INTO users (name, email, password, profile_photo, role, is_verified)
			  VALUES ($1, $2, $3, $4, $5, $6)
			  RETURNING id`

	var err error
	if tx != nil {
		err = tx.QueryRowContext(c, query, pharmacists.Name, pharmacists.Email, pharmacists.Password, profilePhoto, appconstant.ROLE_PHARMACY, true).Scan(&pharmacists.ID)
	} else {
		err = r.db.QueryRowContext(c, query, pharmacists.Name, pharmacists.Email, pharmacists.Password, profilePhoto, appconstant.ROLE_PHARMACY, true).Scan(&pharmacists.ID)
	}
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r pharmacistRepoImpl) AddPharmacistDetail(c context.Context, pharmacists entity.Pharmacist) error {
	tx := transaction.ExtractTx(c)

	query := `INSERT INTO pharmacist_details (pharmacist_id, sipa_number, phone_number, year_of_experience)
			  VALUES ($1, $2, $3, $4)`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(c, query, pharmacists.ID, pharmacists.SipaNumber, pharmacists.PhoneNumber, pharmacists.YearOfExperience)
	} else {
		_, err = r.db.ExecContext(c, query, pharmacists.ID, pharmacists.SipaNumber, pharmacists.PhoneNumber, pharmacists.YearOfExperience)
	}
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return nil
}

func (r pharmacistRepoImpl) UpdatePharmacist(c context.Context, pharmacist entity.Pharmacist) error {
	query := `UPDATE pharmacist_details
				SET pharmacy_id = $2, phone_number = $3, year_of_experience = $4, updated_at = NOW()
				WHERE pharmacist_id = $1 AND deleted_at IS NULL`

	_, err := r.db.Exec(query, pharmacist.ID, pharmacist.PharmacyID, pharmacist.PhoneNumber, pharmacist.YearOfExperience)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return nil
}

func (r pharmacistRepoImpl) UpdatePharmacistPhoto(c context.Context, url string, pharmacistID int) error {
	query := `UPDATE users
				SET profile_photo = $2, updated_at = NOW()
				WHERE id = $1  AND deleted_at IS NULL`

	_, err := r.db.Exec(query, pharmacistID, url)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return nil
}

func (r pharmacistRepoImpl) DeletePharmacist(c context.Context, pharmacyID int) error {
	tx := transaction.ExtractTx(c)

	query := `UPDATE users
				SET deleted_at = NOW()
				WHERE id = $1`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(c, query, pharmacyID)
	} else {
		_, err = r.db.ExecContext(c, query, pharmacyID)
	}
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return nil
}

func (r pharmacistRepoImpl) DeletePharmacistDetail(c context.Context, pharmacyID int) error {
	tx := transaction.ExtractTx(c)

	query := `UPDATE pharmacist_details
				SET deleted_at = NOW()
				WHERE pharmacist_id = $1`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(c, query, pharmacyID)
	} else {
		_, err = r.db.ExecContext(c, query, pharmacyID)
	}
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return nil
}

func (r pharmacistRepoImpl) GetPharmacistByID(c context.Context, id int) (*entity.Pharmacist, error) {
	query := `SELECT u.id, p.id, p.name, u.name, email, sipa_number, phone_number, year_of_experience
				FROM users u 
				JOIN pharmacist_details pd ON pd.pharmacist_id = u.id AND pd.deleted_at IS null
				LEFT JOIN pharmacies p ON p.id = pd.pharmacy_id AND p.deleted_at IS NULL
				WHERE u.id = $1 AND u.deleted_at IS NULL;`

	var pharmacist entity.Pharmacist
	err := r.db.QueryRowContext(c, query, id).Scan(
		&pharmacist.ID,
		&pharmacist.PharmacyID,
		&pharmacist.PharmacyName,
		&pharmacist.Name,
		&pharmacist.Email,
		&pharmacist.SipaNumber,
		&pharmacist.PhoneNumber,
		&pharmacist.YearOfExperience,
	)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return &pharmacist, nil
}

func (r pharmacistRepoImpl) GetPharmacists(c context.Context, queryParams queryparams.QueryParams, totalItems int) ([]entity.Pharmacist, error) {
	pharmacists := []entity.Pharmacist{}

	query := `SELECT u.id, p.id, p.name, u.name, email, sipa_number, phone_number, year_of_experience
				FROM users u 
				JOIN pharmacist_details pd ON pd.pharmacist_id = u.id AND pd.deleted_at IS null
				LEFT JOIN pharmacies p ON p.id = pd.pharmacy_id AND p.deleted_at IS NULL
				WHERE u.deleted_at IS NULL`

	var params []any
	query += queryparams.AddQueryParams(&params, queryParams)

	rows, err := r.db.Query(query, params...)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	defer rows.Close()

	for rows.Next() {
		var pharmacist entity.Pharmacist
		err := rows.Scan(
			&pharmacist.ID,
			&pharmacist.PharmacyID,
			&pharmacist.PharmacyName,
			&pharmacist.Name,
			&pharmacist.Email,
			&pharmacist.SipaNumber,
			&pharmacist.PhoneNumber,
			&pharmacist.YearOfExperience,
		)
		if err != nil {
			return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		}
		pharmacists = append(pharmacists, pharmacist)
	}

	return pharmacists, nil
}

func (r pharmacistRepoImpl) GetTotalPharmacist(c context.Context, queryParams queryparams.QueryParams) (int, error) {
	queryParams.Limit = 0
	queryParams.Page = 0
	queryParams.SortBy = ""
	queryParams.Order = ""

	query := `SELECT count(*)
				FROM users u 
				JOIN pharmacist_details pd ON pd.pharmacist_id = u.id AND pd.deleted_at IS null
				LEFT JOIN pharmacies p ON p.id = pd.pharmacy_id AND p.deleted_at IS NULL
				WHERE u.deleted_at IS NULL`

	var params []any
	query += queryparams.AddQueryParams(&params, queryParams)

	var totalPharmacist int
	err := r.db.QueryRow(query, params...).Scan(&totalPharmacist)
	if err != nil {
		return 0, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return totalPharmacist, nil
}

func (r pharmacistRepoImpl) GetPharmacyIDByPharmacistID(c context.Context, pharmacistID int) (*int, error) {

	query := `SELECT pharmacy_id
				FROM pharmacist_details pd 
				WHERE pharmacist_id = $1 AND deleted_at IS NULL`

	var pharmacyID *int
	err := r.db.QueryRow(query, pharmacistID).Scan(&pharmacyID)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return pharmacyID, nil
}

func (r pharmacistRepoImpl) GetTotalPharmacistByPharmacyID(c context.Context, pharmacyID int) (int, error) {
	tx := transaction.ExtractTx(c)

	query := `SELECT COUNT(pharmacy_id)
				FROM pharmacist_details pd 
				WHERE pharmacy_id = $1
				GROUP BY pharmacy_id`

	var totalPharmacist int
	var err error
	if tx != nil {
		err = tx.QueryRowContext(c, query, pharmacyID).Scan(&totalPharmacist)
	} else {
		err = r.db.QueryRowContext(c, query, pharmacyID).Scan(&totalPharmacist)
	}
	if err != nil {
		return 0, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return totalPharmacist, nil
}
