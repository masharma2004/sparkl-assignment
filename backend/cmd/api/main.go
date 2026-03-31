package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"sparklassignment/backend/internal/config"
	"sparklassignment/backend/internal/db"
	"sparklassignment/backend/internal/middleware"
	"sparklassignment/backend/internal/routes"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	database, err := db.NewPostgresConn(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	if err := db.RunMigrations(database); err != nil {
		log.Fatalf("migration error: %v", err)
	}

	if err := db.Seed(database); err != nil {
		log.Fatalf("error seeding database: %v", err)
	}

	db.StartMaintenanceLoop(context.Background(), database, cfg, log.Default())

	router := gin.Default()
	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatalf("Error configuring trusted proxies: %v", err)
	}
	router.Use(middleware.CORSMiddleware(cfg))

	routes.RegisterRoutes(router, database, cfg)

	log.Printf("Listening on port %s", cfg.Port)
	if err := router.Run(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
