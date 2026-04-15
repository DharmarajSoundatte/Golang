package middleware

import (
	"context"
	"net/http"
	"strings"

	jwtpkg "github.com/DharmarajSoundatte/Golang/backend/pkg/jwt"
	"github.com/DharmarajSoundatte/Golang/backend/pkg/response"
)

// contextKey avoids key collisions in request context.
type contextKey string

const (
	// ContextKeyUserID is the key used to store the authenticated user's ID in ctx.
	ContextKeyUserID contextKey = "user_id"
	// ContextKeyUserRole is the key used to store the authenticated user's role in ctx.
	ContextKeyUserRole contextKey = "user_role"
)

// Authenticate is a JWT validation middleware.
// It expects: Authorization: Bearer <token>
func Authenticate(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.WriteError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				response.WriteError(w, http.StatusUnauthorized, "invalid authorization header format")
				return
			}

			claims, err := jwtpkg.ValidateToken(parts[1], secret)
			if err != nil {
				response.WriteError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
			ctx = context.WithValue(ctx, ContextKeyUserRole, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole builds a middleware that enforces a minimum role requirement.
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value(ContextKeyUserRole).(string)
			if !ok || userRole != role {
				response.WriteError(w, http.StatusForbidden, "insufficient permissions")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// GetUserIDFromCtx is a helper to extract the user ID from the request context.
func GetUserIDFromCtx(ctx context.Context) (uint, bool) {
	id, ok := ctx.Value(ContextKeyUserID).(uint)
	return id, ok
}
