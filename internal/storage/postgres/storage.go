package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"jedyEvgeny/online-music-library/internal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type DataBase struct {
	db *sql.DB
}

const dirMigrations = "migrations"

func MustNew(cfg *config.Config) *DataBase {
	db, err := runDB(cfg)
	if err != nil {
		log.Fatalf(errCreateDB, err)
	}
	return &DataBase{
		db: db,
	}
}

func runDB(cfg *config.Config) (*sql.DB, error) {
	connName := createConnStr(cfg)
	db, err := sql.Open(cfg.Database.Type, connName)
	if err != nil {
		return nil, fmt.Errorf(errCantOpen, connName, err)
	}

	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)

	idleTime, err := time.ParseDuration(cfg.Database.ConnMaxIdleTime)
	if err != nil {
		return nil, fmt.Errorf(errParseNotActiveConn, err)
	}
	db.SetConnMaxIdleTime(idleTime)

	maxLifetime, err := time.ParseDuration(cfg.Database.ConnMaxLifetime)
	if err != nil {
		return nil, fmt.Errorf(errParseLifeConn, err)
	}
	db.SetConnMaxLifetime(maxLifetime)

	tStartPing := time.Now()
	err = db.Ping()
	tEndPing := time.Now()
	if err != nil {
		return nil, fmt.Errorf(errPing, connName, err)
	}
	log.Printf(msgTimePing, connName, tEndPing.Sub(tStartPing))

	err = runMigrate(cfg, db)
	if err != nil {
		return nil, fmt.Errorf(errLaunchMigrations, connName, err)
	}

	return db, nil
}

func (d *DataBase) Close() error {
	return d.db.Close()
}

func createConnStr(cfg *config.Config) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s connect_timeout=%d",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password,
		cfg.Database.Name, cfg.Database.SSLMode, cfg.Database.ConnectTimeout)

}

func runMigrate(cfg *config.Config, db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf(errDriver, err)
	}

	workingDirApp, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	migrationsPathDir := filepath.Join(workingDirApp, dirMigrations)
	if _, err = os.Stat(migrationsPathDir); os.IsNotExist(err) {
		return fmt.Errorf(isntExistMigrations, migrationsPathDir, err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		pathMigrations(),
		cfg.Database.Type, driver,
	)
	if err != nil {
		return fmt.Errorf(errInitMigrations, err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf(errExecMigrations, err)
	}

	if err == migrate.ErrNoChange {
		log.Println(msgMigrationsNotNeed)
	} else {
		log.Println(msgMigrationsDone)
	}

	return nil
}

func pathMigrations() string {
	return "file://" + dirMigrations
}
