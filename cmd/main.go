package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	handler "github.com/strconvitoa/martian-service/internal/adapters/handlers"
	"github.com/strconvitoa/martian-service/internal/adapters/repository"
	"github.com/strconvitoa/martian-service/internal/core/services"
)

func main() {

	ctx := context.Background()
	// 1. Get the connection string from Env
	connStr := "postgres://postgres:@PulsarMap1977@db.kwcgfcibvsxkvchztbnx.supabase.co:5432/postgres"
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// 2. Initialize the pool directly using the connection string
	// You do NOT need pgx.Connect here.
	dbPool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}
	defer dbPool.Close()

	// 3. Optional: Ping the database to ensure the connection is actually alive
	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("Could not ping database: %v", err)
	}

	// Initialize specific repositories
	orgRepo := repository.NewPostgresOrgRepository(dbPool)
	userRepo := repository.NewPostgresUserRepository(dbPool)
	authRepo := repository.NewPostgresAuthRepository(dbPool)
	intakeRepo := repository.NewPostgresIntakeRepository(dbPool)

	// Inject into services
	orgsvc := services.NewOrgService(orgRepo)
	usrsvc := services.NewUserService(userRepo)
	authsvc := services.NewAuthService(authRepo)
	intakesvc := services.NewIntakeService(intakeRepo)

	// Setup handlers
	usrhdl := handler.NewUserHandler(usrsvc, authsvc)
	enthdl := handler.NewOrgHandler(orgsvc, usrsvc)
	authhdl := handler.NewAuthHandler(authsvc, usrsvc)
	intakehhdl := handler.NewIntakeHandler(intakesvc)

	app := fiber.New()

	app.Post("/users", usrhdl.Create)
	app.Post("/orgs", enthdl.Create)
	app.Post("/intakes", intakehhdl.Create)
	app.Post("/auth/reset-password", authhdl.Reset)
	app.Post("/auth/change-password", authhdl.Change)

	log.Fatal(app.Listen(":3001"))
}
