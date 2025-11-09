-- ============================================
-- Health Check System - Database Setup
-- ============================================

USE mito_inventory;

SET FOREIGN_KEY_CHECKS=0;

-- ============================================
-- TABLE 1: hc_niam_users
-- NIAM user pool with session tracking
-- ============================================
DROP TABLE IF EXISTS hc_niam_users;
CREATE TABLE hc_niam_users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user VARCHAR(245) UNIQUE NOT NULL,
    passwd VARCHAR(245) NOT NULL,
    niam_ip VARCHAR(245) NOT NULL,
    niam_port VARCHAR(45) NOT NULL,
    login_status VARCHAR(45) DEFAULT 'Yes',
    current_sessions INT DEFAULT 0,
    max_sessions INT DEFAULT 5,
    active_session_ids JSON,
    is_expired BOOLEAN DEFAULT FALSE,
    expiry_date DATE,
    total_usage_count INT DEFAULT 0,
    last_used_at DATETIME,
    failure_count INT DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user (user),
    INDEX idx_login_status (login_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================
-- TABLE 2: hc_nodes
-- Node inventory with health check config
-- ============================================
DROP TABLE IF EXISTS hc_nodes;
CREATE TABLE hc_nodes (
    id INT PRIMARY KEY AUTO_INCREMENT,
    neId VARCHAR(245) UNIQUE NOT NULL,
    IPAddress VARCHAR(245) NOT NULL,
    Hostname VARCHAR(245) NOT NULL,
    Site VARCHAR(245),
    Circle VARCHAR(245),
    Login_status VARCHAR(50) DEFAULT 'Yes',
    vendor VARCHAR(50),
    node_type VARCHAR(50) DEFAULT 'router',
    environment VARCHAR(50) DEFAULT 'production',
    priority VARCHAR(20) DEFAULT 'medium',
    health_check_enabled BOOLEAN DEFAULT TRUE,
    custom_commands JSON,
    tags JSON,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_neId (neId),
    INDEX idx_login_status (Login_status),
    INDEX idx_circle (Circle)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================
-- TABLE 3: hc_node_status
-- Current status of each node
-- ============================================
DROP TABLE IF EXISTS hc_node_status;
CREATE TABLE hc_node_status (
    neId VARCHAR(245) PRIMARY KEY,
    current_status ENUM('idle','queued','connecting','running','polling','collecting','completed','failed','timeout','cancelled') DEFAULT 'idle',
    current_session_id VARCHAR(100),
    current_username VARCHAR(50),
    last_check_started DATETIME,
    last_check_completed DATETIME,
    last_check_duration INT,
    last_check_result VARCHAR(50),
    health_score INT,
    error_message TEXT,
    retry_count INT DEFAULT 0,
    consecutive_failures INT DEFAULT 0,
    last_successful_check DATETIME,
    total_checks INT DEFAULT 0,
    successful_checks INT DEFAULT 0,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_status (current_status),
    FOREIGN KEY (neId) REFERENCES hc_nodes(neId) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================
-- TABLE 4: hc_history
-- Historical health check records
-- ============================================
DROP TABLE IF EXISTS hc_history;
CREATE TABLE hc_history (
    id INT PRIMARY KEY AUTO_INCREMENT,
    session_id VARCHAR(100) UNIQUE NOT NULL,
    neId VARCHAR(245) NOT NULL,
    node_ip VARCHAR(245),
    hostname VARCHAR(245),
    circle VARCHAR(245),
    username VARCHAR(50),
    mito_proxy_used VARCHAR(100),
    app_server_used VARCHAR(100),
    started_at DATETIME NOT NULL,
    completed_at DATETIME,
    duration INT,
    final_status VARCHAR(50),
    result VARCHAR(50),
    health_score INT,
    metrics JSON,
    error_message TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_neId (neId),
    INDEX idx_started (started_at),
    FOREIGN KEY (neId) REFERENCES hc_nodes(neId) ON DELETE CASCADE,
    FOREIGN KEY (username) REFERENCES hc_niam_users(user) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================
-- TABLE 5: hc_live_updates
-- Real-time progress updates
-- ============================================
DROP TABLE IF EXISTS hc_live_updates;
CREATE TABLE hc_live_updates (
    id INT PRIMARY KEY AUTO_INCREMENT,
    session_id VARCHAR(100) NOT NULL,
    neId VARCHAR(245) NOT NULL,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50),
    message TEXT,
    progress_percentage INT,
    INDEX idx_session (session_id),
    FOREIGN KEY (session_id) REFERENCES hc_history(session_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================
-- TABLE 6: hc_active_sessions
-- Currently active SSH sessions
-- ============================================
DROP TABLE IF EXISTS hc_active_sessions;
CREATE TABLE hc_active_sessions (
    session_id VARCHAR(100) PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    neId VARCHAR(245) NOT NULL,
    node_ip VARCHAR(50),
    mito_proxy_used VARCHAR(100),
    started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_activity DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_username (username),
    FOREIGN KEY (username) REFERENCES hc_niam_users(user) ON DELETE CASCADE,
    FOREIGN KEY (neId) REFERENCES hc_nodes(neId) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================
-- TABLE 7: hc_mito_proxies
-- Mito proxy servers with failover
-- ============================================
DROP TABLE IF EXISTS hc_mito_proxies;
CREATE TABLE hc_mito_proxies (
    id INT PRIMARY KEY AUTO_INCREMENT,
    proxy_name VARCHAR(100) UNIQUE NOT NULL,
    proxy_ip VARCHAR(50) NOT NULL,
    proxy_port INT DEFAULT 22,
    proxy_user VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    is_primary BOOLEAN DEFAULT FALSE,
    priority INT DEFAULT 1,
    current_connections INT DEFAULT 0,
    max_connections INT DEFAULT 100,
    total_attempts INT DEFAULT 0,
    failed_attempts INT DEFAULT 0,
    success_rate DECIMAL(5,2) DEFAULT 100.00,
    last_success DATETIME,
    last_failure DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_active (is_active),
    INDEX idx_priority (priority)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================
-- TABLE 8: hc_app_servers
-- App servers with failover
-- ============================================
DROP TABLE IF EXISTS hc_app_servers;
CREATE TABLE hc_app_servers (
    id INT PRIMARY KEY AUTO_INCREMENT,
    server_name VARCHAR(100) UNIQUE NOT NULL,
    server_ip VARCHAR(50) NOT NULL,
    server_user VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    is_primary BOOLEAN DEFAULT FALSE,
    priority INT DEFAULT 1,
    current_load INT DEFAULT 0,
    max_load INT DEFAULT 50,
    total_requests INT DEFAULT 0,
    failed_requests INT DEFAULT 0,
    last_success DATETIME,
    last_failure DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS=1;

-- ============================================
-- Copy data from original tables
-- ============================================

INSERT INTO hc_niam_users (user, passwd, niam_ip, niam_port, login_status, expiry_date, max_sessions)
SELECT user, passwd, niam_ip, niam_port, login_status, 
       DATE_ADD(CURDATE(), INTERVAL 90 DAY), 5
FROM niam_users;

INSERT INTO hc_nodes (neId, IPAddress, Hostname, Site, Circle, Login_status)
SELECT neId, IPAddress, Hostname, Site, Circle, Login_status
FROM IBM_director_Info;

INSERT INTO hc_node_status (neId, current_status)
SELECT neId, 'idle' FROM hc_nodes WHERE Login_status = 'Yes';

-- ============================================
-- Insert Mito Proxy Servers
-- ============================================
INSERT INTO hc_mito_proxies (proxy_name, proxy_ip, proxy_port, proxy_user, is_primary, priority) VALUES
('mito-proxy-1', '150.236.16.69', 22, 'mitorunner', TRUE, 1),
('mito-proxy-2', '150.236.16.74', 22, 'mitorunner', FALSE, 2),
('mito-proxy-3', '150.236.16.92', 22, 'mitorunner', FALSE, 3),
('mito-proxy-4', '150.236.16.117', 22, 'mitorunner', FALSE, 4);

-- ============================================
-- Insert App Servers
-- ============================================
INSERT INTO hc_app_servers (server_name, server_ip, server_user, is_primary, priority) VALUES
('app-server-1', '103.170.144.33', 'mitorunner', TRUE, 1),
('app-server-2', '103.170.144.37', 'mitorunner', FALSE, 2),
('app-server-3', '103.170.144.39', 'mitorunner', FALSE, 3),
('app-server-4', '103.170.144.41', 'mitorunner', FALSE, 4);

-- ============================================
-- Verification Output
-- ============================================

SELECT ' Database setup complete!' as Status;
SELECT '' as '';
SELECT 'Table Counts:' as Info;
SELECT 'hc_niam_users' as TableName, COUNT(*) as Records FROM hc_niam_users
UNION ALL SELECT 'hc_nodes', COUNT(*) FROM hc_nodes
UNION ALL SELECT 'hc_node_status', COUNT(*) FROM hc_node_status
UNION ALL SELECT 'hc_history', COUNT(*) FROM hc_history
UNION ALL SELECT 'hc_live_updates', COUNT(*) FROM hc_live_updates
UNION ALL SELECT 'hc_active_sessions', COUNT(*) FROM hc_active_sessions
UNION ALL SELECT 'hc_mito_proxies', COUNT(*) FROM hc_mito_proxies
UNION ALL SELECT 'hc_app_servers', COUNT(*) FROM hc_app_servers;

SELECT '' as '';
SELECT 'Mito Proxies:' as Info;
SELECT proxy_name, proxy_ip, is_primary, priority FROM hc_mito_proxies ORDER BY priority;

SELECT '' as '';
SELECT 'App Servers:' as Info;
SELECT server_name, server_ip, is_primary, priority FROM hc_app_servers ORDER BY priority;
