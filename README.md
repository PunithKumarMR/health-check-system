# Health Check System

Automated parallel health check system for network nodes using Go.

## Overview

This system performs health checks on network nodes (routers, switches, servers) in parallel using:
- **User Pool Management**: 10 NIAM users × 5 sessions = 50 concurrent checks
- **Multi-hop SSH**: Mito Proxy → NIAM Proxy → Target Node
- **Failover Support**: 4 Mito proxies + 4 App servers with automatic failover
- **Real-time Status Tracking**: Monitor progress and health of all checks
- **Database-driven**: Dynamic inventory and configuration

## Quick Start

### 1. Database Setup
```bash
mysql -h 103.170.144.21 -u root -pmito mito_inventory < scripts/setup_database.sql
```

### 2. Configuration
```bash
cp .env.example .env
# Edit .env with actual passwords
nano .env
```

### 3. Build & Run
```bash
go mod tidy
go build -o health-check-system cmd/main.go
./health-check-system
```

## Architecture
```
Database → Inventory Manager → User Pool Manager
                             ↓
                    SSH Session Manager
                             ↓
                  Health Check Executor
                             ↓
                     Status Manager → Results
```

## Infrastructure

**Database:** 103.170.144.21 (mito_inventory)

**Mito Proxies (Failover):**
- Primary: 150.236.16.74
- Backup: 150.236.16.75, 150.236.16.76, 150.236.16.77

**App Servers (Failover):**
- Primary: 103.170.144.39
- Backup: 103.170.144.33, 103.170.144.41, 103.170.144.37

## Project Structure
```
├── cmd/main.go           # Application entry point
├── pkg/                  # Core modules
├── config/               # Configuration files
├── scripts/              # Database setup scripts
└── docs/                 # Documentation
```

## Development Status

- [x] Phase 1: Database schema
- [ ] Phase 2: Configuration files
- [ ] Phase 3: Go modules implementation
- [ ] Phase 4: Testing
- [ ] Phase 5: Deployment

## License

Internal use - Airtel
