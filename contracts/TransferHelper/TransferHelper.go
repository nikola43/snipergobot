// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package TransferHelper

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// PancakeABI is the input ABI used to generate the binding from.
const PancakeABI = "[]"

// PancakeBin is the compiled bytecode used for deploying new contracts.
var PancakeBin = "0x60566023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea26469706673582212207f85ac0eaed49fca814d4d4094339934f89277f3c42bbc0227af69c7e7058b6f64736f6c63430006060033"

// DeployPancake deploys a new Ethereum contract, binding an instance of Pancake to it.
func DeployPancake(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Pancake, error) {
	parsed, err := abi.JSON(strings.NewReader(PancakeABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(PancakeBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Pancake{PancakeCaller: PancakeCaller{contract: contract}, PancakeTransactor: PancakeTransactor{contract: contract}, PancakeFilterer: PancakeFilterer{contract: contract}}, nil
}

// Pancake is an auto generated Go binding around an Ethereum contract.
type Pancake struct {
	PancakeCaller     // Read-only binding to the contract
	PancakeTransactor // Write-only binding to the contract
	PancakeFilterer   // Log filterer for contract events
}

// PancakeCaller is an auto generated read-only Go binding around an Ethereum contract.
type PancakeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PancakeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PancakeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PancakeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PancakeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PancakeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PancakeSession struct {
	Contract     *Pancake          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PancakeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PancakeCallerSession struct {
	Contract *PancakeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// PancakeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PancakeTransactorSession struct {
	Contract     *PancakeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// PancakeRaw is an auto generated low-level Go binding around an Ethereum contract.
type PancakeRaw struct {
	Contract *Pancake // Generic contract binding to access the raw methods on
}

// PancakeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PancakeCallerRaw struct {
	Contract *PancakeCaller // Generic read-only contract binding to access the raw methods on
}

// PancakeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PancakeTransactorRaw struct {
	Contract *PancakeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPancake creates a new instance of Pancake, bound to a specific deployed contract.
func NewPancake(address common.Address, backend bind.ContractBackend) (*Pancake, error) {
	contract, err := bindPancake(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Pancake{PancakeCaller: PancakeCaller{contract: contract}, PancakeTransactor: PancakeTransactor{contract: contract}, PancakeFilterer: PancakeFilterer{contract: contract}}, nil
}

// NewPancakeCaller creates a new read-only instance of Pancake, bound to a specific deployed contract.
func NewPancakeCaller(address common.Address, caller bind.ContractCaller) (*PancakeCaller, error) {
	contract, err := bindPancake(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PancakeCaller{contract: contract}, nil
}

// NewPancakeTransactor creates a new write-only instance of Pancake, bound to a specific deployed contract.
func NewPancakeTransactor(address common.Address, transactor bind.ContractTransactor) (*PancakeTransactor, error) {
	contract, err := bindPancake(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PancakeTransactor{contract: contract}, nil
}

// NewPancakeFilterer creates a new log filterer instance of Pancake, bound to a specific deployed contract.
func NewPancakeFilterer(address common.Address, filterer bind.ContractFilterer) (*PancakeFilterer, error) {
	contract, err := bindPancake(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PancakeFilterer{contract: contract}, nil
}

// bindPancake binds a generic wrapper to an already deployed contract.
func bindPancake(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PancakeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Pancake *PancakeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Pancake.Contract.PancakeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Pancake *PancakeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Pancake.Contract.PancakeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Pancake *PancakeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Pancake.Contract.PancakeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Pancake *PancakeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Pancake.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Pancake *PancakeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Pancake.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Pancake *PancakeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Pancake.Contract.contract.Transact(opts, method, params...)
}
