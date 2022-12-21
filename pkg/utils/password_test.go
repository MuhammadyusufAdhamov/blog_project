package utils

import (
	"testing"
	"github.com/stretchr/testify/require"

)

func TestPassword(t *testing.T) {
	password := "12345678"

	hashesPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashesPassword)

	hashesPassword1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashesPassword1)

	err = CheckPassword(password, hashesPassword)
	require.NoError(t, err)
}