package eventshandler

import (
	"net/http"
	"time"

	"thebackendcompany/app/core/events"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const DefaultLocation = "UTC"

type EventCreateRequest struct {
	Name string `json:"name" form:"name"`
	Data string `json:"data" form:"data"`
}

func HandleEventsCreate(svc *events.EventService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req EventCreateRequest

		if err := c.Bind(&req); err != nil {
			log.Error().Err(err).Msg("failed to parse request body")

			c.JSON(
				http.StatusBadRequest,
				`{"success": false, "error": "invalid request body"}`,
			)
			return
		}

		tzLocation := c.Request.Header.Get("location")
		tz, err := time.LoadLocation(tzLocation)
		if err != nil {
			tz, err = time.LoadLocation(DefaultLocation)
		}

		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				`{"success": false, "error": "failed to get timezone"}`,
			)
			return
		}

		// Adapters/events.go
		event := &events.Event{
			Name:      req.Name,
			Data:      []byte(req.Data),
			CreatedAt: time.Now().In(tz),
			UpdatedAt: time.Now().In(tz),
		}

		if err := svc.Create(c.Request.Context(), event); err != nil {
			log.Error().Err(err).Msg("failed to upload to upstash")

			c.JSON(
				http.StatusInternalServerError,
				`{"success": false, "error": "failed to create event"}`,
			)
			return
		}

		log.Info().Msg("pushed to kafka upstash")
		c.JSON(
			http.StatusCreated,
			`{"success": true, "message": "event created"}`,
		)
		return
	}
}
