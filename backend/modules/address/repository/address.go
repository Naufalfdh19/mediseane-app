package repository

import (
	"context"
	"database/sql"
	"montelukast/modules/address/entity"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/transaction"
)

type AddressRepo interface {
	GetProvinces(c context.Context) ([]entity.Location, error)
	GetCitiesByProvinceID(c context.Context, provinceID int) ([]entity.Location, error)
	GetDistrictsByCityID(c context.Context, cityID int) ([]entity.Location, error)
	GetSubDistrictsByDistrictID(c context.Context, districtID int) ([]entity.Location, error)
	IsUserActiveAddressExists(c context.Context, userID int) (bool, error)
	AddUserAddress(c context.Context, userAddress entity.UserAddress) error
	GetCurrentLocationByLongAndLat(c context.Context, longitude string, latitude string) (*entity.UserAddress, error)
	GetAddressesByUserID(c context.Context, userID int, filter entity.AddressFilter) ([]entity.UserAddress, error)
	IsUserAddressExists(c context.Context, addressID int, userID int) (bool, error)
	GetActiveAddressByUserID(c context.Context, userID int) (int, error)
	DeactivateAddressByID(c context.Context, addressID int) error
	UpdateUserAddress(c context.Context, address entity.UserAddress) error
	GetAddressByID(c context.Context, addressID int) (*entity.UserAddress, error)
	DeleteUserAddressByID(c context.Context, addressID int) error
}

type addressRepoImpl struct {
	db *sql.DB
}

func NewAddressRepo(dbConn *sql.DB) *addressRepoImpl {
	return &addressRepoImpl{
		db: dbConn,
	}
}

func (r *addressRepoImpl) GetProvinces(c context.Context) ([]entity.Location, error) {
	query := `SELECT 
							id, 
							name 
						FROM 
							provinces 
						WHERE 
							deleted_at is NULL`

	rows, err := r.db.QueryContext(c, query)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	defer rows.Close()

	provinces := []entity.Location{}

	for rows.Next() {
		var province entity.Location
		err := rows.Scan(
			&province.ID,
			&province.Name,
		)
		if err != nil {
			return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		}
		provinces = append(provinces, province)
	}

	return provinces, nil
}

func (r *addressRepoImpl) GetCitiesByProvinceID(c context.Context, provinceID int) ([]entity.Location, error) {
	query := `SELECT 
							id, 
							name, 
							ST_X(location::geometry), 
							ST_Y(location::geometry) 
						FROM 
							cities 
						WHERE 
							province_id = $1
						AND
							deleted_at is NULL`

	rows, err := r.db.QueryContext(c, query, provinceID)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	defer rows.Close()

	cities := []entity.Location{}

	for rows.Next() {
		var city entity.Location
		err := rows.Scan(
			&city.ID,
			&city.Name,
			&city.Longitude,
			&city.Latitude,
		)
		if err != nil {
			return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		}
		cities = append(cities, city)
	}

	return cities, nil
}

func (r *addressRepoImpl) GetDistrictsByCityID(c context.Context, cityID int) ([]entity.Location, error) {
	query := `SELECT 
							id, 
							name, 
							ST_X(location::geometry), 
							ST_Y(location::geometry) 
						FROM 
							districts 
						WHERE 
							city_id = $1
						AND
							deleted_at is NULL`

	rows, err := r.db.QueryContext(c, query, cityID)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	defer rows.Close()

	districts := []entity.Location{}

	for rows.Next() {
		var district entity.Location
		err := rows.Scan(
			&district.ID,
			&district.Name,
			&district.Longitude,
			&district.Latitude,
		)
		if err != nil {
			return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		}
		districts = append(districts, district)
	}

	return districts, nil
}

func (r *addressRepoImpl) GetSubDistrictsByDistrictID(c context.Context, districtID int) ([]entity.Location, error) {
	query := `SELECT 
							id, 
							name, 
							postal_codes
						FROM 
							sub_districts 
						WHERE 
							district_id = $1
						AND
							deleted_at is NULL`

	rows, err := r.db.QueryContext(c, query, districtID)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	defer rows.Close()

	subDistricts := []entity.Location{}

	for rows.Next() {
		var subDistrict entity.Location
		err := rows.Scan(
			&subDistrict.ID,
			&subDistrict.Name,
			&subDistrict.PostalCodes,
		)
		if err != nil {
			return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		}
		subDistricts = append(subDistricts, subDistrict)
	}

	return subDistricts, nil
}

