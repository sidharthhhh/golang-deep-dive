package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/sidharthhhh/go-auth-service/internal/config"
)

func NewMySQLConnection(cfg *config.Config) *sql.DB {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.MysqlUser,
		cfg.MysqlPassword,
		cfg.MysqlHost,
		cfg.MysqlPort,
		cfg.MysqlDB,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Database connection error:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Database ping error:", err)
	}

	log.Println("MySQL connected successfully")

	return db
}
