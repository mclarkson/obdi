// Copyright 2016 Mark Clarkson
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"log"
	//"fmt"
	"github.com/mclarkson/obdi/external/go.crypto/bcrypt"
)

type Crypt struct {
	Pass []byte
	Hash []byte
}

func (c *Crypt) clear() {
	for i := 0; i < len(c.Pass); i++ {
		c.Pass[i] = 0
	}
}

func (c *Crypt) crypt() ([]byte, error) {
	defer c.clear()
	return bcrypt.GenerateFromPassword(c.Pass, bcrypt.DefaultCost)
}

func (c *Crypt) Crypt() {

	hash, err := c.crypt()

	if err != nil {
		log.Fatal(err)
	}

	c.Hash = hash
}

func (c *Crypt) Check() error {
	defer c.clear()
	return bcrypt.CompareHashAndPassword(c.Hash, c.Pass)
}
