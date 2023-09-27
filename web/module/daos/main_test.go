package daos

import (
	"log"
	"os"
	"testing"

	"github.com/learn-go/web/pkg/dbutils"
	"github.com/learn-go/web/pkg/migration"
)

var sqlPara = dbutils.SQLPara{
	UserName: "root",
	Password: "root",
	Addr:     "127.0.0.1:3306",
	DBName:   "demo",
	Para:     "charset=utf8mb4&parseTime=true",
}

func TestMain(m *testing.M) {
	err := migration.MigrationMysql(sqlPara, "../../migrations")
	if err != nil {
		log.Fatal(err)
	}

	db := dbutils.ConnectDBbyCfg(&sqlPara)
	lockerDao, _ = NewDaoLocker(db)

	code := m.Run()
	os.Exit(code)
}
