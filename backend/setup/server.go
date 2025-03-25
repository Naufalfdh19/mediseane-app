package setup

import (
	"context"
	"log"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

const Timeout = 600

func Run() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	db, err := ConnectDB()
	if err != nil {
		log.Fatalf("unable to connect to the database: %v", err)
	}
	defer db.Close()

	redisDB, err := ConnectRedisDB(context.Background())
	if err != nil {
		log.Fatalf("unable to connect to redis database: %v", err)
	}
	defer redisDB.Close()
	conn, rabbitMQ, err := ConnectRabbitMQ()
	if err != nil {
		log.Fatalf("unable to connect to rabbit mq: %v", err)
	}
	defer conn.Close()
	defer rabbitMQ.Close()

	resendClient := SetupResend()
	logger.SetLogger(logger.NewLogrusLogger())
	apperror.FormatValidatedField()

	router := SetUp(db, redisDB, resendClient, rabbitMQ)

	s := &http.Server{
		Addr:         ":" + os.Getenv("SERVER_PORT"),
		Handler:      router,
		ReadTimeout:  Timeout * time.Second,
		WriteTimeout: Timeout * time.Second,
	}

	isGracefulShutdown := true
	if !isGracefulShutdown {
		s.ListenAndServe()
	} else {
		go func() {
			if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Server error: %v", err)
			}
		}()
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		log.Println("Shutdown server...")

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			log.Fatalf("Server shutdown error: %v", err)
		}

		<-ctx.Done()

		log.Println("Server exited gracefully")
	}

}