func (r *addressRepoImpl) IsUserActiveAddressExists(c context.Context, userID int) (bool, error) {
	query := `SELECT EXISTS (
							SELECT 
								1 
							FROM 
								user_addresses 
							WHERE 
								user_id = $1 
							AND 
								is_active = TRUE
							AND
								deleted_at IS NULL
						)`

	var isExists bool
	err := r.db.QueryRowContext(c, query, userID).Scan(&isExists)
	if err != nil {
		return isExists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return isExists, nil
}

func (r *addressRepoImpl) AddUserAddress(c context.Context, userAddress entity.UserAddress) error {
	query := `INSERT INTO user_addresses(
							user_id,
							name,
							phone_number,
							address, 
							province_id, 
							province, 
							city_id, 
							city, 
							district_id, 
							district, 
							sub_district_id, 
							sub_district, 
							postal_code, 
							location, 
							is_active
						)
						VALUES
						($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, ST_SetSRID(ST_MakePoint($14, $15), 4326), $16)`

	_, err := r.db.ExecContext(c, query, userAddress.UserID, userAddress.Name, userAddress.PhoneNumber, userAddress.Address, userAddress.ProvinceID, userAddress.Province, userAddress.CityID, userAddress.City, userAddress.DistrictID, userAddress.District, userAddress.SubDistrictID, userAddress.SubDistrict, userAddress.PostalCode, userAddress.Longitude, userAddress.Latitude, userAddress.IsActive)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r *addressRepoImpl) GetCurrentLocationByLongAndLat(c context.Context, longitude string, latitude string) (*entity.UserAddress, error) {
	query := `SELECT 
							p.id, 
							c.id
						FROM provinces p
						JOIN cities c
							ON p.id = c.province_id
						WHERE
							p.deleted_at IS NULL
						AND
							c.deleted_at IS NULL
						ORDER BY
							ST_Distance(location, ST_SetSRID(ST_MakePoint($1, $2), 4326))
						LIMIT
							1`

	var address entity.UserAddress
	err := r.db.QueryRowContext(c, query, longitude, latitude).Scan(&address.ProvinceID, &address.CityID)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return &address, nil
}

func (r *addressRepoImpl) GetAddressesByUserID(c context.Context, userID int, filter entity.AddressFilter) ([]entity.UserAddress, error) {
	query := `SELECT
							id,
							name,
							phone_number,
							address, 
							province, 
							city, 
							district, 
							sub_district, 
							postal_code, 
							is_active
						FROM 
							user_addresses
						WHERE
							user_id = $1
						AND
							deleted_at IS NULL `

	order := "ORDER BY is_active DESC"
	queryFilter := AddressesParam(filter)
	rows, err := r.db.QueryContext(c, query+queryFilter+order, userID)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	defer rows.Close()

	addresses := []entity.UserAddress{}

	for rows.Next() {
		var address entity.UserAddress
		err := rows.Scan(
			&address.ID,
			&address.Name,
			&address.PhoneNumber,
			&address.Address,
			&address.Province,
			&address.City,
			&address.District,
			&address.SubDistrict,
			&address.PostalCode,
			&address.IsActive,
		)
		if err != nil {
			return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		}
		addresses = append(addresses, address)
	}

	return addresses, nil
}

func (r *addressRepoImpl) IsUserAddressExists(c context.Context, addressID int, userID int) (bool, error) {
	query := `SELECT EXISTS (
							SELECT
								1
							FROM
								user_addresses
							WHERE
								id = $1
							AND
								user_id = $2
							AND
								deleted_at IS NULL
						)`

	var isExists bool
	err := r.db.QueryRowContext(c, query, addressID, userID).Scan(&isExists)
	if err != nil {
		return isExists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return isExists, nil
}

func (r *addressRepoImpl) GetActiveAddressByUserID(c context.Context, userID int) (int, error) {
	query := `SELECT
							id
						FROM
							user_addresses
						WHERE
							user_id = $1
						AND 
							is_active = TRUE
						AND
							deleted_at IS NULL`

	var addressID int
	err := r.db.QueryRowContext(c, query, userID).Scan(&addressID)
	if err != nil {
		return -1, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return addressID, nil
}

func (r *addressRepoImpl) DeactivateAddressByID(c context.Context, addressID int) error {
	query := `UPDATE user_addresses
						SET is_active = FALSE
						WHERE id = $1`

	tx := transaction.ExtractTx(c)
	_, err := tx.ExecContext(c, query, addressID)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return nil
}

func (r *addressRepoImpl) UpdateUserAddress(c context.Context, address entity.UserAddress) error {
	query := `UPDATE user_addresses
						SET name = $1,
								phone_number = $2,
								address = $3,
								province_id = $4,
								province = $5,
								city_id = $6,
								city = $7,
								district_id = $8,
								district = $9,
								sub_district_id = $10,
								sub_district = $11,
								postal_code = $12,
								location = ST_SetSRID(ST_MakePoint($13, $14), 4326),
								is_active = $15,
								updated_at = NOW()
						WHERE
								id = $16`

	tx := transaction.ExtractTx(c)
	_, err := tx.ExecContext(
		c,
		query,
		address.Name,
		address.PhoneNumber,
		address.Address,
		address.ProvinceID,
		address.Province,
		address.CityID,
		address.City,
		address.DistrictID,
		address.District,
		address.SubDistrictID,
		address.SubDistrict,
		address.PostalCode,
		address.Longitude,
		address.Latitude,
		address.IsActive,
		address.ID,
	)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return nil
}

func (r *addressRepoImpl) GetAddressByID(c context.Context, addressID int) (*entity.UserAddress, error) {
	query := `SELECT
							name,
							phone_number,
							address,
							province_id,
							city_id,
							district_id,
							sub_district_id,
							postal_code,
							ST_X(location::geometry),
							ST_Y(location::geometry),
							is_active
						FROM 
							user_addresses
						WHERE
							id = $1`

	var address entity.UserAddress
	err := r.db.QueryRowContext(c, query, addressID).Scan(
		&address.Name,
		&address.PhoneNumber,
		&address.Address,
		&address.ProvinceID,
		&address.CityID,
		&address.DistrictID,
		&address.SubDistrictID,
		&address.PostalCode,
		&address.Longitude,
		&address.Latitude,
		&address.IsActive,
	)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return &address, nil
}

func (r *addressRepoImpl) DeleteUserAddressByID(c context.Context, addressID int) error {
	query := `UPDATE user_addresses
						SET deleted_at = NOW()
						WHERE id = $1`
	_, err := r.db.ExecContext(c, query, addressID)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return nil
}
