package inventory

import (
    "database/sql"
    "fmt"
)

// Node represents a network node
type Node struct {
    NeID       string
    IPAddress  string
    Hostname   string
    Site       string
    Circle     string
    Vendor     string
    NodeType   string
}

// Manager manages node inventory
type Manager struct {
    db *sql.DB
}

// NewManager creates a new inventory manager
func NewManager(db *sql.DB) *Manager {
    return &Manager{
        db: db,
    }
}

// GetNodesToCheck returns nodes that need health check
func (m *Manager) GetNodesToCheck(limit int) ([]*Node, error) {
    query := `
        SELECT 
            n.neId, 
            n.IPAddress, 
            n.Hostname, 
            n.Site, 
            n.Circle, 
            COALESCE(n.vendor, 'unknown') as vendor,
            COALESCE(n.node_type, 'router') as node_type
        FROM hc_nodes n
        JOIN hc_node_status s ON n.neId = s.neId
        WHERE n.Login_status = 'Yes'
          AND n.health_check_enabled = TRUE
          AND s.current_status = 'idle'
        ORDER BY 
            COALESCE(s.last_check_completed, '2000-01-01') ASC,
            n.priority DESC
        LIMIT ?
    `

    rows, err := m.db.Query(query, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var nodes []*Node
    for rows.Next() {
        node := &Node{}
        err := rows.Scan(
            &node.NeID,
            &node.IPAddress,
            &node.Hostname,
            &node.Site,
            &node.Circle,
            &node.Vendor,
            &node.NodeType,
        )
        if err != nil {
            return nil, err
        }
        nodes = append(nodes, node)
    }

    if len(nodes) == 0 {
        return nil, fmt.Errorf("no nodes available for checking")
    }

    return nodes, nil
}

// GetNodeByID returns a specific node
func (m *Manager) GetNodeByID(neID string) (*Node, error) {
    query := `
        SELECT 
            neId, 
            IPAddress, 
            Hostname, 
            Site, 
            Circle, 
            COALESCE(vendor, 'unknown') as vendor,
            COALESCE(node_type, 'router') as node_type
        FROM hc_nodes
        WHERE neId = ? AND Login_status = 'Yes'
    `

    node := &Node{}
    err := m.db.QueryRow(query, neID).Scan(
        &node.NeID,
        &node.IPAddress,
        &node.Hostname,
        &node.Site,
        &node.Circle,
        &node.Vendor,
        &node.NodeType,
    )
    if err != nil {
        return nil, fmt.Errorf("node not found: %w", err)
    }

    return node, nil
}

// GetNodesByCircle returns nodes in a specific circle
func (m *Manager) GetNodesByCircle(circle string, limit int) ([]*Node, error) {
    query := `
        SELECT 
            n.neId, 
            n.IPAddress, 
            n.Hostname, 
            n.Site, 
            n.Circle, 
            COALESCE(n.vendor, 'unknown') as vendor,
            COALESCE(n.node_type, 'router') as node_type
        FROM hc_nodes n
        JOIN hc_node_status s ON n.neId = s.neId
        WHERE n.Circle = ?
          AND n.Login_status = 'Yes'
          AND n.health_check_enabled = TRUE
          AND s.current_status = 'idle'
        LIMIT ?
    `

    rows, err := m.db.Query(query, circle, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var nodes []*Node
    for rows.Next() {
        node := &Node{}
        err := rows.Scan(
            &node.NeID,
            &node.IPAddress,
            &node.Hostname,
            &node.Site,
            &node.Circle,
            &node.Vendor,
            &node.NodeType,
        )
        if err != nil {
            return nil, err
        }
        nodes = append(nodes, node)
    }

    return nodes, nil
}
