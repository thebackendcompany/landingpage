package limiters

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

type SessionLimiter struct {
	Store           map[string]*LimitMetadata
	mu              *sync.RWMutex
	limit           rate.Limit
	burstSize       int
	cleanupInterval time.Duration
}

type LimitMetadata struct {
	key       string
	createdAt time.Time
	limiter   *rate.Limiter
}

func NewSessionLimiter(limit rate.Limit, bustSize int, cleanupInterval time.Duration) *SessionLimiter {
	return &SessionLimiter{
		Store:           make(map[string]*LimitMetadata),
		mu:              &sync.RWMutex{},
		limit:           limit,
		burstSize:       bustSize,
		cleanupInterval: cleanupInterval,
	}
}

func (sessionlimiter *SessionLimiter) AddToken(token string) *rate.Limiter {
	log.Info().Msg("adding token for limiter ")

	limiter := rate.NewLimiter(sessionlimiter.limit, sessionlimiter.burstSize)
	sessionlimiter.mu.Lock()
	defer sessionlimiter.mu.Unlock()

	sessionlimiter.Store[token] = &LimitMetadata{
		limiter:   limiter,
		key:       token,
		createdAt: time.Now().UTC(),
	}

	return limiter
}

func (sessionlimiter *SessionLimiter) CleanupExpired(metadata *LimitMetadata) (cleanedUp bool) {
	sessionlimiter.mu.RLock()

	if metadata.createdAt.Add(sessionlimiter.cleanupInterval).Before(time.Now().UTC()) {
		log.Info().Msg("cleanup rate limiter for token")

		sessionlimiter.mu.RUnlock()
		sessionlimiter.Deregister(metadata.key)
		return true
	}

	log.Info().Msg("rate limiter yet to expire for token")
	sessionlimiter.mu.RUnlock()
	return false
}

func (sessionlimiter *SessionLimiter) GetLimiter(token string) (*rate.Limiter, bool) {
	log.Info().Msg("getting token for limiter ")

	sessionlimiter.mu.RLock()

	data, ok := sessionlimiter.Store[token]
	if !ok {
		sessionlimiter.mu.RUnlock() // because add token will locak again
		return sessionlimiter.AddToken(token), true
	}

	sessionlimiter.mu.RUnlock()
	// if sessionlimiter.CleanupExpired(data) {
	// 	return nil, false
	// }
	log.Info().Msg("limiter returned from store")
	return data.limiter, true
}

func (sessionlimiter *SessionLimiter) Deregister(token string) {
	sessionlimiter.mu.Lock()
	defer sessionlimiter.mu.Unlock()

	delete(sessionlimiter.Store, token)
}
