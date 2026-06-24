package response

import (
	"testing"

	appvalidator "enterprise-order-management-api/internal/pkg/validator"

	playgroundvalidator "github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

type validationItem struct {
	ProductID int64 `json:"product_id" validate:"required,gt=0"`
}

type validationRequest struct {
	RefreshToken string           `json:"refresh_token" validate:"required"`
	Items        []validationItem `json:"items" validate:"required,min=1,dive"`
}

func TestFormatValidationErrorsUsesJSONFieldNames(t *testing.T) {
	validate := appvalidator.New()
	err := validate.Validate(validationRequest{Items: []validationItem{{}}})
	require.Error(t, err)

	validationErrors, ok := err.(playgroundvalidator.ValidationErrors)
	require.True(t, ok)

	formatted := formatValidationErrors(validationErrors)
	require.Contains(t, formatted, "refresh_token")
	require.Contains(t, formatted, "product_id")
	require.NotContains(t, formatted, "refreshtoken")
	require.NotContains(t, formatted, "productid")
}

func TestMapDatabaseError(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		expectedCode string
	}{
		{name: "unique", code: "23505", expectedCode: "CONFLICT"},
		{name: "foreign key", code: "23503", expectedCode: "BAD_REQUEST"},
		{name: "check", code: "23514", expectedCode: "BAD_REQUEST"},
		{name: "invalid text", code: "22P02", expectedCode: "BAD_REQUEST"},
		{name: "unknown", code: "99999", expectedCode: "INTERNAL_ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appErr := mapDatabaseError(&pgconn.PgError{Code: tt.code})
			require.Equal(t, tt.expectedCode, appErr.Code)
		})
	}
}
