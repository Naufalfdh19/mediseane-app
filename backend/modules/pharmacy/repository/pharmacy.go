package repository

import (
	"context"
	"database/sql"
	"log"
	"montelukast/modules/pharmacy/entity"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
)

type PharmacyRepository interface {
	AddPharmacy(c context.Context, pharmacy entity.Pharmacy) (err error)
	UpdatePharmacy(c context.Context, pharmacy entity.Pharmacy) (err error)
	DeletePharmacy(c context.Context, id int) (err error)
	IsPharmacyExists(c context.Context, id int) (bool, error)
	IsPharmacistExists(c context.Context, id int) (bool, error)
	GetAllPharmacies(c context.Context, filter entity.PharmacyFilter) (pharmacies *entity.PaginatedPharmacies, err error)
	GetTotalItem(c context.Context, filter entity.PharmacyFilterCount) (int, error)
	AddLogo(c context.Context, url string, id int) (err error)
	GetPharmacyByID(c context.Context, id int) (pharmacy entity.Pharmacy, err error)
}

type pharmacyRepository struct {
	db *sql.DB
}

func NewPharmacyRepository(db *sql.DB) PharmacyRepository {
	return &pharmacyRepository{db: db}
}
func (r *pharmacyRepository) GetPharmacyByID(c context.Context, id int) (pharmacy entity.Pharmacy, err error) {
	query := `SELECT p.id,q.id,q.name,p.name,p.address,
	p.province_id,
	p.province,
	p.city_id,
	p.city,
	p.district_id,
	p.district,
	p.sub_district_id,
	p.sub_district,
	ST_X(location::geometry), 
	ST_Y(location::geometry),
	p.postal_code,
	p.is_active,
	p.logo,
	p.updated_at
	FROM pharmacies p
	JOIN partners q
	ON p.partner_id = q.id
	WHERE
		p.id = $1
	AND
		p.deleted_at IS NULL
	AND
		q.deleted_at IS NULL`
	err = r.db.QueryRowContext(c, query, id).Scan(&pharmacy.ID, &pharmacy.PartnerID, &pharmacy.PartnerName, &pharmacy.Name,
		&pharmacy.Address, &pharmacy.ProvinceID, &pharmacy.Province, &pharmacy.CityID, &pharmacy.City,
		&pharmacy.DistrictID, &pharmacy.District, &pharmacy.SubDistrictID,
		&pharmacy.SubDistrict, &pharmacy.Longitude, &pharmacy.Latitude,
		&pharmacy.PostalCode, &pharmacy.IsActive, &pharmacy.Logo,
		&pharmacy.UpdatedAt)
	if err != nil {
		return pharmacy, apperror.NewErrInternalServerError(appconstant.FieldErrPharmacy, apperror.ErrInternalServer, err)
	}
	return pharmacy, nil
}

func (r *pharmacyRepository) GetAllPharmacies(c context.Context, filter entity.PharmacyFilter) (pharmacies *entity.PaginatedPharmacies, err error) {
	query := ` SELECT p.id,q.id,q.name,p.name,p.address,
	p.province_id,
	p.province,
	p.city_id,
	p.city,
	p.district_id,
	p.district,
	p.sub_district_id,
	p.sub_district,
	ST_X(location::geometry), 
	ST_Y(location::geometry),
	p.postal_code,
	p.logo,
	p.updated_at
	FROM pharmacies p
	JOIN partners q
	ON p.partner_id = q.id
	WHERE
		p.deleted_at IS NULL
	AND
		q.deleted_at IS NULL
	`
	result := entity.PaginatedPharmacies{}
	paramQuery, sortPage, args := PharmacyParam(filter)
	finalQuery := query + paramQuery + sortPage
	rows, err := r.db.QueryContext(c, finalQuery, args...)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrPharmacy, apperror.ErrInternalServer, err)

	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal()
		}
	}()
	for rows.Next() {
		var pharmacy entity.Pharmacy
		err := rows.Scan(&pharmacy.ID, &pharmacy.PartnerID, &pharmacy.PartnerName, &pharmacy.Name,
			&pharmacy.Address, &pharmacy.ProvinceID, &pharmacy.Province, &pharmacy.CityID, &pharmacy.City,
			&pharmacy.DistrictID, &pharmacy.District, &pharmacy.SubDistrictID,
			&pharmacy.SubDistrict, &pharmacy.Latitude, &pharmacy.Longitude,
			&pharmacy.PostalCode, &pharmacy.Logo,
			&pharmacy.UpdatedAt)
		if err != nil {
			return nil, apperror.NewErrInternalServerError(appconstant.FieldErrPharmacy, apperror.ErrInternalServer, err)
		}
		result.Pharmacies = append(result.Pharmacies, pharmacy)
	}
	err = rows.Err()
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrPharmacy, apperror.ErrInternalServer, err)
	}
	return &result, nil
}

