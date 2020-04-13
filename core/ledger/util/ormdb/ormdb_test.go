/*
Copyright Zhigui.com. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package ormdb

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewORMDBInstance(t *testing.T) {
	yaml := "---\n" +
		"ledger:\n" +
		"  state:\n" +
		"    ormDBConfig:\n" +
		"      username: test\n" +
		"      sqlite3Config:\n" +
		"        path: /tmp/ormdb\n"

	defer viper.Reset()
	viper.SetConfigType("yaml")

	if err := viper.ReadConfig(bytes.NewReader([]byte(yaml))); err != nil {
		t.Fatalf("Error reading config: %s", err)
	}

	ormDBInstance, err := NewORMDBInstance(nil)
	assert.NoError(t, err)
	assert.NotNil(t, ormDBInstance)
	assert.Equal(t, "test", ormDBInstance.Config.Username)
	assert.Equal(t, "/tmp/ormdb", ormDBInstance.Config.Sqlite3Config.Path)
}

func TestCreateORMDatabase(t *testing.T) {
	yaml := "---\n" +
		"ledger:\n" +
		"  state:\n" +
		"    ormDBConfig:\n" +
		"      username: test\n" +
		"      sqlite3Config:\n" +
		"        path: /tmp/ormdb\n"

	defer viper.Reset()
	viper.SetConfigType("yaml")

	if err := viper.ReadConfig(bytes.NewReader([]byte(yaml))); err != nil {
		t.Fatalf("Error reading config: %s", err)
	}

	ormDBInstance, err := NewORMDBInstance(nil)
	assert.NoError(t, err)

	db, err := CreateORMDatabase(ormDBInstance, "test", "sqlite3")
	assert.NoError(t, err)
	assert.NotNil(t, db)
	err = DeleteORMDatabase(db, "test")
	assert.NoError(t, err)
}
