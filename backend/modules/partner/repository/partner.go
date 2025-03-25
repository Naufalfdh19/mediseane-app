package repository

import (
	"context"
	"database/sql"
	"montelukast/modules/partner/entity"
	queryparams "montelukast/modules/partner/query_params"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/transaction"
)

type PartnerRepo interface {
	IsPartnerExistsByName(c context.Context, name string) (bool, error)
	IsPartnerExistsByID(c context.Context, partnerID int) (bool, error)
	AddPartner(c context.Context, partner entity.Partner) error
	IsPharmacyExistsByPartnerID(c context.Context, partnerID int) (bool, error)
	DeletePartner(c context.Context, partnerID int) error
	UpdatePartner(c context.Context, partner entity.Partner) error
	GetPartners(c context.Context, queryParams queryparams.QueryParams, totalItems int) ([]entity.Partner, error)
	GetTotalPartners(c context.Context, queryParams queryparams.QueryParams) (int, error)
	GetPartner(c context.Context, partnerID int) (*entity.Partner, error)
}

type partnerRepoImpl struct {
	db *sql.DB
}

func NewPartnerRepo(dbConn *sql.DB) partnerRepoImpl {
	return partnerRepoImpl{
		db: dbConn,
	}
}

func (r partnerRepoImpl) IsPartnerExistsByName(c context.Context, name string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM partners WHERE name = $1 AND deleted_at IS NULL) FOR UPDATE `

	var exists bool
	err := r.db.QueryRow(query, name).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return exists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return exists, nil
}

func (r partnerRepoImpl) IsPartnerExistsByID(c context.Context, partnerID int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM partners WHERE id = $1 AND deleted_at IS NULL) FOR UPDATE `

	var exists bool
	err := r.db.QueryRow(query, partnerID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return exists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return exists, nil
}

func (r partnerRepoImpl) AddPartner(c context.Context, partner entity.Partner) error {
	query := `INSERT INTO partners (name, year_founded, active_days, start_hour, end_hour, is_active)
			  VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(query, partner.Name, partner.YearFounded, partner.ActiveDays, partner.StartHour, partner.EndHour, partner.IsActive)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r partnerRepoImpl) IsPharmacyExistsByPartnerID(c context.Context, partnerID int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM pharmacies WHERE partner_id = $1 AND deleted_at IS NULL) FOR UPDATE `

	var exists bool
	err := r.db.QueryRow(query, partnerID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return exists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return exists, nil
}

func (r partnerRepoImpl) UpdatePartner(c context.Context, partner entity.Partner) error {
	tx := transaction.ExtractTx(c)

	query := `UPDATE partners
				SET active_days = $2, start_hour = $3, end_hour = $4, is_active = $5, updated_at = NOW()
				WHERE id = $1`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(c, query, partner.ID, partner.ActiveDays, partner.StartHour, partner.EndHour, partner.IsActive)
	} else {
		_, err = r.db.ExecContext(c, query, partner.ID, partner.ActiveDays, partner.StartHour, partner.EndHour, partner.IsActive)
	}
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r partnerRepoImpl) DeletePartner(c context.Context, partnerID int) error {
	tx := transaction.ExtractTx(c)

	query := `UPDATE partners
				SET deleted_at = NOW()
				WHERE id = $1`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(c, query, partnerID)
	} else {
		_, err = r.db.ExecContext(c, query, partnerID)
	}
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r partnerRepoImpl) GetPartners(c context.Context, queryParams queryparams.QueryParams, totalItems int) ([]entity.Partner, error) {
	partners := []entity.Partner{}

	query := `SELECT id, name, year_founded, active_days, start_hour, end_hour, is_active 
				FROM partners p 
				WHERE deleted_at IS NULL`

	var params []any
	query += queryparams.AddQueryParams(&params, queryParams)

	rows, err := r.db.Query(query, params...)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	defer rows.Close()

	for rows.Next() {
		var partner entity.Partner
		err := rows.Scan(
			&partner.ID,
			&partner.Name,
			&partner.YearFounded,
			&partner.ActiveDays,
			&partner.StartHour,
			&partner.EndHour,
			&partner.IsActive,
		)
		if err != nil {
			return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		}
		partners = append(partners, partner)
	}

	return partners, nil
}

func (r partnerRepoImpl) GetTotalPartners(c context.Context, queryParams queryparams.QueryParams) (int, error) {
	queryParams.Limit = 0
	queryParams.Page = 0
	queryParams.SortBy = ""
	queryParams.Order = ""

	query := `SELECT count(*)
				FROM partners p 
				WHERE deleted_at IS NULL`

	var params []any
	query += queryparams.AddQueryParams(&params, queryParams)

	var totalPartner int
	err := r.db.QueryRow(query, params...).Scan(&totalPartner)
	if err != nil {
		return 0, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return totalPartner, nil
}

func (r partnerRepoImpl) GetPartner(c context.Context, partnerID int) (*entity.Partner, error) {
	query := `select id, name, year_founded, active_days, start_hour, end_hour, is_active
				from partners p 
				where id = $1`

	var partner entity.Partner
	err := r.db.QueryRow(query, partnerID).Scan(
		&partner.ID,
		&partner.Name,
		&partner.YearFounded,
		&partner.ActiveDays,
		&partner.StartHour,
		&partner.EndHour,
		&partner.IsActive,
	)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return &partner, nil
}
