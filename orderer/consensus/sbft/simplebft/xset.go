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
	"reflect"

	sb "github.com/hyperledger/fabric/protos/orderer/sbft"
)

// makeXset returns a request subject that should be proposed as blocks
// for new-view.  If there is no request to select (null request), it
// will return nil for subject.  makeXset always returns a blocks for
// the most recent checkpoint.
func (s *SBFT) makeXset(vcs []*sb.ViewChange) (*sb.Subject, *sb.Batch, bool) {
	// first select base commit (equivalent to checkpoint/low water mark)
	var best *sb.Batch
	for _, vc := range vcs {
		seq := vc.Checkpoint.DecodeHeader().Seq
		if best == nil || seq > best.DecodeHeader().Seq {
			best = vc.Checkpoint
		}
	}

	if best == nil {
		return nil, nil, false
	}

	next := best.DecodeHeader().Seq + 1
	logger.Debugf("replica %d: xset starts at commit %d", s.id, next)

	// now determine which request could have executed for best+1
	var xset *sb.Subject

	// This is according to Castro's TOCS PBFT, Fig. 4
	// find some message m in S,
	emptycount := 0
nextm:
	for _, m := range vcs {
		notfound := true
		// which has <n,d,v> in its Pset
		for _, mtuple := range m.Pset {
			logger.Debugf("replica %d: trying %v", s.id, mtuple)
			if mtuple.Seq.Seq < next {
				continue
			}

			// we found an entry for next
			notfound = false

			// A1. where 2f+1 messages mp from S
			count := 0
		nextmp:
			for _, mp := range vcs {
				// "low watermark" is less than n
				if mp.Checkpoint.DecodeHeader().Seq > mtuple.Seq.Seq {
					continue
				}
				// and all <n,d',v'> in its Pset
				for _, mptuple := range mp.Pset {
					logger.Debugf("replica %d: matching %v", s.id, mptuple)
					if mptuple.Seq.Seq != mtuple.Seq.Seq {
						continue
					}

					// either v' < v or (v' == v and d' == d)
					if mptuple.Seq.View < mtuple.Seq.View ||
						(mptuple.Seq.View == mtuple.Seq.View && reflect.DeepEqual(mptuple.Digest, mtuple.Digest)) {
						continue
					} else {
						continue nextmp
					}
				}
				count += 1
			}
			if count < s.viewChangeQuorum() {
				continue
			}
			logger.Debugf("replica %d: found %d replicas for Pset %d/%d", s.id, count, mtuple.Seq.Seq, mtuple.Seq.View)

			// A2. f+1 messages mp from S
			count = 0
			for _, mp := range vcs {
				// and all <n,d',v'> in its Qset
				for _, mptuple := range mp.Qset {
					if mptuple.Seq.Seq != mtuple.Seq.Seq {
						continue
					}
					if mptuple.Seq.View < mtuple.Seq.View {
						continue
					}
					// d' == d
					if !reflect.DeepEqual(mptuple.Digest, mtuple.Digest) {
						continue
					}
					count += 1
					// there exists one ...
					break
				}
			}
			if count < s.oneCorrectQuorum() {
				continue
			}
			logger.Debugf("replica %d: found %d replicas for Qset %d", s.id, count, mtuple.Seq.Seq)

			logger.Debugf("replica %d: selecting %d with %x", s.id, next, mtuple.Digest)
			xset = &sb.Subject{
				Seq:    &sb.SeqView{Seq: next, View: s.view},
				Digest: mtuple.Digest,
			}
			break nextm
		}

		if notfound {
			emptycount += 1
		}
	}

	// B. otherwise select null request
	// We actually don't select a null request, but report the most recent blocks instead.
	if emptycount >= s.viewChangeQuorum() {
		logger.Debugf("replica %d: no pertinent requests found for %d", s.id, next)
		return nil, best, true
	}

	if xset == nil {
		return nil, nil, false
	}

	return xset, best, true
}
