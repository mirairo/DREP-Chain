package service

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/drep-project/DREP-Chain/chain/store"
	"github.com/drep-project/DREP-Chain/params"
	"github.com/drep-project/DREP-Chain/pkgs/evm/vm"
	"math/big"

	"github.com/drep-project/DREP-Chain/blockmgr"
	"github.com/drep-project/DREP-Chain/common"
	"github.com/drep-project/DREP-Chain/crypto"
	"github.com/drep-project/DREP-Chain/crypto/secp256k1"
	"github.com/drep-project/DREP-Chain/database"
	"github.com/drep-project/DREP-Chain/pkgs/accounts/addrgenerator"
	"github.com/drep-project/DREP-Chain/pkgs/evm"
	"github.com/drep-project/DREP-Chain/types"
)

/*
name: Account RPC interface
usage: Address management and initiate simple transactions
prefix:account
*/
type AccountApi struct {
	EvmService         *evm.EvmService
	Wallet             *Wallet
	accountService     *AccountService
	poolQuery          blockmgr.IBlockMgrPool
	messageBroadCastor blockmgr.ISendMessage
	databaseService    *database.DatabaseService
}

/*
 name: listAddress
 usage: Lists all local addresses
 return: Address of the array
 example:   curl http://localhost:10085 -X POST --data '{"jsonrpc":"2.0","method":"account_listAddress","params":[], "id": 3}' -H "Content-Type:application/json"
 response:
  {"jsonrpc":"2.0","id":3,"result":["0x3296d3336895b5baaa0eca3df911741bd0681c3f","0x3ebcbe7cb440dd8c52940a2963472380afbb56c5"]}
*/
func (accountapi *AccountApi) ListAddress() ([]string, error) {
	if !accountapi.Wallet.IsOpen() {
		return nil, ErrClosedWallet
	}
	return accountapi.Wallet.ListAddress()
}

/*
 name: createAccount
 usage: Create a local account
 params:
	1. password
 return: New account address information
 example:   curl http://localhost:10085 -X POST --data '{"jsonrpc":"2.0","method":"account_createAccount","params":["123456"], "id": 3}' -H "Content-Type:application/json"
 response:
	  {"jsonrpc":"2.0","id":3,"result":"0x2944c15c466fad03ec1282bab579dec5a0cf0fa3"}
*/
func (accountapi *AccountApi) CreateAccount(password string) (*crypto.CommonAddress, error) {
	if !accountapi.Wallet.IsOpen() {
		return nil, ErrClosedWallet
	}
	newAaccount, err := accountapi.Wallet.NewAccount(password)
	if err != nil {
		return nil, err
	}
	return newAaccount.Address, nil
}

/*
 name: createWallet
 usage: Create a local wallet
 params:
	1. The wallet password
 return:  Failure returns the reason for the error, and success returns no information
 example:   curl http://localhost:10085 -X POST --data '{"jsonrpc":"2.0","method":"account_createWallet","params":["123"], "id": 3}' -H "Content-Type:application/json"
 response:
	  {"jsonrpc":"2.0","id":3,"result":null}
*/
func (accountapi *AccountApi) CreateWallet(password string) error {
	err := accountapi.accountService.CreateWallet(password)
	if err != nil {
		return err
	}
	return accountapi.OpenWallet(password)
}

/*
 name: lockAccount
 usage: Lock the account
 params:
 return:  Failure returns the reason for the error, and success returns no information
 example:   curl http://localhost:10085 -X POST --data '{"jsonrpc":"2.0","method":"account_lockAccount","params":["0x518b3fefa3fb9a72753c6ad10a2b68cc034ec391"], "id": 3}' -H "Content-Type:application/json"
 response:
	 {"jsonrpc":"2.0","id":3,"result":null}
*/
func (accountapi *AccountApi) LockAccount(addr crypto.CommonAddress) error {
	if !accountapi.Wallet.IsOpen() {
		return ErrClosedWallet
	}

	return accountapi.Wallet.Lock(&addr)

	return ErrLockedWallet
}

