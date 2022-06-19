package web3utils

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ierc20 "github.com/nikola43/snipergobot/contracts/IERC20"
	"github.com/nikola43/snipergobot/dbutils"
	"github.com/nikola43/snipergobot/menuutils"
	"github.com/nikola43/web3golanghelper/web3helper"
	"gorm.io/gorm"
)

var FactoryAddress = "0xB7926C0430Afb07AA7DEfDE6DA862aE0Bde767bc"

func InitWeb3(pk string) *web3helper.Web3GolangHelper {

	rpcUrl := "https://speedy-nodes-nyc.moralis.io/84a2745d907034e6d388f8d6/bsc/testnet"
	wsUrl := "wss://speedy-nodes-nyc.moralis.io/84a2745d907034e6d388f8d6/bsc/testnet/ws"

	web3GolangHelper := web3helper.NewWeb3GolangHelper(rpcUrl, wsUrl, pk)

	chainID, err := web3GolangHelper.HttpClient().NetworkID(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", menuutils.Green("Web3 Conected"))
	fmt.Printf("%s: %s\n", menuutils.Cyan("Account"), menuutils.Yellow(web3GolangHelper.FromAddress))
	fmt.Printf("%s: %s\n", menuutils.Cyan("Detected Chain Id: "), menuutils.Yellow(chainID))

	//web3GolangHelper.BuyV2("0x7ef95a0FEE0Dd31b22626fA2e10Ee6A223F8a684", big.NewInt(10000000000000000))

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

			tokenAddressA := vLog.Topics[1]
			tokenAddressB := vLog.Topics[2]

			fmt.Println(tokenAddressA)
			fmt.Println(tokenAddressB)

			if tokenAddressA.Hex() != "" && tokenAddressB.Hex() != "" {
				dbutils.InsertNewEvent(db, res, vLog)
			}
		}
	}

	/*
		for {
			select {
			case err := <-sub.Err():
				fmt.Println(err)
				//out <- err.Error()

			case vLog := <-logs:
				fmt.Println("paco")
				fmt.Println("vLog.TxHash: " + vLog.TxHash.Hex())
				fmt.Println(vLog)
				if len(vLog.Data) > 0 {
					switch vLog.Topics[0].Hex() {
					case giftCardPayedSigHash.Hex():
						res, err := nodeAbi.Unpack("GiftCardPayed", vLog.Data)
						if err != nil {
							log.Fatal(err)
						}
						fmt.Println(res)
						services.SetGiftCardIntentPayment(res[2].(string))
					}
				}

			}
		}
	*/

}

func GetTokenInfo(db *gorm.DB, web3GolangHelper *web3helper.Web3GolangHelper, tokenAddress string, lpAddress string) bool {
	// create pancakeRouter pancakeRouterInstance
	tokenContractInstance, instanceErr := ierc20.NewPancake(common.HexToAddress(tokenAddress), web3GolangHelper.HttpClient())
	if instanceErr != nil {
		fmt.Println("instanceErr")
		fmt.Println(instanceErr)
	}

	name, getNameErr := tokenContractInstance.Name(nil)
	if getNameErr != nil {
		fmt.Println("getNameErr")
		fmt.Println(getNameErr)
	}

	symbol, getSymbolErr := tokenContractInstance.Symbol(nil)
	if getSymbolErr != nil {
		fmt.Println("getSymbolErr")
		fmt.Println(getSymbolErr)
	}

	/*
		lpPairAddress, getSymbolErr := tokenContractInstance.LpPair(nil)
		if getSymbolErr != nil {
			fmt.Println(getSymbolErr)
		}
	*/

	reserves := web3GolangHelper.GetReserves(lpAddress)
	fmt.Println("reserves")
	fmt.Println(reserves)
	if reserves.BlockTimestampLast != 0 {
		//dbutils.UpdateLiquidity(db, token.ID)
	}

	fmt.Printf("%s: %s\n", menuutils.Cyan("Token Address"), menuutils.Yellow(tokenAddress))
	fmt.Printf("%s: %s\n", menuutils.Cyan("Name"), menuutils.Yellow(name))
	fmt.Printf("%s: %s\n", menuutils.Cyan("Symbol"), menuutils.Yellow(symbol))
	fmt.Printf("%s:\n", menuutils.Cyan("LP Pairs"))

	fmt.Printf("\t%s: %s\n", menuutils.Cyan("LP Address"), menuutils.Yellow("0x536f0A9fdC03eDcAF78720f6E3855F7bb6fEcA36"))
	fmt.Printf("\t%s: %s\n", menuutils.Cyan("LP TokenA Address"), menuutils.Yellow("0x536f0A9fdC03eDcAF78720f6E3855F7bb6fEcA36"))
	fmt.Printf("\t%s: %s\n", menuutils.Cyan("LP TokenB Address"), menuutils.Yellow("0xae13d989daC2f0dEbFf460aC112a837C89BAa7cd"))
	fmt.Printf("\t%s: %s\n", menuutils.Cyan("LP Has Liquidity"), menuutils.GetPairLiquidityIcon(reserves.Reserve0.Uint64() > 0 && reserves.Reserve1.Uint64() > 0))
	fmt.Printf("\t%s: %s\n", menuutils.Cyan("Trading Active"), menuutils.GetPairLiquidityIcon(false))

	//UpdateTradingActive(db, token.ID)

	return (reserves.Reserve0.Uint64() > 0 && reserves.Reserve1.Uint64() > 0)
}

func checkTradingActive(tokenAddress string, web3GolangHelper *web3helper.Web3GolangHelper) bool {
	return true
}

/*
func GetTokenInfoV2(db *gorm.DB, web3GolangHelper *web3helper.Web3GolangHelper, tokenAddres string) {

	// create pancakeRouter pancakeRouterInstance
	tokenContractInstance, instanceErr := ierc20.NewPancake(common.HexToAddress(tokenAddres), web3GolangHelper.HttpClient())
	if instanceErr != nil {
		fmt.Println(instanceErr)
	}

	name, getNameErr := tokenContractInstance.Name(nil)
	if getNameErr != nil {
		fmt.Println(getNameErr)
	}

	symbol, getSymbolErr := tokenContractInstance.Symbol(nil)
	if getSymbolErr != nil {
		fmt.Println(getSymbolErr)
	}

	lpPairAddress, getSymbolErr := tokenContractInstance.LpPair(nil)
	if getSymbolErr != nil {
		fmt.Println(getSymbolErr)
	}

	reserves := web3GolangHelper.GetReserves("0x")
	if reserves.BlockTimestampLast != 0 {
		dbutils.UpdateLiquidity(db, token.ID)
	}

	fmt.Printf("%s: %s\n", menuutils.Cyan("Token Name: "), menuutils.Yellow(name))
	fmt.Printf("%s: %s\n", menuutils.Cyan("Token Symbol: "), menuutils.Yellow(symbol))
	fmt.Printf("%s: %s\n", menuutils.Cyan("Token Symbol: "), menuutils.Yellow(symbol))

	tradingActive := checkTradingActive(token.TokenAddress, web3GolangHelper)
	fmt.Println(tradingActive)
	//UpdateTradingActive(db, token.ID)
}
*/
