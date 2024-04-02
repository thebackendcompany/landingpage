package emailleads

import (
	"net/http"
	"regexp"
	coreleads "thebackendcompany/app/core/emailleads"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type EmailLeadsHandler struct {
	svc *coreleads.EmailLeadsSheets
}

type EmailLeadRequest struct {
	Email     string `form:"email" json:"email"`
	CsrfToken string `form:"csrf_token" json:"csrf_token"`
	TimeZone  string `form:"-" json:"-"`
}

func NewEmailLeadsHandler(svc *coreleads.EmailLeadsSheets) EmailLeadsHandler {
	return EmailLeadsHandler{svc}
}

func (e EmailLeadsHandler) HandlerFunc(ctx *gin.Context) {
	scheme := "http://"
	if ctx.Request.TLS != nil {
		scheme = "https://"
	}

	redirectToURL := scheme + ctx.Request.Host

	session := sessions.Default(ctx)
	_token := session.Get("csrfToken")

	csrfTokenInSession, ok := _token.(string)
	if !ok {
		log.Error().Msg("invalid csrf token in session")

		ctx.Redirect(http.StatusFound, redirectToURL)
		return
	}

	var req EmailLeadRequest
	if err := ctx.Bind(&req); err != nil {
		log.Error().Err(err).Msg("failed to parse request")

		// TODO: redirect to proper error pages
		ctx.Redirect(http.StatusFound, redirectToURL)
		return
	}

	if req.CsrfToken == "" || req.Email == "" {
		log.Error().Msg("missing values in req")

		// TODO: redirect to proper error pages
		ctx.Redirect(http.StatusFound, redirectToURL)
		return
	}

	log.Debug().Str(
		"req csrf", req.CsrfToken,
	).Str(
		"session csrf", csrfTokenInSession,
	).Msg("mismatch csrf tokens")

	if req.CsrfToken != csrfTokenInSession {
		log.Error().Msg("csrf token did not match")

		// TODO: redirect to proper error pages
		ctx.Redirect(http.StatusFound, redirectToURL)
		return
	}

	if match, _ := regexp.MatchString(`[\S]+@[\S]+\.[\S]{2,}`, req.Email); !match {
		log.Error().Msg("invalid email provided")
		ctx.Redirect(http.StatusFound, redirectToURL)
		return
	}

	updateCount, err := e.svc.Write(ctx.Request.Context(), &coreleads.SheetsEmailLeads{
		Email:       req.Email,
		Body:        "Get in Touch",
		ContactedAt: time.Now(), // try to get location
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to write to email leads sheet")

		// TODO: redirect to proper error pages
		ctx.Redirect(http.StatusFound, redirectToURL)
		return
	}

	log.Info().Int64("update_count", updateCount).Msg("email leads inserted in table")
	ctx.Redirect(http.StatusFound, redirectToURL)
}
