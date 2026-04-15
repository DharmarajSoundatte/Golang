package routes

import (
	"net/http"

	"github.com/DharmarajSoundatte/Golang/backend/internal/config"
	"github.com/DharmarajSoundatte/Golang/backend/internal/handlers"
	"github.com/DharmarajSoundatte/Golang/backend/internal/middleware"
	"github.com/DharmarajSoundatte/Golang/backend/internal/repository"
	"github.com/DharmarajSoundatte/Golang/backend/internal/services"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Setup wires all dependencies and returns a fully configured chi router.
func Setup(db *gorm.DB, cfg *config.Config, log *zap.Logger) http.Handler {
	// ── Repositories ──────────────────────────────────────────────────────────
	userRepo := repository.NewUserRepository(db)

	// ── Services ──────────────────────────────────────────────────────────────
	authSvc := services.NewAuthService(userRepo, cfg)
	userSvc := services.NewUserService(userRepo)

	// ── Handlers ──────────────────────────────────────────────────────────────
	authHandler := handlers.NewAuthHandler(authSvc, log)
	userHandler := handlers.NewUserHandler(userSvc, log)

	// ── Router ────────────────────────────────────────────────────────────────
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.CORS(cfg.AllowedOrigins))
	r.Use(middleware.RequestLogger(log))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// API v1
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Authenticate(cfg.JWTSecret))

			// Users — any authenticated user
			r.Route("/users", func(r chi.Router) {
				r.Get("/", userHandler.GetAll)
				r.Get("/{id}", userHandler.GetByID)

				// Admin-only mutations
				r.Group(func(r chi.Router) {
					r.Use(middleware.RequireRole("admin"))
					r.Put("/{id}", userHandler.UpdateUser)
					r.Delete("/{id}", userHandler.DeleteUser)
				})
			})
		})
	})

	return r
}
