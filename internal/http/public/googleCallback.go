package http

import (
	"log/slog"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)


func (h *Handlers) CallbackWithGoogle(c *gin.Context, oauthConfig *oauth2.Config, verifier *oidc.IDTokenVerifier) {
	slog.Info("callback with google")	

	state := c.Query("state")
	slog.Info("state", "state", state)
	code := c.Query("code")
	cookie, err := c.Cookie("oauth_state")
	if err != nil {
		slog.Error("error getting cookie", "error", err)
		c.JSON(400, gin.H{
			"error": "error getting cookie",
		})
		return
	}

	session, err := h.userUsecase.LoginWithGoogle(c, oauthConfig, verifier, state, code, cookie)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	}
	
	c.JSON(200, gin.H{
		"session": session,
	})	
}