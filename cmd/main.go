package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	handler "github.com/strconvitoa/martian-service/internal/adapters/handlers"
	"github.com/strconvitoa/martian-service/internal/adapters/repository"
	"github.com/strconvitoa/martian-service/internal/core/services"
)

func main() {

	ctx := context.Background()
	// 1. Get the connection string from Env
	connStr := "postgresql://postgres.kwcgfcibvsxkvchztbnx:PostgresInBrooklyn@aws-1-us-west-2.pooler.supabase.com:6543/postgres"
	// if connStr == "" {
	// 	log.Fatal("DATABASE_URL environment variable is not set")
	// }

	// // 2. Initialize the pool directly using the connection string
	// // You do NOT need pgx.Connect here.
	// dbPool, err := pgxpool.New(ctx, connStr)
	// if err != nil {
	// 	log.Fatalf("Unable to create connection pool: %v", err)
	// }
	// defer dbPool.Close()

	// // 3. Optional: Ping the database to ensure the connection is actually alive
	// if err := dbPool.Ping(ctx); err != nil {
	// 	log.Fatalf("Could not ping database: %v", err)
	// }

	// 1. Parse the connection string into a config object instead of connecting immediately
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		// Handle error (e.g., log.Fatalf("Unable to parse config: %v", err))
	}

	// 2. FORCE pgx to use Simple Protocol (This stops it from making prepared statements)
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	// 3. Pass the config object into NewWithConfig
	dbPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		// Handle error (e.g., log.Fatalf("Unable to create connection pool: %v", err))
	}

	// Your dbPool is now perfectly configured for Supabase Port 6543!

	// Initialize specific repositories
	orgRepo := repository.NewPostgresOrgRepository(dbPool)
	userRepo := repository.NewPostgresUserRepository(dbPool)
	authRepo := repository.NewPostgresAuthRepository(dbPool)
	LeadRepo := repository.NewPostgresLeadRepository(dbPool)

	// Inject into services
	orgsvc := services.NewOrgService(orgRepo)
	usrsvc := services.NewUserService(userRepo)
	authsvc := services.NewAuthService(authRepo)
	Leadsvc := services.NewLeadService(LeadRepo)
	emailsvc := services.NewEmailService()
	// Setup handlers
	usrhdl := handler.NewUserHandler(usrsvc, authsvc, emailsvc, orgsvc)
	enthdl := handler.NewOrgHandler(orgsvc, usrsvc, authsvc)
	authhdl := handler.NewAuthHandler(authsvc, usrsvc, emailsvc)
	leadhhdl := handler.NewLeadHandler(Leadsvc)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		// Ensure the exact protocol (http vs https) and port match your frontend
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowCredentials: true,
	}))
	app.Static("/static", "../public")
	app.Post("/users/login", usrhdl.Login)
	app.Post("/users", usrhdl.Create)
	app.Get("/users", usrhdl.Get)
	app.Delete("/users", usrhdl.Remove)
	app.Post("/orgs", enthdl.Create)
	app.Post("/leads", leadhhdl.Create)
	app.Put("/leads/:id", leadhhdl.Update)
	app.Get("/leads", leadhhdl.Get)
	app.Post("/auth/reset-password", authhdl.Reset)
	app.Post("/auth/change-password", authhdl.Change)
	//log.Fatal(app.Listen(":8443"))
	log.Fatal(app.ListenTLS(":8443", "./cert.pem", "./key.pem"))
}
