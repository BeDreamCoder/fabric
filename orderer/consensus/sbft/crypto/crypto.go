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

package crypto

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func ParseCertPEM(certFile string) ([]byte, error) {
	certBytes, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, err
	}

	var b *pem.Block
	for {
		b, certBytes = pem.Decode(certBytes)
		if b == nil {
			break
		}
		if b.Type == "CERTIFICATE" {
			break
		}
	}

	if b == nil {
		return nil, fmt.Errorf("no certificate found")
	}

	return b.Bytes, nil
}

func LoadX509KeyPair(dir string) (*tls.Certificate, error) {
	//ReadDir returns sorted list of files in dir
	signPath := filepath.Join(dir, "signcerts")
	fis, err := ioutil.ReadDir(signPath)
	if err != nil {
		return nil, fmt.Errorf("ReadDir signcerts failed %s\n", err)
	}
	if len(fis) != 1 {
		return nil, fmt.Errorf("Invalid signcerts dir there are %v files. ", len(fis))
	}
	signFile := filepath.Join(signPath, fis[0].Name())

	keyPath := filepath.Join(dir, "keystore")
	keyFis, err := ioutil.ReadDir(keyPath)
	if err != nil {
		return nil, fmt.Errorf("ReadDir keystore failed %s\n", err)
	}
	if len(keyFis) != 1 {
		return nil, fmt.Errorf("Invalid keystore dir there are %v files. ", len(keyFis))
	}
	keyFile := filepath.Join(keyPath, keyFis[0].Name())

	cert, err := tls.LoadX509KeyPair(signFile, keyFile)
	if err != nil {
		return nil, err
	}
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, err
	}
	return &cert, nil
}
