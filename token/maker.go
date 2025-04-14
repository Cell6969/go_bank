package token

import "time"

// Maker is an interface for managing token
type Maker interface {
	// Create Token for specific username and duration
	CreateToken(username string, duration time.Duration) (string, error)

	VerifyToken(token string) (*Payload, error)
}
