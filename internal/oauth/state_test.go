package oauth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGenerateAndParseState(t *testing.T) {
	state, err := GenerateState("secret", "google", 5*time.Minute)
	require.NoError(t, err)

	claims, err := ParseState(state, "secret", "google")
	require.NoError(t, err)
	require.Equal(t, "google", claims.Provider)
}

func TestParseStateRejectsWrongProvider(t *testing.T) {
	state, err := GenerateState("secret", "google", 5*time.Minute)
	require.NoError(t, err)

	claims, err := ParseState(state, "secret", "github")
	require.Error(t, err)
	require.Nil(t, claims)
}
