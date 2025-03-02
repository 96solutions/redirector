package registry

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/lroman242/redirector/config"
	"github.com/lroman242/redirector/domain/interactor"
	"github.com/lroman242/redirector/domain/repository"
	"github.com/lroman242/redirector/domain/service"
	"github.com/lroman242/redirector/infrastructure/logger"
	"github.com/lroman242/redirector/infrastructure/server"
	serviceImpl "github.com/lroman242/redirector/infrastructure/service"
	"github.com/lroman242/redirector/infrastructure/storage"
	"github.com/lroman242/redirector/infrastructure/transport/http"
	"github.com/oschwald/geoip2-golang"
	"github.com/redis/go-redis/v9"
)

// Registry provides factory methods for creating the main application components.
// It handles dependency injection and ensures proper initialization of services.
type Registry interface {
	// NewService creates a new RedirectInteractor instance
	NewService() interactor.RedirectInteractor
	// NewServer creates and configures the HTTP server
	NewServer() *server.Server
	// NewIPAddressParser creates a service for parsing IP addresses
	NewIPAddressParser() service.IPAddressParserInterface
	// NewUserAgentParser creates a service for parsing User-Agent strings
	NewUserAgentParser() service.UserAgentParserInterface
	// NewTrackingLinksRepository creates a repository for managing tracking links
	NewTrackingLinksRepository() repository.TrackingLinksRepositoryInterface
	// NewRedisClient creates a new Redis client
	NewRedisClient() *redis.Client
	// NewDB initializes the database connection
	NewDB() *sql.DB
}

// registry implements Registry interface and manages application component initialization
type registry struct {
	conf *config.AppConfig
}

// NewRegistry function initialize new Registry instance.
func NewRegistry(conf *config.AppConfig) Registry {
	slog.Info("Initializing Registry ...", "config", conf)

	r := &registry{
		conf: conf,
	}

	r.NewLogger()

	return r
}

// NewService func creates redirect interactor (interactor.RedirectInteractor) implementation.
func (r *registry) NewService() interactor.RedirectInteractor {
	slog.Info("initializing RedirectInteractor....")
	clickHandlers := make([]interactor.ClickHandlerInterface, 0)

	//clickHandlers = append(clickHandlers, click_handlers.NewClickHandlerWithMetrics(/* ... create some click handler*/))

	redirectInteractor := interactor.NewRedirectInteractor(
		r.NewTrackingLinksRepository(),
		r.NewIPAddressParser(),
		r.NewUserAgentParser(),
		clickHandlers,
	)

	return serviceImpl.NewRedirectWithMetrics(redirectInteractor)
}

// NewServer func creates an instance of new Server (HTTP).
func (r *registry) NewServer() *server.Server {
	slog.Info("initializing Server....")
	return server.NewServer(r.conf.HTTPServerConf, http.NewHandler(r.NewService()))
}

// NewDB func creates mysql session.
func (r *registry) NewDB() *sql.DB {
	slog.Info("initializing sql connection ...", slog.String("DSN", r.conf.DBConf.DSN()))

	db, err := sql.Open("postgres", fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		r.conf.DBConf.User, r.conf.DBConf.Password, r.conf.DBConf.Database, r.conf.DBConf.Host, r.conf.DBConf.Port))
	if err != nil {
		slog.Error("Couldn't open connection to postgres database", logger.ErrAttr(err))
		panic(err)
	}

	if err = db.Ping(); err != nil {
		slog.Error("Couldn't ping postgres database", logger.ErrAttr(err))
		panic(err)
	}

	db.SetConnMaxLifetime(r.conf.DBConf.ConnectionMaxLifeDuration())
	db.SetMaxIdleConns(r.conf.DBConf.MaxIdleConnections)
	db.SetMaxOpenConns(r.conf.DBConf.MaxOpenConnections)

	// Add periodic health check
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		for range ticker.C {
			if err := db.Ping(); err != nil {
				slog.Error("Database health check failed", logger.ErrAttr(err))
				//TODO: Implement reconnection logic
			}
		}
	}()

	return db
}

// NewIPAddressParser creates service.IPAddressParserInterface implementation.
func (r *registry) NewIPAddressParser() service.IPAddressParserInterface {
	slog.Info("initializing geoip2 db", "path", r.conf.GeoIP2DBPath)
	db, err := geoip2.Open(r.conf.GeoIP2DBPath)

	if err != nil {
		panic(err)
	}

	return serviceImpl.NewGeoIP2(db)
}

// NewUserAgentParser creates service.UserAgentParserInterface implementation.
func (r *registry) NewUserAgentParser() service.UserAgentParserInterface {
	slog.Info("initializing user agent parser...")
	return serviceImpl.NewUserAgentParser()
}

// NewTrackingLinksRepository creates repository.TrackingLinksRepositoryInterface implementation.
func (r *registry) NewTrackingLinksRepository() repository.TrackingLinksRepositoryInterface {
	slog.Info("initializing tracking links repository...")
	return storage.NewSQLStorage(r.NewDB())
}

// NewLogger creates pointer to *slog.Logger instance (which might be set as default logger).
func (r *registry) NewLogger() *slog.Logger {
	slog.Info("initializing logger...")
	return logger.NewLogger(r.conf.LogConf)
}

// NewRedisClient creates a new Redis client
func (r *registry) NewRedisClient() *redis.Client {
	return redis.NewClient(
		// Options contains Redis client options based on the configuration
		&redis.Options{
			Addr:            r.conf.RedisConf.Addr(),
			Password:        r.conf.RedisConf.Password,
			DB:              r.conf.RedisConf.DB,
			MaxRetries:      r.conf.RedisConf.MaxRetries,
			MinRetryBackoff: time.Duration(r.conf.RedisConf.MinRetryBackoff) * time.Millisecond,
			MaxRetryBackoff: time.Duration(r.conf.RedisConf.MaxRetryBackoff) * time.Millisecond,
			DialTimeout:     time.Duration(r.conf.RedisConf.DialTimeout) * time.Second,
			ReadTimeout:     time.Duration(r.conf.RedisConf.ReadTimeout) * time.Second,
			WriteTimeout:    time.Duration(r.conf.RedisConf.WriteTimeout) * time.Second,
			PoolSize:        r.conf.RedisConf.PoolSize,
			PoolTimeout:     time.Duration(r.conf.RedisConf.PoolTimeout) * time.Second,
		})
}
