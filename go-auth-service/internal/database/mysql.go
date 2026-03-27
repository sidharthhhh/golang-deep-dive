package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/sidharthhhh/go-auth-service/internal/config"
)

func NewMySQLConnection(cfg *config.Config) (*sql.DB, error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		cfg.MysqlUser,
		cfg.MysqlPassword,
		cfg.MysqlHost,
		cfg.MysqlPort,
		cfg.MysqlDB,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("database ping error: %w", err)
	}

	log.Println("MySQL connected successfully")

	return db, nil
}
