package usecase

import (
	"context"
	"math"
	"montelukast/modules/category/entity"
	"montelukast/modules/category/repository"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
)

type CategoryUsecase interface {
	AddCategory(c context.Context, category entity.Category) error
	UpdateCategory(c context.Context, category entity.Category) error
	DeleteCategory(c context.Context, id int) error
	GetCategoryDetail(c context.Context, id int) (*entity.Category, error)
	GetCategories(c context.Context, filter entity.CategoryFilter) (*entity.PaginatedCategories, error)
}

type categoryUsecaseImpl struct {
	r repository.CategoryRepo
}

func NewCategoryUsecase(r repository.CategoryRepo) categoryUsecaseImpl {
	return categoryUsecaseImpl{
		r: r,
	}
}

func (u categoryUsecaseImpl) AddCategory(c context.Context, category entity.Category) error {
	isAlreadyExists, err := u.r.IsCategoryExistsByName(c, category.Name)
	if err != nil {
		return err
	}
	if isAlreadyExists {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrAddCategory, apperror.ErrCategoryAlreadyExists, apperror.ErrCategoryAlreadyExists)
	}
	err = u.r.AddCategory(c, category)
	if err != nil {
		return err
	}
	return nil
}

func (u categoryUsecaseImpl) UpdateCategory(c context.Context, category entity.Category) error {
	isExists, err := u.r.IsCategoryExistsByID(c, category.ID)
	if err != nil {
		return err
	}
	if !isExists {
		return apperror.NewErrStatusNotFound(appconstant.FieldErrUpdateCategory, apperror.ErrCategoryNotExists, apperror.ErrCategoryNotExists)
	}
	
	isNameExists, err := u.r.IsCategoryExistsByName(c, category.Name)
	if err != nil {
		return err
	}
	if isNameExists {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrUpdateCategory, apperror.ErrCategoryAlreadyExists, apperror.ErrCategoryAlreadyExists)
	}

	err = u.r.UpdateCategory(c, category)
	if err != nil {
		return err
	}

	return nil
}

func (u categoryUsecaseImpl) DeleteCategory(c context.Context, id int) error {
	isExists, err := u.r.IsCategoryExistsByID(c, id)
	if err != nil {
		return err
	}
	if !isExists {
		return apperror.NewErrStatusNotFound(appconstant.FieldErrDeleteCategory, apperror.ErrCategoryNotExists, apperror.ErrCategoryNotExists)
	}

	hasProductAlready, err := u.r.IsProductExistsByCategory(c, id)
	if err != nil {
		return err
	}
	if hasProductAlready {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrDeleteCategory, apperror.ErrCategoryHasProductAlready, apperror.ErrCategoryHasProductAlready)
	}

	err = u.r.DeleteCategoryByID(c, id)
	if err != nil {
		return err
	}

	return nil
}

func (u categoryUsecaseImpl) GetCategories(c context.Context, filter entity.CategoryFilter) (*entity.PaginatedCategories, error) {
	categories, err := u.r.GetCategories(c, filter)
	if err != nil {
		return nil, err
	}
	var categoryCount entity.CategoryFilterCount
	categoryCount.Name = filter.Name
	totalItem, err := u.r.GetTotalItem(c, categoryCount)
	if err != nil {
		return nil, err
	}
	if filter.Page < 1 {
		filter.Page = 1
	}
	totalPage := math.Ceil(float64(totalItem) / float64(filter.GetLimit()))
	categories.Pagination.TotalItem = totalItem
	categories.Pagination.CurrentPage = filter.Page
	categories.Pagination.TotalPage = int(totalPage)
	return categories, nil
}

func (u categoryUsecaseImpl) GetCategoryDetail(c context.Context, id int) (*entity.Category, error) {
	isExists, err := u.r.IsCategoryExistsByID(c, id)
	if err != nil {
		return nil, err
	}
	if !isExists {
		return nil, apperror.NewErrStatusNotFound(appconstant.FieldErrGetCategory, apperror.ErrCategoryNotExists, apperror.ErrCategoryNotExists)
	}

	category, err := u.r.GetCategoryByID(c, id)
	if err != nil {
		return nil, err
	}

	return category, nil
}
