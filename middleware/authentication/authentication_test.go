package authentication

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthentication_Success(t *testing.T) {
	claims := jwt.MapClaims{
		"sub": float64(userId),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token.Raw)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.Set("user", token)
	c.SetParamNames("userId")
	c.SetParamValues(userIdStr)

	authMiddleware := NewAuthenticationMiddleware(secret, nil)
	authenticate := authMiddleware.Authenticate(nextFunc)

	err := authenticate(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Success", rec.Body.String())
}

func TestAuthentication_Errors(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	authMiddleware := NewAuthenticationMiddleware(secret, nil)
	authenticate := authMiddleware.Authenticate(nextFunc)

	for _, test := range testCases {
		if test.token != nil {
			req.Header.Set("Authorization", "Bearer "+test.token.Raw)
			c.Set("user", test.token)
		}
		c.SetParamNames("userId")
		c.SetParamValues(test.pathParamValue)

		err := authenticate(c)

		assert.Error(t, err)
		assert.IsType(t, test.expectedError, err)
	}
}
