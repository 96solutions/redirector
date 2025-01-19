package registry

import (
	"database/sql"
	"fmt"
	"log/slog"

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
)

// Registry interface describe struct which should build main application parts.
type Registry interface {
	NewService() interactor.RedirectInteractor
	NewServer() *server.Server
	NewDB() *sql.DB
	NewIPAddressParser() service.IpAddressParserInterface
	NewUserAgentParser() service.UserAgentParserInterface
	NewTrackingLinksRepository() repository.TrackingLinksRepositoryInterface
}

type registry struct {
	conf *config.AppConfig
}

// NewRegistry function initialize new Registry instance.
func NewRegistry(conf *config.AppConfig) Registry {
	r := &registry{
		conf: conf,
	}

	return r
}

// NewService func creates redirect interactor (interactor.RedirectInteractor) implementation.
func (r *registry) NewService() interactor.RedirectInteractor {
	slog.Info("initializing RedirectInteractor....")
	clickHandlers := make([]interactor.ClickHandlerInterface, 0)

	return interactor.NewRedirectInteractor(
		r.NewTrackingLinksRepository(),
		r.NewIPAddressParser(),
		r.NewUserAgentParser(),
		clickHandlers,
	)
}

// NewServer func creates an instance of new Server (HTTP).
func (r *registry) NewServer() *server.Server {
	slog.Info("initializing Server....")
	return server.NewServer(r.conf.HttpServerConf, http.NewHandler(r.NewService()))
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

	return db
}

// NewIPAddressParser creates service.IpAddressParserInterface implementation.
func (r *registry) NewIPAddressParser() service.IpAddressParserInterface {
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
