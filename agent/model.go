package agent

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/openzipkin/zipkin-go/model"
)

type ModelSpanModel model.SpanModel

type Span struct {
	ModelSpanModel
	Type    string `json:"type,omitempty"`
	Service string `json:"service,omitempty"`
	System  string `json:"system,omitempty"`
}

// MarshalJSON exports our Model into the correct format for the Zipkin V2 API.
func (s Span) MarshalJSON() ([]byte, error) {
	type Alias Span

	var timestamp int64
	if !s.Timestamp.IsZero() {
		if s.Timestamp.Unix() < 1 {
			// Zipkin does not allow Timestamps before Unix epoch
			return nil, model.ErrValidTimestampRequired
		}
		timestamp = s.Timestamp.Round(time.Microsecond).UnixNano() / 1e3
	}

	if s.Duration < time.Microsecond {
		if s.Duration < 0 {
			// negative duration is not allowed and signals a timing logic error
			return nil, model.ErrValidDurationRequired
		} else if s.Duration > 0 {
			// sub microsecond durations are reported as 1 microsecond
			s.Duration = 1 * time.Microsecond
		}
	} else {
		// Duration will be rounded to nearest microsecond representation.
		//
		// NOTE: Duration.Round() is not available in Go 1.8 which we still support.
		// To handle microsecond resolution rounding we'll add 500 nanoseconds to
		// the duration. When truncated to microseconds in the call to marshal, it
		// will be naturally rounded. See TestSpanDurationRounding in span_test.go
		s.Duration += 500 * time.Nanosecond
	}

	s.Name = strings.ToLower(s.Name)

	if s.LocalEndpoint.Empty() {
		s.LocalEndpoint = nil
	}

	if s.RemoteEndpoint.Empty() {
		s.RemoteEndpoint = nil
	}

	return json.Marshal(&struct {
		T int64 `json:"timestamp,omitempty"`
		D int64 `json:"duration,omitempty"`
		Alias
	}{
		T:     timestamp,
		D:     s.Duration.Nanoseconds() / 1e3,
		Alias: (Alias)(s),
	})
}

// UnmarshalJSON imports our Model from a Zipkin V2 API compatible span
// representation.
func (s *Span) UnmarshalJSON(b []byte) error {
	type Alias Span
	span := &struct {
		T uint64 `json:"timestamp,omitempty"`
		D uint64 `json:"duration,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(b, &span); err != nil {
		return err
	}
	if s.ID < 1 {
		return model.ErrValidIDRequired
	}
	if span.T > 0 {
		s.Timestamp = time.Unix(0, int64(span.T)*1e3)
	}
	s.Duration = time.Duration(span.D*1e3) * time.Nanosecond
	if s.LocalEndpoint.Empty() {
		s.LocalEndpoint = nil
	}

	if s.RemoteEndpoint.Empty() {
		s.RemoteEndpoint = nil
	}
	return nil
}
