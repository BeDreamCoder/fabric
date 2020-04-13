/*
Copyright Zhigui.com. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package config

import (
	"github.com/hyperledger/fabric/common/viperutil"
	"github.com/pkg/errors"
)

// ORMDBConfig
type ORMDBConfig struct {
	Username      string         `mapstructure:"username" json:"username" yaml:"username"`
	Password      string         `mapstructure:"password" json:"password" yaml:"password"`
	Host          string         `mapstructure:"host" json:"host" yaml:"host"`
	Port          int            `mapstructure:"port" json:"port" yaml:"port"`
	Sqlite3Config *Sqlite3Config `mapstructure:"sqlite3Config" json:"sqlite3Config" yaml:"sqlite3Config"`
}

// Sqlite3Config
type Sqlite3Config struct {
	Path string `mapstructure:"path" json:"path" yaml:"path"`
}

//GetORMDBConfig converts ormdb yaml config to struct
func GetORMDBConfig() (*ORMDBConfig, error) {
	config := &ORMDBConfig{Sqlite3Config: &Sqlite3Config{}}
	err := viperutil.EnhancedExactUnmarshalKey("ledger.state.ormDBConfig", config)
	if err != nil {
		return nil, errors.WithMessage(err, "load ormdb config from yaml failed")
	}
	return config, nil
}
