package main

import (
	// "context"
	"log"
	"net/http"

	// "time"

	"github.com/prajwal-huggi/backend_go/internal/auth"
	"github.com/prajwal-huggi/backend_go/internal/config"
	"github.com/prajwal-huggi/backend_go/internal/repository"
	"github.com/prajwal-huggi/backend_go/internal/service"

	db "github.com/prajwal-huggi/backend_go/internal/db"
	httptransport "github.com/prajwal-huggi/backend_go/internal/transport/http"

	_ "github.com/prajwal-huggi/backend_go/docs"
)

// @title Backend Go API
// @version 1.0
// @description Learning Backend Project
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	cfg := config.LoadConfig() // reads the environment file

	pool := db.NewPostgresPool(cfg) // initialize the database connection

	userRepo := repository.NewUserRepository(pool)

	// means if the operation> 5 sec then cancel it(in this case db query)
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	jwtService := auth.NewJWTService(cfg.JWTSecret)
	userService := service.NewUserService(userRepo, jwtService)

	userHandler := httptransport.NewUserHandler(userService)

	router := httptransport.NewRouter(userHandler, jwtService)

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
