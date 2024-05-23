package error

import (
	"errors"
	"github.com/google/uuid"
	common "github.com/khivuksergey/portmonetka.common"
	"github.com/labstack/echo/v4"
	"net/http"
)

var (
	authorizationError       = &common.AuthorizationError{}
	validationError          = &common.ValidationError{}
	unprocessableEntityError = &common.UnprocessableEntityError{}
	echoHttpError            *echo.HTTPError
)

type ErrorHandlingMiddleware struct{}

func NewErrorHandlingMiddleware() *ErrorHandlingMiddleware {
	return &ErrorHandlingMiddleware{}
}

func (e *ErrorHandlingMiddleware) HandleError(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error

		requestUuid, ok := c.Get(common.RequestUuidKey).(string)
		if ok {
			err = uuid.Validate(requestUuid)
		}

		if !ok || err != nil {
			requestUuid = uuid.New().String()
		}

		c.Set(common.RequestUuidKey, requestUuid)

		err = next(c)

		if err == nil {
			return nil
		}

		switch {
		case errors.As(err, authorizationError):
			return c.JSON(http.StatusUnauthorized, common.Response{
				Message:     err.Error(),
				RequestUuid: requestUuid,
			})
		case errors.As(err, validationError):
			return c.JSON(http.StatusBadRequest, common.Response{
				Message:     err.Error(),
				RequestUuid: requestUuid,
			})
		case errors.As(err, unprocessableEntityError):
			return c.JSON(http.StatusUnprocessableEntity, common.Response{
				Message:     err.Error(),
				RequestUuid: requestUuid,
			})
		case errors.As(err, &echoHttpError):
			echoErr := err.(*echo.HTTPError)
			return c.JSON(echoErr.Code, common.Response{
				Message:     echoErr.Message.(string),
				RequestUuid: requestUuid,
			})
		default:
			return c.JSON(http.StatusInternalServerError, common.Response{
				Message:     "internal server error",
				Data:        err,
				RequestUuid: requestUuid,
			})
		}
	}
}
