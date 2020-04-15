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
	sb "github.com/hyperledger/fabric/protos/orderer/sbft"
)

// Request proposes a new request to the BFT network.ls
func (s *SBFT) Request(req []byte) {
	logger.Infof("Replica %d: broadcasting a request", s.id)
	s.broadcast(&sb.Msg{Type: &sb.Msg_Request{Request: &sb.Request{Payload: req}}})
}

func (s *SBFT) handleRequest(req *sb.Request, src uint64) {
	logger.Infof("Replica %d: inserting request", s.id)
	if s.isPrimary() && s.activeView {
		blocks, valid := s.sys.Validate(s.chainId, req)
		if !valid {
			logger.Errorf("Validate request invalid")
			// this one is problematic, lets skip it
			//delete(s.pending, key)
			return
		}
		if len(blocks) == 0 {
			s.startBatchTimer()
		} else {
			for _, v := range blocks {
				key := hash2str(hash(v[0].Payload))
				logger.Infof("Replica %d: inserting %x into pending", s.id, key)
				s.pending[key] = v[0]
				s.validated[key] = valid
			}
			s.blocks = append(s.blocks, blocks...)
			s.maybeSendNextBatch()
		}
	}
}

////////////////////////////////////////////////

func (s *SBFT) startBatchTimer() {
	if s.batchTimer == nil {
		s.batchTimer = s.sys.Timer(s.support.SharedConfig().BatchTimeout(), s.cutAndMaybeSend)
	}
}

func (s *SBFT) cutAndMaybeSend() {
	batch := s.sys.Cut(s.chainId)
	s.blocks = append(s.blocks, batch)
	s.maybeSendNextBatch()
}

func (s *SBFT) batchSize() uint64 {
	size := uint64(0)
	if len(s.blocks) == 0 {
		return size
	}
	for _, req := range s.blocks[0] {
		size += uint64(len(req.Payload))
	}
	return size
}

func (s *SBFT) maybeSendNextBatch() {
	if s.batchTimer != nil {
		s.batchTimer.Cancel()
		s.batchTimer = nil
	}

	if !s.isPrimary() || !s.activeView {
		return
	}

	if !s.cur.checkpointDone {
		return
	}

	if len(s.blocks) == 0 {
		hasPending := len(s.pending) != 0
		for k, req := range s.pending {
			if s.validated[k] == false {
				batches, valid := s.sys.Validate(s.chainId, req)
				s.blocks = append(s.blocks, batches...)
				if !valid {
					logger.Panicf("Replica %d: one of our own pending requests is erroneous.", s.id)
					delete(s.pending, k)
					continue
				}
				s.validated[k] = true
			}
		}
		if len(s.blocks) == 0 {
			// if we have no pending, every req was included in blocks
			if !hasPending {
				return
			}
			// we have pending reqs that were just sent for validation or
			// were already sent (they are in s.validated)
			batch := s.sys.Cut(s.chainId)
			s.blocks = append(s.blocks, batch)
		}
	}

	block := s.blocks[0]
	s.blocks = s.blocks[1:]
	logger.Infof("Send Preprepare for chainId: %s", s.chainId)
	s.sendPreprepare(block)
}
