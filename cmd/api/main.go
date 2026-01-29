package main

import (
	"backend/cmd/models"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "backend/docs" // Swagger docs - se genera automáticamente en Dockerfile
)

const version = "1.0.0"

// @title           Accounting Users API
// @version         1.0.0
// @description     API REST para gestión de usuarios de Accounting A-Z
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@accounting-a-z.ch

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8000
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	jwt struct {
		secret string
	}
}

type AppStatus struct {
	Status     string `json: "status"`
	Enviroment string `json: "environment"`
	Version    string `json: "version"`
}
type application struct {
	config config
	logger *log.Logger
	models models.Models
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		// .env opcional: si no existe, se usan las variables de entorno del sistema
		godotenv.Load(".env.dev")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	dsnDefault := os.Getenv("DATABASE_DSN")
	if dsnDefault == "" {
		dsnDefault = os.Getenv("DSN")
	}
	if dsnDefault == "" {
		dsnDefault = "preview_usr:TU8uPynebAqadady@tcp(server29.hostfactory.ch:3306)/preview?parseTime=true"
	}

	var cfg config
	flag.StringVar(&cfg.env, "env", "Development", "Application enviroment (Development|Production)")
	flag.StringVar(&cfg.db.dsn, "dsn", dsnDefault, "MySQL connection string")
	flag.StringVar(&cfg.jwt.secret, "jwt-secret", "2dce505d96a53c5768052ee90f3df2055657518dad489160df9913f66042e160", "secret")
	flag.Parsed()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	app := &application{
		config: cfg,
		logger: logger,
		models: models.NewModels(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	logger.Println("Starting server on port ", port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
