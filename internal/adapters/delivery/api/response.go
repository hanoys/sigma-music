package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hanoys/sigma-music/internal/ports"
)

const (
	ErrBadRequest          = "bad request"
	ErrNotFound            = "not found"
	ErrUnauthorized        = "unauthorized"
	ErrForbidden           = "forbidden"
	ErrInternalServerError = "internal server error"
	ErrRequestTimeout      = "request timeout"
)

var (
	BadRequestError         = errors.New("bad request")
	NotFoundError           = errors.New("not Found")
	UnauthorizedError       = errors.New("unauthorized")
	ForbiddenError          = errors.New("forbidden")
	InternalServerError     = errors.New("internal server error")
	PathIDNotFoundError     = errors.New("id not found in query path")
	InvalidPathIDError      = errors.New("query path id is invalid")
	ParseGenreIDError       = errors.New("cannot parse genre id")
	ParseAlbumIDError       = errors.New("cannot parse album id")
	UnexpectedFileExtension = errors.New("unexpecte file extension")
)

var errorStatusMap = map[error]int{
	ports.ErrAlbumDuplicate:    http.StatusBadRequest,
	ports.ErrAlbumIDNotFound:   http.StatusNotFound,
	ports.ErrAlbumPublish:      http.StatusInternalServerError,
	ports.ErrInternalAlbumRepo: http.StatusInternalServerError,

	ports.ErrCommentDuplicate:         http.StatusBadRequest,
	ports.ErrCommentIDNotFound:        http.StatusNotFound,
	ports.ErrCommentByTrackIDNotFound: http.StatusNotFound,
	ports.ErrCommentByUserIDNotFound:  http.StatusNotFound,
	ports.ErrInternalCommentRepo:      http.StatusInternalServerError,

	ports.ErrGenreIDNotFound:   http.StatusNotFound,
	ports.ErrGenreNotFound:     http.StatusNotFound,
	ports.ErrInternalGenreRepo: http.StatusInternalServerError,

	ports.ErrTrackDuplicate:    http.StatusBadRequest,
	ports.ErrTrackIDNotFound:   http.StatusNotFound,
	ports.ErrTrackDelete:       http.StatusBadRequest,
	ports.ErrInternalTrackRepo: http.StatusInternalServerError,

	ports.ErrUserDuplicate:      http.StatusBadRequest,
	ports.ErrUserIDNotFound:     http.StatusNotFound,
	ports.ErrUserNameNotFound:   http.StatusNotFound,
	ports.ErrUserEmailNotFound:  http.StatusNotFound,
	ports.ErrUserPhoneNotFound:  http.StatusNotFound,
	ports.ErrUserUnknownCountry: http.StatusBadRequest,
	ports.ErrInternalUserRepo:   http.StatusInternalServerError,

	ports.ErrUserWithSuchNameAlreadyExists:  http.StatusConflict,
	ports.ErrUserWithSuchEmailAlreadyExists: http.StatusConflict,
	ports.ErrUserWithSuchPhoneAlreadyExists: http.StatusConflict,

	ports.ErrMusicianDuplicate:      http.StatusBadRequest,
	ports.ErrMusicianIDNotFound:     http.StatusNotFound,
	ports.ErrMusicianNameNotFound:   http.StatusNotFound,
	ports.ErrMusicianEmailNotFound:  http.StatusNotFound,
	ports.ErrMusicianUnknownCountry: http.StatusBadRequest,
	ports.ErrInternalMusicianRepo:   http.StatusInternalServerError,

	ports.ErrMusicianWithSuchNameAlreadyExists:  http.StatusConflict,
	ports.ErrMusicianWithSuchEmailAlreadyExists: http.StatusConflict,

	ports.ErrIncorrectName:     http.StatusUnauthorized,
	ports.ErrIncorrectPassword: http.StatusUnauthorized,
	ports.ErrUnexpectedRole:    http.StatusUnauthorized,
	ports.ErrInternalAuthRepo:  http.StatusUnauthorized,
	ports.ErrInvalidToken:      http.StatusUnauthorized,

	PathIDNotFoundError: http.StatusBadRequest,
	InvalidPathIDError:  http.StatusBadRequest,
}

type RestErr interface {
	Status() int
	Error() string
}

