package registry

import (
	"database/sql"
	"fmt"
	"github.com/lroman242/redirector/infrastructure/logger"
	"log/slog"

	"github.com/lroman242/redirector/config"
	"github.com/lroman242/redirector/domain/interactor"
	"github.com/lroman242/redirector/domain/repository"
	"github.com/lroman242/redirector/domain/service"
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
	return server.NewServer(r.conf.HttpServerConf, http.NewHandler(r.NewService()))
}

// NewDB func creates mysql session.
func (r *registry) NewDB() *sql.DB {
	slog.Info("initializing mysql connection ...")

	db, err := sql.Open("mysql", r.conf.DBConf.DSN())
	if err != nil {
		panic(fmt.Sprintf("cannot connect to mysql %s", err))
	}

	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("unsuccessfull ping database. error: %w", err))
	}

	db.SetConnMaxLifetime(r.conf.DBConf.ConnectionMaxLifeDuration())
	db.SetMaxIdleConns(r.conf.DBConf.MaxIdleConnections)
	db.SetMaxOpenConns(r.conf.DBConf.MaxOpenConnections)

	return db
}

// NewIPAddressParser creates service.IpAddressParserInterface implementation.
func (r *registry) NewIPAddressParser() service.IpAddressParserInterface {
	slog.Info(" geoip2 db ", "path", r.conf.GeoIP2DBPath)
	db, err := geoip2.Open(r.conf.GeoIP2DBPath)

	if err != nil {
		panic(err)
	}

	return serviceImpl.NewGeoIP2(db)
}

// NewUserAgentParser creates service.UserAgentParserInterface implementation.
func (r *registry) NewUserAgentParser() service.UserAgentParserInterface {
	return serviceImpl.NewUserAgentParser()
}

// NewTrackingLinksRepository creates repository.TrackingLinksRepositoryInterface implementation.
func (r *registry) NewTrackingLinksRepository() repository.TrackingLinksRepositoryInterface {
	return storage.NewMySQLStorage(r.NewDB())
}

// NewLogger creates pointer to *slog.Logger instance (which might be set as default logger).
func (r *registry) NewLogger() *slog.Logger {
	return logger.NewLogger(r.conf.LogConf)
}
