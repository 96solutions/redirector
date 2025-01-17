package registry

import (
	"database/sql"
	"fmt"
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
}

type registry struct {
	conf               *config.AppConfig
	redirectInteractor interactor.RedirectInteractor
	mysql              *sql.DB
}

// NewRegistry function initialize new Registry instance.
func NewRegistry(conf *config.AppConfig) Registry {
	r := &registry{
		conf: conf,
	}

	r.mysql = r.NewDB()

	return r
}

// NewService func creates redirect interactor (interactor.RedirectInteractor) implementation.
func (r *registry) NewService() interactor.RedirectInteractor {
	clickHandlers := make([]interactor.ClickHandlerInterface, 0)

	r.redirectInteractor = interactor.NewRedirectInteractor(
		r.NewTrackingLinksRepository(),
		r.NewIPAddressParser(),
		r.NewUserAgentParser(),
		clickHandlers,
	)

	return r.redirectInteractor
}

func (r *registry) NewServer() *server.Server {
	return server.NewServer(r.conf.HttpServerConf, http.NewHandler(r.NewService()))
}

// NewDB func creates mysql session.
func (r *registry) NewDB() *sql.DB {
	if r.mysql != nil {
		return r.mysql
	}

	slog.Info("initializing mysql connection ...")

	mysqlSession, err := sql.Open("mysql", r.conf.DBConf.DSN())
	if err != nil {
		panic(fmt.Sprintf("cannot connect to mysql %s", err))
	}

	mysqlSession.SetConnMaxLifetime(r.conf.DBConf.ConnectionMaxLifeDuration())
	mysqlSession.SetMaxIdleConns(r.conf.DBConf.MaxIdleConnections)
	mysqlSession.SetMaxOpenConns(r.conf.DBConf.MaxOpenConnections)

	r.mysql = mysqlSession

	return r.mysql
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
