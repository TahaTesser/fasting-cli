package model

import "time"

// FastingSession represents a single fasting session.
type FastingSession struct {
	ID           int
	StartTime    time.Time
	Duration     time.Duration // New field for the total duration of the fast
	ProtocolName string        // Renamed from Protocol for clarity
	IsActive     bool          // New field to indicate if the session is currently active
}

// UserProfile represents the user's settings and data.

type UserProfile struct {
	Name            string
	ActiveProtocol  string
	FastingHistory  []FastingSession
}

// Timer represents the fasting timer's state.
type Timer struct {
	Duration  time.Duration
	Remaining time.Duration
	Running   bool
}
