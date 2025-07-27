package main

import (
	"database/sql"
	"exp_tracker/config"
	"exp_tracker/database"
	"exp_tracker/routers"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load("config/.env")
	if err != nil {
		panic(err)
	}

	psqlInfo := fmt.Sprintf(
		`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`,
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGDATABASE"),
	)

	config.Db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer config.Db.Close()

	err = config.Db.Ping()
	if err != nil {
		panic(err)
	}

	database.DBMigrate(config.Db)

	routers.StartServer(config.Db)
}
