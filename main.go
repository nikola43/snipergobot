package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/nikola43/snipergobot/dbutils"
	"github.com/nikola43/snipergobot/licenceutils"
	"github.com/nikola43/snipergobot/menuutils"
	"github.com/nikola43/snipergobot/models"

	//ierc20 "buytokenspancakegolang/contracts/IERC20"
	pancakeFactory "github.com/nikola43/snipergobot/contracts/IPancakeFactory"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/go-cmd/cmd"
	"github.com/nikola43/web3golanghelper/web3helper"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

// Create SprintXxx functions to mix strings with other non-colorized strings:

func main() {
	web3GolangHelper := initWeb3()
	db := dbutils.InitDatabase()
	//migrate(db)
	factoryAddress := "0xB7926C0430Afb07AA7DEfDE6DA862aE0Bde767bc"
	factoryAbi, _ := abi.JSON(strings.NewReader(string(pancakeFactory.PancakeABI)))

	//checkLicense()
	createWalletFolder("wallets")
	menuutils.PrintWelcome()

	// check tokens on other goroutine each 5 seconds
	go func() {
		for {
			checkTokens(db, web3GolangHelper)
			time.Sleep(time.Second * 5)
		}
	}()
	proccessEvents(db, web3GolangHelper, factoryAddress, factoryAbi)

	selectedMainMenuOption := mainProgram()
	if selectedMainMenuOption == "1" {
		// LOGIC -----------------------------------------------------------

	} else if selectedMainMenuOption == "2" {

	} else if selectedMainMenuOption == "2" {

	} else {

	}

	//printWelcome()
	//os.Exit(0)
}

func proccessEvents(db *gorm.DB, web3GolangHelper *web3helper.Web3GolangHelper, contractAddress string, contractAbi abi.ABI) {

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

func initWeb3() *web3helper.Web3GolangHelper {
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

func checkTokens(db *gorm.DB, web3GolangHelper *web3helper.Web3GolangHelper) {
	events := make([]*models.EventsCatched, 0)
	db.Joins("INNER JOIN lp_pairs ON lp_pairs.events_catched_id = events_catcheds.id").Where("lp_pairs.has_liquidity = ?", 0).Preload("LPPairs").Find(&events)
	lo.ForEach(events, func(element *models.EventsCatched, _ int) {
		menuutils.PrintTokenStatus(element)
		//UpdateTokenStatus(db, web3GolangHelper, element)
	})
}

func checkTradingActive(tokenAddress string, web3GolangHelper *web3helper.Web3GolangHelper) bool {
	return true
}

const DefaultTimeoutTime = "1m"

func RunCMD(name string, args ...string) (err error, stdout, stderr []string) {
	c := cmd.NewCmd(name, args...)
	s := <-c.Start()
	stdout = s.Stdout
	stderr = s.Stderr
	return
}

func mainProgram() string {
	valid := false
	mainMenuOption := "1"

	for ok := true; ok; ok = !valid {
		menuutils.PrintMainMenu()
		mainMenuOption = readFromKeyBoard("Select any option: ")
		valid = mainMenuOption == "1" || mainMenuOption == "2" || mainMenuOption == "3" || mainMenuOption == "4"
		if !valid {
			fmt.Printf("\n%s\n", menuutils.Red("Invalid option"))
		} else if mainMenuOption == "1" {
			menuutils.PrintAccounts()
		} else if mainMenuOption == "2" {
			menuutils.PrintLoginMenu()
		} else if mainMenuOption == "4" {
			licenceutils.ShowPaymentQr()
		}
	}
	fmt.Println("You select " + mainMenuOption)
	mainMenuOption = "0"
	return mainMenuOption
}

func readFromKeyBoard(text string) string {
	var data string

	fmt.Print(text + ": ")
	fmt.Scanf("%s", &data)

	return data
}

func createWalletFolder(dirname string) bool {
	_, error := os.Stat(dirname)
	if os.IsNotExist(error) {
		if err := os.Mkdir(dirname, os.ModePerm); err != nil {
			log.Fatal(err)
		}
		return true
	} else {
		return true
	}
}
