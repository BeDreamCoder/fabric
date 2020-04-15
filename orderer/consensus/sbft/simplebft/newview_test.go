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
	"testing"

	sb "github.com/hyperledger/fabric/protos/orderer/sbft"
)

func TestXsetNoByz(t *testing.T) {
	s := &SBFT{config: sb.Options{N: 4, F: 1}, view: 3}
	vcs := []*sb.ViewChange{
		{
			View: 3,
			Pset: nil,
			Qset: []*sb.Subject{{Seq: &sb.SeqView{View: 1, Seq: 2}, Digest: []byte("val1")},
				{Seq: &sb.SeqView{View: 2, Seq: 2}, Digest: []byte("val2")}},
			Checkpoint: s.makeBatch(1, []byte("prev"), nil),
		},
		{
			View:       3,
			Pset:       []*sb.Subject{{Seq: &sb.SeqView{View: 1, Seq: 2}, Digest: []byte("val1")}},
			Qset:       []*sb.Subject{{Seq: &sb.SeqView{View: 1, Seq: 2}, Digest: []byte("val1")}},
			Checkpoint: s.makeBatch(1, []byte("prev"), nil),
		},
		{
			View: 3,
			Pset: []*sb.Subject{{Seq: &sb.SeqView{View: 2, Seq: 2}, Digest: []byte("val2")}},
			Qset: []*sb.Subject{{Seq: &sb.SeqView{View: 1, Seq: 2}, Digest: []byte("val1")},
				{Seq: &sb.SeqView{View: 2, Seq: 2}, Digest: []byte("val2")}},
			Checkpoint: s.makeBatch(1, []byte("prev"), nil),
		},
	}

	xset, _, ok := s.makeXset(vcs)
	if !ok {
		t.Fatal("no xset")
	}

	if !reflect.DeepEqual(xset, &sb.Subject{Seq: &sb.SeqView{View: 3, Seq: 2}, Digest: []byte("val2")}) {
		t.Error(xset)
	}
}

func TestXsetNoNew(t *testing.T) {
	s := &SBFT{config: sb.Options{N: 4, F: 1}, view: 3}
	prev := s.makeBatch(2, []byte("prev"), nil)
	vcs := []*sb.ViewChange{
		{
			View:       3,
			Pset:       []*sb.Subject{{Seq: &sb.SeqView{View: 2, Seq: 2}, Digest: []byte("val2")}},
			Qset:       []*sb.Subject{{Seq: &sb.SeqView{View: 2, Seq: 2}, Digest: []byte("val2")}},
			Checkpoint: prev,
		},
		{
			View:       3,
			Pset:       []*sb.Subject{{Seq: &sb.SeqView{View: 2, Seq: 2}, Digest: []byte("val2")}},
			Qset:       []*sb.Subject{{Seq: &sb.SeqView{View: 2, Seq: 2}, Digest: []byte("val2")}},
			Checkpoint: prev,
		},
		{
			View:       3,
			Pset:       []*sb.Subject{{Seq: &sb.SeqView{View: 2, Seq: 2}, Digest: []byte("val2")}},
			Qset:       []*sb.Subject{{Seq: &sb.SeqView{View: 2, Seq: 2}, Digest: []byte("val2")}},
			Checkpoint: prev,
		},
	}

	xset, prevBatch, ok := s.makeXset(vcs)
	if !ok {
		t.Fatal("no xset")
	}

	if xset != nil {
		t.Errorf("should have null request")
	}

	if !reflect.DeepEqual(prevBatch, prev) {
		t.Errorf("blocks don't match: %v, %v", prevBatch.DecodeHeader(), prev.DecodeHeader())
	}
}

