package ormdb

import (
	"github.com/hyperledger/fabric/common/viperutil"
	"github.com/pkg/errors"
)

type ORMDBConfig struct {
	Username      string
	Password      string
	Host          string
	Port          int
	DBName        string
	Sqlite3Config *Sqlite3Config
}

// Sqlite3Config
type Sqlite3Config struct {
	Path string
}

//ORMDBInstance represents a ORMDB instance
type ORMDBInstance struct {
	Config *ORMDBConfig
}

func NewORMDBInstance() (*ORMDBInstance, error) {
	config := &ORMDBConfig{}
	err := viperutil.EnhancedExactUnmarshalKey("ledger.state.ormDBConfig", config)
	if err != nil {
		return nil, errors.WithMessage(err, "load ormdb config from yaml failed")
	}

	ormDBInstance := &ORMDBInstance{Config: config}
	return ormDBInstance, nil
}
