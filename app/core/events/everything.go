package events

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-batteries/snowflake"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

const (
	CreateSchema = `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY UNIQUE
		,event_name VARCHAR(50) NOT NULL
		,event_data TEXT
		,created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		,updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`

	InsertEventQuery = `
	INSERT INTO events (id, event_name, event_data) 
	VALUES (:id, :event_name, :event_data)
	`

	SelectEventQuery = `
	SELECT
		id
		,event_name
		,event_data
		,created_at
	FROM	
	  events
	`
)

type Event struct {
	ID   int64  `db:"id"`
	Name string `db:"event_name"`
	Data []byte `db:"event_data"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type EventRepo struct {
	db    *sqlx.DB
	idgen *snowflake.SequenceGenerator
}

func NewEventRepo(db *sqlx.DB) *EventRepo {
	return &EventRepo{
		db:    db,
		idgen: snowflake.NewSequenceGenerator(),
	}
}

func (e EventRepo) Migrate(ctx context.Context) error {
	_, err := e.db.ExecContext(ctx, CreateSchema)
	return err
}

func (e EventRepo) Create(ctx context.Context, event *Event) error {
	event.ID = e.idgen.NextID()

	_, err := e.db.NamedExecContext(ctx, InsertEventQuery, event)
	return err
}

type EventQuery struct {
	Name string `db:"event_name"`
	ID   int64  `db:"id"`
}

func (e EventRepo) Where(ctx context.Context, eq *EventQuery) ([]*Event, error) {
	query := SelectEventQuery

	if eq != nil {
		query = fmt.Sprintf("%s WHERE name=:name", SelectEventQuery)
	}

	events := []*Event{}

	err := e.db.SelectContext(ctx, &events, query, eq)
	return events, err
}

type EventService struct {
	// We need a kafka producer to send events
	// repo to save events to database
	repo     *EventRepo
	producer *EventsKafkaMq
}

func NewEventService(repo *EventRepo, producer *EventsKafkaMq) *EventService {
	return &EventService{
		repo:     repo,
		producer: producer,
	}
}

var ErrKafkaEventProducerFailed = errors.New("failed_to_produce_kafka_events")

func (es *EventService) Create(
	ctx context.Context,
	ev *Event,
) error {
	log.Info().Msg("creating event in kofka")

	// hasError := false
	//
	// // Handle partial failures better, by returning
	// for result := range es.producer.GetResults(ctx) {
	// 	if result.Err != nil {
	// 		hasError = true
	// 		log.Error().Err(result.Err).Msg("failed to push to kafka")
	// 	}
	// }
	//
	// if hasError {
	// 	return ErrKafkaEventProducerFailed
	// }

	if err := es.repo.Create(ctx, ev); err != nil {
		log.Error().Err(err).Msg("failed to save event to db")
		return err
	}

	if err := es.producer.PushEvent(ctx, ev); err != nil {
		log.Error().Err(err).Msg("failed to push events to producer")
		return err
	}

	return nil
}
