package worker

import (
	"context"
	"hrsync-backend/internal/service"
	"log"
	"time"
)

type LeaveWorker interface {
	Start(ctx context.Context)
}

type leaveWorker struct {
	srv service.LeaveService
}

func NewLeaveWorker(srv service.LeaveService) LeaveWorker {
	return &leaveWorker{srv: srv}
}

func (w *leaveWorker) Start(ctx context.Context) {
	// Run once on startup
	log.Println("[LeaveWorker] Initial check for leave status updates...")
	if err := w.srv.AutoUpdateStatus(ctx); err != nil {
		log.Printf("[LeaveWorker] Error updating leave status: %v", err)
	}

	for {
		now := time.Now()
		// Calculate time until next midnight
		nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		durationUntilMidnight := nextMidnight.Sub(now)

		select {
		case <-ctx.Done():
			log.Println("[LeaveWorker] Stopping leave worker...")
			return
		case <-time.After(durationUntilMidnight):
			log.Println("[LeaveWorker] Running daily leave status update...")
			if err := w.srv.AutoUpdateStatus(ctx); err != nil {
				log.Printf("[LeaveWorker] Error updating leave status: %v", err)
			}
		}
	}
}
