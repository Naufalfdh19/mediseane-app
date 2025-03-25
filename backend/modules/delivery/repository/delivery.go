package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"montelukast/modules/delivery/entity"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"net/http"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type DeliveryRepository interface {
	IsPostalCodeExist(c context.Context, postalCode string) (exists bool, err error)
	GetLocationID(c context.Context, postalCode string) (locationID int, err error)
	GetOngkirLocationID(c context.Context, postal_code string) (locationID int, err error)
	GetOngkirCost(c context.Context, ongkir entity.CalculateOngkir) (result entity.UserCostResponse, err error)
	GetListOngkir(c context.Context, addressID int, pharmacyID int) (listOngkir []entity.OngkirData, err error)
	AddLocationID(c context.Context, postalCode string, locationID int) (err error)
	StoreListOngkir(c context.Context, listOngkir []entity.OngkirData, userID int, pharmacyID int) (err error)
	GetSameDayPrice(c context.Context) (price float64, err error)
	GetInstantPrice(c context.Context) (price float64, err error)
	CalculateDistance(c context.Context, pharmacyId int, addressID int) (distance float64, err error)
	GetUserPostalCode(c context.Context, userID int) (postalCode *string, addressID int, err error)
	GetPharmacyPostalCode(c context.Context, userID int) (postalCode *string, err error)
}

type deliveryRepository struct {
	redisDB *redis.Client
	db      *sql.DB
}

func NewDeliveryRepository(db *sql.DB, redisDB *redis.Client) DeliveryRepository {
	return &deliveryRepository{
		db:      db,
		redisDB: redisDB,
	}
}

func (r *deliveryRepository) CalculateDistance(c context.Context, pharmacyId int, addressID int) (distance float64, err error) {
	query := `SELECT ST_DISTANCE(
	(SELECT location FROM pharmacies WHERE id = $1 AND deleted_at IS NULL AND is_active = TRUE),
	(SELECT location FROM user_addresses WHERE id = $2 AND deleted_at IS NULL AND is_active=TRUE)
	) AS distance`
	err = r.db.QueryRowContext(c, query, pharmacyId, addressID).Scan(&distance)
	if err != nil {
		return -1, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return distance, nil
}

func (r *deliveryRepository) GetPharmacyPostalCode(c context.Context, pharmacyID int) (postalCode *string, err error) {
	query := `SELECT postal_code FROM pharmacies
		WHERE id = $1
		AND deleted_at IS NULL`
	err = r.db.QueryRowContext(c, query, pharmacyID).Scan(&postalCode)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return postalCode, nil
}

func (r *deliveryRepository) GetUserPostalCode(c context.Context, userID int) (postalCode *string, addressID int, err error) {
	query := `SELECT id,postal_code FROM user_addresses
		WHERE user_id = $1
		AND deleted_at IS NULL AND is_active = TRUE`
	err = r.db.QueryRowContext(c, query, userID).Scan(&addressID, &postalCode)
	if err != nil {
		return nil, -1, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return postalCode, addressID, nil
}

func (r *deliveryRepository) GetSameDayPrice(c context.Context) (price float64, err error) {
	query := `SELECT price FROM logistics
		WHERE name = $1
		AND deleted_at IS NULL`
	err = r.db.QueryRowContext(c, query, "same day").Scan(&price)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, apperror.NewErrInternalServerError(appconstant.FieldErrAddToCart, apperror.ErrDataNotExists, err)
		}
		return -1, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return price, nil
}

func (r *deliveryRepository) GetInstantPrice(c context.Context) (price float64, err error) {
	query := `SELECT price FROM logistics
		WHERE name = $1 
		AND deleted_at IS NULL`
	err = r.db.QueryRowContext(c, query, "instant").Scan(&price)
	if err != nil {
		return -1, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return price, nil
}

func (r *deliveryRepository) StoreListOngkir(c context.Context, listOngkir []entity.OngkirData, userID int, pharmacyID int) (err error) {
	key := fmt.Sprintf(appconstant.OngkirRedisKey, userID, pharmacyID)
	json, err := json.Marshal(listOngkir)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	_, err = r.redisDB.Set(c, key, json, appconstant.OngkirTimeExpiration).Result()
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r *deliveryRepository) GetListOngkir(c context.Context, addressID int, pharmacyID int) (listOngkir []entity.OngkirData, err error) {
	key := fmt.Sprintf(appconstant.OngkirRedisKey, addressID, pharmacyID)
	value, err := r.redisDB.Get(c, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, err
		}
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	json.Unmarshal([]byte(value), &listOngkir)
	return listOngkir, nil
}

func (r *deliveryRepository) IsPostalCodeExist(c context.Context, postalCode string) (exists bool, err error) {
	checkQuery := `SELECT EXISTS (SELECT 1 FROM postal_ongkir WHERE postal_code = $1 AND deleted_at IS NULL) `
	err = r.db.QueryRowContext(c, checkQuery, postalCode).Scan(&exists)
	if err != nil {
		return exists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return exists, nil
}

func (r *deliveryRepository) GetLocationID(c context.Context, postalCode string) (locationID int, err error) {
	query := `SELECT id_location
			  FROM postal_ongkir
			  WHERE postal_code = $1 AND deleted_at IS NULL`
	err = r.db.QueryRowContext(c, query, postalCode).Scan(&locationID)
	if err != nil {
		return -1, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return locationID, nil
}

func (r *deliveryRepository) AddLocationID(c context.Context, postalCode string, locationID int) (err error) {
	query := `INSERT INTO postal_ongkir(
				id_location,postal_code)
				VALUES ($1,$2)`
	_, err = r.db.ExecContext(c, query, locationID, postalCode)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r *deliveryRepository) GetOngkirLocationID(c context.Context, postal_code string) (locationID int, err error) {
	baseUrl := appconstant.URLOngkirLocationID
	req, _ := http.NewRequest("GET", baseUrl, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("key", os.Getenv("RAJA_ONGKIR_API_KEY"))
	query := req.URL.Query()
	query.Set("search", postal_code)
	req.URL.RawQuery = query.Encode()
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return locationID, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return locationID, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	var response entity.OngkirLocation
	err = json.Unmarshal(body, &response)
	if err != nil {
		return locationID, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	if len(response.Data) == 0 {
		return locationID, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, apperror.ErrRajaOngkirNoResult)
	}
	locationID = response.Data[0].ID
	return locationID, nil
}

func (r *deliveryRepository) GetOngkirCost(c context.Context, ongkir entity.CalculateOngkir) (result entity.UserCostResponse, err error) {
	baseUrl := ongkir.BaseUrl
	req, _ := http.NewRequest("POST", baseUrl, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("key", os.Getenv("RAJA_ONGKIR_API_KEY"))
	query := req.URL.Query()

	query.Set("origin", strconv.Itoa(ongkir.OriginID))
	query.Set("destination", strconv.Itoa(ongkir.DestinationID))
	query.Set("weight", strconv.Itoa(ongkir.Weight))
	query.Set("price", ongkir.SortingPrice)
	query.Set("courier", ongkir.Courier)

	req.URL.RawQuery = query.Encode()
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	var response entity.OngkirCostResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return result, err
	}
	result.Name = response.Data[0].Name
	result.Cost = response.Data[0].Cost
	result.Etd = response.Data[0].Etd
	return result, nil
}