type RestError struct {
	ErrStatus  int       `json:"status,omitempty"`
	ErrMessage string    `json:"error,omitempty"`
	Timestamp  time.Time `json:"timestamp,omitempty"`
}

type RestErrorBadRequest struct {
	ErrStatus  int       `json:"status,omitempty" example:"400"`
	ErrMessage string    `json:"error,omitempty" example:"bad request"`
	Timestamp  time.Time `json:"timestamp,omitempty" example:"2020-11-10T23:00:00+00:00"`
}

type RestErrorUnauthorized struct {
	ErrStatus  int       `json:"status,omitempty" example:"401"`
	ErrMessage string    `json:"error,omitempty" example:"unauthorized"`
	Timestamp  time.Time `json:"timestamp,omitempty" example:"2020-11-10T23:00:00+00:00"`
}

type RestErrorForbidden struct {
	ErrStatus  int       `json:"status,omitempty" example:"403"`
	ErrMessage string    `json:"error,omitempty" example:"forbidden"`
	Timestamp  time.Time `json:"timestamp,omitempty" example:"2020-11-10T23:00:00+00:00"`
}

type RestErrorNotFound struct {
	ErrStatus  int       `json:"status,omitempty" example:"404"`
	ErrMessage string    `json:"error,omitempty" example:"not found"`
	Timestamp  time.Time `json:"timestamp,omitempty" example:"2020-11-10T23:00:00+00:00"`
}

type RestErrorConflict struct {
	ErrStatus  int       `json:"status,omitempty" example:"409"`
	ErrMessage string    `json:"error,omitempty" example:"conflict"`
	Timestamp  time.Time `json:"timestamp,omitempty" example:"2020-11-10T23:00:00+00:00"`
}

type RestErrorInternalError struct {
	ErrStatus  int       `json:"status,omitempty" example:"500"`
	ErrMessage string    `json:"error,omitempty" example:"internal server error"`
	Timestamp  time.Time `json:"timestamp,omitempty" example:"2020-11-10T23:00:00+00:00"`
}

func (e RestError) Error() string {
	return fmt.Sprintf("status: %d, error: %s", e.ErrStatus, e.ErrMessage)
}

func (e RestError) Status() int {
	return e.ErrStatus
}

func NewRestError(status int, err string) RestErr {
	return RestError{
		ErrStatus:  status,
		ErrMessage: err,
		Timestamp:  time.Now().UTC(),
	}
}

func getValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("field '%s' must be not empty", strings.ToLower(err.Field()))
	case "email":
		return fmt.Sprintf("invalid email")
	case "url":
		return fmt.Sprintf("field '%s' must be URL", strings.ToLower(err.Field()))
	case "oneof":
		return fmt.Sprintf("field '%s' must be enum type", err.Field())
	default:
		return "json validation error"
	}
}

func ParseError(err error) RestErr {
	var validationErrors validator.ValidationErrors

	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return NewRestError(http.StatusRequestTimeout, ErrRequestTimeout)
	case errors.Is(err, UnauthorizedError):
		return NewRestError(http.StatusUnauthorized, ErrUnauthorized)
	case errors.Is(err, BadRequestError):
		return NewRestError(http.StatusBadRequest, ErrBadRequest)
	case errors.Is(err, ForbiddenError):
		return NewRestError(http.StatusForbidden, ErrForbidden)
	case errors.Is(err, ports.ErrInvalidToken):
		return NewRestError(http.StatusUnauthorized, err.Error())
	case errors.As(err, &validationErrors):
		return NewRestError(http.StatusBadRequest, getValidationMessage(validationErrors[0]))
	default:
		for errType, errResponseCode := range errorStatusMap {
			if errors.Is(err, errType) {
				return NewRestError(errResponseCode, errType.Error())
			}
		}

		if restErr, ok := err.(*RestError); ok {
			return restErr
		}
		return NewRestError(http.StatusInternalServerError, ErrInternalServerError)
	}
}

func errorResponse(context *gin.Context, err error) {
	debug.PrintStack()
	restErr := ParseError(err)
	context.AbortWithStatusJSON(restErr.Status(), restErr)
}

func successResponse(context *gin.Context, data interface{}) {
	context.JSON(http.StatusOK, data)
}

func createdResponse(context *gin.Context, data interface{}) {
	context.JSON(http.StatusCreated, data)
}
