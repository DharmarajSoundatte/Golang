package routes

import (
	"net/http"

	"github.com/DharmarajSoundatte/Golang/backend/internal/handlers"
	"github.com/DharmarajSoundatte/Golang/backend/internal/middleware"
	"github.com/DharmarajSoundatte/Golang/backend/internal/repository"
	"github.com/DharmarajSoundatte/Golang/backend/internal/services"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// requireAdminOrTeacher allows requests from users with role "admin" or "teacher".
func requireAdminOrTeacher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(middleware.ContextKeyUserRole).(string)
		if !ok || (role != "admin" && role != "teacher") {
			http.Error(w, `{"error":"insufficient permissions"}`, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RegisterSanjanaRoutes wires up all of Sanjana's module routes.
func RegisterSanjanaRoutes(r chi.Router, db *gorm.DB, log *zap.Logger) {
	// ── Repositories ──────────────────────────────────────────────────────────
	teacherRepo      := repository.NewTeacherRepository(db)
	classRepo        := repository.NewClassRepository(db)
	attendanceRepo   := repository.NewAttendanceRepository(db)
	gradeRepo        := repository.NewGradeRepository(db)
	timetableRepo    := repository.NewTimetableRepository(db)
	announcementRepo := repository.NewAnnouncementRepository(db)

	// ── Services ──────────────────────────────────────────────────────────────
	teacherSvc      := services.NewTeacherService(teacherRepo)
	classSvc        := services.NewClassService(classRepo)
	attendanceSvc   := services.NewAttendanceService(attendanceRepo)
	gradeSvc        := services.NewGradeService(gradeRepo)
	timetableSvc    := services.NewTimetableService(timetableRepo)
	announcementSvc := services.NewAnnouncementService(announcementRepo)

	// ── Handlers ──────────────────────────────────────────────────────────────
	teacherH      := handlers.NewTeacherHandler(teacherSvc, log)
	classH        := handlers.NewClassHandler(classSvc, log)
	attendanceH   := handlers.NewAttendanceHandler(attendanceSvc, log)
	gradeH        := handlers.NewGradeHandler(gradeSvc, log)
	timetableH    := handlers.NewTimetableHandler(timetableSvc, log)
	announcementH := handlers.NewAnnouncementHandler(announcementSvc, log)

	// ── Teachers ──────────────────────────────────────────────────────────────
	r.Route("/teachers", func(r chi.Router) {
		r.Get("/", teacherH.List)
		r.Get("/{id}", teacherH.GetByID)
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole("admin"))
			r.Post("/", teacherH.Add)
			r.Put("/{id}", teacherH.Update)
		})
	})

	// ── Classes ───────────────────────────────────────────────────────────────
	r.Route("/classes", func(r chi.Router) {
		r.Get("/", classH.List)
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole("admin"))
			r.Post("/", classH.Create)
			r.Post("/{id}/students", classH.AssignStudent)
		})
	})

	// ── Attendance ────────────────────────────────────────────────────────────
	r.Route("/attendance", func(r chi.Router) {
		r.Get("/", attendanceH.Get)
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole("teacher"))
			r.Post("/", attendanceH.Mark)
			r.Put("/{id}", attendanceH.Edit)
		})
	})

	// ── Grades ────────────────────────────────────────────────────────────────
	r.Route("/grades", func(r chi.Router) {
		r.Get("/", gradeH.List)
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole("teacher"))
			r.Post("/", gradeH.Enter)
		})
	})

	// ── Timetable ─────────────────────────────────────────────────────────────
	r.Route("/timetable", func(r chi.Router) {
		r.Get("/{class_id}", timetableH.GetByClass)
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole("admin"))
			r.Post("/", timetableH.Create)
		})
	})

	// ── Announcements ─────────────────────────────────────────────────────────
	r.Route("/announcements", func(r chi.Router) {
		r.Get("/", announcementH.List)
		r.Group(func(r chi.Router) {
			r.Use(requireAdminOrTeacher)
			r.Post("/", announcementH.Post)
		})
	})
}
