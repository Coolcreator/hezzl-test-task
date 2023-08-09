package v1

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"

	"goods-service/internal/good/domain"
)

var (
	notFoundError = &errorResponse{
		Code:    3,
		Message: "errors.good.notFound",
	}

	badRequest = &errorResponse{
		Code:    4,
		Message: "errors.badRequest",
	}

	internalServerError = &errorResponse{
		Code:    5,
		Message: "errors.internalServerError",
	}
)

type (
	errorResponse struct {
		Code    int32  `json:"code"`
		Message string `json:"message"`
		Details any    `json:"details"`
	}

	errorHandlerFunc func(w http.ResponseWriter, r *http.Request) (err error)

	errorHandler struct{}
)

func (h *errorHandler) wrap(f errorHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrGoodNotFound):
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, notFoundError)
			case errors.Is(err, domain.ErrBadRequest):
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, badRequest)
			default:
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, internalServerError)
			}
		}
	}
}
