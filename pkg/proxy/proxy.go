package proxy

import (
    "database/sql"
    "fmt"
    "sync"
)

// Proxy represents a Mito proxy server
type Proxy struct {
    Name     string
    IP       string
    Port     int
    User     string
    Priority int
    IsPrimary bool
}

// Manager manages Mito proxy pool
type Manager struct {
    db *sql.DB
    mu sync.Mutex
}

// NewManager creates a new proxy manager
func NewManager(db *sql.DB) *Manager {
    return &Manager{
        db: db,
    }
}

// GetProxy gets the best available proxy (failover support)
func (m *Manager) GetProxy() (*Proxy, error) {
    m.mu.Lock()
    defer m.mu.Unlock()

    query := `
        SELECT proxy_name, proxy_ip, proxy_port, proxy_user, priority, is_primary
        FROM hc_mito_proxies
        WHERE is_active = TRUE
        ORDER BY priority ASC
        LIMIT 1
    `

    proxy := &Proxy{}
    err := m.db.QueryRow(query).Scan(
        &proxy.Name,
        &proxy.IP,
        &proxy.Port,
        &proxy.User,
        &proxy.Priority,
        &proxy.IsPrimary,
    )
    if err != nil {
        return nil, fmt.Errorf("no available proxy: %w", err)
    }

    return proxy, nil
}

// RecordSuccess records successful proxy usage
func (m *Manager) RecordSuccess(proxyName string) error {
    _, err := m.db.Exec(`
        UPDATE hc_mito_proxies
        SET total_attempts = total_attempts + 1,
            last_success = NOW(),
            success_rate = ((total_attempts - failed_attempts) * 100.0) / (total_attempts + 1)
        WHERE proxy_name = ?
    `, proxyName)
    return err
}

// RecordFailure records failed proxy usage
func (m *Manager) RecordFailure(proxyName string) error {
    _, err := m.db.Exec(`
        UPDATE hc_mito_proxies
        SET total_attempts = total_attempts + 1,
            failed_attempts = failed_attempts + 1,
            last_failure = NOW(),
            success_rate = ((total_attempts - failed_attempts - 1) * 100.0) / (total_attempts + 1)
        WHERE proxy_name = ?
    `, proxyName)
    return err
}

// GetAllProxies returns all active proxies in priority order
func (m *Manager) GetAllProxies() ([]*Proxy, error) {
    query := `
        SELECT proxy_name, proxy_ip, proxy_port, proxy_user, priority, is_primary
        FROM hc_mito_proxies
        WHERE is_active = TRUE
        ORDER BY priority ASC
    `

    rows, err := m.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var proxies []*Proxy
    for rows.Next() {
        proxy := &Proxy{}
        err := rows.Scan(
            &proxy.Name,
            &proxy.IP,
            &proxy.Port,
            &proxy.User,
            &proxy.Priority,
            &proxy.IsPrimary,
        )
        if err != nil {
            return nil, err
        }
        proxies = append(proxies, proxy)
    }

    return proxies, nil
}
