package ormdb

import (
	"github.com/hyperledger/fabric/common/viperutil"
	"github.com/pkg/errors"
)

// ORMDBConfig
type ORMDBConfig struct {
	Username      string
	Password      string
	Host          string
	Port          int
	Sqlite3Config *Sqlite3Config
}

// Sqlite3Config
type Sqlite3Config struct {
	Path string
}

//GetORMDBConfig converts ormdb yaml config to struct
func GetORMDBConfig() (*ORMDBConfig, error) {
	config := &ORMDBConfig{}
	err := viperutil.EnhancedExactUnmarshalKey("ledger.state.ormDBConfig", config)
	if err != nil {
		return nil, errors.WithMessage(err, "load ormdb config from yaml failed")
	}
	return config, nil
}
