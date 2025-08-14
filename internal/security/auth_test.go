package security

import (
	"testing"
	"time"
)

func TestNewSessionManager(t *testing.T) {
	timeout := 30 * time.Minute
	sm := NewSessionManager(timeout)

	if sm == nil {
		t.Fatal("Expected session manager to be created")
	}

	if sm.sessionTimeout != timeout {
		t.Errorf("Expected timeout %v, got %v", timeout, sm.sessionTimeout)
	}

	if sm.sessions == nil {
		t.Error("Expected sessions map to be initialized")
	}
}

func TestCreateSession(t *testing.T) {
	sm := NewSessionManager(30 * time.Minute)
	userID := "test-user-123"

	session, err := sm.CreateSession(userID)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	if session == nil {
		t.Fatal("Expected session to be created")
	}

	if session.UserID != userID {
		t.Errorf("Expected user ID %s, got %s", userID, session.UserID)
	}

	if session.ID == "" {
		t.Error("Expected session ID to be set")
	}

	if session.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}

	if session.ExpiresAt.IsZero() {
		t.Error("Expected ExpiresAt to be set")
	}

	if session.ExpiresAt.Before(session.CreatedAt) {
		t.Error("Expected ExpiresAt to be after CreatedAt")
	}

	// Check that session is stored in manager
	storedSession, exists := sm.sessions[session.ID]
	if !exists {
		t.Error("Expected session to be stored in manager")
	}

	if storedSession != session {
		t.Error("Expected stored session to match created session")
	}
}

func TestValidateSession(t *testing.T) {
	sm := NewSessionManager(30 * time.Minute)
	userID := "test-user-123"

	// Create a session
	session, err := sm.CreateSession(userID)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// Test valid session
	validatedSession, err := sm.ValidateSession(session.ID)
	if err != nil {
		t.Fatalf("Failed to validate session: %v", err)
	}

	if validatedSession.ID != session.ID {
		t.Errorf("Expected session ID %s, got %s", session.ID, validatedSession.ID)
	}

	// Test non-existent session
	_, err = sm.ValidateSession("non-existent-session")
	if err != ErrInvalidCredentials {
		t.Errorf("Expected ErrInvalidCredentials, got %v", err)
	}

	// Test expired session
	expiredSession, _ := sm.CreateSession("expired-user")
	expiredSession.ExpiresAt = time.Now().Add(-1 * time.Hour) // Set to past
	sm.sessions[expiredSession.ID] = expiredSession

	_, err = sm.ValidateSession(expiredSession.ID)
	if err != ErrSessionExpired {
		t.Errorf("Expected ErrSessionExpired, got %v", err)
	}

	// Check that expired session was removed
	_, exists := sm.sessions[expiredSession.ID]
	if exists {
		t.Error("Expected expired session to be removed")
	}
}

func TestRefreshSession(t *testing.T) {
	sm := NewSessionManager(30 * time.Minute)
	userID := "test-user-123"

	// Create a session
	session, err := sm.CreateSession(userID)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	originalExpiry := session.ExpiresAt

	// Wait a bit to ensure time difference
	time.Sleep(10 * time.Millisecond)

	// Refresh the session
	err = sm.RefreshSession(session.ID)
	if err != nil {
		t.Fatalf("Failed to refresh session: %v", err)
	}

	// Check that expiry time was updated
	if !session.ExpiresAt.After(originalExpiry) {
		t.Error("Expected session expiry to be extended")
	}

	// Test refreshing non-existent session
	err = sm.RefreshSession("non-existent-session")
	if err != ErrInvalidCredentials {
		t.Errorf("Expected ErrInvalidCredentials, got %v", err)
	}
}

func TestDestroySession(t *testing.T) {
	sm := NewSessionManager(30 * time.Minute)
	userID := "test-user-123"

	// Create a session
	session, err := sm.CreateSession(userID)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// Verify session exists
	_, err = sm.ValidateSession(session.ID)
	if err != nil {
		t.Fatalf("Session should exist: %v", err)
	}

	// Destroy the session
	sm.DestroySession(session.ID)

	// Verify session no longer exists
	_, err = sm.ValidateSession(session.ID)
	if err != ErrInvalidCredentials {
		t.Errorf("Expected ErrInvalidCredentials after destroying session, got %v", err)
	}

	// Test destroying non-existent session (should not panic)
	sm.DestroySession("non-existent-session")
}