/*
 name: account_unlockAccount
 usage: Unlock the account
 params:
	1. The account address
	2. password
 return: Failure returns the reason for the error, and success returns no information
 example:   curl http://localhost:10085 -X POST --data '{"jsonrpc":"2.0","method":"account_unlockAccount","params":["0x518b3fefa3fb9a72753c6ad10a2b68cc034ec391", "123456"], "id": 3}' -H "Content-Type:application/json"
 response:
	 {"jsonrpc":"2.0","id":3,"result":null}
*/
func (accountapi *AccountApi) UnlockAccount(addr crypto.CommonAddress, password string) error {
	if !accountapi.Wallet.IsOpen() {
		return ErrClosedWallet
	}

	return accountapi.Wallet.UnLock(&addr, password)
}

/*
 name: openWallet
 usage: Open my wallet
 params:
	1. The wallet password
 return: error or none
 example:   curl http://localhost:10085 -X POST --data '{"jsonrpc":"2.0","method":"account_openWallet","params":["123"], "id": 3}' -H "Content-Type:application/json"
 response:
	 {"jsonrpc":"2.0","id":3,"result":null}
*/
func (accountapi *AccountApi) OpenWallet(password string) error {
	return accountapi.Wallet.OpenWallet(password)
}

/*
 name: closeWallet
 usage: close wallet
 params:
 return: none
 example:   curl http://localhost:10085 -X POST --data '{"jsonrpc":"2.0","method":"account_closeWallet","params":[], "id": 3}' -H "Content-Type:application/json"
 response:
	 {"jsonrpc":"2.0","id":3,"result":null}
*/
func (accountapi *AccountApi) CloseWallet() {
	accountapi.Wallet.Close()
}

/*
 name: transfer
 usage: transfer
 params:
	1. The address at which the transfer was initiated
	2. Recipient's address
	3. Mount
	4. gas price
	5. gas limit
	6. commit
 return: transaction hash
 example:   curl -H "Content-Type: application/json" -X post --data '{"jsonrpc":"2.0","method":"account_transfer","params":["0x3ebcbe7cb440dd8c52940a2963472380afbb56c5","0x3ebcbe7cb440dd8c52940a2963472380afbb56c5","0x111","0x110","0x30000",""],"id":1}' http://127.0.0.1:10085
 response:
	 {"jsonrpc":"2.0","id":1,"result":"0x3a3b59f90a21c2fd1b690aa3a2bc06dc2d40eb5bdc26fdd7ecb7e1105af2638e"}
*/
func (accountapi *AccountApi) Transfer(from crypto.CommonAddress, to crypto.CommonAddress, amount, gasprice, gaslimit *common.Big, data common.Bytes) (string, error) {
	if gasprice.ToInt().Uint64() < blockmgr.DefaultGasPrice {
		gasprice.SetMathBig(*new(big.Int).SetUint64(blockmgr.DefaultGasPrice))
	}

	nonce := accountapi.poolQuery.GetTransactionCount(&from)
	tx := types.NewTransaction(to, (*big.Int)(amount), (*big.Int)(gasprice), (*big.Int)(gaslimit), nonce)
	sig, err := accountapi.Wallet.Sign(&from, tx.TxHash().Bytes())
	if err != nil {
		return "", err
	}
	tx.Sig = sig
	err = accountapi.messageBroadCastor.SendTransaction(tx, true)
	if err != nil {
		return "", err
	}
	return tx.TxHash().String(), nil
}

