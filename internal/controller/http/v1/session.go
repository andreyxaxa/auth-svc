package v1

import (
	"errors"
	"net/http"
	"strings"

	"github.com/andreyxaxa/auth-svc/internal/controller/http/v1/request"
	"github.com/andreyxaxa/auth-svc/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var (
	ErrMissingAuthHeader = errors.New("missing 'Authorization' header")
	ErrInvalidAuthFormat = errors.New("invalid 'Authorization' format")
)

// @Summary     Create token-pair
// @Description Create token-pair for user
// @ID          create
// @Tags            sessions
// @Produce 	json
// @Param 		user_id query string true "Create token-pair for user with specified GUID"
// @Success 	200 {object} entity.Token
// @Failure 	400 {object} response.Error
// @Failure 	500 {object} response.Error
// @Router 		/v1/session/token [post]
func (r *V1) create(ctx *fiber.Ctx) error {
	uidStr := ctx.Query("user_id")
	userID, err := uuid.Parse(uidStr)
	if err != nil {
		r.l.Error(err, "http - v1 - create - uuid.Parse")

		return errorResponse(ctx, http.StatusBadRequest, "invalid or missing user_id")
	}

	tokens, err := r.s.Create(ctx.Context(), userID, ctx.Get("User-Agent"), ctx.IP())
	if err != nil {
		r.l.Error(err, "http - v1 - create - r.s.Create")

		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(tokens)
}

// @Summary 	Refresh token-pair
// @Description Refresh token-pair for user. The previous one becomes invalid
// @ID 			refresh
// @Tags 			sessions
// @Accept 		json
// @Produce 	json
// @Security 	TokenAuth
// @Param 		request body request.RefreshRequest true "Refresh token request body"
// @Success 	200 {object} entity.Token
// @Failure 	400 {object} response.Error
// @Failure 	401 {object} response.Error
// @Router 		/v1/session/refresh [post]
func (r *V1) refresh(ctx *fiber.Ctx) error {
	var req request.RefreshRequest

	if err := ctx.BodyParser(&req); err != nil {
		r.l.Error(err, "http - v1 - refresh - ctx.BodyParser")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	userID, err := extractUserIDFromHeader(ctx, r.s)
	if err != nil {
		r.l.Error(err, "http - v1 - refresh - extractUserIDFromHeader")

		return errorResponse(ctx, http.StatusUnauthorized, err.Error())
	}

	tokens, err := r.s.Refresh(ctx.Context(), userID, req.RefreshToken, ctx.Get("User-Agent"), ctx.IP())
	if err != nil {
		r.l.Error(err, "http - v1 - refresh - r.s.Refresh")

		return errorResponse(ctx, http.StatusUnauthorized, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(tokens)
}

// @Summary 	Logout
// @Description All sessions become invalid
// @ID 			logout
// @Tags 			sessions
// @Produce 	json
// @Security 	TokenAuth
// @Success 	200
// @Failure 	401 {object} response.Error
// @Failure 	500 {object} response.Error
// @Router 		/v1/session/logout [post]
func (r *V1) logout(ctx *fiber.Ctx) error {
	userID, err := extractUserIDFromHeader(ctx, r.s)
	if err != nil {
		r.l.Error(err, "http - v1 - logout - extractUserIDFromHeader")

		return errorResponse(ctx, http.StatusUnauthorized, err.Error())
	}

	if err := r.s.Logout(ctx.Context(), userID); err != nil {
		r.l.Error(err, "http - v1 - logout - r.s.Logout")

		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.SendStatus(fiber.StatusOK)
}

// @Summary  Show user_id
// @ID 		 me
// @Tags 	 	user
// @Produce  json
// @Security TokenAuth
// @Success  200 {object} response.UserID
// @Failure  401 {object} response.Error
// @Router   /v1/user/me [get]
func (r *V1) me(ctx *fiber.Ctx) error {
	userID, err := extractUserIDFromHeader(ctx, r.s)
	if err != nil {
		r.l.Error(err, "http - v1 - me - extractUserIDFromHeader")

		return errorResponse(ctx, http.StatusUnauthorized, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"user_id": userID.String()})
}

func extractUserIDFromHeader(ctx *fiber.Ctx, s usecase.Sessions) (uuid.UUID, error) {
	auth := ctx.Get("Authorization")
	if auth == "" {
		return uuid.Nil, ErrMissingAuthHeader
	}

	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return uuid.Nil, ErrInvalidAuthFormat
	}

	return s.Me(ctx.Context(), parts[1])
}
