/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sbft

import (
	"strconv"
	"sync"

	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric/core/comm"
	"github.com/hyperledger/fabric/orderer/common/localconfig"
	"github.com/hyperledger/fabric/orderer/consensus"
	"github.com/hyperledger/fabric/orderer/consensus/sbft/backend"
	"github.com/hyperledger/fabric/orderer/consensus/sbft/connection"
	"github.com/hyperledger/fabric/orderer/consensus/sbft/persist"
	"github.com/hyperledger/fabric/orderer/consensus/sbft/simplebft"
	cb "github.com/hyperledger/fabric/protos/common"
	sb "github.com/hyperledger/fabric/protos/orderer/sbft"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
)

type consensusStack struct {
	persist *persist.Persist
	backend *backend.Backend
}

var logger = logging.MustGetLogger("orderer/consensus/sbft")
var once sync.Once

// Consenter interface implementation for new main application
type consenter struct {
	cert            []byte
	localMsp        string
	config          *sb.ConsensusConfig
	consensusStack  *consensusStack
	sbftStackConfig *backend.StackConfig
	sbftPeers       map[string]*simplebft.SBFT
}

type chain struct {
	chainID        string
	exitChan       chan struct{}
	consensusStack *consensusStack
}

// New creates a new consenter for the SBFT consensus scheme.
// It accepts messages being delivered via Enqueue, orders them, and then uses the blockcutter to form the messages
// into blocks before writing to the given ledger.
func New(conf *localconfig.TopLevel, srvConf comm.ServerConfig) consensus.Consenter {
	sc := &backend.StackConfig{ListenAddr: conf.SbftLocal.PeerCommAddr,
		CertFile: conf.SbftLocal.CertFile,
		KeyFile:  conf.SbftLocal.KeyFile,
		WALDir:   conf.SbftLocal.WALDir,
		SnapDir:  conf.SbftLocal.SnapDir}

	return &consenter{cert: srvConf.SecOpts.Certificate, localMsp: conf.General.LocalMSPDir, sbftStackConfig: sc}
}

func (sbft *consenter) HandleChain(support consensus.ConsenterSupport, metadata *cb.Metadata) (consensus.Chain, error) {
	logger.Infof("Starting a chain: %d", support.ChainID())

	m := &sb.ConfigMetadata{}
	if err := proto.Unmarshal(support.SharedConfig().SbftMetadata(), m); err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal consensus metadata")
	}
	if m.Options == nil {
		return nil, errors.New("Sbft options have not been provided")
	}
	if len(m.Consenters) == 0 {
		return nil, errors.New("Sbft consenters have not been provided")
	}
	peers := make(map[string]*sb.Consenter)
	for _, consenter := range m.Consenters {
		endpoint := consenter.Host + ":" + strconv.FormatUint(uint64(consenter.Port), 10)
		peers[endpoint] = consenter
	}

	sbft.config = &sb.ConsensusConfig{
		Consensus: m.Options,
		Peers:     peers,
	}

	if sbft.sbftPeers == nil {
		sbft.consensusStack = createConsensusStack(sbft)
		sbft.sbftPeers = make(map[string]*simplebft.SBFT)
	}
	sbft.sbftPeers[support.ChainID()] = initSbftPeer(sbft, support)

	return &chain{
		chainID:        support.ChainID(),
		exitChan:       make(chan struct{}),
		consensusStack: sbft.consensusStack,
	}, nil
}

func createConsensusStack(sbft *consenter) *consensusStack {
	conn, err := connection.New(sbft.sbftStackConfig.ListenAddr, sbft.sbftStackConfig.CertFile, sbft.sbftStackConfig.KeyFile)
	if err != nil {
		logger.Errorf("Error when trying to connect: %s", err)
		panic(err)
	}
	pPersist := persist.New(sbft.sbftStackConfig.WALDir)
	pBackend, err := backend.NewBackend(sbft.config.Peers, conn, pPersist, sbft.cert, sbft.localMsp)
	if err != nil {
		logger.Errorf("Backend instantiation error: %v", err)
		panic(err)
	}

	return &consensusStack{
		backend: pBackend,
		persist: pPersist,
	}
}

func initSbftPeer(sbft *consenter, support consensus.ConsenterSupport) *simplebft.SBFT {
	sbftPeer, err := sbft.consensusStack.backend.AddSbftPeer(support.ChainID(), support, sbft.config.Consensus)
	if err != nil {
		logger.Errorf("SBFT peer instantiation error.")
		panic(err)
	}
	return sbftPeer
}

// Chain interface implementation:

// Start allocates the necessary resources for staying up to date with this Chain.
// It implements the multichain.Chain interface. It is called by multichain.NewManagerImpl()
// which is invoked when the ordering process is launched, before the call to NewServer().
func (ch *chain) Start() {
	once.Do(ch.consensusStack.backend.StartAndConnectWorkers)
}

// Halt frees the resources which were allocated for this Chain
func (ch *chain) Halt() {
	// panic("There is no way to halt SBFT")
	select {
	case <-ch.exitChan:
		// Allow multiple halts without panic
	default:
		close(ch.exitChan)
	}
}

// Enqueue accepts a message and returns true on acceptance, or false on shutdown
//func (ch *chain) Enqueue(env *cb.Envelope) bool {
//	return ch.consensusStack.backend.Enqueue(ch.chainID, env)
//}

func (ch *chain) WaitReady() error {
	return nil
}

// Order accepts normal messages for ordering
func (ch *chain) Order(env *cb.Envelope, configSeq uint64) error {
	return ch.consensusStack.backend.Enqueue(ch.chainID, env)
}

// Configure accepts configuration update messages for ordering
func (ch *chain) Configure(config *cb.Envelope, configSeq uint64) error {
	return ch.consensusStack.backend.Enqueue(ch.chainID, config)
}

// Errored only closes on exit
func (ch *chain) Errored() <-chan struct{} {
	return ch.exitChan
}
