package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashPasswordAndCheckPassword(t *testing.T) {
	hashedPassword, err := HashPassword("secret123")

	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
	require.True(t, CheckPassword("secret123", hashedPassword))
	require.False(t, CheckPassword("wrong-password", hashedPassword))
}
