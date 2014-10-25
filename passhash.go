// Obdi - a REST interface and GUI for deploying software
// Copyright (C) 2014  Mark Clarkson
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
