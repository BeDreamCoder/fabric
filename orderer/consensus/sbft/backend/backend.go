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

package backend

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ecdsa"
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/asn1"
	"encoding/gob"
	"fmt"
	"io"
	"math/big"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/orderer/consensus"
	"github.com/hyperledger/fabric/orderer/consensus/sbft/connection"
	cry "github.com/hyperledger/fabric/orderer/consensus/sbft/crypto"
	"github.com/hyperledger/fabric/orderer/consensus/sbft/persist"
	"github.com/hyperledger/fabric/orderer/consensus/sbft/simplebft"
	cb "github.com/hyperledger/fabric/protos/common"
	sb "github.com/hyperledger/fabric/protos/orderer/sbft"
	"github.com/hyperledger/fabric/protos/utils"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

const headerIndex = 0
const signaturesIndex = 1
const metadataLen = 2

var logger = logging.MustGetLogger("orderer.consensus.sbft.backend")

type Backend struct {
	conn        *connection.Manager
	lock        sync.Mutex
	peers       map[uint64]chan<- *sb.MultiChainMsg
	queue       chan Executable
	persistence *persist.Persist

	self *PeerInfo
	// address to PeerInfo mapping
	peerInfo map[string]*PeerInfo
	// msg sign and verify
	signCert    *tls.Certificate
	verifyCerts map[uint64]*x509.Certificate // replica number to cert

	// chainId to instance mapping
	consensus   map[string]simplebft.Receiver
	bc          map[string]*blockCreator
	lastBatches map[string]*sb.Batch
	supports    map[string]consensus.ConsenterSupport
}

type consensusConn Backend

type StackConfig struct {
	ListenAddr string
	CertFile   string
	KeyFile    string
	WALDir     string
	SnapDir    string
}

type PeerInfo struct {
	info connection.PeerInfo
	id   uint64
}

type peerInfoSlice []*PeerInfo

func (pi peerInfoSlice) Len() int {
	return len(pi)
}

func (pi peerInfoSlice) Less(i, j int) bool {
	return strings.Compare(pi[i].info.Fingerprint(), pi[j].info.Fingerprint()) == -1
}

func (pi peerInfoSlice) Swap(i, j int) {
	pi[i], pi[j] = pi[j], pi[i]
}

func NewBackend(peers map[string]*sb.Consenter, conn *connection.Manager, persist *persist.Persist, tlsCert []byte, localMsp string) (*Backend, error) {
	c := &Backend{
		conn:        conn,
		peers:       make(map[uint64]chan<- *sb.MultiChainMsg),
		peerInfo:    make(map[string]*PeerInfo),
		verifyCerts: make(map[uint64]*x509.Certificate),
		supports:    make(map[string]consensus.ConsenterSupport),
		consensus:   make(map[string]simplebft.Receiver),
		bc:          make(map[string]*blockCreator),
		lastBatches: make(map[string]*sb.Batch),
	}
	if signCert, err := cry.LoadX509KeyPair(localMsp); err != nil {
		return nil, err
	} else {
		c.signCert = signCert
	}

	var peerInfo []*PeerInfo
	for addr, consenter := range peers {
		logger.Infof("New peer connect addr: %s", addr)
		pi, err := connection.NewPeerInfo(addr, consenter.ServerTlsCert, consenter.ClientSignCert)
		if err != nil {
			return nil, err
		}
		cpi := &PeerInfo{info: pi}
		if bytes.Equal(tlsCert, consenter.ServerTlsCert) {
			c.self = cpi
		}
		peerInfo = append(peerInfo, cpi)
		c.peerInfo[pi.Fingerprint()] = cpi
	}

	sort.Sort(peerInfoSlice(peerInfo))
	for i, pi := range peerInfo {
		pi.id = uint64(i)
		c.verifyCerts[pi.id] = pi.info.VerifyCert()
		logger.Infof("replica %d: %s", i, pi.info.Fingerprint())
	}

	if c.self == nil {
		return nil, fmt.Errorf("peer list does not contain local node")
	}

	logger.Infof("Current peer is replica %d (%s)", c.self.id, c.self.info)

	sb.RegisterConsensusServer(conn.Server, (*consensusConn)(c))
	c.persistence = persist
	c.queue = make(chan Executable)
	return c, nil
}

