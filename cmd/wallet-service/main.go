package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wallet-service/pkg/db"

	"wallet-service/internal/config"
	"wallet-service/internal/server"
)

func main() {

	dsn := config.DSNFromEnv()

	conn, err := db.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer conn.Close(context.Background())

	srv := server.NewServer(conn)

	go func() {
		log.Printf("starting wallet service on %s", config.GetAppPort())
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Stop(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
}
