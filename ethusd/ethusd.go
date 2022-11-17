// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ethusd

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// EthusdMetaData contains all meta data concerning the Ethusd contract.
var EthusdMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// EthusdABI is the input ABI used to generate the binding from.
// Deprecated: Use EthusdMetaData.ABI instead.
var EthusdABI = EthusdMetaData.ABI

// Ethusd is an auto generated Go binding around an Ethereum contract.
type Ethusd struct {
	EthusdCaller     // Read-only binding to the contract
	EthusdTransactor // Write-only binding to the contract
	EthusdFilterer   // Log filterer for contract events
}

// EthusdCaller is an auto generated read-only Go binding around an Ethereum contract.
type EthusdCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthusdTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EthusdTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthusdFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EthusdFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthusdSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EthusdSession struct {
	Contract     *Ethusd           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EthusdCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EthusdCallerSession struct {
	Contract *EthusdCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// EthusdTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EthusdTransactorSession struct {
	Contract     *EthusdTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EthusdRaw is an auto generated low-level Go binding around an Ethereum contract.
type EthusdRaw struct {
	Contract *Ethusd // Generic contract binding to access the raw methods on
}

// EthusdCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EthusdCallerRaw struct {
	Contract *EthusdCaller // Generic read-only contract binding to access the raw methods on
}

// EthusdTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EthusdTransactorRaw struct {
	Contract *EthusdTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEthusd creates a new instance of Ethusd, bound to a specific deployed contract.
func NewEthusd(address common.Address, backend bind.ContractBackend) (*Ethusd, error) {
	contract, err := bindEthusd(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Ethusd{EthusdCaller: EthusdCaller{contract: contract}, EthusdTransactor: EthusdTransactor{contract: contract}, EthusdFilterer: EthusdFilterer{contract: contract}}, nil
}

// NewEthusdCaller creates a new read-only instance of Ethusd, bound to a specific deployed contract.
func NewEthusdCaller(address common.Address, caller bind.ContractCaller) (*EthusdCaller, error) {
	contract, err := bindEthusd(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EthusdCaller{contract: contract}, nil
}

// NewEthusdTransactor creates a new write-only instance of Ethusd, bound to a specific deployed contract.
func NewEthusdTransactor(address common.Address, transactor bind.ContractTransactor) (*EthusdTransactor, error) {
	contract, err := bindEthusd(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EthusdTransactor{contract: contract}, nil
}

// NewEthusdFilterer creates a new log filterer instance of Ethusd, bound to a specific deployed contract.
func NewEthusdFilterer(address common.Address, filterer bind.ContractFilterer) (*EthusdFilterer, error) {
	contract, err := bindEthusd(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EthusdFilterer{contract: contract}, nil
}

// bindEthusd binds a generic wrapper to an already deployed contract.
func bindEthusd(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EthusdABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ethusd *EthusdRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Ethusd.Contract.EthusdCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ethusd *EthusdRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ethusd.Contract.EthusdTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ethusd *EthusdRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ethusd.Contract.EthusdTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ethusd *EthusdCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Ethusd.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ethusd *EthusdTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ethusd.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ethusd *EthusdTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ethusd.Contract.contract.Transact(opts, method, params...)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Ethusd *EthusdCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Ethusd.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Ethusd *EthusdSession) Decimals() (uint8, error) {
	return _Ethusd.Contract.Decimals(&_Ethusd.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Ethusd *EthusdCallerSession) Decimals() (uint8, error) {
	return _Ethusd.Contract.Decimals(&_Ethusd.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_Ethusd *EthusdCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Ethusd.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_Ethusd *EthusdSession) Description() (string, error) {
	return _Ethusd.Contract.Description(&_Ethusd.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_Ethusd *EthusdCallerSession) Description() (string, error) {
	return _Ethusd.Contract.Description(&_Ethusd.CallOpts)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_Ethusd *EthusdCaller) GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _Ethusd.contract.Call(opts, &out, "getRoundData", _roundId)

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Answer = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.StartedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.AnsweredInRound = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_Ethusd *EthusdSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _Ethusd.Contract.GetRoundData(&_Ethusd.CallOpts, _roundId)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_Ethusd *EthusdCallerSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _Ethusd.Contract.GetRoundData(&_Ethusd.CallOpts, _roundId)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_Ethusd *EthusdCaller) LatestRoundData(opts *bind.CallOpts) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _Ethusd.contract.Call(opts, &out, "latestRoundData")

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Answer = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.StartedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.AnsweredInRound = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_Ethusd *EthusdSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _Ethusd.Contract.LatestRoundData(&_Ethusd.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_Ethusd *EthusdCallerSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _Ethusd.Contract.LatestRoundData(&_Ethusd.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_Ethusd *EthusdCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Ethusd.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_Ethusd *EthusdSession) Version() (*big.Int, error) {
	return _Ethusd.Contract.Version(&_Ethusd.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_Ethusd *EthusdCallerSession) Version() (*big.Int, error) {
	return _Ethusd.Contract.Version(&_Ethusd.CallOpts)
}
