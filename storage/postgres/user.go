package postgres

import (
	"fmt"

	"blog_project/storage/repo"

	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUser(db *sqlx.DB) repo.UserStorageI {
	return &userRepo{
		db: db,
	}
}

func (ur *userRepo) Create(user *repo.User) (*repo.User, error) {
	query := `
		insert into users(
			first_name,
			last_name,
			phone_number,
			email,
			gender,
			password,
			username,
			profile_image_url,
			type
		) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		returning id, created_at
	`

	row := ur.db.QueryRow(
		query,
		user.FirstName,
		user.LastName,
		user.PhoneNumber,
		user.Email,
		user.Gender,
		user.Password,
		user.Username,
		user.ProfileImageUrl,
		user.Type,
	)

	err := row.Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *userRepo) Activate(id int64) error {
	query := `update users set is_active=true where id =$1`
	_, err := ur.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepo) Get(id int64) (*repo.User, error) {
	var result repo.User

	query := `
		SELECT
			id,
			first_name,
			last_name,
			phone_number,
			email,
			gender,
			password,
			username,
			profile_image_url,
			type,
			created_at
		FROM users
		WHERE id=$1
	`

	row := ur.db.QueryRow(query, id)
	err := row.Scan(
		&result.ID,
		&result.FirstName,
		&result.LastName,
		&result.PhoneNumber,
		&result.Email,
		&result.Gender,
		&result.Password,
		&result.Username,
		&result.ProfileImageUrl,
		&result.Type,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *userRepo) GetAll(params *repo.GetAllUsersParams) (*repo.GetAllUsersResult, error) {
	result := repo.GetAllUsersResult{
		Users: make([]*repo.User, 0),
	}

	offset := (params.Page - 1) * params.Limit

	limit := fmt.Sprintf(" limit %d offset %d ", params.Limit, offset)

	filter := ""
	if params.Search != "" {
		str := "%" + params.Search + "%"
		filter += fmt.Sprintf(`
		where first_name ilike '%s' or last_name ilike '%s' or email ilike '%s'
		or username ilike '%s' or phone_number ilike '%s'`,
			str, str, str, str, str,
		)
	}

	query := `
		SELECT
			id,
			first_name,
			last_name,
			phone_number,
			email,
			gender,
			password,
			username,
			profile_image_url,
			type,
			created_at
		FROM users
		` + filter + `
		ORDER BY created_at desc
		` + limit

	rows, err := ur.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var u repo.User

		err := rows.Scan(
			&u.ID,
			&u.FirstName,
			&u.LastName,
			&u.PhoneNumber,
			&u.Email,
			&u.Gender,
			&u.Password,
			&u.Username,
			&u.ProfileImageUrl,
			&u.Type,
			&u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result.Users = append(result.Users, &u)
	}

	queryCount := `SELECT count(1) FROM users ` + filter
	err = ur.db.QueryRow(queryCount).Scan(&result.Count)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *userRepo) GetByEmail(email string) (*repo.User, error) {
	var result repo.User

	query := `
		select
			id,
			first_name,
			last_name,
			phone_number,
			email,
			gender,
			password,
			username,
			profile_image_url,
			type,
			created_at
		from users
		where email=$1
	`

	row := ur.db.QueryRow(query, email)
	err := row.Scan(
		&result.ID,
		&result.FirstName,
		&result.LastName,
		&result.PhoneNumber,
		&result.Email,
		&result.Gender,
		&result.Password,
		&result.Username,
		&result.ProfileImageUrl,
		&result.Type,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *userRepo) UpdatePassword(req *repo.UpdatePassword) error {
	query := `UPDATE users SET password=$1 WHERE id=$2`

	_, err := ur.db.Exec(query, req.Password, req.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (ur *userRepo) UpdateUser(user *repo.User) (*repo.User, error) {
	var result repo.User

	// UpdatedAt := time.Now()
	query := `update users set
				first_name=$1,
				last_name=$2,
				phone_number=$3,
				email=$4,
				gender=$5,
				password=$6,
				username=$7,
				profile_image_url=$8,
				type=$9
			where id=$10
			returning id
			`

	row := ur.db.QueryRow(
		query,
		user.FirstName,
		user.LastName,
		user.PhoneNumber,
		user.Email,
		user.Gender,
		user.Password,
		user.Username,
		user.ProfileImageUrl,
		user.Type,
		user.ID,
	)

	err := row.Scan(
		&result.ID,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *userRepo) DeleteUser(user *repo.User) (*repo.User, error) {
	var result repo.User
	// DeletedAt := time.Now()

	query := `delete from users
			where id=$1
			returning id`

	row := ur.db.QueryRow(
		query,
		user.ID,
	)

	err := row.Scan(
		&result.ID,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
