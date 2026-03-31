package db

import (
	"context"
	"log"
	"time"

	"gorm.io/gorm"

	"sparklassignment/backend/internal/config"
	"sparklassignment/backend/internal/models"
)

func StartMaintenanceLoop(ctx context.Context, database *gorm.DB, cfg *config.Config, logger *log.Logger) {
	interval := time.Duration(cfg.MaintenanceIntervalMinutes) * time.Minute
	if interval <= 0 {
		interval = 30 * time.Minute
	}

	revokedRefreshRetention := time.Duration(cfg.RevokedRefreshRetentionHours) * time.Hour
	if revokedRefreshRetention <= 0 {
		revokedRefreshRetention = 24 * time.Hour
	}

	runCleanup := func() {
		if err := cleanupSecurityState(database, revokedRefreshRetention); err != nil {
			logger.Printf("security maintenance cleanup failed: %v", err)
		}
	}

	go func() {
		runCleanup()

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				runCleanup()
			}
		}
	}()
}

func cleanupSecurityState(database *gorm.DB, revokedRefreshRetention time.Duration) error {
	now := time.Now().UTC()
	revokedCutoff := now.Add(-revokedRefreshRetention)

	return database.Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Where("reset_at <= ?", now).
			Delete(&models.RateLimitEntry{}).Error; err != nil {
			return err
		}

		if err := tx.
			Where("expires_at <= ?", now).
			Delete(&models.RefreshSession{}).Error; err != nil {
			return err
		}

		if err := tx.
			Where("revoked_at IS NOT NULL AND revoked_at <= ?", revokedCutoff).
			Delete(&models.RefreshSession{}).Error; err != nil {
			return err
		}

		return nil
	})
}
