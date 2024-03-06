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
	// "github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/heroku/x/hmetrics/onload"
)

const version = "1.0.0"

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
	// os.Setenv("PORT", "8080")
	// Find .env file
	// err := godotenv.Load(".env")
	// if err != nil{
	//  log.Fatalf("Error loading .env file: %s", err)
	// }
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	
	var cfg config
	// flag.IntVar(&cfg.port, "port", port, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "Development", "Application enviroment (Development|Production)")
	flag.StringVar(&cfg.db.dsn, "dsn", "preview_usr:TU8uPynebAqadady@tcp(server29.hostfactory.ch:3306)/preview?parseTime=true", "MySQL connection string")
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
