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
	"crypto/sha256"

	"github.com/golang/protobuf/proto"
)

func hash(data []byte) []byte {
	h := sha256.Sum256(data)
	return h[:]
}

// Hash returns the hash of the Batch.
func (m *Batch) Hash() []byte {
	return hash(m.Header)
}

func (m *Batch) DecodeHeader() *BatchHeader {
	batchheader := &BatchHeader{}
	err := proto.Unmarshal(m.Header, batchheader)
	if err != nil {
		panic(err)
	}

	return batchheader
}