/*
 name: transferWithNonce
 usage: transfer with nonce
 params:
	1. The address at which the transfer was initiated
	2. Recipient's address
	3. Mount
	4. gas price
	5. gas limit
	6. commit
    7. nonce
 return: transaction hash
 example:   curl -H "Content-Type: application/json" -X post --data '{"jsonrpc":"2.0","method":"account_transferWithNonce","params":["0x3ebcbe7cb440dd8c52940a2963472380afbb56c5","0x3ebcbe7cb440dd8c52940a2963472380afbb56c5","0x111","0x110","0x30000","",1],"id":1}' http://127.0.0.1:10085
 response:
	 {"jsonrpc":"2.0","id":1,"result":"0x3a3b59f90a21c2fd1b690aa3a2bc06dc2d40eb5bdc26fdd7ecb7e1105af2638e"}
*/
func (accountapi *AccountApi) TransferWithNonce(from crypto.CommonAddress, to crypto.CommonAddress, amount, gasprice, gaslimit *common.Big, data common.Bytes, nonce uint64) (string, error) {
	//nonce := accountapi.poolQuery.GetTransactionCount(&from)
	tx := types.NewTransaction(to, (*big.Int)(amount), (*big.Int)(gasprice), (*big.Int)(gaslimit), nonce)
	sig, err := accountapi.Wallet.Sign(&from, tx.TxHash().Bytes())
	if err != nil {
		return "", err
	}
	tx.Sig = sig
	err = accountapi.messageBroadCastor.SendTransaction(tx, true)
	if err != nil {
		return "", err
	}
	return tx.TxHash().String(), nil
}

/*
 name: setAlias
 usage: Set an alias
 params:
	1. address
	2. alias
	3. gas price
	4. gas lowLimit
 return: transaction hash
 example:
	curl -H "Content-Type: application/json" -X post --data '{"jsonrpc":"2.0","method":"account_setAlias","params":["0x3ebcbe7cb440dd8c52940a2963472380afbb56c5","AAAAA","0x110","0x30000"],"id":1}' http://127.0.0.1:10085
response:
	{"jsonrpc":"2.0","id":1,"result":"0x5adb248f2943e12fb91c140bd3d0df6237712061e9abae97345b0869c3daa749"}
*/
func (accountapi *AccountApi) SetAlias(srcAddr crypto.CommonAddress, alias string, gasprice, gaslimit *common.Big) (string, error) {
	nonce := accountapi.poolQuery.GetTransactionCount(&srcAddr)
	t := types.NewAliasTransaction(alias, (*big.Int)(gasprice), (*big.Int)(gaslimit), nonce)
	sig, err := accountapi.Wallet.Sign(&srcAddr, t.TxHash().Bytes())
	if err != nil {
		return "", err
	}
	t.Sig = sig
	fmt.Println(hex.EncodeToString(t.AsPersistentMessage()))
	fmt.Println(t.TxHash().String())
	err = accountapi.messageBroadCastor.SendTransaction(t, true)
	if err != nil {
		return "", err
	}
	return t.TxHash().String(), nil
}

/*
 name: VoteCredit
 usage: vote credit to candidate
 params:
	1. address of voter
	2. address of candidate
	3. amount
	4. gas price
	5. gas uplimit of transaction
 return: transaction hash
 example:   curl -H "Content-Type: application/json" -X post --data '{"jsonrpc":"2.0","method":"account_voteCredit","params":["0x3ebcbe7cb440dd8c52940a2963472380afbb56c5","0x3ebcbe7cb440dd8c52940a2963472380afbb56c5","0x111","0x110","0x30000"],"id":1}' http://127.0.0.1:10085
 response:
	 {"jsonrpc":"2.0","id":1,"result":"0x3a3b59f90a21c2fd1b690aa3a2bc06dc2d40eb5bdc26fdd7ecb7e1105af2638e"}
*/
func (accountapi *AccountApi) VoteCredit(from crypto.CommonAddress, to crypto.CommonAddress, amount, gasprice, gaslimit *common.Big) (string, error) {
	nonce := accountapi.poolQuery.GetTransactionCount(&from)
	tx := types.NewVoteTransaction(to, (*big.Int)(amount), (*big.Int)(gasprice), (*big.Int)(gaslimit), nonce)
	sig, err := accountapi.Wallet.Sign(&from, tx.TxHash().Bytes())
	if err != nil {
		return "", err
	}
	tx.Sig = sig
	err = accountapi.messageBroadCastor.SendTransaction(tx, true)
	if err != nil {
		return "", err
	}
	return tx.TxHash().String(), nil
}

