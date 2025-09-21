package dbhelper

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/quokka2020/gohelpers/util"
)

var dbhost = flag.String(
	"dbhost", util.GetEnv("DB_HOST", "192.168.10.4"), "the db-host")
var dbname = flag.String(
	"dbname", util.GetEnv("DB_NAME", ""), "the dbname")
var dbuser = flag.String(
	"dbuser", util.GetEnv("DB_USER", ""), "the dbuser")
var dbpassword = flag.String(
	"dbpassword", util.GetEnv("DB_PASS", ""), "the dbpassword")

type Db_Helper struct {
	pool_config *pgxpool.Config
}

func CreateDbHelper(prefix string) *Db_Helper {
	var err error
	helper := Db_Helper{}
	pgconfig := fmt.Sprintf("host=%s database=%s user=%s password=%s", *dbhost, *dbname, *dbuser, *dbpassword)
	helper.pool_config, err = pgxpool.ParseConfig(pgconfig)
	if err != nil {
		log.Panicf("Unable to parse DATABASE_URL:[%s] %v", pgconfig, err)
		// os.Exit(1)
	}
	return &helper
}

func (helper *Db_Helper) SingleQuery(sql func(db *pgxpool.Pool) error) (error) {
	ctx := context.Background()
	db, err := pgxpool.ConnectConfig(ctx, helper.pool_config)
	if err != nil {
		log.Printf("Unable to create connection pool %v", err)
		return fmt.Errorf("unable to create connection pool %v", err)
		// os.Exit(1)
	}
	defer db.Close()

	_, err = db.Exec(ctx, "set timezone = 'UTC'")
	if err != nil {
		log.Printf("failed to set timezone to UTC err:%v",err)
		return fmt.Errorf("failed to set timezone to UTC err:%v",err)
	}

	return sql(db)
}
