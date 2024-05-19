package authorization

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/khivuksergey/portmonetka.middleware"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

const (
	secret = "verystrongjwtsecret"
	userId = 111
)

var (
	userIdStr    = strconv.FormatUint(userId, 10)
	invalidIdStr = userIdStr + "1"
)

type testCase struct {
	token          *jwt.Token
	pathParamValue string
	expectedError  error
}

var testCases = []testCase{
	{
		token:          nil,
		pathParamValue: "",
		expectedError:  middleware.AuthorizationError{},
	},
	{
		token:          jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{}),
		pathParamValue: "",
		expectedError:  middleware.AuthorizationError{},
	},
	{
		token:          jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": userId}),
		pathParamValue: "",
		expectedError:  middleware.AuthorizationError{},
	},
	{
		token:          jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": userId}),
		pathParamValue: invalidIdStr,
		expectedError:  middleware.AuthorizationError{},
	},
}

var nextFunc = func(c echo.Context) error {
	id, ok := c.Get("userId").(uint64)
	if !ok || id != userId {
		return c.String(http.StatusUnauthorized, "Invalid userId in context in next handler")
	}
	return c.String(http.StatusOK, "Success")
}