/*
 name: CancelVoteCredit
 usage:
 params:
	1. address of voter
	2. address of candidate
	3. amount
	4. gas price
	5. gas limit
	6. 备注
 return: transaction hash
 example:   curl -H "Content-Type: application/json" -X post --data '{"jsonrpc":"2.0","method":"account_cancelVoteCredit","params":["0x3ebcbe7cb440dd8c52940a2963472380afbb56c5","0x3ebcbe7cb440dd8c52940a2963472380afbb56c5","0x111","0x110","0x30000"],"id":1}' http://127.0.0.1:10085
 response:
	 {"jsonrpc":"2.0","id":1,"result":"0x3a3b59f90a21c2fd1b690aa3a2bc06dc2d40eb5bdc26fdd7ecb7e1105af2638e"}
*/
func (accountapi *AccountApi) CancelVoteCredit(from crypto.CommonAddress, to crypto.CommonAddress, amount, gasprice, gaslimit *common.Big) (string, error) {
	nonce := accountapi.poolQuery.GetTransactionCount(&from)
	tx := types.NewCancelVoteTransaction(to, (*big.Int)(amount), (*big.Int)(gasprice), (*big.Int)(gaslimit), nonce)
	sig, err := accountapi.Wallet.Sign(&from, tx.TxHash().Bytes())
	if err != nil {
		return "", err
	}
	tx.Sig = sig
	err = accountapi.messageBroadCastor.SendTransaction(tx, true)
	if err != nil {
		return "", err
	}
	return tx.TxHash().String(), nil
}

/*
 name: CandidateCredit
 usage: Candidate node pledge
 params:
	1. The address of the pledger
	2. The pledge amount
	3. gas price
	4. gas limit
	5. The pubkey corresponding to the address of the pledger, and the P2p information of the pledger
 return: transaction hash
 example:   curl -H "Content-Type: application/json" -X post --data '{"jsonrpc":"2.0","method":"account_candidateCredit","params":["0x3ebcbe7cb440dd8c52940a2963472380afbb56c5","0x111","0x110","0x30000","{\"Pubkey\":\"0x020e233ebaed5ade5e48d7ee7a999e173df054321f4ddaebecdb61756f8a43e91c\",\"Node\":\"enode://3f05da2475bf09ce20b790d76b42450996bc1d3c113a1848be1960171f9851c0@149.129.172.91:44444\"}"],"id":1}' http://127.0.0.1:10085
 response:
	 {"jsonrpc":"2.0","id":1,"result":"0x3a3b59f90a21c2fd1b690aa3a2bc06dc2d40eb5bdc26fdd7ecb7e1105af2638e"}
*/
func (accountapi *AccountApi) CandidateCredit(from crypto.CommonAddress, amount, gasprice, gaslimit *common.Big, data string) (string, error) {
	cd := types.CandidateData{}
	err := cd.Unmarshal([]byte(data))
	if err != nil {
		return "", err
	}

	if !bytes.Equal(crypto.PubkeyToAddress(cd.Pubkey).Bytes(), from.Bytes()) {
		return "", nil
	}

	nonce := accountapi.poolQuery.GetTransactionCount(&from)
	tx := types.NewCandidateTransaction((*big.Int)(amount), (*big.Int)(gasprice), (*big.Int)(gaslimit), nonce, []byte(data))
	sig, err := accountapi.Wallet.Sign(&from, tx.TxHash().Bytes())
	if err != nil {
		return "", err
	}
	tx.Sig = sig
	err = accountapi.messageBroadCastor.SendTransaction(tx, true)
	if err != nil {
		return "", err
	}
	return tx.TxHash().String(), nil
}

