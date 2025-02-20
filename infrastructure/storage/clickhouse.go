package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/lroman242/redirector/domain/entity"
)

// TODO:
const clickhouseInsertClickQuery = "INSERT INTO clicks (id) VALUES (?)"

type ClickhouseStorage struct {
	session *sql.DB
}

// TODO:
func NewClickHouseStorage(host, port, database, username, password string) *ClickhouseStorage {
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", host, port)},
		Auth: clickhouse.Auth{
			Database: database,
			Username: username,
			Password: password,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: 5 * time.Second,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
	})
	conn.SetMaxIdleConns(5)
	conn.SetMaxOpenConns(10)
	conn.SetConnMaxLifetime(time.Hour)

	if err := conn.PingContext(context.Background()); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Catch exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}

		panic(err)
	}

	return &ClickhouseStorage{session: conn}
}

// Save function performs storing of a Click record into clickhouse database.
func (c *ClickhouseStorage) Save(ctx context.Context, click *entity.Click) error {
	scope, err := c.session.Begin()
	if err != nil {
		return err
	}

	batch, err := scope.Prepare("INSERT INTO clicks")
	if err != nil {
		return err
	}

	for i := 0; i < 1000; i++ {
		_, err = batch.Exec(
			click.ID,
			//TODO:
		)
		if err != nil {
			return err
		}
	}
	if err = scope.Commit(); err != nil {
		return err
	}

	return nil
}
