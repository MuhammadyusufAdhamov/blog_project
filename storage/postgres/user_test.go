package postgres_test

import (
	"testing"

	"blog_project/storage/repo"

	"github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/require"
)

func createUser(t *testing.T) *repo.User {
	user, err := strg.User().Create(&repo.User{
		FirstName: faker.FirstName(),
		LastName:  faker.LastName(),
		Email:     faker.Email(),
		Password:  faker.Password(),
		Type:      repo.UserTypeUser,
	})
	require.NoError(t, err)
	require.NotEmpty(t, user)

	return user
}

func TestGetUser(t *testing.T) {
	c := createUser(t)

	user, err := strg.User().Get(c.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user)
}

func TestCreateUser(t *testing.T) {
	createUser(t)
}

func TestDeleteUser(t *testing.T) {
	user, err := strg.User().DeleteUser(&repo.User{})
	require.NoError(t, err)
	require.NotEmpty(t, user)
}

func TestGetAllUsers(t *testing.T) {
	user, err := strg.User().GetAll(&repo.GetAllUsersParams{
		Limit:  3,
		Page:   1,
		Search: "ab",
	})
	require.NoError(t, err)
	require.NotEmpty(t, user)
}
