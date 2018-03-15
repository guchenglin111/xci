package ipfs

import (
	"github.com/xcareteam/xci/common"
	"github.com/xcareteam/xci/contracts/ipfs/contract"
	"github.com/xcareteam/xci/accounts/abi/bind"
	"github.com/xcareteam/xci/core/types"
	"github.com/xcareteam/xci/accounts"
)

//go:generate abigen --sol contract/ipfs.sol --pkg contract --out contract/ipfs.go

var (
	MainNetAddress = common.HexToAddress("0x314159265dD8dbb310642f98f50C066173C1259b")
	TestNetAddress = common.HexToAddress("0x112234455c3a32fd11230c42e7bccd4a84e02010")
)

type IPFS struct {
	*contract.IpfsSession
	contractBackend bind.ContractBackend
}

func NewIpfs(transactOpts *bind.TransactOpts, callOpts *bind.CallOpts, contractAddr common.Address, contractBackend bind.ContractBackend) (*IPFS, error) {
	ipfs, err := contract.NewIpfs(contractAddr, contractBackend)
	if err != nil {
		return nil, err
	}

	return &IPFS{
		&contract.IpfsSession{
			Contract:     ipfs,
			CallOpts:	  *callOpts,
			TransactOpts: *transactOpts,
		},
		contractBackend,
	}, nil
}

func DeployIPFS(transactOpts *bind.TransactOpts, callOpts *bind.CallOpts, contractBackend bind.ContractBackend) (common.Address, *IPFS, error) {
	ipfsAddr, _, _, err := contract.DeployIpfs(transactOpts, contractBackend)
	if err != nil {
		return ipfsAddr, nil, err
	}

	ipfs, err := NewIpfs(transactOpts, callOpts, ipfsAddr, contractBackend)
	if err != nil {
		return ipfsAddr, nil, err
	}

	return ipfsAddr, ipfs, nil
}

func (self *IPFS) GetIpfsUrl(date,fileName string) (string, error) {
	return self.Contract.GetIpfsUrl(&self.CallOpts, date, fileName);
}
/*
func (self *IPFS) GetFileListByDate(date string) (string, error) {
	return self.Contract.GetFileListByDate(&self.CallOpts, date);
}*/

func (self *IPFS) GetOwner() (common.Address, error) {
	return self.Contract.GetOwner(&self.CallOpts);
}

func (self *IPFS) AddNewIpfsUrl(date,fileName, url string) (*types.Transaction, error) {
	return self.Contract.AddNewIpfsUrl(&self.TransactOpts, date,fileName, url);
}

func GetNewIPFS(accMng *accounts.Manager, backend bind.ContractBackend, address common.Address, passphrase string) (*IPFS, error) {

	account := accounts.Account{Address: address}
	wallet, err := accMng.Find(account)
	if err != nil {
		return nil, err
	}

	transactOpts, err := wallet.NewKeyedTransactor(account, passphrase)
	if err != nil {
		return nil, err
	}

	callOpts := &bind.CallOpts{
		From:	address,
	}

	contract, err := NewIpfs(transactOpts, callOpts, TestNetAddress, backend)
	if err != nil {
		return nil, err
	}

	return contract, nil
}