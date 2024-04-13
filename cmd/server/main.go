package main

import (
	"context"
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"thebackendcompany/app/web/emailleads"
	"thebackendcompany/pkg/config"
	"thebackendcompany/pkg/limiters"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	coreleads "thebackendcompany/app/core/emailleads"

	"github.com/rs/zerolog/log"

	"github.com/jmoiron/sqlx"
)

func RunMigrations(db *sqlx.DB) {
	log.Info().Msg("running migrations")
}

var (
	//go:embed static/*.html
	indexHtmlFs embed.FS

	//go:embed static/favicon.ico
	//go:embed all:static/_next
	nextFs embed.FS

	//go:embed all:static/images
	imageFs embed.FS
)

const (
	CSRFTokenHeader = "X-CSRF-Token"
)

func main() {
	env := os.Getenv("ENVIRONMENT")
	env = strings.ToLower(env)

	// runMigrateFlag := flag.Bool("migrate", false, "run migrations")
	// flag.Parse()

	if env == "" {
		env = "local"
	}

	cfg := config.BuildAppConfig(env)

	// db, err := config.ConnectDB("sqlite3", "sink.db", "sink.db")
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("failed to connect to db")
	// }

	// if *runMigrateFlag {
	// 	RunMigrations(db.GetDB())
	// }

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = ":8080"
	}

	distImgFS, err := fs.Sub(imageFs, "static/images")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load static images")
	}

	distFS, err := fs.Sub(nextFs, "static/_next")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load static templates")
	}

	baseLayout := template.Must(template.New("layout").ParseFS(indexHtmlFs, "static/*.html"))

	r := gin.Default()
	r.SetHTMLTemplate(baseLayout)

	r.StaticFS("/_next", http.FS(distFS))
	r.StaticFS("/images", http.FS(distImgFS))

	sessionKey := os.Getenv("SESSION_SECRET")
	if sessionKey == "" {
		sessionKey = emailleads.GenerateToken(64)
	}

	emailLeadsSvc := coreleads.NewEmailLeadsSheets(
		[]byte(cfg.GoogleCreds),
		cfg.EmailLeadsDbName,
	)

	// 2 request, every 10 seconds
	csrfLimiter := limiters.NewSessionLimiter(rate.Every(10*time.Minute), 2, 15*time.Minute)

	store := cookie.NewStore([]byte(sessionKey))
	shouldSecure := true

	if env == "local" {
		shouldSecure = false
	}

	store.Options(sessions.Options{
		HttpOnly: true,
		Secure:   shouldSecure,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   1200,
		Domain:   cfg.DomainName,
	})

	leadsHandler := emailleads.NewEmailLeadsHandler(emailLeadsSvc, csrfLimiter)

	r.Use(sessions.Sessions("thebackendcompany", store))

	r.GET("/", leadsHandler.LandingHandleFunc)

	r.GET("favicon.ico", func(ctx *gin.Context) {
		file, err := nextFs.ReadFile("static/favicon.ico")
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Oops! Server Crashed")
			ctx.Error(err)
			ctx.Abort()
			return
		}

		ctx.Data(
			http.StatusOK,
			"image/x-icon",
			file,
		)
	})

	r.GET("/tbc/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.POST("/tbc/emails/leads", leadsHandler.CreateLeadHandlerFunc)

	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	go func() {
		// service connections
		log.Info().Str("port", port).Msg("starting server")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	ctx := context.Background()

	appCtx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	<-appCtx.Done()

	if err := srv.Shutdown(appCtx); err != nil {
		log.Fatal().Err(err).Msg("Server Shutdown")
	}

	log.Info().Msg("Server exiting")
}
