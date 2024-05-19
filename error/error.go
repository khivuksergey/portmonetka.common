package error

import (
	"errors"
	"github.com/google/uuid"
	middleware "github.com/khivuksergey/portmonetka.middleware"
	"github.com/labstack/echo/v4"
	"net/http"
)

const RequestUuidKey = "request_uuid"

var (
	authorizationError       = &middleware.AuthorizationError{}
	validationError          = &middleware.ValidationError{}
	unprocessableEntityError = &middleware.UnprocessableEntityError{}
	echoHttpError            *echo.HTTPError
)

type ErrorHandlingMiddleware struct{}

func NewErrorHandlingMiddleware() *ErrorHandlingMiddleware {
	return &ErrorHandlingMiddleware{}
}

func (e *ErrorHandlingMiddleware) HandleError(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error

		requestUuid, ok := c.Get(RequestUuidKey).(string)
		if ok {
			err = uuid.Validate(requestUuid)
		}

		if !ok || err != nil {
			requestUuid = uuid.New().String()
		}

		c.Set(RequestUuidKey, requestUuid)

		err = next(c)

		if err == nil {
			return nil
		}

		switch {
		case errors.As(err, authorizationError):
			return c.JSON(http.StatusUnauthorized, middleware.Response{
				Message:     err.Error(),
				RequestUuid: requestUuid,
			})
		case errors.As(err, validationError):
			return c.JSON(http.StatusBadRequest, middleware.Response{
				Message:     err.Error(),
				RequestUuid: requestUuid,
			})
		case errors.As(err, unprocessableEntityError):
			return c.JSON(http.StatusUnprocessableEntity, middleware.Response{
				Message:     err.Error(),
				RequestUuid: requestUuid,
			})
		case errors.As(err, &echoHttpError):
			echoErr := err.(*echo.HTTPError)
			return c.JSON(echoErr.Code, middleware.Response{
				Message:     echoErr.Message.(string),
				RequestUuid: requestUuid,
			})
		default:
			return c.JSON(http.StatusInternalServerError, middleware.Response{
				Message:     "internal server error",
				Data:        err,
				RequestUuid: requestUuid,
			})
		}
	}
}
