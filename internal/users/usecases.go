package users

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type UsersUsecase struct {
	repo Repository
}

func NewUsersUsecase(repo Repository) UsersUsecase {
	return UsersUsecase{repo: repo}
}

func (u *UsersUsecase) LoginWithGoogle(ctx context.Context, oauthConfig *oauth2.Config, verifier *oidc.IDTokenVerifier, state string, code string, cookie string) (*Session, error) {
	if state != cookie {
		slog.Error("state mismatch", "state", state, "cookie", cookie)
		return nil, errors.New("state mismatch")
	}

	token, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		slog.Error("error exchanging token", "error", err, "state", state)
		return nil, err
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		slog.Error("error getting id token", "state", state)
		return nil, errors.New("error getting id token",)
	}

	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		slog.Error("error verifying id token", "error", err)
		return nil, err
	}

	var claims Claims
	

	if err := idToken.Claims(&claims); err != nil {
		slog.Error("error parsing id token claims", "error", err)
		return nil, err
	}
	
	expires_at := time.Now().Add(30 * 24 * time.Hour)
	return u.repo.LoginWithGoogle(ctx, claims, expires_at)
}

