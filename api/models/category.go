package models

import "time"

type Category struct {
	ID int64 `json:"id"`
	Title string `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type UpdateCategory struct {
	Title string `json:"title"`
}

type CreateCategoryRequest struct {
	Title string `json:"title" binding:"required,max=100"`
}

type GetAllCategoriesResponse struct {
	Categories []*Category `json:"categories"`
	Count int32 `json:"count"`
}
