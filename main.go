package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/stripe/stripe-go"

	"github.com/amborle/featmap/migrations"
	"github.com/amborle/featmap/webapp"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
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
	StripeKey           string `json:"stripeKey"`
	StripeWebhookSecret string `json:"stripeWebhookSecret"`
	StripeBasicPlan     string `json:"stripeBasicPlan"`
	StripeProPlan       string `json:"stripeProPlan"`
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
		log.Fatalln("no conf.json found")
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
		log.Fatalln("database error:" + err.Error())
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// Apply migrations
	d, err := iofs.New(migrations.FS, ".")
	if err != nil {
		log.Fatalln(err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, config.DbConnectionString)
	if err != nil {
		log.Fatalln(err)
	}

	m.Up()

	// Create JWTAuth object
	auth := jwtauth.New("HS256", []byte(config.JWTSecret), nil)

	r.Use(jwtauth.Verifier(auth))
	r.Use(ContextSkeleton(config))

	r.Use(Transaction(db))
	r.Use(Auth(auth))

	r.Use(User())

	stripe.Key = config.StripeKey

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1/users", usersAPI)               // Nothing is needed
	r.Route("/v1/link", linkAPI)                 // Nothing is needed
	r.Route("/v1/subscription", subscriptionAPI) // Nothing is needed

	r.Route("/v1/account", accountAPI) // Account needed
	r.Route("/v1/", workspaceAPI)      // Account + workspace is needed

	buildFS, err := fs.Sub(webapp.FS, "build")
	if err != nil {
		log.Fatalln(err)
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

	fmt.Println("Serving on port " + config.Port)
	err = http.ListenAndServe(":"+config.Port, r)
	if err != nil {
		log.Fatalln(err)
	}

}

func readConfiguration() (Configuration, error) {
	file, err := os.Open("conf.json")

	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()

	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)

	if configuration.SMTPPort == "" {
		configuration.SMTPPort = "587"
	}

	return configuration, err
}
