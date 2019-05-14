package crypto

import (
	"encoding/hex"
	"errors"
	"github.com/drep-project/binary"
	"github.com/drep-project/drep-chain/common"
	"github.com/drep-project/drep-chain/crypto/secp256k1"
	"github.com/drep-project/drep-chain/crypto/sha3"
	"math/big"
	"reflect"
)

const (
	HashLength    = 32
	AddressLength = 20
)

var (
	ErrExceedHashLength = errors.New("bytes length exceed maximum hash length of 32")
	hashT               = reflect.TypeOf(Hash{})
	addressT            = reflect.TypeOf(CommonAddress{})
)

type CommonAddress [AddressLength]byte

func (addr CommonAddress) IsEmpty() bool {
	return addr == CommonAddress{}
}

func String2Address(s string) CommonAddress {
	if s == "" {
		return CommonAddress{}
	}
	bytes := common.MustDecode(s)
	return Bytes2Address(bytes)
}

func Bytes2Address(b []byte) CommonAddress {
	if b == nil {
		return CommonAddress{}
	}
	var addr CommonAddress
	addr.SetBytes(b)
	return addr
}

func (addr *CommonAddress) SetBytes(b []byte) {
	if len(b) > len(addr) {
		copy(addr[:], b[len(b)-AddressLength:])
	} else {
		copy(addr[AddressLength-len(b):], b)
	}
}

func (addr CommonAddress) Bytes() []byte {
	return addr[:]
}

func Hex2Address(s string) CommonAddress {
	if s == "" {
		return CommonAddress{}
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		return CommonAddress{}
	}
	return Bytes2Address(b)
}

func (addr CommonAddress) Hex() string {
	return hex.EncodeToString(addr.Bytes())
}

func Big2Address(x *big.Int) CommonAddress {
	return Bytes2Address(x.Bytes())
}

func (addr CommonAddress) Big() *big.Int {
	return new(big.Int).SetBytes(addr.Bytes())
}

// MarshalText returns the hex representation of a.
func (addr CommonAddress) MarshalText() ([]byte, error) {
	return common.Bytes(addr[:]).MarshalText()
}

// UnmarshalText parses a hash in hex syntax.
func (addr *CommonAddress) UnmarshalText(input []byte) error {
	return common.UnmarshalFixedText("Address", input, addr[:])
}

// UnmarshalJSON parses a hash in hex syntax.
func (addr *CommonAddress) UnmarshalJSON(input []byte) error {
	return common.UnmarshalFixedJSON(addressT, input, addr[:])
}

func PubKey2Address(pubKey *secp256k1.PublicKey) CommonAddress {
	return Bytes2Address(sha3.Keccak256(pubKey.Serialize()))
}

type ByteCode []byte

func GetByteCodeHash(byteCode ByteCode) Hash {
	return Bytes2Hash(sha3.Keccak256(byteCode))
}

func GetByteCodeAddress(callerAddr CommonAddress, nonce uint64) CommonAddress {
	b, err := binary.Marshal(
		struct {
			CallerAddr CommonAddress
			Nonce      uint64
		}{
			CallerAddr: callerAddr,
			Nonce:      nonce,
		})
	if err != nil {
		return CommonAddress{}
	}
	return Bytes2Address(sha3.Keccak256(b))
}