// GetMyId returns the ID of the backend in the SFTT network (1..N)
func (b *Backend) GetMyId() uint64 {
	return b.self.id
}

// Enqueue enqueues an Envelope for a chainId for ordering, marshalling it first
func (b *Backend) Enqueue(chainID string, env *cb.Envelope) error {
	requestBytes, err := proto.Marshal(env)
	if err != nil {
		return err
	}
	b.enqueueRequest(chainID, requestBytes)
	return nil
}

func (b *Backend) connectWorker(peer *PeerInfo) {
	timeout := 1 * time.Second

	delay := time.After(0)
	for {
		// pace reconnect attempts
		<-delay

		// set up for next
		delay = time.After(timeout)

		logger.Infof("connecting to replica %d (%s)", peer.id, peer.info)
		conn, err := b.conn.DialPeer(peer.info, grpc.WithBlock(), grpc.WithTimeout(timeout))
		if err != nil {
			logger.Warningf("could not connect to replica %d (%s): %s", peer.id, peer.info, err)
			continue
		}

		ctx := context.TODO()

		client := sb.NewConsensusClient(conn)
		consensus, err := client.Consensus(ctx, &sb.Handshake{})
		if err != nil {
			logger.Warningf("could not establish consensus stream with replica %d (%s): %s", peer.id, peer.info, err)
			continue
		}
		logger.Noticef("connection to replica %d (%s) established", peer.id, peer.info)

		for {
			msg, err := consensus.Recv()
			//if err == io.EOF || err == transport.ErrConnClosing {
			//	break
			//}
			if err == io.EOF {
				break
			}
			if err != nil {
				logger.Warningf("consensus stream with replica %d (%s) broke: %v", peer.id, peer.info, err)
				break
			}
			b.enqueueForReceive(msg.ChainID, msg.Msg, peer.id)
		}
	}
}

func (b *Backend) enqueueConnection(chainID string, peerid uint64) {
	go func() {
		b.queue <- &connectionEvent{chainID: chainID, peerid: peerid}
	}()
}

func (b *Backend) enqueueRequest(chainID string, request []byte) {
	go func() {
		b.queue <- &requestEvent{chainId: chainID, req: request}
	}()
}

func (b *Backend) enqueueForReceive(chainID string, msg *sb.Msg, src uint64) {
	go func() {
		b.queue <- &msgEvent{chainId: chainID, msg: msg, src: src}
	}()
}

func (b *Backend) initTimer(t *Timer, d time.Duration) {
	send := func() {
		if t.execute {
			b.queue <- t
		}
	}
	time.AfterFunc(d, send)
}

func (b *Backend) run() {
	for {
		e := <-b.queue
		e.Execute(b)
	}
}

// AddSbftPeer adds a new SBFT peer for the given chainId using the given support and configuration
func (b *Backend) AddSbftPeer(chainID string, support consensus.ConsenterSupport, config *sb.Options) (*simplebft.SBFT, error) {
	b.supports[chainID] = support
	return simplebft.New(b.GetMyId(), chainID, support, config, b)
}