/*
 name: CancelCandidateCredit
 usage: To cancel the candidate
 params:
	1. The address at which the transfer was cancel
	2. address of candidate
	3. amount
	4. gas price
	5. gas limit

 return: transaction hash
 example:   curl -H "Content-Type: application/json" -X post --data '{"jsonrpc":"2.0","method":"account_cancelCandidateCredit","params":["0x3ebcbe7cb440dd8c52940a2963472380afbb56c5","0x111","0x110","0x30000",""],"id":1}' http://127.0.0.1:10085
 response:
	 {"jsonrpc":"2.0","id":1,"result":"0x3a3b59f90a21c2fd1b690aa3a2bc06dc2d40eb5bdc26fdd7ecb7e1105af2638e"}
*/
func (accountapi *AccountApi) CancelCandidateCredit(from crypto.CommonAddress, amount, gasprice, gaslimit *common.Big) (string, error) {
	nonce := accountapi.poolQuery.GetTransactionCount(&from)
	tx := types.NewCancleCandidateTransaction((*big.Int)(amount), (*big.Int)(gasprice), (*big.Int)(gaslimit), nonce)
	sig, err := accountapi.Wallet.Sign(&from, tx.TxHash().Bytes())
	if err != nil {
		return "", err
	}
	tx.Sig = sig
	err = accountapi.messageBroadCastor.SendTransaction(tx, true)
	if err != nil {
		return "", err
	}
	return tx.TxHash().String(), nil
}

/*
 name: readContract
 usage: Read smart contract (no data modified)
 params:
    1. The account address of the transaction
	2. Contract address
	3. Contract api
 return: The query results
 example:
	curl -H "Content-Type: application/json" -X post --data '{"jsonrpc":"2.0","method":"account_readContract","params":["0xec61c03f719a5c214f60719c3f36bb362a202125","0xecfb51e10aa4c146bf6c12eee090339c99841efc","0x6d4ce63c"],"id":1}' http://127.0.0.1:10085
 response:
	 {"jsonrpc":"2.0","id":1,"result":""}
*/
func (accountapi *AccountApi) ReadContract(from, to crypto.CommonAddress, input common.Bytes) (common.Bytes, error) {
	header := accountapi.EvmService.Chain.GetCurrentHeader()
	tx := types.NewTransaction(to, new(big.Int).SetUint64(0), &big.Int{}, new(big.Int).SetUint64(params.MinGasLimit), 0)
	tx.Data.Data = input

	sig, err := accountapi.Wallet.Sign(&from, tx.TxHash().Bytes())
	if err != nil {
		return nil, err
	}
	tx.Sig = sig

	trieStore, err := store.TrieStoreFromStore(accountapi.databaseService.LevelDb(), header.StateRoot)
	if err != nil {
		return nil, err
	}

	ret, err := accountapi.EvmService.Call(trieStore, tx, header)
	fmt.Println(string(common.Bytes(ret)))
	fmt.Println(new(big.Int).SetBytes(ret))
	fmt.Println(common.Bytes(ret))

	return common.Bytes(ret), err
}

/*
 name: estimateGas
 usage: Estimate how much gas is needed for the transaction
 params:
	1. The address at which the transfer was initiated
	2. amount
	3. commit
	4. Address of recipient
 return: Evaluate the result, failure returns an error
 example:
	curl -H "Content-Type: application/json" -X post --data '{"jsonrpc":"2.0","method":"account_estimateGas","params":["0xec61c03f719a5c214f60719c3f36bb362a202125","0xecfb51e10aa4c146bf6c12eee090339c99841efc","0x6d4ce63c","0x110","0x30000"],"id":1}' http://127.0.0.1:10085
 response:
	 {"jsonrpc":"2.0","id":1,"result":"0x5d74aba54ace5f01a5f0057f37bfddbbe646ea6de7265b368e2e7d17d9cdeb9c"}
*/
func (accountapi *AccountApi) EstimateGas(from crypto.CommonAddress, amount *common.Big, data common.Bytes, to *crypto.CommonAddress) (uint64, error) {
	if amount.ToInt().Uint64() != 0 {
		return params.MinGasLimit, nil
	}

	header := accountapi.EvmService.Chain.GetCurrentHeader()
	tx := types.NewTransaction(*to, amount.ToInt(), new(big.Int).SetUint64(blockmgr.DefaultGasPrice), new(big.Int).SetUint64(params.MinGasLimit), 0)
	tx.Data.Data = data

	sig, err := accountapi.Wallet.Sign(&from, tx.TxHash().Bytes())
	if err != nil {
		return 0, err
	}
	tx.Sig = sig

	trieStore, err := store.TrieStoreFromStore(accountapi.databaseService.LevelDb(), header.StateRoot)
	if err != nil {
		return 0, err
	}

	state := vm.NewState(trieStore, header.Height)

	gl := new(big.Int).SetUint64(params.MinGasLimit)
	var (
		fail bool
	)

	for {
		_, _, _, fail, err = accountapi.EvmService.Eval(state, tx, header, gl.Uint64(), amount.ToInt())
		if err != nil || fail {
			if err == vm.ErrCodeStoreOutOfGas || err == vm.ErrOutOfGas {
				gl = gl.Add(gl, new(big.Int).SetUint64(1))
			} else {
				return 0, fmt.Errorf("err:%v or fail:%v", err, fail)
			}
		} else {
			tx.Data.GasLimit = *(*common.Big)(gl)
			return tx.Data.GasLimit.ToInt().Uint64(), err
		}
	}

	return 0, err
}

