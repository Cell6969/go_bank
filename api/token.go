package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type renewTokenRequest struct {
	RefreshToken string `json:"_refresh_token" binding:"required"`
}

type renewTokenResponse struct {
	Token          string    `json:"_token"`
	TokenExpiredAt time.Time `json:"token_expired_at"`
}

func (server *Server) renewToken(ctx *gin.Context) {
	var req renewTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	refresh_payload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.store.GetSession(ctx, refresh_payload.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("session blocked")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.Username != refresh_payload.Username {
		err := fmt.Errorf("incorrect user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("session not valid")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if time.Now().After(session.ExpiredAt) {
		err := fmt.Errorf("session expired")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	token, token_payload, err := server.tokenMaker.CreateToken(refresh_payload.Username, server.config.TokenDuration)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := renewTokenResponse{
		Token:          token,
		TokenExpiredAt: token_payload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, response)
}