func (b *Backend) Validate(chainID string, req *sb.Request) ([][]*sb.Request, bool) {
	// ([][]*cb.Envelope, bool)
	// If the message is a valid normal message and fills a batch, the batch, true is returned
	// If the message is a valid special message (like a config message) it terminate the current batch
	// and returns the current batch, plus a second batch containing the special transaction and true
	logger.Debugf("Validate request %v for chain %s", req, chainID)
	env := &cb.Envelope{}
	err := proto.Unmarshal(req.Payload, env)
	if err != nil {
		logger.Panicf("Request format error: %s", err)
	}

	batches, pending, err := b.ordered(chainID, env)
	if err != nil {
		logger.Errorf("Failed to order message: %s", err)
		return nil, false
	}
	logger.Infof("Envelope order pending: %t", pending)

	if !pending {
		blocks := b.propose(chainID, batches...)

		if len(blocks) == 1 {
			rb1 := toRequestBlock(blocks[0])
			return [][]*sb.Request{rb1}, true
		}
		if len(blocks) == 2 {
			rb1 := toRequestBlock(blocks[0])
			rb2 := toRequestBlock(blocks[1])
			return [][]*sb.Request{rb1, rb2}, true
		}

	}

	return nil, true
}

func (b *Backend) propose(chainID string, batches ...[]*cb.Envelope) (blocks []*cb.Block) {
	for _, batch := range batches {
		block := b.bc[chainID].createNextBlock(batch)
		logger.Infof("Created block [%d]", block.Header.Number)
		blocks = append(blocks, block)
	}

	return
}

func (b *Backend) isConfig(env *cb.Envelope) bool {
	h, err := utils.ChannelHeader(env)
	if err != nil {
		logger.Panicf("failed to extract channel header from envelope")
	}

	return h.Type == int32(cb.HeaderType_CONFIG) || h.Type == int32(cb.HeaderType_ORDERER_TRANSACTION)
}

// Orders the envelope in the `req` content. Request.
// Returns
//   -- batches [][]*cb.Envelope; the batches cut,
//   -- pending bool; if there are envelopes pending to be ordered,
//   -- err error; the error encountered, if any.
// It takes care of config messages as well as the revalidation of messages if the config sequence has advanced.
func (b *Backend) ordered(chainID string, payload *cb.Envelope) (batches [][]*cb.Envelope, pending bool, err error) {
	if b.isConfig(payload) {
		// ConfigMsg
		logger.Infof("Config message was validated")
		payload, _, err = b.supports[chainID].ProcessConfigMsg(payload)
		if err != nil {
			return nil, true, errors.Errorf("bad config message: %s", err)
		}

		batch := b.supports[chainID].BlockCutter().Cut()
		batches = [][]*cb.Envelope{}
		if len(batch) != 0 {
			batches = append(batches, batch)
		}
		batches = append(batches, []*cb.Envelope{payload})
		return batches, false, nil
	}
	// it is a normal message
	logger.Infof("Normal message was validated")
	if _, err := b.supports[chainID].ProcessNormalMsg(payload); err != nil {
		return nil, true, errors.Errorf("bad normal message: %s", err)
	}

	batches, pending = b.supports[chainID].BlockCutter().Ordered(payload)
	return batches, pending, nil
}

func (b *Backend) Cut(chainID string) []*sb.Request {
	envBatch := b.supports[chainID].BlockCutter().Cut()
	block := b.bc[chainID].createNextBlock(envBatch)
	return toRequestBlock(block)
}

func toRequestBlock(block *cb.Block) []*sb.Request {
	rqs := make([]*sb.Request, 0, 1)
	requestBytes, err := utils.Marshal(block)
	if err != nil {
		logger.Panicf("Cannot marshal envelope: %s", err)
	}
	rq := &sb.Request{Payload: requestBytes}
	rqs = append(rqs, rq)

	return rqs
}

