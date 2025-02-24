package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/lroman242/redirector/domain/entity"
)

const clickhouseInsertClickQuery = `
	INSERT INTO clicks (
		id, target_url, referer, trk_url, slug, parent_slug,
		source_id, campaign_id, affiliate_id, advertiser_id, is_parallel,
		landing_id, gclid,
		user_agent, agent, platform, browser, device,
		ip, country_code,
		p1, p2, p3, p4,
		created_at
	) VALUES (
		?, ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?,
		?, ?,
		?, ?, ?, ?, ?,
		?, ?,
		?, ?, ?, ?,
		?
	)
`

// ClickhouseStorage implements ClicksRepository interface using Clickhouse as the underlying storage.
// It provides methods to store and manage click tracking data.
type ClickhouseStorage struct {
	session *sql.DB
}

// NewClickHouseStorage creates and initializes a new ClickhouseStorage instance.
// It establishes a connection to the Clickhouse database using the provided credentials
// and returns a configured storage instance ready for use.
// Panics if unable to establish a connection to the database.
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

// Save stores a Click record in the Clickhouse database.
// It handles all click-related fields including ID, LinkID, IPAddress, UserAgent,
// Referer, and timestamp fields.
// Returns an error if the database operation fails.
func (c *ClickhouseStorage) Save(ctx context.Context, click *entity.Click) error {
	scope, err := c.session.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer scope.Rollback()

	batch, err := scope.Prepare(clickhouseInsertClickQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer batch.Close()

	_, err = batch.Exec(
		click.ID,
		click.TargetURL,
		click.Referer,
		click.TrkURL,
		click.Slug,
		click.ParentSlug,
		click.SourceID,
		click.CampaignID,
		click.AffiliateID,
		click.AdvertiserID,
		click.IsParallel,
		click.LandingID,
		click.GCLID,
		click.UserAgent.SrcString,
		click.Agent,
		click.Platform,
		click.Browser,
		click.Device,
		click.IP.String(),
		click.CountryCode,
		click.P1,
		click.P2,
		click.P3,
		click.P4,
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	if err = scope.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
