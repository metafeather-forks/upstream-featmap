package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/amborle/featmap/migrations"
	"github.com/amborle/featmap/webapp"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
)

// Configuration ...
type Configuration struct {
	Environment         string `json:"environment"`
	Mode                string `json:"mode"`
	AppSiteURL          string `json:"appSiteURL"`
	DbConnectionString  string `json:"dbConnectionString"`
	JWTSecret           string `json:"jwtSecret"`
	Port                string `json:"port"`
	EmailFrom           string `json:"emailFrom"`
	SMTPServer          string `json:"smtpServer"`
	SMTPPort            string `json:"smtpPort"`
	SMTPUser            string `json:"smtpUser"`
	SMTPPass            string `json:"smtpPass"`
}

func main() {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// r.Use(middleware.SetHeader("Content-Type", "application/json"))

	config, err := readConfiguration()
	if err != nil {
		slog.Error("configuration error", "error", err)
		os.Exit(1)
	}

	// CORS
	corsConfiguration := cors.New(cors.Options{
		AllowedOrigins:   []string{config.AppSiteURL, "http://localhost:3000"}, // localhost is for development work
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Workspace", "X-CSRF-Token"},
		ExposedHeaders:   []string{""},
		AllowCredentials: true,
		MaxAge:           300,
	})

	r.Use(corsConfiguration.Handler)

	db, err := sqlx.Connect("postgres", config.DbConnectionString)
	if err != nil {
		slog.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("database close error", "error", err)
		}
	}()

	// Apply migrations
	d, err := iofs.New(migrations.FS, ".")
	if err != nil {
		slog.Error("migration source error", "error", err)
		os.Exit(1)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, config.DbConnectionString)
	if err != nil {
		slog.Error("migration init error", "error", err)
		os.Exit(1)
	}

	m.Up()

	// Create JWTAuth object
	auth := jwtauth.New("HS256", []byte(config.JWTSecret), nil)

	r.Use(jwtauth.Verifier(auth))
	r.Use(ContextSkeleton(config))

	r.Use(Transaction(db))
	r.Use(Auth(auth))

	r.Use(User())

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1/users", usersAPI)               // Nothing is needed
	r.Route("/v1/link", linkAPI)                 // Nothing is needed

	r.Get("/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if db == nil || db.Ping() != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"unhealthy"}`))
			return
		}
		w.Write([]byte(`{"status":"healthy"}`))
	})

	r.Route("/v1/account", accountAPI) // Account needed
	r.Route("/v1/", workspaceAPI)      // Account + workspace is needed

	buildFS, err := fs.Sub(webapp.FS, "build")
	if err != nil {
		slog.Error("embedded filesystem error", "error", err)
		os.Exit(1)
	}

	// Static files + SPA fallback: serve files from build/ if they exist,
	// otherwise serve index.html for client-side routing.
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		f, err := buildFS.Open(path)
		if err != nil {
			index, _ := webapp.FS.ReadFile("build/index.html")
			http.ServeContent(w, r, "index.html", time.Now(), strings.NewReader(string(index)))
			return
		}
		f.Close()
		http.FileServer(http.FS(buildFS)).ServeHTTP(w, r)
	})

	slog.Info("starting server", "port", config.Port)

	srv := &http.Server{
		Addr:    ":" + config.Port,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("forced shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped")

}

func readConfiguration() (Configuration, error) {
	c := Configuration{}

	// Read from environment variables first, fall back to conf.json
	if env := os.Getenv("PORT"); env != "" {
		c.Port = env
	}
	if env := os.Getenv("DATABASE_URL"); env != "" {
		c.DbConnectionString = env
	}
	if env := os.Getenv("JWT_SECRET"); env != "" {
		c.JWTSecret = env
	}
	if env := os.Getenv("APP_SITE_URL"); env != "" {
		c.AppSiteURL = env
	}
	if env := os.Getenv("ENVIRONMENT"); env != "" {
		c.Environment = env
	}
	if env := os.Getenv("MODE"); env != "" {
		c.Mode = env
	}
	if env := os.Getenv("EMAIL_FROM"); env != "" {
		c.EmailFrom = env
	}
	if env := os.Getenv("SMTP_SERVER"); env != "" {
		c.SMTPServer = env
	}
	if env := os.Getenv("SMTP_PORT"); env != "" {
		c.SMTPPort = env
	}
	if env := os.Getenv("SMTP_USER"); env != "" {
		c.SMTPUser = env
	}
	if env := os.Getenv("SMTP_PASS"); env != "" {
		c.SMTPPass = env
	}

	// Fall back to conf.json for any unset values
	file, err := os.Open("conf.json")
	if err == nil {
		defer file.Close()
		var fileConfig Configuration
		if decodeErr := json.NewDecoder(file).Decode(&fileConfig); decodeErr == nil {
			if c.Port == "" {
				c.Port = fileConfig.Port
			}
			if c.DbConnectionString == "" {
				c.DbConnectionString = fileConfig.DbConnectionString
			}
			if c.JWTSecret == "" {
				c.JWTSecret = fileConfig.JWTSecret
			}
			if c.AppSiteURL == "" {
				c.AppSiteURL = fileConfig.AppSiteURL
			}
			if c.Environment == "" {
				c.Environment = fileConfig.Environment
			}
			if c.Mode == "" {
				c.Mode = fileConfig.Mode
			}
			if c.EmailFrom == "" {
				c.EmailFrom = fileConfig.EmailFrom
			}
			if c.SMTPServer == "" {
				c.SMTPServer = fileConfig.SMTPServer
			}
			if c.SMTPPort == "" {
				c.SMTPPort = fileConfig.SMTPPort
			}
			if c.SMTPUser == "" {
				c.SMTPUser = fileConfig.SMTPUser
			}
			if c.SMTPPass == "" {
				c.SMTPPass = fileConfig.SMTPPass
			}
		}
	}

	// Defaults
	if c.SMTPPort == "" {
		c.SMTPPort = "587"
	}
	if c.Port == "" {
		c.Port = "5000"
	}

	// Validate required fields
	if c.DbConnectionString == "" {
		return c, fmt.Errorf("DATABASE_URL or conf.json required")
	}
	if c.JWTSecret == "" {
		return c, fmt.Errorf("JWT_SECRET or conf.json required")
	}

	return c, nil
}
