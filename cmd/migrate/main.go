package main

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	acceptedArgs := map[string]string{
		"up":    "migrate up",
		"down":  "migrate down",
		"force": "force a specific version",
	}

	acceptedKeys := slices.Collect(maps.Keys(acceptedArgs))

	arg := os.Args[1]

	if !slices.Contains(acceptedKeys, arg) {
		fmt.Printf("Invalid argument: %s\n", arg)
		for _, key := range acceptedKeys {
			fmt.Printf("  %s: %s\n", key, acceptedArgs[key])
		}
		os.Exit(1)
	}

	m, err := migrate.New(
		"file://internal/infra/database/migration",
		loadDBURL(),
	)

	if err != nil {
		fmt.Printf("failed to create migration object", err)
		os.Exit(1)
	}

	switch arg {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			fmt.Printf("failed to run migration up: %v\n", err)
			os.Exit(1)
		}
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			fmt.Printf("failed to run migration down: %v\n", err)
			os.Exit(1)
		}
	case "force":
		version := os.Args[2]
		versionInt, err := strconv.Atoi(version)
		if err != nil {
			fmt.Printf("invalid version number: %s\n", version)
			os.Exit(1)
		}
		if err := m.Force(versionInt); err != nil {
			fmt.Printf("failed to force migration version %s: %v\n", version, err)
			os.Exit(1)
		}
	}

	fmt.Println("migrations sucessfully applied")
}

func loadDBURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
}
