package chain

import (
	"errors"
)

var (
	ErrInvalidateTimestamp       = errors.New("timestamp equals parent's")
	ErrInvalidateBlockNumber     = errors.New("invalid block number")
	ErrBlockNotFound             = errors.New("block not exist")
	ErrTxIndexOutOfRange         = errors.New("tx index out of range")
	ErrReachGasLimit             = errors.New("gas limit reached")
	ErrInvalidateBlockMultisig   = errors.New("verify multisig error")
	ErrUnsupportTxType           = errors.New("not support transaction type")
	ErrNegativeAmount            = errors.New("negative amount in tx")
	ErrExceedGasLimit            = errors.New("gas limit in tx has exceed block limit")
	ErrNonceTooHigh              = errors.New("nonce too high")
	ErrNonceTooLow               = errors.New("nonce too low")
	ErrTxPool                    = errors.New("transaction pool full")
	ErrInitStateFail             = errors.New("initChainState")
	ErrNotMathcedStateRoot       = errors.New("state root not matched")
	ErrGasUsed                   = errors.New("gas used not matched")
	ErrChainId                   = errors.New("chainId not matched")
	ErrVersion                   = errors.New("version not matched")
	ErrPreHash                   = errors.New("previous hash not matched")
	ErrBpNotInList               = errors.New("bp node not in local list")
	ErrBlockExsist               = errors.New("already have block")
	ErrBalance                   = errors.New("not enough balance")
	ErrGas                       = errors.New("not enough gas")
	ErrInsufficientBalanceForGas = errors.New("insufficient balance to pay for gas")
	ErrOrphanBlockExsist         = errors.New("already have block (orphan)")
	ErrGenesisPkNotFound         = errors.New("genesisi pubkey not found")
	ErrBlockProducerNotFound     = errors.New("block producer not found")
	ErrNotSupportRenameAlias     = errors.New("not suppport rename alias")
	ErrTooShortAlias			 = errors.New("alias too short")
	ErrTooLongAlias			  	 = errors.New("alias too long")
	ErrUnsupportAliasChar		 = errors.New("alias only support number and letter")
)