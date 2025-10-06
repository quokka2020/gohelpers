package dbhelper

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/quokka2020/gohelpers/util"
)

type Db_Helper struct {
	pool_config *pgxpool.Config
}

func CreateDbHelper() *Db_Helper {
	var err error
	helper := Db_Helper{}
	pgconfig := fmt.Sprintf("host=%s database=%s user=%s password=%s", util.GetEnv("DB_HOST", "192.168.10.4"), util.GetEnv("DB_NAME", ""), util.GetEnv("DB_USER", ""), util.GetEnv("DB_PASS", ""))
	helper.pool_config, err = pgxpool.ParseConfig(pgconfig)
	if err != nil {
		log.Panicf("Unable to parse DATABASE_URL:[%s] %v", pgconfig, err)
		// os.Exit(1)
	}
	return &helper
}

func (helper *Db_Helper) SingleQuery(sql func(ctx context.Context, db *pgxpool.Pool) error) error {
	ctx := context.Background()
	db, err := pgxpool.ConnectConfig(ctx, helper.pool_config)
	if err != nil {
		log.Printf("Unable to create connection pool %v", err)
		return fmt.Errorf("unable to create connection pool %v", err)
	}
	defer db.Close()

	_, err = db.Exec(ctx, "set timezone = 'UTC'")
	if err != nil {
		log.Printf("failed to set timezone to UTC err:%v", err)
		return fmt.Errorf("failed to set timezone to UTC err:%v", err)
	}

	return sql(ctx, db)
}