/*
 name: executeContract
 usage: Execute smart contract (cause data to be modified)
 params:
	1. The address of the caller
	2. Contract address
	3. Contract code
	3. gas price
	4. gas limit
 return: transaction hash
 example:
	curl -H "Content-Type: application/json" -X post --data '{"jsonrpc":"2.0","method":"account_executeContract","params":["0xec61c03f719a5c214f60719c3f36bb362a202125","0xecfb51e10aa4c146bf6c12eee090339c99841efc","0x6d4ce63c","0x110","0x30000"],"id":1}' http://127.0.0.1:10085
 response:
	 {"jsonrpc":"2.0","id":1,"result":"0x5d74aba54ace5f01a5f0057f37bfddbbe646ea6de7265b368e2e7d17d9cdeb9c"}
*/
func (accountapi *AccountApi) ExecuteContract(from crypto.CommonAddress, to crypto.CommonAddress, input common.Bytes, gasprice, gaslimit *common.Big) (string, error) {
	nonce := accountapi.poolQuery.GetTransactionCount(&from)
	t := types.NewCallContractTransaction(to, input, &big.Int{}, (*big.Int)(gasprice), (*big.Int)(gaslimit), nonce)
	sig, err := accountapi.Wallet.Sign(&from, t.TxHash().Bytes())
	if err != nil {
		return "", err
	}
	t.Sig = sig
	accountapi.messageBroadCastor.SendTransaction(t, true)
	return t.TxHash().String(), nil
}

/*
 name: createCode
 usage: Deployment of contract
 params:
	1. The account address of the deployment contract
	2. Content of the contract
	3. gas price
	4. gas limit
 return: transaction hash
 example:
 	curl -H "Content-Type: application/json" -X post --data '{"jsonrpc":"2.0","method":"account_createCode","params":["0x3ebcbe7cb440dd8c52940a2963472380afbb56c5","0x608060405234801561001057600080fd5b5061018c806100206000396000f3fe608060405260043610610051576000357c0100000000000000000000000000000000000000000000000000000000900480634f2be91f146100565780636d4ce63c1461006d578063db7208e31461009e575b600080fd5b34801561006257600080fd5b5061006b6100dc565b005b34801561007957600080fd5b5061008261011c565b604051808260070b60070b815260200191505060405180910390f35b3480156100aa57600080fd5b506100da600480360360208110156100c157600080fd5b81019080803560070b9060200190929190505050610132565b005b60016000808282829054906101000a900460070b0192506101000a81548167ffffffffffffffff021916908360070b67ffffffffffffffff160217905550565b60008060009054906101000a900460070b905090565b806000806101000a81548167ffffffffffffffff021916908360070b67ffffffffffffffff1602179055505056fea165627a7a723058204b651e4313ab6bc4eda61084cac1f805699cefbb979ddfd3a2d7f970903307cd0029","0x111","0x110","0x30000"],"id":1}' http://127.0.0.1:10085
 response:
	 {"jsonrpc":"2.0","id":1,"result":"0x9a8d8d5d7d00bbe0eb1b9431a13a7219008e352241b751b177bfb29e4e75b0d1"}
*/
func (accountapi *AccountApi) CreateCode(from crypto.CommonAddress, byteCode common.Bytes, gasprice, gaslimit *common.Big) (string, error) {
	nonce := accountapi.poolQuery.GetTransactionCount(&from)
	t := types.NewContractTransaction(byteCode, (*big.Int)(gasprice), (*big.Int)(gaslimit), nonce)
	sig, err := accountapi.Wallet.Sign(&from, t.TxHash().Bytes())
	if err != nil {
		return "", err
	}
	t.Sig = sig
	err = accountapi.messageBroadCastor.SendTransaction(t, true)
	if err != nil {
		return "", err
	}
	return t.TxHash().String(), nil
}

