package main

import (
	"database/sql"
	"fmt"
	"github.com/TeenBanner/Inventory_system/database"
	"os"
)

var envDbVarsRequired = []string{"DB_USER", "DB_PASSWORD", "DB_PORT", "DB_HOST", "DB_NAME"}

func connectDB() *sql.DB {
	database.CreateConnection(createStrConn())

	db := database.Pool()

	return db
}

func createStrConn() string {
	vars := getEnvDbVars()
	var (
		user     = ""
		password = ""
		host     = ""
		port     = ""
		dbName   = ""
	)

	required := []string{user, password, host, port, dbName}
	for i := 0; i < len(vars); i++ {
		for j := 0; j < len(required); j++ {
			required[j] = vars[i]
		}
	}

	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", user, password, host, port, dbName)
}

func getEnvDbVars() []string {
	values := []string{}
	for _, envVar := range envDbVarsRequired {
		value := os.Getenv(envVar)

		values = append(values, value)
	}

	return values
}
