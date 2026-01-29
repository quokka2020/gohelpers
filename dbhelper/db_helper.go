package dbhelper

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/quokka2020/gohelpers/util"
)

type Db_Helper struct {
	pool_config *pgxpool.Config
}

type Db_Session struct {
	session_context context.Context
	db_conn         *pgxpool.Pool
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

// You have to defer a close
func (helper *Db_Helper) CreateSession(ctx context.Context) (*Db_Session, error) {
	var err error
	session := Db_Session{
		session_context: ctx,
	}
	session.db_conn, err = pgxpool.ConnectConfig(ctx, helper.pool_config)
	if err != nil {
		log.Printf("Unable to create connection pool %v", err)
		return nil, fmt.Errorf("unable to create connection pool %v", err)
	}

	_, err = session.db_conn.Exec(ctx, "set timezone = 'UTC'")
	if err != nil {
		log.Printf("failed to set timezone to UTC err:%v", err)
		session.db_conn.Close()
		return nil, fmt.Errorf("failed to set timezone to UTC err:%v", err)
	}
	return &session, nil
}

func (helper *Db_Helper) SingleQuery(sql func(ctx context.Context, db *pgxpool.Pool) error) error {
	ctx := context.Background()
	session, err := helper.CreateSession(ctx)
	if err != nil {
		return err
	}
	defer session.Close()
	return session.Run(sql)
}

func (session *Db_Session) Run(sql func(ctx context.Context, db *pgxpool.Pool) error) error {
	return sql(session.session_context, session.db_conn)
}

func (session *Db_Session) Exec(query string, args ...any) (pgconn.CommandTag, error) {
	return session.db_conn.Exec(session.session_context, query, args)
}

func (session *Db_Session) Close() {
	session.db_conn.Close()
}
