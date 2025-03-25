package handler

import (
	"montelukast/modules/category/converter"
	"montelukast/modules/category/dto"
	"montelukast/modules/category/usecase"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/pagination"
	"montelukast/pkg/wrapper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	u usecase.CategoryUsecase
}

func NewCategoryHandler(u usecase.CategoryUsecase) CategoryHandler {
	return CategoryHandler{
		u: u,
	}
}

func (h *CategoryHandler) AddCategoryHandler(c *gin.Context) {
	err := apperror.JsonValidator(c)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrAddCategory, apperror.ErrInvalidJSON, err)
		c.Error(err)
		return
	}
	categoryReq := dto.CategoryAddRequest{}

	err = c.ShouldBindJSON(&categoryReq)
	if err != nil {
		c.Error(err)
		return
	}

	err = h.u.AddCategory(c, converter.CategoryAddConverter{}.ToEntity(categoryReq))
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(nil, "add category success!", nil)
	c.JSON(http.StatusCreated, response)
}

func (h *CategoryHandler) UpdateCategoryHandler(c *gin.Context) {
	err := apperror.JsonValidator(c)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrUpdateCategory, apperror.ErrInvalidJSON, err)
		c.Error(err)
		return
	}
	categoryReq := dto.CategoryUpdateRequest{}

	err = c.ShouldBindJSON(&categoryReq)
	if err != nil {
		c.Error(err)
		return
	}

	err = h.u.UpdateCategory(c, converter.CategoryUpdateConverter{}.ToEntity(categoryReq))
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(nil, "update category success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h *CategoryHandler) DeleteCategoryHandler(c *gin.Context) {
	id := c.Param("id")
	categoryID, err := strconv.Atoi(id)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrDeleteCategory, apperror.ErrConvertVariableType, err)
		c.Error(err)
		return
	}

	err = h.u.DeleteCategory(c, categoryID)
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(nil, "delete category success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h *CategoryHandler) GetCategoryDetailHandler(c *gin.Context) {
	id := c.Param("id")
	categoryID, err := strconv.Atoi(id)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrGetCategory, apperror.ErrConvertVariableType, err)
		c.Error(err)
		return
	}

	res, err := h.u.GetCategoryDetail(c, categoryID)
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(converter.GetCategoriesConverter{}.ToDto(*res), "get product detail success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h *CategoryHandler) GetCategoriesHandler(c *gin.Context) {
	filterDTO := dto.CategoryFilterRequest{}
	if err := c.ShouldBindQuery(&filterDTO); err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrGetCategories, apperror.ErrConvertVariableType, nil)
		c.Error(err)
		return
	}

	var filterConverter converter.FilterCategoriesConverter
	paginatedCategories, err := h.u.GetCategories(c, filterConverter.ToEntity(filterDTO))
	if err != nil {
		c.Error(err)
		return
	}
	categories := []dto.CategoryResponse{}

	var categoryConverter converter.GetCategoriesConverter
	for _, category := range paginatedCategories.Categories {
		categories = append(categories, categoryConverter.ToDto(category))
	}
	var finalResult dto.PaginatedCategoriesResponse
	finalResult.Categories = categories
	var paginationInfo pagination.PaginationConverter
	finalResult.Pagination = paginationInfo.ToDto(paginatedCategories.Pagination)

	response := wrapper.ResponseData(finalResult, "get categories success", nil)
	c.JSON(http.StatusOK, response)
}
