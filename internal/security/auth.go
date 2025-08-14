package security

import (
	"crypto/subtle"
	"errors"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrSessionExpired     = errors.New("session expired")
)

// Session represents an application session
type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// SessionManager handles application sessions
type SessionManager struct {
	sessions       map[string]*Session
	sessionTimeout time.Duration
}

// NewSessionManager creates a new session manager
func NewSessionManager(timeout time.Duration) *SessionManager {
	return &SessionManager{
		sessions:       make(map[string]*Session),
		sessionTimeout: timeout,
	}
}

// CreateSession creates a new session for a user
func (sm *SessionManager) CreateSession(userID string) (*Session, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, err
	}

	session := &Session{
		ID:        sessionID,
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(sm.sessionTimeout),
	}

	sm.sessions[sessionID] = session
	return session, nil
}

// ValidateSession validates a session and returns the session if valid
func (sm *SessionManager) ValidateSession(sessionID string) (*Session, error) {
	session, exists := sm.sessions[sessionID]
	if !exists {
		return nil, ErrInvalidCredentials
	}

	if time.Now().After(session.ExpiresAt) {
		delete(sm.sessions, sessionID)
		return nil, ErrSessionExpired
	}

	return session, nil
}

// RefreshSession extends the session expiration time
func (sm *SessionManager) RefreshSession(sessionID string) error {
	session, err := sm.ValidateSession(sessionID)
	if err != nil {
		return err
	}

	session.ExpiresAt = time.Now().Add(sm.sessionTimeout)
	return nil
}

// DestroySession removes a session
func (sm *SessionManager) DestroySession(sessionID string) {
	delete(sm.sessions, sessionID)
}

// CleanupExpiredSessions removes expired sessions
func (sm *SessionManager) CleanupExpiredSessions() {
	now := time.Now()
	for id, session := range sm.sessions {
		if now.After(session.ExpiresAt) {
			delete(sm.sessions, id)
		}
	}
}

// generateSessionID generates a secure session ID
func generateSessionID() (string, error) {
	key, err := GenerateKey()
	if err != nil {
		return "", err
	}

	// Convert to hex string for session ID
	sessionID := ""
	for _, b := range key[:16] { // Use first 16 bytes for session ID
		sessionID += string(rune('a' + (b % 26)))
	}

	return sessionID, nil
}

// SecureCompare performs a constant-time comparison of two strings
func SecureCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