/*
 name: dumpPrivkey
 usage: The private key corresponding to the export address
 params:
	1.address
 return: private key
 example:   curl http://localhost:10085 -X POST --data '{"jsonrpc":"2.0","method":"account_dumpPrivkey","params":["0x3ebcbe7cb440dd8c52940a2963472380afbb56c5"], "id": 3}' -H "Content-Type:application/json"
 response:
	 {"jsonrpc":"2.0","id":3,"result":"0x270f4b122603999d1c07aec97e972a2ddf7bd8b5bfe3543c10814e6a19f13aaf"}
*/
func (accountapi *AccountApi) DumpPrivkey(address *crypto.CommonAddress) (*secp256k1.PrivateKey, error) {
	if !accountapi.Wallet.IsOpen() {
		return nil, ErrClosedWallet
	}
	if accountapi.Wallet.IsLock() {
		return nil, ErrLockedWallet
	}

	node, err := accountapi.Wallet.GetAccountByAddress(address)
	if err != nil {
		return nil, err
	}
	return node.PrivateKey, nil
}

/*
 name: DumpPubkey
 usage: Export the public key corresponding to the address
 params:
	1.address
 return: public key
 example:   curl http://localhost:10085 -X POST --data '{"jsonrpc":"2.0","method":"account_dumpPubkey","params":["0x3ebcbe7cb440dd8c52940a2963472380afbb56c5"], "id": 3}' -H "Content-Type:application/json"
 response:
	 {"jsonrpc":"2.0","id":3,"result":"0x270f4b122603999d1c07aec97e972a2ddf7bd8b5bfe3543c10814e6a19f13aaf"}
*/
func (accountapi *AccountApi) DumpPubkey(address *crypto.CommonAddress) (*secp256k1.PublicKey, error) {
	if !accountapi.Wallet.IsOpen() {
		return nil, ErrClosedWallet
	}
	if accountapi.Wallet.IsLock() {
		return nil, ErrLockedWallet
	}

	node, err := accountapi.Wallet.GetAccountByAddress(address)
	if err != nil {
		return nil, err
	}
	return node.PrivateKey.PubKey(), nil
}

/*
 name: sign
 usage: Signature transaction
 params:
	1.account of sig
	2.msg for sig
 return: private key
 example:
	curl http://localhost:10085 -X POST --data '{"jsonrpc":"2.0","method":"account_sign","params":["0x3ebcbe7cb440dd8c52940a2963472380afbb56c5", "0x00001c9b8c8fdb1f53faf02321f76253704123e2b56cce065852bab93e526ae2"], "id": 3}' -H "Content-Type:application/json"

response:
	 {"jsonrpc":"2.0","id":3,"result":"0x1f1d16412468dd9b67b568d31839ac608bdfddf2580666db4d364eefbe285fdaed569a3c8fa1decfebbfa0ed18b636059dbbf4c2106c45fc8846909833ef2cb1de"}
*/
func (accountapi *AccountApi) Sign(address crypto.CommonAddress, hash common.Bytes) (common.Bytes, error) {
	sig, err := accountapi.Wallet.Sign(&address, hash)
	if err != nil {
		return nil, err
	}
	return sig, nil
}

