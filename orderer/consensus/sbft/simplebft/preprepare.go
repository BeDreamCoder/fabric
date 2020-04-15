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

package simplebft

import (
	"bytes"
	"time"

	sb "github.com/hyperledger/fabric/protos/orderer/sbft"
)

func (s *SBFT) sendPreprepare(batch []*sb.Request) {
	seq := s.nextSeq()

	if len(batch) != 1 {
		logger.Panicf("replica %d: sendPreprepare batch length %d", s.id, len(batch))
	}

	data := make([][]byte, len(batch))
	for i, req := range batch {
		data[i] = req.Payload
	}

	lasthash := hash(s.sys.LastBatch(s.chainId).Header)

	m := &sb.Preprepare{
		Seq:   &seq,
		Batch: s.makeBatch(seq.Seq, lasthash, data),
	}

	s.sys.Persist(s.chainId, preprepared, m)
	s.broadcast(&sb.Msg{Type: &sb.Msg_Preprepare{Preprepare: m}})
	logger.Infof("replica %d: sendPreprepare", s.id)
	s.handleCheckedPreprepare(m)
}

func (s *SBFT) handlePreprepare(pp *sb.Preprepare, src uint64) {
	if src == s.id {
		logger.Infof("replica %d: ignoring preprepare from self: %d", s.id, src)
		return
	}
	if src != s.primaryID() {
		logger.Infof("replica %d: preprepare from non-primary %d", s.id, src)
		return
	}
	nextSeq := s.nextSeq()
	if (*pp.Seq).Seq != nextSeq.Seq || (*pp.Seq).View != nextSeq.View {
		logger.Infof("replica %d: preprepare does not match expected %v, got %v", s.id, nextSeq, *pp.Seq)
		return
	}
	if s.cur.subject.Seq.Seq == pp.Seq.Seq {
		logger.Infof("replica %d: duplicate preprepare for %v", s.id, *pp.Seq)
		return
	}
	if pp.GetBatch() == nil {
		logger.Infof("replica %d: preprepare without blocks", s.id)
		return
	}

	batchheader, err := s.checkBatch(pp.Batch, true, false)
	if err != nil || batchheader.Seq != pp.Seq.Seq {
		logger.Infof("replica %d: preprepare %v blocks head inconsistent from %d: %s", s.id, pp.Seq, src, err)
		return
	}

	prevhash := s.sys.LastBatch(s.chainId).Hash()
	if !bytes.Equal(batchheader.PrevHash, prevhash) {
		logger.Infof("replica %d: preprepare blocks prev hash does not match expected %s, got %s", s.id, hash2str(batchheader.PrevHash), hash2str(prevhash))
		return
	}

	//blockOK, committers := s.getCommittersFromBatch(pp.Batch)
	//if !blockOK {
	//	logger.Debugf("Replica %d found Byzantine block in preprepare, Seq: %d View: %d", s.id, pp.Seq.Seq, pp.Seq.View)
	//	s.sendViewChange()
	//	return
	//}
	logger.Infof("replica %d: handlePrepare", s.id)
	s.handleCheckedPreprepare(pp)
}

func (s *SBFT) acceptPreprepare(pp *sb.Preprepare) {
	sub := sb.Subject{Seq: pp.Seq, Digest: pp.Batch.Hash()}

	logger.Infof("replica %d: accepting preprepare for %v, %x", s.id, sub.Seq, sub.Digest)
	s.sys.Persist(s.chainId, preprepared, pp)

	s.cur = reqInfo{
		subject:    sub,
		timeout:    s.sys.Timer(time.Duration(s.config.RequestTimeoutNsec)*time.Nanosecond, s.requestTimeout),
		preprep:    pp,
		prep:       make(map[uint64]*sb.Subject),
		commit:     make(map[uint64]*sb.Subject),
		checkpoint: make(map[uint64]*sb.Checkpoint),
	}
}

func (s *SBFT) handleCheckedPreprepare(pp *sb.Preprepare) {
	s.acceptPreprepare(pp)
	if !s.isPrimary() {
		s.sendPrepare()
		s.processBacklog()
	}

	s.maybeSendCommit()
}

////////////////////////////////////////////////

func (s *SBFT) requestTimeout() {
	logger.Infof("replica %d: request timed out: %s", s.id, s.cur.subject.Seq)
	s.sendViewChange()
}
