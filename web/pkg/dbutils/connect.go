package dbutils

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

type SQLPara struct {
	DriverName string `yaml:"driver_name"`
	UserName   string `yaml:"username"`
	Password   string `yaml:"password"`
	NetMode    string `yaml:"net_mode"`
	Addr       string `yaml:"addr"`
	DBName     string `yaml:"db_name"`
	Para       string `yaml:"para"`
}

// SQLDriverName sql driver name (only support mysql)
const SQLDriverName = "mysql"

// JSONVersion json version
const JSONVersion = "1.1.0"

// DbSQLTimeout Database sql execution timeout
const DbSQLTimeout = 50

// GrpcConnBackoffMaxDelay grpc conn reconnect max delay
const GrpcConnBackoffMaxDelay = 10

// GrpcClientTimeout is timeout for calling grpc interface
const GrpcClientTimeout = 5

const maxStatementTimeStr = "max_statement_time"
const maxExecutionTimeStr = "max_execution_time"

// DeadlockRetryDelay is the delay time(ms) before retry when occur deadlock
const DeadlockRetryDelay = 20

var (
	// SQLTimeoutVarName is the db variable to set sql execution timeout
	SQLTimeoutVarName = "max_statement_time"
	// DBType is the type of db
	DBType = "mariadb"
)

// ConnectDBbyCfg create use
func ConnectDBbyCfg(SQL *SQLPara) *sql.DB {
	SQLDataSourceName := SQL.UserName + ":" + SQL.Password + "@" + SQL.NetMode + "(" + SQL.Addr + ")/" +
		SQL.DBName + "?" + SQL.Para

	return ConnectDB(SQLDataSourceName)
}

// ConnectDB create sql connection
func ConnectDB(dataSourceName string) *sql.DB {
	retry := 0
	var db *sql.DB
	var err error

	for {
		if retry > 10 {
			log.Panic("max retry exceeded")
		}

		db, err = sql.Open(SQLDriverName, dataSourceName)
		if err != nil {
			log.Panic("open error: ", err)
		}
		log.Debug("after sql.Open")

		now := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		err = db.PingContext(ctx)
		if err != nil {
			if strings.Contains(err.Error(), "connection timed out") && time.Since(now) > time.Second*60 {
				log.Fatalf("%v", err)
			}
			db.Close() // nolint
			log.Warn("Retrying, ", err)
			time.Sleep(5 * time.Second)
			retry++
			continue
		}
		break
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(50)

	setSQLTimeoutTime(db)

	return db
}

// DisconnectDB close sql connection
func DisconnectDB(db *sql.DB) {
	db.Close() // nolint
}

func setSQLTimeoutTime(db *sql.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*DbSQLTimeout)
	defer cancel()

	verStr := ""
	err := db.QueryRowContext(ctx, "select version()").Scan(&verStr)
	if err != nil {
		log.Errorf("[setSQLTimeoutTime] select version() failed, err:%v", err)
		return
	}

	if strings.Contains(verStr, "MariaDB") {
		SQLTimeoutVarName = maxStatementTimeStr
		DBType = "mariadb"
	} else {
		DBType = "mysql"
		splits := strings.Split(verStr, ".")
		if len(splits) >= 2 {
			i, err := strconv.Atoi(splits[0])
			if err != nil {
				log.Errorf("[setSQLTimeoutTime] type (string -> int) cast failed, err:%v, verStr:%v", err, verStr)
				return
			}
			if i <= 4 {
				SQLTimeoutVarName = maxStatementTimeStr
			} else if i >= 6 {
				SQLTimeoutVarName = maxExecutionTimeStr
			} else {
				j, err := strconv.Atoi(splits[1])
				if err != nil {
					log.Errorf("[setSQLTimeoutTime] type (string -> int) cast failed, err:%v, verStr:%v", err, verStr)
					return
				}
				if j <= 6 {
					SQLTimeoutVarName = maxStatementTimeStr
				} else {
					SQLTimeoutVarName = maxExecutionTimeStr
				}
			}
		}
	}
}
