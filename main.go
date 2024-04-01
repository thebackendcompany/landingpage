package main

import (
	"context"
	"crypto/rand"
	"embed"
	"encoding/base64"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var (
	//go:embed static/index.html
	indexHtmlFs embed.FS

	//go:embed static/404.html
	//go:embed static/favicon.ico
	//go:embed static/index.txt
	//go:embed all:static/_next
	nextFs embed.FS
)

func main() {
	distFS, err := fs.Sub(nextFs, "static/_next")
	if err != nil {
		log.Fatal(err)
	}

	baseLayout := template.Must(template.New("layout").ParseFS(indexHtmlFs, "static/index.html"))

	r := gin.Default()
	r.SetHTMLTemplate(baseLayout)

	r.StaticFS("/_next", http.FS(distFS))

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		csrfToken := generateCSRFToken()

		session.Set("csrfToken", csrfToken)

		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"csrf_token": csrfToken,
		})
	})

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

	// r.POST("/transaction", func(c *gin.Context) {
	// 	session := sessions.Default(c)
	// 	publicKey := c.GetHeader("Authorization")
	// 	csrfToken := session.Get("csrfToken")
	//
	// 	if publicKey != "" && csrfToken != nil && c.GetHeader("X-CSRF-Token") == csrfToken {
	// 		// process the transaction
	// 	} else {
	// 		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
	// 	}
	// })
	//
	// r.GET("/login", func(c *gin.Context) {
	// 	session := sessions.Default(c)
	// 	csrfToken := generateCSRFToken()
	// 	session.Set("csrfToken", csrfToken)
	// 	session.Save()
	// })

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		// service connections
		log.Println("starting server at port", ":8080")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	ctx := context.Background()

	appCtx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	<-appCtx.Done()

	if err := srv.Shutdown(appCtx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}

func generateCSRFToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