// Consensus implements the SBFT consensus gRPC interface
func (c *consensusConn) Consensus(_ *sb.Handshake, srv sb.Consensus_ConsensusServer) error {
	pi := connection.GetPeerInfo(srv)
	peer, ok := c.peerInfo[pi.Fingerprint()]

	if !ok || !peer.info.Cert().Equal(pi.Cert()) {
		logger.Infof("rejecting connection from unknown replica %s", pi)
		return fmt.Errorf("unknown peer certificate")
	}
	logger.Infof("connection from replica %d (%s)", peer.id, pi)

	ch := make(chan *sb.MultiChainMsg)
	c.lock.Lock()
	if oldCh, ok := c.peers[peer.id]; ok {
		logger.Debugf("replacing connection from replica %d", peer.id)
		close(oldCh)
	}
	c.peers[peer.id] = ch
	c.lock.Unlock()

	for chainID := range c.supports {
		((*Backend)(c)).enqueueConnection(chainID, peer.id)
	}

	var err error
	for msg := range ch {
		err = srv.Send(msg)
		if err != nil {
			c.lock.Lock()
			delete(c.peers, peer.id)
			c.lock.Unlock()

			logger.Infof("lost connection from replica %d (%s): %s", peer.id, pi, err)
		}
	}

	return err
}

// Unicast sends to all external SBFT peers
func (b *Backend) Broadcast(msg *sb.MultiChainMsg) error {
	b.lock.Lock()
	for _, ch := range b.peers {
		ch <- msg
	}
	b.lock.Unlock()
	return nil
}

// Unicast sends to a specific external SBFT peer identified by chainId and dest
func (b *Backend) Unicast(chainID string, msg *sb.Msg, dest uint64) error {
	b.lock.Lock()
	ch, ok := b.peers[dest]
	b.lock.Unlock()

	if !ok {
		err := fmt.Errorf("peer not found: %v", dest)
		logger.Debug(err)
		return err
	}
	ch <- &sb.MultiChainMsg{Msg: msg, ChainID: chainID}
	return nil
}

// AddReceiver adds a receiver instance for a given chainId
func (b *Backend) AddReceiver(chainId string, recv simplebft.Receiver) {
	b.consensus[chainId] = recv
	block := b.supports[chainId].Block(b.supports[chainId].Height() - 1)
	b.bc[chainId] = &blockCreator{
		hash:   block.Header.Hash(),
		number: block.Header.Number,
	}
	b.lastBatches[chainId] = &sb.Batch{Header: nil, Signatures: nil, Payloads: [][]byte{}}
}

// Send sends to a specific SBFT peer identified by chainId and dest
func (b *Backend) Send(chainID string, msg *sb.Msg, dest uint64) {
	if dest == b.self.id {
		b.enqueueForReceive(chainID, msg, b.self.id)
		return
	}
	b.Unicast(chainID, msg, dest)
}

// Timer starts a timer
func (b *Backend) Timer(d time.Duration, tf func()) simplebft.Canceller {
	tm := &Timer{tf: tf, execute: true}
	b.initTimer(tm, d)
	return tm
}

// Deliver writes a block
func (b *Backend) Deliver(chainId string, batch *sb.Batch) {
	block := utils.UnmarshalBlockOrPanic(batch.Payloads[0])
	b.lastBatches[chainId] = batch
	// TODO SBFT needs to use Rawledger's structures and signatures over the Block.
	// This a quick and dirty solution to make it work.
	blockMetadata := &cb.BlockMetadata{}
	metadata := make([][]byte, metadataLen)
	metadata[headerIndex] = batch.Header
	metadata[signaturesIndex] = encodeSignatures(batch.Signatures)
	blockMetadata.Metadata = metadata
	m := utils.MarshalOrPanic(blockMetadata)
	if utils.IsConfigBlock(block) {
		logger.Infof("Writing config block [%d] for chainId: %s to ledger", block.Header.Number, chainId)
		b.supports[chainId].WriteConfigBlock(block, m)
	} else {
		logger.Infof("Writing block [%d] for chainId: %s to ledger", block.Header.Number, chainId)
		b.supports[chainId].WriteBlock(block, m)
	}
}

