package repo

import "time"

type Comment struct {
	ID int64
	PostID int64
	UserID int64
	Description string
	CreatedAt time.Time
	UpdatedAt *time.Time
	User struct {
		FirstName string
		LastName string
		Email string
		ProfileIMageUrl *string
	}
}

type GetAllCommentsParams struct {
	Limit int32
	Page int32
	UserID int64
	PostID int64
}

type GetAllCommentsResult struct {
	Comments []*Comment
	Count int32
}

type CommentStorageI interface {
	Create(c *Comment) (*Comment, error)
	GetAll(params *GetAllCommentsParams) (*GetAllCommentsResult, error)
	Get(id int64) (*Comment, error)
	UpdateComment(c *Comment) (*Comment, error)
	DeleteComment(id int64) error
}