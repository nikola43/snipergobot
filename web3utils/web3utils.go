package web3utils

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/nikola43/snipergobot/dbutils"
	"github.com/nikola43/web3golanghelper/web3helper"
	"gorm.io/gorm"
)

func InitWeb3() *web3helper.Web3GolangHelper {
	pk := "3062979ebcda3efb3bae3919e003f8e3a3597d9244244a13e4a9ff7776221501"
	rpcUrl := "https://speedy-nodes-nyc.moralis.io/84a2745d907034e6d388f8d6/bsc/testnet"
	wsUrl := "wss://speedy-nodes-nyc.moralis.io/84a2745d907034e6d388f8d6/bsc/testnet/ws"
	web3GolangHelper := web3helper.NewWeb3GolangHelper(rpcUrl, wsUrl, pk)

	chainID, err := web3GolangHelper.HttpClient().NetworkID(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Chain Id: " + chainID.String())

	web3GolangHelper.BuyV2("0x7ef95a0FEE0Dd31b22626fA2e10Ee6A223F8a684", big.NewInt(10000000000000000))

	return web3GolangHelper
}

func ProccessContractEvents(db *gorm.DB, web3GolangHelper *web3helper.Web3GolangHelper, contractAddress string, contractAbi abi.ABI) {

	logs := make(chan types.Log)
	sub := web3GolangHelper.BuildContractEventSubscription(contractAddress, logs)
	fmt.Println("sniping....")
	for {
		select {
		case err := <-sub.Err():
			fmt.Println(err)
			//out <- err.Error()

		case vLog := <-logs:
			fmt.Println("vLog.TxHash: " + vLog.TxHash.Hex())

			res, err := contractAbi.Unpack("PairCreated", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(res)

			dbutils.InsertNewEvent(db, res, vLog)
		}
	}
}
