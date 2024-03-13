package repository

import "errors"

var (
	ErrUserExist    = errors.New("user already exist")
	ErrUserNotFound = errors.New("user not found")

	ErrCategoryExist  = errors.New("category already exist")
	ErrCategoryDelete = errors.New("error deleting category")
	ErrUpdateCategory = errors.New("error updating category name")
	ErrAllCategoryies = errors.New("error getting categories from database")

	ErrSaveProduct           = errors.New("product is not saved")
	ErrDeleteProduct         = errors.New("error deleting a product")
	ErrDeleteProductCategory = errors.New("error deleting a product category")
	ErrUpdateProduct         = errors.New("error updating product name")
	ErrSaveProductCategory   = errors.New("error save product category")
	ErrProductNotFound       = errors.New("product not found")
)
