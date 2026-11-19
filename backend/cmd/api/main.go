package main

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"

	"job-portal/internal/db"
	"job-portal/internal/db/dbgen"
	"job-portal/internal/email"
	"job-portal/internal/handler"
	"job-portal/internal/logger"
	appmw "job-portal/internal/middleware"
	"job-portal/internal/service"
	"job-portal/internal/storage"
)

func main() {
	logger.Init()

	if err := godotenv.Load(); err != nil {
		slog.Info("no .env file found, using environment")
	}

	if err := db.RunMigrations(); err != nil {
		slog.Error("migrations failed", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	pool, err := db.NewPool(ctx)
	if err != nil {
		slog.Error("db connection failed", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	queries := dbgen.New(pool)
	jwtSvc := service.NewJWTService()
	refreshSvc := service.NewRefreshTokenService(queries)
	authTokens := service.NewAuthTokenService(queries)
	emailSender := email.NewSender()
	notificationSvc := service.NewNotificationService(queries)

	appName := os.Getenv("APP_NAME")
	if appName == "" {
		appName = "Job Portal"
	}
	appBaseURL := os.Getenv("APP_BASE_URL")
	if appBaseURL == "" {
		appBaseURL = "http://localhost:5173"
		slog.Warn("APP_BASE_URL not set, defaulting to localhost")
	}
	uploadsDir := os.Getenv("UPLOADS_DIR")
	if uploadsDir == "" {
		uploadsDir = "./uploads"
	}

	fileStorage, err := storage.NewLocalStorage(uploadsDir, "/uploads")
	if err != nil {
		slog.Error("storage init failed", "error", err)
		os.Exit(1)
	}

	// Background cleanup of expired tokens every hour.
	go func() {
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			_ = queries.DeleteExpiredRefreshTokens(ctx)
			_ = queries.DeleteExpiredTokens(ctx)
		}
	}()

	// Handlers
	authHandler := handler.NewAuthHandler(queries, jwtSvc, refreshSvc, authTokens, emailSender, appBaseURL, appName)
	authVerificationHandler := handler.NewAuthVerificationHandler(queries, authTokens, emailSender, appBaseURL, appName)
	authPasswordHandler := handler.NewAuthPasswordHandler(queries, authTokens, emailSender, appBaseURL, appName)
	meHandler := handler.NewMeHandler(queries, refreshSvc)
	seekerProfileHandler := handler.NewSeekerProfileHandler(queries)
	seekerResumeHandler := handler.NewSeekerResumeHandler(queries, fileStorage)
	seekerAppsHandler := handler.NewSeekerApplicationsHandler(queries, notificationSvc)
	seekerSavedJobsHandler := handler.NewSavedJobsHandler(queries)
	seekerSavedSearchesHandler := handler.NewSavedSearchesHandler(queries)
	companyHandler := handler.NewCompanyHandler(queries)
	companyLogoHandler := handler.NewCompanyLogoHandler(queries, fileStorage)
	jobsHandler := handler.NewJobsHandler(queries)
	employerAppsHandler := handler.NewEmployerApplicationsHandler(queries, notificationSvc)
	adminStatsHandler := handler.NewAdminStatsHandler(queries)
	adminUsersHandler := handler.NewAdminUsersHandler(queries)
	adminJobsHandler := handler.NewAdminJobsHandler(queries, notificationSvc)
	adminAppsHandler := handler.NewAdminApplicationsHandler(queries)
	publicJobsHandler := handler.NewPublicJobsHandler(queries)
	publicCompaniesHandler := handler.NewPublicCompaniesHandler(queries)
	notificationsHandler := handler.NewNotificationsHandler(queries)

	// Rate limiters
	registerLimiter := appmw.NewIPRateLimiter(appmw.Limit{Rate: rate.Every(10 * time.Second), Burst: 5, TTL: 10 * time.Minute})
	loginLimiter := appmw.NewIPRateLimiter(appmw.Limit{Rate: rate.Every(2 * time.Second), Burst: 10, TTL: 10 * time.Minute})
	forgotLimiter := appmw.NewIPRateLimiter(appmw.Limit{Rate: rate.Every(5 * time.Minute), Burst: 3, TTL: 15 * time.Minute})
	resendLimiter := appmw.NewIPRateLimiter(appmw.Limit{Rate: rate.Every(20 * time.Minute), Burst: 3, TTL: 1 * time.Hour})
	uploadLimiter := appmw.NewIPRateLimiter(appmw.Limit{Rate: rate.Every(6 * time.Second), Burst: 10, TTL: 10 * time.Minute})

	// Echo
	e := echo.New()
	e.Use(appmw.RequestLogger())
	e.Use(echomw.Recover())

	originsRaw := os.Getenv("CORS_ORIGINS")
	if originsRaw == "" {
		originsRaw = "http://localhost:5173"
	}
	allowedOrigins := make([]string, 0)
	for _, o := range strings.Split(originsRaw, ",") {
		if o = strings.TrimSpace(o); o != "" {
			allowedOrigins = append(allowedOrigins, o)
		}
	}
	slog.Info("CORS configured", "origins", allowedOrigins)

	e.Use(echomw.CORSWithConfig(echomw.CORSConfig{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "X-CSRF-Token"},
		AllowCredentials: true, // ← critical for cookies
	}))

	e.Static("/uploads", uploadsDir)

	api := e.Group("/api")

	api.POST("/register", authHandler.Register, registerLimiter.Middleware())
	api.POST("/login", authHandler.Login, loginLimiter.Middleware())
	api.POST("/auth/logout", authHandler.Logout)
	api.POST("/auth/refresh", authHandler.Refresh)

	api.POST("/auth/verify-email", authVerificationHandler.VerifyEmail)
	api.POST("/auth/forgot-password", authPasswordHandler.ForgotPassword, forgotLimiter.Middleware())
	api.POST("/auth/reset-password", authPasswordHandler.ResetPassword)

	api.GET("/jobs", publicJobsHandler.List)
	api.GET("/jobs/:id", publicJobsHandler.Detail)
	api.GET("/companies", publicCompaniesHandler.List)
	api.GET("/companies/:id", publicCompaniesHandler.Detail)

	// Protected (auth + CSRF for state-changing methods)
	protected := api.Group("", appmw.JWTAuth(jwtSvc), appmw.CSRFProtect())
	protected.GET("/me", authHandler.Me)
	protected.POST("/me/password", meHandler.ChangePassword)
	protected.POST("/onboarding/role", authHandler.OnboardingRole)
	protected.POST("/auth/resend-verification", authVerificationHandler.ResendVerification, resendLimiter.Middleware())

	protected.GET("/notifications", notificationsHandler.List)
	protected.GET("/notifications/recent", notificationsHandler.Recent)
	protected.GET("/notifications/unread-count", notificationsHandler.Count)
	protected.POST("/notifications/read-all", notificationsHandler.MarkAllRead)
	protected.POST("/notifications/:id/read", notificationsHandler.MarkRead)

	seekerOnly := protected.Group("/seeker", appmw.RequireRole("job_seeker"))
	seekerOnly.GET("/profile", seekerProfileHandler.Get)
	seekerOnly.PUT("/profile", seekerProfileHandler.Upsert)
	seekerOnly.POST("/profile/resume", seekerResumeHandler.Upload, uploadLimiter.Middleware())
	seekerOnly.DELETE("/profile/resume", seekerResumeHandler.Delete)
	seekerOnly.POST("/applications", seekerAppsHandler.Create)
	seekerOnly.GET("/applications", seekerAppsHandler.List)
	seekerOnly.GET("/applications/status", seekerAppsHandler.HasApplied)
	seekerOnly.POST("/applications/:id/withdraw", seekerAppsHandler.Withdraw)
	seekerOnly.POST("/saved-jobs/:job_id", seekerSavedJobsHandler.Save)
	seekerOnly.DELETE("/saved-jobs/:job_id", seekerSavedJobsHandler.Unsave)
	seekerOnly.GET("/saved-jobs/:job_id", seekerSavedJobsHandler.Check)
	seekerOnly.GET("/saved-jobs", seekerSavedJobsHandler.List)
	seekerOnly.POST("/saved-searches", seekerSavedSearchesHandler.Create)
	seekerOnly.GET("/saved-searches", seekerSavedSearchesHandler.List)
	seekerOnly.DELETE("/saved-searches/:id", seekerSavedSearchesHandler.Delete)

	employerOnly := protected.Group("/employer", appmw.RequireRole("employer"))
	employerOnly.GET("/company", companyHandler.Get)
	employerOnly.PUT("/company", companyHandler.Upsert)
	employerOnly.POST("/company/logo", companyLogoHandler.Upload, uploadLimiter.Middleware())
	employerOnly.DELETE("/company/logo", companyLogoHandler.Delete)
	employerOnly.POST("/jobs", jobsHandler.Create)
	employerOnly.GET("/jobs", jobsHandler.List)
	employerOnly.GET("/jobs/:id", jobsHandler.Get)
	employerOnly.PUT("/jobs/:id", jobsHandler.Update)
	employerOnly.POST("/jobs/:id/submit", jobsHandler.Submit)
	employerOnly.POST("/jobs/:id/close", jobsHandler.Close)
	employerOnly.POST("/jobs/:id/archive", jobsHandler.Archive)
	employerOnly.GET("/applications", employerAppsHandler.List)
	employerOnly.GET("/jobs/:id/applications", employerAppsHandler.ListByJob)
	employerOnly.GET("/applications/:id", employerAppsHandler.Detail)
	employerOnly.POST("/applications/:id/status", employerAppsHandler.UpdateStatus)

	adminOnly := protected.Group("/admin", appmw.RequireRole("admin"))
	adminOnly.GET("/stats", adminStatsHandler.Stats)
	adminOnly.GET("/jobs", adminJobsHandler.List)
	adminOnly.GET("/jobs/:id", adminJobsHandler.Get)
	adminOnly.POST("/jobs/:id/approve", adminJobsHandler.Approve)
	adminOnly.POST("/jobs/:id/reject", adminJobsHandler.Reject)
	adminOnly.GET("/users", adminUsersHandler.List)
	adminOnly.GET("/users/:id", adminUsersHandler.Detail)
	adminOnly.GET("/applications", adminAppsHandler.List)
	adminOnly.GET("/applications/:id", adminAppsHandler.Detail)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("server starting", "port", port)
	e.Logger.Fatal(e.Start(":" + port))
}