func (r pharmacyRepository) GetTotalItem(c context.Context, filter entity.PharmacyFilterCount) (int, error) {
	query := `SELECT COUNT(*)
				FROM pharmacies p
				WHERE p.deleted_at IS NULL`
	paramQuery, args := PharmacyParamCount(filter)
	var totalItem int
	err := r.db.QueryRowContext(c, query+paramQuery, args...).Scan(&totalItem)
	if err != nil {
		return 0, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return totalItem, nil
}

func (r *pharmacyRepository) AddPharmacy(c context.Context, pharmacy entity.Pharmacy) (err error) {
	query := `INSERT INTO pharmacies (partner_id, name, address,location,
			province, province_id, city, city_id, district, 
			district_id, sub_district, sub_district_id,postal_code,
			logo)
			VALUES ($1, $2,$3,ST_SetSRID(ST_MakePoint($4, $5),4326),
			$6,$7,$8,$9,$10,$11,$12,$13,$14,$15) RETURNING id`
	_, err = r.db.ExecContext(c, query,
		pharmacy.PartnerID, pharmacy.Name, pharmacy.Address,
		pharmacy.Longitude, pharmacy.Latitude,
		pharmacy.Province, pharmacy.ProvinceID,
		pharmacy.City, pharmacy.CityID,
		pharmacy.District, pharmacy.DistrictID,
		pharmacy.SubDistrict, pharmacy.SubDistrictID, pharmacy.PostalCode,
		appconstant.DefaultProfileIMG)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r *pharmacyRepository) UpdatePharmacy(c context.Context, pharmacy entity.Pharmacy) (err error) {
	query := `UPDATE 
				pharmacies 
			SET
				partner_id = $2, name = $3, 
				address = $4 , 
				location = ST_SetSRID(ST_MakePoint($5, $6),4326),
					province = $7, province_id = $8 , 
					city = $9, city_id = $10 , 
					district = $11, district_id = $12, 
					sub_district = $13, sub_district_id = $14 ,
				postal_code = $15,
				is_active = $16,
				updated_at = $17
			WHERE id = $1
				AND
			deleted_at IS NULL`
	_, err = r.db.ExecContext(c, query,
		pharmacy.ID,
		pharmacy.PartnerID,
		pharmacy.Name,
		pharmacy.Address,
		pharmacy.Longitude,
		pharmacy.Latitude,
		pharmacy.Province,
		pharmacy.ProvinceID,
		pharmacy.City,
		pharmacy.CityID,
		pharmacy.District,
		pharmacy.DistrictID,
		pharmacy.SubDistrict,
		pharmacy.SubDistrictID,
		pharmacy.PostalCode,
		pharmacy.IsActive,
		"NOW()",
	)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r *pharmacyRepository) IsPharmacistExists(c context.Context, id int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM pharmacist_details WHERE pharmacy_id = $1 AND deleted_at IS NULL)`
	var exists bool
	err := r.db.QueryRowContext(c, query, id).Scan(&exists)
	if err != nil {
		return exists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return exists, nil
}

func (r *pharmacyRepository) IsPharmacyExists(c context.Context, id int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM pharmacies WHERE id = $1 AND deleted_at IS NULL) FOR UPDATE `
	var exists bool
	err := r.db.QueryRowContext(c, query, id).Scan(&exists)
	if err != nil {
		return exists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return exists, nil
}

func (r *pharmacyRepository) DeletePharmacy(c context.Context, id int) (err error) {
	query := `UPDATE pharmacies 
			SET
				updated_at = NOW(),
				deleted_at = NOW()
			WHERE
				id = $1`
	_, err = r.db.ExecContext(c, query, id)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r *pharmacyRepository) AddLogo(c context.Context, url string, id int) (err error) {
	query := `UPDATE pharmacies
				SET 
					logo = $2,
					updated_at = NOW()
				WHERE id = $1
					AND 
				deleted_at IS NULL`
	_, err = r.db.ExecContext(c, query, id, url)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}
