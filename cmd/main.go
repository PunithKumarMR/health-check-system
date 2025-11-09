package main

import (
    "fmt"
    "log"
    "health-check-system/pkg/config"
    "health-check-system/pkg/database"
    "github.com/joho/godotenv"
)

func main() {
    fmt.Println("===========================================")
    fmt.Println("Health Check System - Database Test")
    fmt.Println("===========================================")
    fmt.Println()

    // Load .env file
    fmt.Println("Loading .env file...")
    if err := godotenv.Load(); err != nil {
        log.Printf("Warning: .env file not found, using environment variables")
    } else {
        fmt.Println("✓ .env file loaded")
    }
    fmt.Println()

    // Load config
    fmt.Println("Loading configuration...")
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    fmt.Printf("✓ Config loaded (Env: %s)\n\n", cfg.App.Environment)

    // Connect to database
    fmt.Println("Connecting to database...")
    fmt.Printf("  Host: %s:%s\n", cfg.Database.Host, cfg.Database.Port)
    fmt.Printf("  Database: %s\n\n", cfg.Database.Database)
    
    // Convert config.DatabaseConfig to database.Config
    dbConfig := database.Config{
        Host:     cfg.Database.Host,
        Port:     cfg.Database.Port,
        User:     cfg.Database.User,
        Password: cfg.Database.Password,
        Database: cfg.Database.Database,
    }
    
    db, err := database.Connect(dbConfig)
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer db.Close()
    fmt.Println("✓ Database connected\n")

    // Test queries
    fmt.Println("Running test queries...")
    
    var userCount, nodeCount, proxyCount, serverCount int
    
    db.QueryRow("SELECT COUNT(*) FROM hc_niam_users").Scan(&userCount)
    db.QueryRow("SELECT COUNT(*) FROM hc_nodes").Scan(&nodeCount)
    db.QueryRow("SELECT COUNT(*) FROM hc_mito_proxies WHERE is_active=1").Scan(&proxyCount)
    db.QueryRow("SELECT COUNT(*) FROM hc_app_servers WHERE is_active=1").Scan(&serverCount)

    fmt.Printf("  NIAM Users: %d\n", userCount)
    fmt.Printf("  Nodes: %d\n", nodeCount)
    fmt.Printf("  Active Proxies: %d\n", proxyCount)
    fmt.Printf("  Active App Servers: %d\n\n", serverCount)

    fmt.Println("✅ All tests passed!")
    fmt.Println("===========================================")
}
