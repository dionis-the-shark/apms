package httpauth

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/dionis-the-shark/apms-task-tracker/internal/authz"
	authusecase "github.com/dionis-the-shark/apms-task-tracker/internal/modules/auth/usecase"
	roleusecase "github.com/dionis-the-shark/apms-task-tracker/internal/modules/role/usecase"
)

func RequireAuth(authSvc authusecase.API, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractBearerToken(r.Header.Get("Authorization"))
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		if _, err := authSvc.ValidateToken(ctx, token); err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func RequirePermission(
	authSvc authusecase.API,
	roleSvc roleusecase.API,
	action authz.Action,
	next http.HandlerFunc,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractBearerToken(r.Header.Get("Authorization"))
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		claims, err := authSvc.ValidateToken(ctx, token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		allowed, err := roleSvc.ActionAllowed(ctx, claims.RoleID, action)
		if err != nil {
			http.Error(w, "Failed to authorize request", http.StatusInternalServerError)
			return
		}
		if !allowed {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}

func extractBearerToken(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