func TestCleanupExpiredSessions(t *testing.T) {
	sm := NewSessionManager(30 * time.Minute)

	// Create some sessions
	validSession, _ := sm.CreateSession("valid-user")
	expiredSession1, _ := sm.CreateSession("expired-user-1")
	expiredSession2, _ := sm.CreateSession("expired-user-2")

	// Set some sessions to expired
	expiredSession1.ExpiresAt = time.Now().Add(-1 * time.Hour)
	expiredSession2.ExpiresAt = time.Now().Add(-2 * time.Hour)

	// Verify initial state
	if len(sm.sessions) != 3 {
		t.Errorf("Expected 3 sessions, got %d", len(sm.sessions))
	}

	// Run cleanup
	sm.CleanupExpiredSessions()

	// Verify only valid session remains
	if len(sm.sessions) != 1 {
		t.Errorf("Expected 1 session after cleanup, got %d", len(sm.sessions))
	}

	_, exists := sm.sessions[validSession.ID]
	if !exists {
		t.Error("Expected valid session to remain after cleanup")
	}

	_, exists = sm.sessions[expiredSession1.ID]
	if exists {
		t.Error("Expected expired session 1 to be removed")
	}

	_, exists = sm.sessions[expiredSession2.ID]
	if exists {
		t.Error("Expected expired session 2 to be removed")
	}
}

func TestGenerateSessionID(t *testing.T) {
	id1, err := generateSessionID()
	if err != nil {
		t.Fatalf("Failed to generate session ID: %v", err)
	}

	if id1 == "" {
		t.Error("Expected non-empty session ID")
	}

	id2, err := generateSessionID()
	if err != nil {
		t.Fatalf("Failed to generate second session ID: %v", err)
	}

	// IDs should be different
	if id1 == id2 {
		t.Error("Generated session IDs should be different")
	}

	// Check that ID contains only lowercase letters
	for _, char := range id1 {
		if char < 'a' || char > 'z' {
			t.Errorf("Session ID should contain only lowercase letters, found: %c", char)
		}
	}
}

func TestSecureCompare(t *testing.T) {
	testCases := []struct {
		name     string
		a        string
		b        string
		expected bool
	}{
		{"identical strings", "hello", "hello", true},
		{"different strings", "hello", "world", false},
		{"empty strings", "", "", true},
		{"one empty", "hello", "", false},
		{"different lengths", "hello", "hello world", false},
		{"case sensitive", "Hello", "hello", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := SecureCompare(tc.a, tc.b)
			if result != tc.expected {
				t.Errorf("Expected %v for comparing '%s' and '%s', got %v", tc.expected, tc.a, tc.b, result)
			}
		})
	}
}

func TestSessionTimeout(t *testing.T) {
	shortTimeout := 100 * time.Millisecond
	sm := NewSessionManager(shortTimeout)

	// Create a session
	session, err := sm.CreateSession("test-user")
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// Session should be valid immediately
	_, err = sm.ValidateSession(session.ID)
	if err != nil {
		t.Fatalf("Session should be valid immediately: %v", err)
	}

	// Wait for session to expire
	time.Sleep(shortTimeout + 50*time.Millisecond)

	// Session should now be expired
	_, err = sm.ValidateSession(session.ID)
	if err != ErrSessionExpired {
		t.Errorf("Expected ErrSessionExpired, got %v", err)
	}
}

func BenchmarkCreateSession(b *testing.B) {
	sm := NewSessionManager(30 * time.Minute)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := sm.CreateSession("benchmark-user")
		if err != nil {
			b.Fatalf("Failed to create session: %v", err)
		}
	}
}

func BenchmarkValidateSession(b *testing.B) {
	sm := NewSessionManager(30 * time.Minute)
	session, err := sm.CreateSession("benchmark-user")
	if err != nil {
		b.Fatalf("Failed to create session: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := sm.ValidateSession(session.ID)
		if err != nil {
			b.Fatalf("Failed to validate session: %v", err)
		}
	}
}
