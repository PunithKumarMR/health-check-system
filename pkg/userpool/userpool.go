package userpool

import (
    "database/sql"
    "fmt"
    "sync"
    "time"
)

// User represents a NIAM user
type User struct {
    Username       string
    Password       string
    NiamIP         string
    NiamPort       string
    CurrentSessions int
    MaxSessions    int
}

// Pool manages NIAM user pool
type Pool struct {
    db              *sql.DB
    mu              sync.Mutex
    maxWaitTime     time.Duration
    checkInterval   time.Duration
}

// NewPool creates a new user pool
func NewPool(db *sql.DB) *Pool {
    return &Pool{
        db:            db,
        maxWaitTime:   5 * time.Minute,
        checkInterval: 2 * time.Second,
    }
}

// AcquireUser gets an available user from the pool
func (p *Pool) AcquireUser(sessionID string) (*User, error) {
    timeout := time.After(p.maxWaitTime)
    ticker := time.NewTicker(p.checkInterval)
    defer ticker.Stop()

    for {
        select {
        case <-timeout:
            return nil, fmt.Errorf("timeout waiting for available user")
        case <-ticker.C:
            user, err := p.tryAcquireUser(sessionID)
            if err == nil {
                return user, nil
            }
        }
    }
}

// tryAcquireUser attempts to acquire a user
func (p *Pool) tryAcquireUser(sessionID string) (*User, error) {
    p.mu.Lock()
    defer p.mu.Unlock()

    tx, err := p.db.Begin()
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()

    // Find available user
    query := `
        SELECT user, passwd, niam_ip, niam_port, current_sessions, max_sessions
        FROM hc_niam_users
        WHERE login_status = 'Yes'
          AND is_expired = FALSE
          AND current_sessions < max_sessions
        ORDER BY current_sessions ASC, last_used_at ASC
        LIMIT 1
        FOR UPDATE
    `

    user := &User{}
    err = tx.QueryRow(query).Scan(
        &user.Username,
        &user.Password,
        &user.NiamIP,
        &user.NiamPort,
        &user.CurrentSessions,
        &user.MaxSessions,
    )
    if err != nil {
        return nil, fmt.Errorf("no available users: %w", err)
    }

    // Update user session count
    _, err = tx.Exec(`
        UPDATE hc_niam_users
        SET current_sessions = current_sessions + 1,
            active_session_ids = JSON_ARRAY_APPEND(
                COALESCE(active_session_ids, JSON_ARRAY()),
                '$',
                ?
            ),
            last_used_at = NOW(),
            total_usage_count = total_usage_count + 1
        WHERE user = ?
    `, sessionID, user.Username)
    if err != nil {
        return nil, err
    }

    if err := tx.Commit(); err != nil {
        return nil, err
    }

    return user, nil
}

// ReleaseUser releases a user back to the pool
func (p *Pool) ReleaseUser(username, sessionID string) error {
    p.mu.Lock()
    defer p.mu.Unlock()

    _, err := p.db.Exec(`
        UPDATE hc_niam_users
        SET current_sessions = GREATEST(current_sessions - 1, 0),
            active_session_ids = JSON_REMOVE(
                active_session_ids,
                JSON_UNQUOTE(JSON_SEARCH(active_session_ids, 'one', ?))
            )
        WHERE user = ?
    `, sessionID, username)

    return err
}

// GetPoolStatus returns current pool status
func (p *Pool) GetPoolStatus() (map[string]interface{}, error) {
    var totalUsers, activeUsers, totalCapacity, usedCapacity int

    err := p.db.QueryRow(`
        SELECT 
            COUNT(*) as total_users,
            SUM(CASE WHEN login_status='Yes' THEN 1 ELSE 0 END) as active_users,
            SUM(max_sessions) as total_capacity,
            SUM(current_sessions) as used_capacity
        FROM hc_niam_users
    `).Scan(&totalUsers, &activeUsers, &totalCapacity, &usedCapacity)

    if err != nil {
        return nil, err
    }

    return map[string]interface{}{
        "total_users":      totalUsers,
        "active_users":     activeUsers,
        "total_capacity":   totalCapacity,
        "used_capacity":    usedCapacity,
        "available_capacity": totalCapacity - usedCapacity,
    }, nil
}
