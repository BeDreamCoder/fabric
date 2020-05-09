package dfs

import (
	"github.com/colinmarc/hdfs"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/core/ledger/ledgerconfig"
	"github.com/pkg/errors"
)

var logger = flogging.MustGetLogger("archive.dfs")

func NewHDFSClient() (*hdfs.Client, error) {
	if len(ledgerconfig.GetHDFSNameNodes()) == 0 {
		errMsg := "Archive service can't be initialized, due to no namenode address in configuration"
		logger.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	if len(ledgerconfig.GetHDFSUser()) == 0 {
		errMsg := "Archive service can't be initialized, due to no HDFS user in configuration"
		logger.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	options := hdfs.ClientOptions{
		Addresses: ledgerconfig.GetHDFSNameNodes(),
		User:      ledgerconfig.GetHDFSUser(),
	}
	client, err := hdfs.NewClient(options)

	logger.Infof("Created a dfs client with options: %+v, error: %+v", options, err)
	return client, err
}
