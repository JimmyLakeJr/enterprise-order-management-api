package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGenerateProducesUniqueTokens(t *testing.T) {
	first, _, err := Generate(1, "user@example.com", "admin", "secret", 15*time.Minute)
	require.NoError(t, err)

	second, _, err := Generate(1, "user@example.com", "admin", "secret", 15*time.Minute)
	require.NoError(t, err)

	require.NotEqual(t, first, second)
}
