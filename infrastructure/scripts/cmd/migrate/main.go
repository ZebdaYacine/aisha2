package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	path := flag.String("path", "", "migration directory")
	databaseURL := flag.String("database", "", "PostgreSQL connection URL")
	flag.Parse()
	if *path == "" || *databaseURL == "" || flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: migrate -path DIR -database URL up|down")
		os.Exit(2)
	}

	runner, err := migrate.New("file://"+*path, *databaseURL)
	if err != nil {
		fail(err)
	}
	defer runner.Close()

	switch flag.Arg(0) {
	case "up":
		err = runner.Up()
	case "down":
		err = runner.Down()
	default:
		fmt.Fprintln(os.Stderr, "direction must be up or down")
		os.Exit(2)
	}
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		fail(err)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
