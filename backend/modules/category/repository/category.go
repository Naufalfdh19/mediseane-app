package repository

import (
	"context"
	"database/sql"
	"montelukast/modules/category/entity"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
)

type CategoryRepo interface {
	AddCategory(c context.Context, category entity.Category) error
	UpdateCategory(c context.Context, category entity.Category) error
	IsCategoryExistsByName(c context.Context, name string) (bool, error)
	DeleteCategoryByID(c context.Context, id int) error
	IsCategoryExistsByID(c context.Context, id int) (bool, error)
	IsProductExistsByCategory(c context.Context, categoryID int) (bool, error)
	GetTotalItem(c context.Context, filter entity.CategoryFilterCount) (int, error)
	GetCategoryByID(c context.Context, id int) (*entity.Category, error)
	GetCategories(c context.Context, filter entity.CategoryFilter) (*entity.PaginatedCategories, error)
}

type categoryRepoImpl struct {
	db *sql.DB
}

func NewCategoryRepo(dbConn *sql.DB) categoryRepoImpl {
	return categoryRepoImpl{
		db: dbConn,
	}
}

func (r categoryRepoImpl) AddCategory(c context.Context, category entity.Category) error {
	query := `INSERT INTO product_categories (name)
			  VALUES ($1)`

	_, err := r.db.ExecContext(c, query, category.Name)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r categoryRepoImpl) UpdateCategory(c context.Context, category entity.Category) error {
	query := `UPDATE 
				product_categories
			SET 
				name = $2,
				updated_at = NOW() 
			WHERE 
				id = $1`

	_, err := r.db.ExecContext(c, query, category.ID, category.Name)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r categoryRepoImpl) IsCategoryExistsByName(c context.Context, name string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM product_categories WHERE name = $1 AND deleted_at IS NULL)`

	var isExists bool
	err := r.db.QueryRowContext(c, query, name).Scan(&isExists)
	if err != nil && err != sql.ErrNoRows {
		return isExists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return isExists, nil
}

func (r categoryRepoImpl) DeleteCategoryByID(c context.Context, id int) error {
	query := `UPDATE 
				product_categories
			SET 
				deleted_at = NOW()
			WHERE 
				id = $1`

	_, err := r.db.ExecContext(c, query, id)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r categoryRepoImpl) IsCategoryExistsByID(c context.Context, id int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM product_categories WHERE id = $1 AND deleted_at IS NULL)`

	var isExists bool
	err := r.db.QueryRowContext(c, query, id).Scan(&isExists)
	if err != nil && err != sql.ErrNoRows {
		return isExists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return isExists, nil
}

func (r categoryRepoImpl) IsProductExistsByCategory(c context.Context, categoryID int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM product_multi_categories WHERE product_category_id = $1 AND deleted_at IS NULL)`

	var isExists bool
	err := r.db.QueryRowContext(c, query, categoryID).Scan(&isExists)
	if err != nil && err != sql.ErrNoRows {
		return isExists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return isExists, nil
}

func (r categoryRepoImpl) GetCategoryByID(c context.Context, id int) (*entity.Category, error) {
	query := `SELECT 
				c.id, c.name, c.updated_at
			FROM 
				product_categories c
			WHERE 
				c.id  = $1 AND c.deleted_at IS NULL`

	var category entity.Category
	err := r.db.QueryRow(query, id).Scan(
		&category.ID,
		&category.Name,
		&category.UpdatedAt,
	)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return &category, nil
}

func (r categoryRepoImpl) GetCategories(c context.Context, filter entity.CategoryFilter) (*entity.PaginatedCategories, error) {
	query := `SELECT 
				c.id, c.name, c.updated_at
			FROM 
				product_categories c
			WHERE 
				c.deleted_at IS NULL`

	result := entity.PaginatedCategories{}
	paramQuery, sortPage, args := CategoryParam(filter)
	finalQuery := query + paramQuery + sortPage
	rows, err := r.db.QueryContext(c, finalQuery, args...)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	defer rows.Close()

	for rows.Next() {
		var category entity.Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		}
		result.Categories = append(result.Categories, category)
	}

	err = rows.Err()
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return &result, nil
}

func (r categoryRepoImpl) GetTotalItem(c context.Context, filter entity.CategoryFilterCount) (int, error) {
	query := `SELECT 
				COUNT(*) 
			FROM 
				product_categories c
			WHERE 
				c.deleted_at IS NULL`

	paramQuery, args := CategoryParamCount(filter)
	var totalItem int
	err := r.db.QueryRowContext(c, query+paramQuery, args...).Scan(&totalItem)
	if err != nil {
		return 0, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return totalItem, nil
}
