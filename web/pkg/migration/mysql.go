package migration

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/learn-go/web/pkg/dbutils"
)

func MigrationMysql(config dbutils.SQLPara, sourcePath string) error {
	return migrateMysqlUp(config, sourcePath)
}

func migrateMysqlUp(config dbutils.SQLPara, sourcePath string) error {
	var db *sql.DB

	uri := fmt.Sprintf("%s:%s@tcp(%s)/", config.UserName, config.Password, config.Addr)

	db = dbutils.ConnectDB(uri)

	defer db.Close() // nolint
	_, err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARSET utf8mb4 COLLATE=utf8mb4_unicode_ci;", config.DBName))
	if err != nil {
		return fmt.Errorf("failed to create db: %v", err)
	}

	m, err := migrate.New("file://"+sourcePath, "mysql://"+uri+config.DBName)
	if err != nil {
		return fmt.Errorf("migrate error: %v", err)
	}
	defer m.Close() // nolint
	err = m.Up()

	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %v", err)
	}

	return nil
}
