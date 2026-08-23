package main

import (
	"context"
	"log"
	"os"
	"taskflow/internal/handler"
	"taskflow/internal/repository"
	"taskflow/internal/service"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using environment variables")
	}
	ctx := context.Background()
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatalf("parse database config: %v", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnIdleTime = 5 * time.Minute

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("create database pool: %v", err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("ping database: %v", err)
	}

	repo := repository.NewRepository(db)

	userService := service.NewUserService(
		repo.Users,
	)

	projectService := service.NewProjectService(
		repo.Projects,
		repo.ProjectMembers,
		repo.Users,
	)

	projectMemberService := service.NewProjectMemberService(
		repo.ProjectMembers,
		repo.Projects,
		repo.Users,
	)

	taskService := service.NewTaskService(
		repo.Tasks,
		repo.Projects,
		repo.Users,
	)

	commentService := service.NewCommentService(
		repo.Comments,
		repo.Tasks,
		repo.Users,
	)
	taskHistoryService := service.NewTaskHistoryService(
		repo.TaskHistory,
		repo.Tasks,
		repo.Projects,
	)
	router := handler.SetupRouter(
		userService,
		projectService,
		taskService,
		projectMemberService,
		commentService,
		taskHistoryService,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server started on :%s", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("start server: %v", err)
	}
}
