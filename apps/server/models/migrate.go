package models

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(dbUrl string) {
	m, err := migrate.New("file://db/migrations", dbUrl)
	if err != nil {
		log.Fatal("error intializing migrations: ", err)
	}
	if err := m.Up(); err != nil {
		fmt.Println("running migrations: ", err)
	}
}
