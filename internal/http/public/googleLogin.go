package http

import (
	"crypto/rand"
	"encoding/base64"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)


func (h *Handlers) LoginWithGoogle(c *gin.Context, oauth *oauth2.Config,  ) {
	state := randomString(32)

	c.SetCookie("oauth_state", state, 3600, "/", "localhost", false, true)
	url := oauth.AuthCodeURL(state)
	slog.Info("redirecting to google", "url", url)
	c.Redirect(http.StatusFound, url)
}

func randomString(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)	
	return base64.RawURLEncoding.EncodeToString(b) 
}