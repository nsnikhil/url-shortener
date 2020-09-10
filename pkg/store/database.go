package store

import (
	"database/sql"
	_ "github.com/newrelic/go-agent/v3/integrations/nrpq"
	"go.uber.org/zap"
	"time"
	"urlshortner/pkg/config"
)

type DBHandler interface {
	GetDB() (*sql.DB, error)
}

type defaultDBHandler struct {
	cfg config.DatabaseConfig
	lgr *zap.Logger
}

func (dbh *defaultDBHandler) GetDB() (*sql.DB, error) {
	db, err := sql.Open(dbh.cfg.GetDriverName(), dbh.cfg.GetSource())
	if err != nil {
		dbh.lgr.Error(err.Error())
		return nil, err
	}

	db.SetMaxOpenConns(dbh.cfg.GetMaxOpenConnections())
	db.SetMaxIdleConns(dbh.cfg.GetIdleConnections())
	db.SetConnMaxLifetime(time.Minute * time.Duration(dbh.cfg.GetConnectionMaxLifetime()))

	if err := db.Ping(); err != nil {
		dbh.lgr.Error(err.Error())
		return nil, err
	}

	return db, nil
}

func NewDBHandler(cfg config.DatabaseConfig, lgr *zap.Logger) DBHandler {
	return &defaultDBHandler{
		cfg: cfg,
		lgr: lgr,
	}
}
