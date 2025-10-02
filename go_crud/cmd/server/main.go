package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prajwal-huggi/backend_go/internal/config"
	"github.com/prajwal-huggi/backend_go/internal/http/handlers/student"
	"github.com/prajwal-huggi/backend_go/internal/storage/sqlite"
)

func main(){
	// 1) load config
	cfg:= config.MustLoad()

	// 2) database setup
	storage, err:= sqlite.New(cfg)
	if err!= nil{
		log.Fatalf("Unable to create the database: %s", err.Error())
	}

	slog.Info("Storage Initialized", slog.String("env", cfg.Env))


	// 3) setup router
	router:= http.NewServeMux()

	//The below is our first endpoint
	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students", student.GetList(storage))

	// 4) setup server
	server:= http.Server{
		Addr: cfg.HTTPServer.Addr,
		Handler: router,
	}

	slog.Info("Server Started: ", slog.String("address", cfg.Addr) )

	done:= make(chan os.Signal, 1)

	signal.Notify(done , os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func(){
		err:= server.ListenAndServe()
		if err!= nil && err != http.ErrServerClosed{
			log.Fatal("Failed to start the server ")
		}
	}()

	<- done

	slog.Info("Shutting down the server")

	ctx, cancel:=context.WithTimeout(context.Background(), 5* time.Second)
	defer cancel()


	err= server.Shutdown(ctx)
	if err!= nil{
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server Shutdown Successfully ")

	// The above shutdown method is known as the graceful shutdown.
}
