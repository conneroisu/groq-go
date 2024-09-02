package groq

import (
	"net/http"
	"strconv"
	"time"
)

// RateLimitHeaders struct represents Groq rate limits headers.
type RateLimitHeaders struct {
	LimitRequests     int       `json:"x-ratelimit-limit-requests"`
	LimitTokens       int       `json:"x-ratelimit-limit-tokens"`
	RemainingRequests int       `json:"x-ratelimit-remaining-requests"`
	RemainingTokens   int       `json:"x-ratelimit-remaining-tokens"`
	ResetRequests     ResetTime `json:"x-ratelimit-reset-requests"`
	ResetTokens       ResetTime `json:"x-ratelimit-reset-tokens"`
}

// ResetTime is a time.Time wrapper for the rate limit reset time.
// string
type ResetTime string

// String returns the string representation of the ResetTime.
func (r ResetTime) String() string {
	return string(r)
}

// Time returns the time.Time representation of the ResetTime.
func (r ResetTime) Time() time.Time {
	d, _ := time.ParseDuration(string(r))
	return time.Now().Add(d)
}

// newRateLimitHeaders creates a new RateLimitHeaders from an http.Header.
func newRateLimitHeaders(h http.Header) RateLimitHeaders {
	limitReq, _ := strconv.Atoi(h.Get("x-ratelimit-limit-requests"))
	limitTokens, _ := strconv.Atoi(h.Get("x-ratelimit-limit-tokens"))
	remainingReq, _ := strconv.Atoi(h.Get("x-ratelimit-remaining-requests"))
	remainingTokens, _ := strconv.Atoi(h.Get("x-ratelimit-remaining-tokens"))
	return RateLimitHeaders{
		LimitRequests:     limitReq,
		LimitTokens:       limitTokens,
		RemainingRequests: remainingReq,
		RemainingTokens:   remainingTokens,
		ResetRequests:     ResetTime(h.Get("x-ratelimit-reset-requests")),
		ResetTokens:       ResetTime(h.Get("x-ratelimit-reset-tokens")),
	}
}
