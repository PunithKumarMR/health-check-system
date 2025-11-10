package main

import (
    "fmt"
    "log"
    
    "health-check-system/pkg/config"
    "health-check-system/pkg/database"
    "health-check-system/pkg/userpool"
    "health-check-system/pkg/proxy"
    "health-check-system/pkg/inventory"
    "health-check-system/pkg/status"
    
    "github.com/joho/godotenv"
)

func main() {
    fmt.Println("===========================================")
    fmt.Println("Health Check System - Module Tests")
    fmt.Println("===========================================")
    fmt.Println()

    // Load .env
    godotenv.Load()

    // Load config
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Connect to database
    dbConfig := database.Config{
        Host:     cfg.Database.Host,
        Port:     cfg.Database.Port,
        User:     cfg.Database.User,
        Password: cfg.Database.Password,
        Database: cfg.Database.Database,
    }
    
    db, err := database.Connect(dbConfig)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    fmt.Println("✓ Database connected\n")

    // Test User Pool
    fmt.Println("=== Testing User Pool ===")
    userPool := userpool.NewPool(db.DB)
    
    poolStatus, err := userPool.GetPoolStatus()
    if err != nil {
        log.Fatalf("Failed to get pool status: %v", err)
    }
    
    fmt.Printf("Total Users: %v\n", poolStatus["total_users"])
    fmt.Printf("Active Users: %v\n", poolStatus["active_users"])
    fmt.Printf("Total Capacity: %v\n", poolStatus["total_capacity"])
    fmt.Printf("Used Capacity: %v\n", poolStatus["used_capacity"])
    fmt.Printf("Available Capacity: %v\n\n", poolStatus["available_capacity"])

    // Test Proxy Manager
    fmt.Println("=== Testing Proxy Manager ===")
    proxyMgr := proxy.NewManager(db.DB)
    
    px, err := proxyMgr.GetProxy()
    if err != nil {
        log.Fatalf("Failed to get proxy: %v", err)
    }
    fmt.Printf("Primary Proxy: %s (%s:%d)\n", px.Name, px.IP, px.Port)
    
    allProxies, _ := proxyMgr.GetAllProxies()
    fmt.Printf("Total Proxies Available: %d\n\n", len(allProxies))

    // Test Inventory Manager
    fmt.Println("=== Testing Inventory Manager ===")
    invMgr := inventory.NewManager(db.DB)
    
    nodes, err := invMgr.GetNodesToCheck(5)
    if err != nil {
        log.Printf("Warning: %v\n", err)
    } else {
        fmt.Printf("Nodes Ready for Check: %d\n", len(nodes))
        for i, node := range nodes {
            fmt.Printf("  %d. %s (%s) - %s\n", i+1, node.Hostname, node.IPAddress, node.Circle)
        }
    }
    fmt.Println()

    // Test Status Manager
    fmt.Println("=== Testing Status Manager ===")
    statusMgr := status.NewManager(db.DB)
    
    activeChecks, err := statusMgr.GetActiveChecks()
    if err != nil {
        log.Fatalf("Failed to get active checks: %v", err)
    }
    fmt.Printf("Currently Active Checks: %d\n\n", activeChecks)

    fmt.Println("✅ All module tests passed!")
    fmt.Println("===========================================")
}
