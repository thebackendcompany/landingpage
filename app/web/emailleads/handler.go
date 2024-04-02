package emailleads

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"regexp"
	coreleads "thebackendcompany/app/core/emailleads"
	"thebackendcompany/pkg/limiters"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type EmailLeadsHandler struct {
	svc     *coreleads.EmailLeadsSheets
	limiter *limiters.SessionLimiter
}

type EmailLeadRequest struct {
	Email     string `form:"email" json:"email"`
	CsrfToken string `form:"csrf_token" json:"csrf_token"`
	TimeZone  string `form:"-" json:"-"`
}

func NewEmailLeadsHandler(svc *coreleads.EmailLeadsSheets, limiter *limiters.SessionLimiter) EmailLeadsHandler {
	// 2 requests every 10 minutes per token
	return EmailLeadsHandler{svc: svc, limiter: limiter}
}

type EmailLeadResponse struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
}

func (e EmailLeadsHandler) CreateLeadHandlerFunc(ctx *gin.Context) {
	// scheme := "http://"
	// if ctx.Request.TLS != nil {
	// 	scheme = "https://"
	// }

	// redirectToURL := scheme + ctx.Request.Host

	// splits := strings.Split(ctx.Request.URL.Path, "#")
	// if len(splits) > 1 {
	// 	redirectToURL = fmt.Sprintf("%s#%s", redirectToURL, splits[1])
	// }

	session := sessions.Default(ctx)
	_token := session.Get("csrfToken")

	csrfTokenInSession, ok := _token.(string)
	if !ok {
		log.Error().Msg("invalid csrf token in session")

		ctx.JSON(http.StatusRequestTimeout, EmailLeadResponse{
			Code: 101,
			Msg:  "Please refresh the page",
		})
		return
	}

	limiter, ok := e.limiter.GetLimiter(csrfTokenInSession)
	if ok && !limiter.Allow() {
		ctx.JSON(http.StatusTooManyRequests, EmailLeadResponse{
			Code: 107,
			Msg:  "We've already received your request. Please give us a moment to check back.",
		})
		return
	}

	var req EmailLeadRequest
	if err := ctx.Bind(&req); err != nil {
		log.Error().Err(err).Msg("failed to parse request")

		// TODO: redirect to proper error pages
		ctx.JSON(http.StatusInternalServerError, EmailLeadResponse{
			Code: 102,
			Msg:  "Oops! Something went wrong, meanwhile please use the email! To your left.",
		})
		return
	}

	if req.CsrfToken == "" || req.Email == "" {
		log.Error().Msg("missing values in req")

		// TODO: redirect to proper error pages
		ctx.JSON(http.StatusBadRequest, EmailLeadResponse{
			Code: 103,
			Msg:  "Please provide your email. If you did, please stop messing around.",
		})
		return
	}

	log.Debug().Str(
		"req csrf", req.CsrfToken,
	).Str(
		"session csrf", csrfTokenInSession,
	).Msg("compare csrf tokens")

	if req.CsrfToken != csrfTokenInSession {
		log.Error().Msg("csrf token did not match")

		// TODO: redirect to proper error pages
		ctx.JSON(http.StatusRequestTimeout, EmailLeadResponse{
			Code: 104,
			Msg:  "Please refresh the page. For security reasons.",
		})
		return
	}

	if match, _ := regexp.MatchString(`[\S]+@[\S]+\.[\S]{2,}`, req.Email); !match {
		log.Error().Msg("invalid email provided")
		ctx.JSON(http.StatusBadRequest, EmailLeadResponse{
			Code: 106,
			Msg:  "Seems like your email is a bit unique. Try mailing us? From what's left.",
		})
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
		ctx.JSON(http.StatusInternalServerError, EmailLeadResponse{
			Code: 107,
			Msg:  "Oops! Seems like something went offline. Help to your left.",
		})
		return
	}

	log.Info().Int64("update_count", updateCount).Msg("email leads inserted in table")
	ctx.JSON(http.StatusCreated, EmailLeadResponse{
		Code:    200,
		Msg:     "We will reach out shortly!",
		Success: true,
	})
}

func (e EmailLeadsHandler) LandingHandleFunc(ctx *gin.Context) {
	session := sessions.Default(ctx)

	maskedToken := session.Get("csrfToken")
	fmt.Println("dsfdsfdsfsdf ", maskedToken)

	var csrfToken string
	if maskedToken == nil {
		csrfToken = GenerateToken(32)
		session.Set("csrfToken", csrfToken)

		if err := session.Save(); err != nil {
			fmt.Println("session save error ", err)
		} else {
			e.limiter.GetLimiter(csrfToken)
		}

	} else if _t, ok := maskedToken.(string); ok {
		csrfToken = _t
	} else {
		log.Error().Msg("something went wrong trying to set csrf token. using blank")
	}

	fmt.Println("klkllkk ", csrfToken)

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"csrf_token": csrfToken,
	})
}

func GenerateToken(tokenLen int) string {
	b := make([]byte, tokenLen)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
