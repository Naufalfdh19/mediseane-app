package repository

import (
	"context"
	"database/sql"
	"montelukast/modules/user/entity"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/transaction"
)

type UserRepo interface {
	AddUser(c context.Context, user entity.User) error
	GetUserByID(c context.Context, userID int) (*entity.User, error)
	GetUserPasswordByEmail(c context.Context, email string) (string, error)
	GetUserByEmail(c context.Context, email string) (*entity.User, error)
	AddResetPassToken(c context.Context, userID int, token string) error
	GetUserIDBYResetPassToken(c context.Context, resetToken string) (int, error)
	UpdateUserPassword(c context.Context, newPassword string, userID int) error
	DeleteResetPasswordToken(c context.Context, token string) error
	IsUserExistsByEmail(c context.Context, email string) (bool, error)
	IsUserExistsByID(c context.Context, userID int) (bool, error)
	UpdateEmailVerifyStatus(c context.Context, userID int) error
	IsResetPasswordTokenExists(c context.Context, token string) (bool, error)
	IsVerifyEmailTokenExists(c context.Context, token string) (bool, error)
	UpdateTokenForVerifyEmail(c context.Context, userID int, token string) error
	DeleteVerifyEmailToken(c context.Context, token string) error
	UpdateNameByID(c context.Context, user entity.User) error
}

type UserRepoImpl struct {
	db *sql.DB
}

func NewUserRepo(dbConn *sql.DB) UserRepoImpl {
	return UserRepoImpl{
		db: dbConn,
	}
}

func (r UserRepoImpl) AddUser(c context.Context, user entity.User) error {
	query := `INSERT INTO users (name, email, password, profile_photo, role, is_verified)
			  VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(query, user.Name, user.Email, user.Password, appconstant.DefaultProfileIMG, user.Role, user.IsVerified)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r UserRepoImpl) GetUserByID(c context.Context, userID int) (*entity.User, error) {
	query := `SELECT id, name, email, profile_photo, role, is_verified
			  FROM users 
			  WHERE id = $1 AND deleted_at IS NULL`

	var userResp entity.User
	err := r.db.QueryRowContext(c, query, userID).Scan(
		&userResp.ID,
		&userResp.Name,
		&userResp.Email,
		&userResp.ProfilePhoto,
		&userResp.Role,
		&userResp.IsVerified,
	)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return &userResp, nil
}

func (r UserRepoImpl) GetUserPasswordByEmail(c context.Context, email string) (string, error) {
	query := `SELECT password
			  FROM users 
			  WHERE email = $1 AND deleted_at IS NULL`

	var hashPass string
	err := r.db.QueryRow(query, email).Scan(
		&hashPass,
	)
	if err != nil {
		return "", apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return hashPass, nil
}

func (r UserRepoImpl) IsUserExistsByEmail(c context.Context, email string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL) FOR UPDATE `

	var exists bool
	err := r.db.QueryRow(query, email).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return exists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return exists, nil
}

func (r UserRepoImpl) IsUserExistsByID(c context.Context, userID int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1 AND deleted_at IS NULL) FOR UPDATE `

	var exists bool
	err := r.db.QueryRowContext(c, query, userID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return exists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return exists, nil
}

func (r UserRepoImpl) GetUserByEmail(c context.Context, email string) (*entity.User, error) {
	query := `SELECT id, name, email, password, profile_photo, role, is_verified
			  FROM users
			  WHERE email = $1 AND deleted_at IS NULL`

	var userResp entity.User
	err := r.db.QueryRow(query, email).Scan(
		&userResp.ID,
		&userResp.Name,
		&userResp.Email,
		&userResp.Password,
		&userResp.ProfilePhoto,
		&userResp.Role,
		&userResp.IsVerified,
	)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return &userResp, nil
}

func (r UserRepoImpl) AddResetPassToken(c context.Context, userID int, token string) error {
	query := `INSERT INTO reset_password_tokens (user_id, token)
			  VALUES ($1, $2)`

	var err error
	tx := transaction.ExtractTx(c)
	if tx != nil {
		_, err = tx.ExecContext(c, query, userID, token)
	} else {
		_, err = r.db.ExecContext(c, query, userID, token)
	}
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r UserRepoImpl) IsResetPasswordTokenExists(c context.Context, token string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM reset_password_tokens WHERE token = $1 and deleted_at is null) FOR UPDATE `

	var exists bool
	var err error
	tx := transaction.ExtractTx(c)
	if tx != nil {
		err = tx.QueryRowContext(c, query, token).Scan(&exists)
	} else {
		err = r.db.QueryRowContext(c, query, token).Scan(&exists)
	}
	if err != nil && err != sql.ErrNoRows {
		return exists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return exists, nil
}