// Persist persists data identified by a chainId and a key
func (b *Backend) Persist(chainId string, key string, data proto.Message) {
	compk := fmt.Sprintf("chain-%s-%s", chainId, key)
	if data == nil {
		b.persistence.DelState(compk)
	} else {
		bytes, err := proto.Marshal(data)
		if err != nil {
			panic(err)
		}
		b.persistence.StoreState(compk, bytes)
	}
}

// Restore loads persisted data identified by chainId and key
func (b *Backend) Restore(chainId string, key string, out proto.Message) bool {
	compk := fmt.Sprintf("chain-%s-%s", chainId, key)
	val, err := b.persistence.ReadState(compk)
	if err != nil {
		return false
	}
	err = proto.Unmarshal(val, out)
	return err == nil
}

// LastBatch returns the last batch for a given chain identified by its ID
func (b *Backend) LastBatch(chainId string) *sb.Batch {
	return b.lastBatches[chainId]
}

// Sign signs a given data
func (b *Backend) Sign(data []byte) []byte {
	return Sign(b.signCert.PrivateKey, data)
}

// CheckSig checks a signature
func (b *Backend) CheckSig(data []byte, src uint64, sig []byte) error {
	cert, ok := b.verifyCerts[src]
	if !ok {
		panic("No public key found: certificate leaf is nil.")
	}
	return CheckSig(cert.PublicKey, data, sig)
}

// Reconnect requests connection to a replica identified by its ID and chainId
func (b *Backend) Reconnect(chainId string, replica uint64) {
	b.enqueueConnection(chainId, replica)
}

func (b *Backend) StartAndConnectWorkers() {
	go func() {
		if err := b.conn.Server.Serve(b.conn.Listener); err != nil {
			logger.Errorf("Backend Server start error: %v", err)
			panic(err)
		}
	}()

	for _, peer := range b.peerInfo {
		if peer == b.self {
			continue
		}
		go b.connectWorker(peer)
	}

	go b.run()
}

// Sign signs a given data
func Sign(privateKey crypto.PrivateKey, data []byte) []byte {
	var err error
	var encsig []byte
	hash := sha256.Sum256(data)
	switch pvk := privateKey.(type) {
	case *rsa.PrivateKey:
		encsig, err = pvk.Sign(crand.Reader, hash[:], crypto.SHA256)
		if err != nil {
			panic(err)
		}
	case *ecdsa.PrivateKey:
		r, s, err := ecdsa.Sign(crand.Reader, pvk, hash[:])
		if err != nil {
			panic(err)
		}
		encsig, _ = asn1.Marshal(struct{ R, S *big.Int }{r, s})
	default:
		panic("Unsupported private key type given.")
	}
	if err != nil {
		panic(err)
	}
	return encsig
}

// CheckSig checks a signature
func CheckSig(publicKey crypto.PublicKey, data []byte, sig []byte) error {
	hash := sha256.Sum256(data)
	switch p := publicKey.(type) {
	case *ecdsa.PublicKey:
		s := struct{ R, S *big.Int }{}
		rest, err := asn1.Unmarshal(sig, &s)
		if err != nil {
			return err
		}
		if len(rest) != 0 {
			return fmt.Errorf("invalid signature (problem with asn unmarshalling for ECDSA)")
		}
		ok := ecdsa.Verify(p, hash[:], s.R, s.S)
		if !ok {
			return fmt.Errorf("invalid signature (problem with verification)")
		}
		return nil
	case *rsa.PublicKey:
		err := rsa.VerifyPKCS1v15(p, crypto.SHA256, hash[:], sig)
		return err
	default:
		return fmt.Errorf("Unsupported public key type.")
	}
}

func encodeSignatures(signatures map[uint64][]byte) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(signatures)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func decodeSignatures(encodedSignatures []byte) map[uint64][]byte {
	if len(encodedSignatures) == 0 {
		return nil
	}
	buf := bytes.NewBuffer(encodedSignatures)
	var r map[uint64][]byte
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&r)
	if err != nil {
		panic(err)
	}
	return r
}