func TestXsetByz0(t *testing.T) {
	s := &SBFT{config: sb.Options{N: 4, F: 1}, view: 3}
	vcs := []*sb.ViewChange{
		{
			View:       3,
			Pset:       nil,
			Qset:       nil,
			Checkpoint: s.makeBatch(1, []byte("prev"), nil),
		},
		{
			View:       3,
			Pset:       []*sb.Subject{{Seq: &sb.SeqView{View: 1, Seq: 2}, Digest: []byte("val1")}},
			Qset:       []*sb.Subject{{Seq: &sb.SeqView{View: 1, Seq: 2}, Digest: []byte("val1")}},
			Checkpoint: s.makeBatch(1, []byte("prev"), nil),
		},
		{
			View: 3,
			Pset: []*sb.Subject{{Seq: &sb.SeqView{View: 2, Seq: 2}, Digest: []byte("val2")}},
			Qset: []*sb.Subject{{Seq: &sb.SeqView{View: 1, Seq: 2}, Digest: []byte("val1")},
				{Seq: &sb.SeqView{View: 2, Seq: 2}, Digest: []byte("val2")}},
			Checkpoint: s.makeBatch(1, []byte("prev"), nil),
		},
	}

	_, _, ok := s.makeXset(vcs)
	if ok {
		t.Error("should not have received an xset")
	}

	vcs = append(vcs, &sb.ViewChange{
		View: 3,
		Pset: []*sb.Subject{{Seq: &sb.SeqView{View: 2, Seq: 2}, Digest: []byte("val2")}},
		Qset: []*sb.Subject{{Seq: &sb.SeqView{View: 1, Seq: 2}, Digest: []byte("val1")},
			{Seq: &sb.SeqView{View: 2, Seq: 2}, Digest: []byte("val2")}},
		Checkpoint: s.makeBatch(2, []byte("prev"), nil),
	})

	xset, _, ok := s.makeXset(vcs)
	if !ok {
		t.Error("no xset")
	}
	if xset != nil {
		t.Error("expected null request")
	}
}

func TestXsetByz2(t *testing.T) {
	s := &SBFT{config: sb.Options{N: 4, F: 1}, view: 3}
	vcs := []*sb.ViewChange{
		{
			View:       3,
			Pset:       nil,
			Qset:       []*sb.Subject{{Seq: &sb.SeqView{View: 1, Seq: 2}, Digest: []byte("val1")}},
			Checkpoint: s.makeBatch(1, []byte("prev"), nil),
		},
		{
			View:       3,
			Pset:       []*sb.Subject{{Seq: &sb.SeqView{View: 1, Seq: 2}, Digest: []byte("val1")}},
			Qset:       []*sb.Subject{{Seq: &sb.SeqView{View: 1, Seq: 2}, Digest: []byte("val1")}},
			Checkpoint: s.makeBatch(1, []byte("prev"), nil),
		},
		{
			View: 3,
			Pset: []*sb.Subject{{Seq: &sb.SeqView{View: 2, Seq: 2}, Digest: []byte("val2")}},
			Qset: []*sb.Subject{{Seq: &sb.SeqView{View: 1, Seq: 2}, Digest: []byte("val1")},
				{Seq: &sb.SeqView{View: 2, Seq: 2}, Digest: []byte("val2")}},
			Checkpoint: s.makeBatch(1, []byte("prev"), nil),
		},
	}

	_, _, ok := s.makeXset(vcs)
	if ok {
		t.Error("should not have received an xset")
	}

	vcs = append(vcs, &sb.ViewChange{
		View:       3,
		Pset:       []*sb.Subject{{Seq: &sb.SeqView{View: 1, Seq: 2}, Digest: []byte("val1")}},
		Qset:       []*sb.Subject{{Seq: &sb.SeqView{View: 1, Seq: 2}, Digest: []byte("val1")}},
		Checkpoint: s.makeBatch(1, []byte("prev"), nil),
	})

	xset, _, ok := s.makeXset(vcs)
	if !ok {
		t.Error("no xset")
	}
	if !reflect.DeepEqual(xset, &sb.Subject{Seq: &sb.SeqView{View: 3, Seq: 2}, Digest: []byte("val1")}) {
		t.Error(xset)
	}
}
