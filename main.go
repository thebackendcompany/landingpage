package main

import (
	"crypto/rand"
	"embed"
	"encoding/base64"
	"html/template"
	"io/fs"
	"log"
	"net/http"

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

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.StaticFS("/_next", http.FS(distFS))
	r.SetHTMLTemplate(baseLayout)

	r.GET("/", func(ctx *gin.Context) {
		csrfToken := generateCSRFToken()

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

	r.POST("/transaction", func(c *gin.Context) {
		session := sessions.Default(c)
		publicKey := c.GetHeader("Authorization")
		csrfToken := session.Get("csrfToken")

		if publicKey != "" && csrfToken != nil && c.GetHeader("X-CSRF-Token") == csrfToken {
			// process the transaction
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		}
	})

	r.GET("/login", func(c *gin.Context) {
		session := sessions.Default(c)
		csrfToken := generateCSRFToken()
		session.Set("csrfToken", csrfToken)
		session.Save()
	})

	r.Run()
}

func generateCSRFToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
