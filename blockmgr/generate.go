package blockmgr

import (
	"github.com/drep-project/drep-chain/chain"
	"github.com/drep-project/drep-chain/chain/store"
	"github.com/drep-project/drep-chain/common"
	"github.com/drep-project/drep-chain/crypto"
	"github.com/drep-project/drep-chain/params"
	"github.com/drep-project/drep-chain/types"
	"math/big"
	"time"
)

func (blockMgr *BlockMgr) GenerateTemplate(trieStore store.StoreInterface, leaderAddr crypto.CommonAddress) (*types.Block, *big.Int, error) {
	parent, err := blockMgr.ChainService.GetHighestBlock()
	if err != nil {
		return nil, nil, err
	}
	newGasLimit := blockMgr.ChainService.CalcGasLimit(parent.Header, params.MinGasLimit, params.MaxGasLimit)
	height := blockMgr.ChainService.BestChain().Height() + 1
	txs := blockMgr.transactionPool.GetPending(newGasLimit)
	previousHash := blockMgr.ChainService.BestChain().Tip().Hash
	timestamp := uint64(time.Now().Unix())

	blockHeader := &types.BlockHeader{
		Version:      common.Version,
		PreviousHash: *previousHash,
		ChainId:      blockMgr.ChainService.ChainID(),
		GasLimit:     *newGasLimit,
		Timestamp:    timestamp,
		Height:       height,
		StateRoot:    []byte{},
		TxRoot:       []byte{},
	}

	block := &types.Block{
		Header: blockHeader,
		Data: &types.BlockData{
			TxCount: uint64(len(txs)),
			TxList:  txs,
		},
	}

	gp := new(chain.GasPool).AddGas(newGasLimit.Uint64())
	//process transaction
	chainStore := &chain.ChainStore{blockMgr.DatabaseService.LevelDb()}
	context := chain.NewBlockExecuteContext(trieStore, gp, chainStore, block)

	templateValidator := NewTemplateBlockValidator(blockMgr.ChainService)
	err = templateValidator.ExecuteBlock(context)
	if err != nil {
		return nil, nil, err
	}
	return context.Block, context.GasFee, nil
}
