package channelconfig

import (
	cb "github.com/hyperledger/fabric/protos/common"
	ab "github.com/hyperledger/fabric/protos/orderer"
	"github.com/pkg/errors"
)

// Add by ztl

const (
	// ConsensusGroupKey is the group name for the consensus config.
	ConsensusGroupKey = "Consensus"
)

type ConsensusProtos struct {
	ConsensusType    *ab.ConsensusType
	OrdererAddresses *cb.OrdererAddresses
}

type ConsensusConfig struct {
	protos *ConsensusProtos
}

// NewConsensusConfig creates a new instance of the consensus config
func NewConsensusConfig(consensusGroup *cb.ConfigGroup) (*ConsensusConfig, error) {
	cc := &ConsensusConfig{
		protos: &ConsensusProtos{},
	}

	if err := DeserializeProtoValuesFromGroup(consensusGroup, cc.protos); err != nil {
		return nil, errors.Wrap(err, "failed to deserialize values")
	}

	return cc, nil
}

func (cc *ConsensusConfig) ConsensusType() string {
	return cc.protos.ConsensusType.Type
}

func (cc *ConsensusConfig) OrdererAddresses() []string {
	return cc.protos.OrdererAddresses.Addresses
}