/*
 name: generateAddresses
 usage: Generate the addresses of the other chains
 params:
	1. address of drep
 return: {BTCaddress, ethAddress, neoAddress}
 example:
	curl http://localhost:10085 -X POST --data '{"jsonrpc":"2.0","method":"account_generateAddresses","params":["0x3ebcbe7cb440dd8c52940a2963472380afbb56c5"], "id": 3}' -H "Content-Type:application/json"

response:
	 {"jsonrpc":"2.0","id":3,"result":""}
*/
func (accountapi *AccountApi) GenerateAddresses(address crypto.CommonAddress) (*RpcAddresses, error) {
	privkey, err := accountapi.Wallet.DumpPrivateKey(&address)
	if err != nil {
		return nil, err
	}
	generator := &addrgenerator.AddrGenerate{
		PrivateKey: privkey,
	}
	return &RpcAddresses{
		BtcAddress:      generator.ToBtc(),
		EthAddress:      generator.ToEth(),
		NeoAddress:      generator.ToNeo(),
		RippleAddress:   generator.ToRipple(),
		DashAddress:     generator.ToDash(),
		DogeCoinAddress: generator.ToDogecoin(),
		LiteCoinAddress: generator.ToLiteCoin(),
		CosmosAddress:   generator.ToAtom(),
		TronAddress:     generator.ToTron(),
	}, nil
}

/*
 name: importKeyStore
 usage: import keystore
 params:
	1.path
	2.password
 return: address list
 example:
	 curl http://localhost:10085 -X POST --data '{"jsonrpc":"2.0","method":"account_importKeyStore","params":["path","123"], "id": 3}' -H "Content-Type:application/json"
response:
	 {"jsonrpc":"2.0","id":3,"result":["0x4082c96e38def8f3851831940485066234fe07b8"]}
*/
func (accountapi *AccountApi) ImportKeyStore(path, password string) ([]*crypto.CommonAddress, error) {
	return accountapi.Wallet.ImportKeyStore(path, password)
}

/*
 name: importPrivkey
 usage: import private key
 params:
	1.privkey(compress hex)
	2.password
 return: address
 example:
	curl http://localhost:10085 -X POST --data '{"jsonrpc":"2.0","method":"account_importPrivkey","params":["0xe5510b32854ca52e7d7d41bb3196fd426d551951e2fd5f6b559a62889d87926c"], "id": 3}' -H "Content-Type:application/json"
response:
	 {"jsonrpc":"2.0","id":3,"result":"0x748eb65493a964e568800c3c2885c63a0de9f9ae"}
*/
func (accountapi *AccountApi) ImportPrivkey(privBytes common.Bytes, password string) (*crypto.CommonAddress, error) {
	priv, _ := secp256k1.PrivKeyFromScalar(privBytes)
	node, err := accountapi.Wallet.ImportPrivKey(priv, password)
	if err != nil {
		return nil, err
	}
	return node.Address, nil
}

/*
 name: getKeyStores
 usage: get ketStores path
 params:

 return: path of keystore
 example:
	curl http://localhost:10085 -X POST --data '{"jsonrpc":"2.0","method":"account_getKeyStores","params":[], "id": 3}' -H "Content-Type:application/json"
response:
	 {"jsonrpc":"2.0","id":3,"result":"'path of keystores is: C:\\Users\\Kun\\AppData\\Local\\Drep\\keystore'"}
*/
func (accountapi *AccountApi) GetKeyStores() string {
	return "path of keystores is: " + accountapi.Wallet.config.KeyStoreDir
}

type RpcAddresses struct {
	BtcAddress      string
	EthAddress      string
	NeoAddress      string
	RippleAddress   string
	DashAddress     string
	DogeCoinAddress string
	LiteCoinAddress string
	CosmosAddress   string
	TronAddress     string
}

type RpcAccount struct {
	Addr   *crypto.CommonAddress
	Pubkey string
}
