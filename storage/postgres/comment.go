package postgres

import (
	"database/sql"
	"fmt"

	"github.com/MuhammadyusufAdhamov/blog_project/storage/repo"
	"github.com/jmoiron/sqlx"
)

type commentRepo struct {
	db *sqlx.DB
}

func NewComment(db *sqlx.DB) repo.CommentStorageI {
	return &commentRepo {
		db: db,
	}
}

func (pr *commentRepo) Create(comment *repo.Comment)(*repo.Comment, error) {
	query := `
		insert into comments(
			user_id,
			post_id,
			description
		) values ($1,$2,$3)
		returning id, created_at
	`

	row := pr.db.QueryRow(
		query,
		comment.UserID,
		comment.PostID,
		comment.Description,
	)

	err := row.Scan(
		&comment.ID,
		&comment.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (ur *commentRepo) Get(id int64) (*repo.Comment, error) {
	var result repo.Comment

	query := `
		SELECT
			id,
			user_id,
			post_id,
			description,
			created_at
		FROM comments
		WHERE id=$1
	`

	row := ur.db.QueryRow(query, id)
	err := row.Scan(
		&result.ID,
		&result.UserID,
		&result.PostID,
		&result.Description,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func(pr *commentRepo) GetAll(params *repo.GetAllCommentsParams) (*repo.GetAllCommentsResult, error) {
	result := repo.GetAllCommentsResult{
		Comments: make([]*repo.Comment, 0),
	}

	offset := (params.Page - 1) * params.Limit

	limit := fmt.Sprintf(" limit %d offset %d ", params.Limit, offset)

	filter := " where true "
	if params.UserID != 0 {
		filter += fmt.Sprintf(" and c.user_id=%d ",params.UserID)
	}

	if params.PostID != 0 {
		filter += fmt.Sprintf(" and c.post_id=%d ", params.PostID)
	}

	query := `
		select
			c.id,
			c.user_id,
			c.post_id,
			c.description,
			c.created_at,
			c.updated_at,
			u.first_name,
			u.last_name,
			u.email,
			u.profile_image_url
		from comments c
		inner join users u on u.id=c.user_id
	` + filter + `order by c.created_at desc` + limit

	rows, err := pr.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var c repo.Comment
		err := rows.Scan(
			&c.ID,
			&c.UserID,
			&c.PostID,
			&c.Description,
			&c.CreatedAt,
			&c.UpdatedAt,
			&c.User.FirstName,
			&c.User.LastName,
			&c.User.Email,
			&c.User.ProfileIMageUrl,
		)
		if err != nil {
			return nil, err
		}

		result.Comments = append(result.Comments, &c)
	}

	queryCount := `
		SELECT count(1) FROM comments c
		INNER JOIN users u ON u.id=c.user_id ` + filter
	err = pr.db.QueryRow(queryCount).Scan(&result.Count)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *commentRepo) UpdateComment(comment *repo.Comment) (*repo.Comment, error) {
	query := `update comments set description=$1 where id=$2
			returning created_at
			`
	
	err := ur.db.QueryRow(query,comment.Description,comment.ID).Scan(&comment.CreatedAt)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func(ur *commentRepo) DeleteComment(id int64) error {
	query := `delete from comments where id=$1`

	result, err := ur.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}