func (r UserRepoImpl) GetUserIDBYResetPassToken(c context.Context, resetToken string) (int, error) {
	tx := transaction.ExtractTx(c)

	query := `SELECT user_id
			  FROM reset_password_tokens
			  WHERE token = $1 and deleted_at is null`

	var userID int
	var err error
	if tx != nil {
		err = tx.QueryRowContext(c, query, resetToken).Scan(&userID)
	} else {
		err = r.db.QueryRowContext(c, query, resetToken).Scan(&userID)
	}
	if err != nil {
		return -1, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return userID, nil
}

func (r UserRepoImpl) UpdateUserPassword(c context.Context, newPassword string, userID int) error {
	tx := transaction.ExtractTx(c)

	query := `UPDATE users
				SET password = $1, updated_at = NOW()
				WHERE id = $2 AND deleted_at IS NULL`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(c, query, newPassword, userID)
	} else {
		_, err = r.db.ExecContext(c, query, newPassword, userID)
	}
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r UserRepoImpl) DeleteResetPasswordToken(c context.Context, token string) error {
	tx := transaction.ExtractTx(c)

	query := `UPDATE reset_password_tokens
				SET deleted_at = NOW()
				WHERE token = $1 AND deleted_at IS NULL`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(c, query, token)
	} else {
		_, err = r.db.ExecContext(c, query, token)
	}
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r UserRepoImpl) IsVerifyEmailTokenExists(c context.Context, token string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM verify_email_tokens WHERE token = $1 and deleted_at is null) FOR UPDATE `

	var exists bool
	var err error
	tx := transaction.ExtractTx(c)
	if tx != nil {
		err = tx.QueryRowContext(c, query, token).Scan(&exists)
	} else {
		err = r.db.QueryRowContext(c, query, token).Scan(&exists)
	}
	if err != nil && err != sql.ErrNoRows {
		return exists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return exists, nil
}

func (r UserRepoImpl) UpdateTokenForVerifyEmail(c context.Context, userID int, token string) error {
	query := `INSERT INTO verify_email_tokens (user_id, token)
				VALUES ($1, $2)`

	var err error
	tx := transaction.ExtractTx(c)
	if tx != nil {
		_, err = tx.ExecContext(c, query, userID, token)
	} else {
		_, err = r.db.ExecContext(c, query, userID, token)
	}
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r UserRepoImpl) UpdateEmailVerifyStatus(c context.Context, userID int) error {
	tx := transaction.ExtractTx(c)

	query := `UPDATE users
				SET is_verified = true
				WHERE id = $1 AND deleted_at IS NULL`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(c, query, userID)
	} else {
		_, err = r.db.ExecContext(c, query, userID)
	}
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r UserRepoImpl) DeleteVerifyEmailToken(c context.Context, token string) error {
	tx := transaction.ExtractTx(c)

	query := `UPDATE verify_email_tokens
				SET deleted_at = NOW()
				WHERE token = $1 AND deleted_at IS NULL`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(c, query, token)
	} else {
		_, err = r.db.ExecContext(c, query, token)
	}
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r UserRepoImpl) UpdateNameByID(c context.Context, user entity.User) error {
	query := `UPDATE 
				users
			SET 
				name = $2,
				updated_at = NOW() 
			WHERE 
				id = $1 AND deleted_at IS NULL`

	_, err := r.db.ExecContext(c, query, user.ID, user.Name)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}
