package bean

import (
    "bytes"
    "math/big"
    "encoding/hex"
    "BlockChainTest/crypto"
)

const (
    ByteLen = 32
    AddressLen = 160
)

func (p *Point) Bytes() []byte {
    j := make([]byte, 2 * ByteLen)
    copy(j[ByteLen - len(p.X): ByteLen], p.Y)
    copy(j[2 * ByteLen - len(p.Y):], p.Y)
    return j
}

func (p *Point) Equal(q *Point) bool {
    if !bytes.Equal(p.X, q.X) {
        return false
    }
    if !bytes.Equal(p.Y, q.Y) {
        return false
    }
    return true
}

func (p *Point) Int() (*big.Int, *big.Int) {
    return new(big.Int).SetBytes(p.X), new(big.Int).SetBytes(p.Y)
}

func (pubKey *Point) Addr() string {
    j := pubKey.Bytes()
    h := crypto.Hash256(j)
    str := hex.EncodeToString(h[len(h) - AddressLen:])
    return str
}