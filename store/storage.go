package store

import (
    "math/big"
    "sync"
    "BlockChainTest/bean"
    "log"
    "fmt"
    "BlockChainTest/database"
)

var (
    balances           = make(map[bean.Address]*big.Int)//map[bean.Address]*big.Int
    nonces             map[bean.Address]int64
    //blocks             []*bean.Block
    accountLock        sync.Mutex
    //currentBlockHeight int64 = -1
)

func init()  {
    //balances = make(map[bean.Address]*big.Int)
    nonces = make(map[bean.Address]int64)
    //blocks = make([]*bean.Block, 0)
}

func GetBalance(addr bean.Address) *big.Int {
    accountLock.Lock()
    defer accountLock.Unlock()
    balance, exists := balances[addr]
    if exists {
        // TODO if map is nil what the fuck
        return balance
    } else {
        balance = big.NewInt(0)
        balances[addr] = balance
        return balance
    }

}

func GetNonce(addr bean.Address) int64 {
    return nonces[addr]
}

func addNonce(addr bean.Address) {
    accountLock.Lock()
    value, exists := nonces[addr]
    if exists {
        if value >= 0 {
            nonces[addr]++
        } else {
            nonces[addr] = 1
        }
    } else {
        nonces[addr] = 1
    }
    accountLock.Unlock()
}

func ExecuteTransactions(b *bean.Block) *big.Int {
    if b == nil || b.Header == nil { // || b.Data == nil || b.Data.TxList == nil {
        fmt.Errorf("error block nil or header nil")
        return nil
    }
    height := database.GetMaxHeight()
    if height + 1 != b.Header.Height {
        fmt.Println("error", height, b.Header.Height)
        return nil
    }
    // TODO check height
    height = b.Header.Height
    database.PutMaxHeight(height)
    //blocks = append(blocks, b)
    database.PutBlock(b)
    total := big.NewInt(0)
    if b.Data == nil || b.Data.TxList == nil {
        return total
    }
    for _, t := range b.Data.TxList {
        gasFee := execute(t)
        fmt.Println("Delete transaction ", *t)
        fmt.Println(removeTransaction(t))
        if gasFee != nil {
            total.Add(total, gasFee)
        }
    }
    return total
}

func execute(t *bean.Transaction) *big.Int {
    addr := bean.Hex2Address(t.Addr().String())
    nonce := t.Data.Nonce
    curN := database.GetNonce(addr)
    if curN + 1 != nonce {
        return nil
    }
    database.PutNonce(addr, curN + 1)
    gasPrice := big.NewInt(0).SetBytes(t.Data.GasPrice)
    gasLimit := big.NewInt(0).SetBytes(t.Data.GasLimit)
    gasFee := big.NewInt(0).Mul(gasLimit, gasPrice)
    balance := database.GetBalance(addr)
    if gasFee.Cmp(balance) > 0 {
        log.Fatal("Error, gas not right")
        return nil
    }
    switch t.Data.Type {
    case TransferType:
        {
            if gasLimit.Cmp(TransferGas) < 0 {
                balance.Sub(balance, gasFee)
            } else {
                amount := big.NewInt(0).SetBytes(t.Data.Amount)
                total := big.NewInt(0).Add(amount, gasFee)
                if balance.Cmp(total) >= 0 {
                    balance.Sub(balance, total)
                    to := bean.Address(t.Data.To)
                    balance2 := database.GetBalance(bean.Hex2Address(to.String()))
                    balance2.Add(balance2, amount)
                } else {
                    balance.Sub(balance, gasFee)
                }
            }
        }
    case MinerType:
        {
            // TODO if not the admin
            if gasLimit.Cmp(MinerGas) < 0 {
                balance.Sub(balance, gasFee)
            } else {
                balance.Sub(balance, gasFee)
                AddMiner(bean.Address(t.Data.Data))
            }
        }
    }
    return gasFee
}

func GetCurrentBlockHeight() int64 {
    if height := database.GetMaxHeight(); height != -1 {
        return int64(height)
    } else {
        fmt.Println("ERROR!!! height is -1")
        return -1
    }
}

func GetBlocks(from int64, number int64) []*bean.Block {
    bs := database.GetAllBlocks()
    l := int64(len(bs))
    if l - 1 < from {
        return []*bean.Block{}
    }
    end := from + number - 1
    r := make([]*bean.Block, 0)
    for i := from; i < l && i <= end; i++ {
        r = append(r, bs[i])
    }
    return r
}