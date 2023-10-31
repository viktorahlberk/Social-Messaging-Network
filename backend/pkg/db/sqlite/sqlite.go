package db

import (
	"database/sql"
	"log"
	"social-network/pkg/models"

	"fmt"

	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
)

// initialize db
func InitDB() *sql.DB {

	db, err := sql.Open("sqlite3", "./tempDB.db")
	if err != nil {
		log.Fatal(err)
	}

	err = Migrations(db)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// InitRepositories should be called in server.go
// creates connection to db for all rep
func InitRepositories(db *sql.DB) *models.Repositories {
	return &models.Repositories{
		UserRepo:    &UserRepository{DB: db},
		SessionRepo: &SessionRepository{DB: db},
		GroupRepo:   &GroupRepository{DB: db},
		PostRepo:    &PostRepository{DB: db},
		CommentRepo: &CommentRepository{DB: db},
		NotifRepo:   &NotifRepository{DB: db},
		EventRepo:   &EventRepository{DB: db},
		MsgRepo:     &MsgRepository{DB: db},
	}
}

func Migrations(db *sql.DB) error {
	migrations := &migrate.FileMigrationSource{
		Dir: "../backend/pkg/db/migration/sqlite",
	}

	// function arguments: DB connection pool, DB management system, 
	// migration location, migration direction
	// for single migration, use ExecMax
	n, err := migrate.Exec(db, "sqlite3", migrations, migrate.Up)
	if err != nil {
		return err
	}

	fmt.Printf("Applied %d migrations to database.db!\n", n)

	return nil
}