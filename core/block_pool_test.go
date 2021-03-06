// Copyright (C) 2017 go-nebulas authors
//
// This file is part of the go-nebulas library.
//
// the go-nebulas library is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// the go-nebulas library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with the go-nebulas library.  If not, see <http://www.gnu.org/licenses/>.
//

package core

import (
	"testing"

	"time"

	"github.com/nebulasio/go-nebulas/crypto"
	"github.com/nebulasio/go-nebulas/crypto/keystore"
	"github.com/nebulasio/go-nebulas/crypto/keystore/secp256k1"
	"github.com/nebulasio/go-nebulas/storage"
	"github.com/nebulasio/go-nebulas/util"
	"github.com/stretchr/testify/assert"
)

type MockConsensus int

func (c MockConsensus) VerifyBlock(block *Block) error {
	return nil
}

func TestBlockPool(t *testing.T) {
	storage, _ := storage.NewMemoryStorage()
	bc, _ := NewBlockChain(0, storage)
	var cons MockConsensus
	bc.SetConsensusHandler(cons)
	pool := bc.bkPool
	assert.Equal(t, pool.blockCache.Len(), 0)

	ks := keystore.DefaultKS
	priv, _ := secp256k1.GeneratePrivateKey()
	pubdata, _ := priv.PublicKey().Encoded()
	from, _ := NewAddressFromPublicKey(pubdata)
	to := &Address{from.address}
	coinbase := &Address{from.address}
	ks.SetKey(from.ToHex(), priv, []byte("passphrase"))
	ks.Unlock(from.ToHex(), []byte("passphrase"), time.Second*60*60*24*365)

	key, _ := ks.GetUnlocked(from.ToHex())
	signature, _ := crypto.NewSignature(keystore.SECP256K1)
	signature.InitSign(key.(keystore.PrivateKey))

	tx1 := NewTransaction(0, from, to, util.NewUint128(), 1, TxPayloadBinaryType, []byte("nas"), util.NewUint128(), util.NewUint128())
	tx1.Sign(signature)
	tx2 := NewTransaction(0, from, to, util.NewUint128(), 2, TxPayloadBinaryType, []byte("nas"), util.NewUint128(), util.NewUint128())
	tx2.Sign(signature)
	tx3 := NewTransaction(0, from, to, util.NewUint128(), 3, TxPayloadBinaryType, []byte("nas"), util.NewUint128(), util.NewUint128())
	tx3.Sign(signature)
	bc.txPool.Push(tx1)
	bc.txPool.Push(tx2)
	bc.txPool.Push(tx3)

	block1 := NewBlock(0, coinbase, bc.tailBlock, bc.txPool, storage)
	block1.CollectTransactions(1)
	block1.Seal()

	block2 := NewBlock(0, coinbase, block1, bc.txPool, storage)
	block2.CollectTransactions(1)
	block2.Seal()

	block3 := NewBlock(0, coinbase, block2, bc.txPool, storage)
	block3.CollectTransactions(1)
	block3.Seal()

	block4 := NewBlock(0, coinbase, block3, bc.txPool, storage)
	block4.CollectTransactions(1)
	block4.Seal()

	pool.Push(block3)
	assert.Equal(t, pool.blockCache.Len(), 1)
	pool.Push(block4)
	assert.Equal(t, pool.blockCache.Len(), 2)
	pool.Push(block2)
	assert.Equal(t, pool.blockCache.Len(), 3)
	pool.Push(block1)
	assert.Equal(t, pool.blockCache.Len(), 0)
	err := pool.Push(block1)
	assert.Equal(t, err, nil)
	bc.SetTailBlock(block1)
	nonce := bc.tailBlock.GetNonce(from.Bytes())
	assert.Equal(t, nonce, uint64(1))
	bc.SetTailBlock(block2)
	nonce = bc.tailBlock.GetNonce(from.Bytes())
	assert.Equal(t, nonce, uint64(2))
	bc.SetTailBlock(block3)
	nonce = bc.tailBlock.GetNonce(from.Bytes())
	assert.Equal(t, nonce, uint64(3))

}
