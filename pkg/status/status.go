package status

import (
    "database/sql"
)

// Status represents node status
type Status string

const (
    StatusIdle       Status = "idle"
    StatusQueued     Status = "queued"
    StatusConnecting Status = "connecting"
    StatusRunning    Status = "running"
    StatusPolling    Status = "polling"
    StatusCompleted  Status = "completed"
    StatusFailed     Status = "failed"
    StatusTimeout    Status = "timeout"
)

// Manager manages node status
type Manager struct {
    db *sql.DB
}

// NewManager creates a new status manager
func NewManager(db *sql.DB) *Manager {
    return &Manager{
        db: db,
    }
}

// UpdateStatus updates node status
func (m *Manager) UpdateStatus(neID string, status Status, sessionID, username string) error {
    _, err := m.db.Exec(`
        UPDATE hc_node_status
        SET current_status = ?,
            current_session_id = ?,
            current_username = ?,
            last_check_started = CASE WHEN ? = 'running' THEN NOW() ELSE last_check_started END,
            updated_at = NOW()
        WHERE neId = ?
    `, status, sessionID, username, status, neID)

    return err
}

// RecordCompletion records health check completion
func (m *Manager) RecordCompletion(neID, sessionID string, success bool, duration int, errorMsg string) error {
    status := StatusCompleted
    result := "success"
    if !success {
        status = StatusFailed
        result = "failed"
    }

    _, err := m.db.Exec(`
        UPDATE hc_node_status
        SET current_status = ?,
            last_check_completed = NOW(),
            last_check_duration = ?,
            last_check_result = ?,
            total_checks = total_checks + 1,
            successful_checks = successful_checks + CASE WHEN ? THEN 1 ELSE 0 END,
            consecutive_failures = CASE WHEN ? THEN 0 ELSE consecutive_failures + 1 END,
            last_successful_check = CASE WHEN ? THEN NOW() ELSE last_successful_check END,
            error_message = ?,
            current_session_id = NULL,
            current_username = NULL
        WHERE neId = ?
    `, status, duration, result, success, success, success, errorMsg, neID)

    return err
}

// GetNodeStatus returns current status of a node
func (m *Manager) GetNodeStatus(neID string) (Status, error) {
    var status string
    err := m.db.QueryRow(`
        SELECT current_status
        FROM hc_node_status
        WHERE neId = ?
    `, neID).Scan(&status)

    if err != nil {
        return "", err
    }

    return Status(status), nil
}

// GetActiveChecks returns count of currently running checks
func (m *Manager) GetActiveChecks() (int, error) {
    var count int
    err := m.db.QueryRow(`
        SELECT COUNT(*)
        FROM hc_node_status
        WHERE current_status IN ('queued', 'connecting', 'running', 'polling')
    `).Scan(&count)

    return count, err
}

// AddLiveUpdate adds a progress update
func (m *Manager) AddLiveUpdate(sessionID, neID, status, message string, progress int) error {
    _, err := m.db.Exec(`
        INSERT INTO hc_live_updates (session_id, neId, status, message, progress_percentage)
        VALUES (?, ?, ?, ?, ?)
    `, sessionID, neID, status, message, progress)

    return err
}
