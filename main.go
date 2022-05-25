package main

import (
	"github.com/nikola43/snipergobot/models"
	"github.com/nikola43/snipergobot/menuhelper"
	"github.com/nikola43/snipergobot/licencehelper"
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	//ierc20 "buytokenspancakegolang/contracts/IERC20"
	pancakeFactory "github.com/nikola43/snipergobot/contracts/IPancakeFactory"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ccolor "github.com/fatih/color"
	"github.com/go-cmd/cmd"

	//"github.com/mattn/go-colorable"

	"github.com/nikola43/web3golanghelper/web3helper"
	"github.com/samber/lo"

	//"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Create SprintXxx functions to mix strings with other non-colorized strings:
var Yellow = ccolor.New(ccolor.FgYellow).SprintFunc()
var Red = ccolor.New(ccolor.FgRed).SprintFunc()
var Cyan = ccolor.New(ccolor.FgCyan).SprintFunc()
var Green = ccolor.New(ccolor.FgGreen).SprintFunc()

const (
	keyFile = "aes.key"
)

var encryptionKey = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func main() {
	web3GolangHelper := initWeb3()
	db := InitDatabase()
	//migrate(db)
	factoryAddress := "0xB7926C0430Afb07AA7DEfDE6DA862aE0Bde767bc"
	factoryAbi, _ := abi.JSON(strings.NewReader(string(pancakeFactory.PancakeABI)))

	//checkLicense()
	createWalletFolder("wallets")
	menuhelper.PrintWelcome()

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

			InsertNewEvent(db, res, vLog)
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

func InitDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func migrate(db *gorm.DB) {
	db.Migrator().DropTable(&models.EventsCatched{})
	db.Migrator().DropTable(&models.LpPair{})
	db.Migrator().CreateTable(&models.LpPair{})
	db.Migrator().CreateTable(&models.EventsCatched{})
}

func InsertNewEvent(db *gorm.DB, newEvent []interface{}, vLog types.Log) bool {
	wBnbContractAddress := "0xae13d989daC2f0dEbFf460aC112a837C89BAa7cd"
	tokenAddressA := vLog.Topics[1]
	tokenAddressB := vLog.Topics[2]
	event := new(models.EventsCatched)
	lpPairs := make([]*models.LpPair, 0)

	event.TxHash = vLog.TxHash.Hex()
	lpPairs = append(lpPairs, &models.LpPair{
		LPAddress:    newEvent[0].(common.Address).Hex(),
		LPPairA:      common.HexToAddress(tokenAddressA.Hex()).Hex(),
		LPPairB:      common.HexToAddress(tokenAddressB.Hex()).Hex(),
		HasLiquidity: false,
	})

	event.LPPairs = lpPairs
	if common.HexToAddress(tokenAddressA.Hex()).Hex() != wBnbContractAddress {
		event.TokenAddress = common.HexToAddress(tokenAddressA.Hex()).Hex()
	} else {
		event.TokenAddress = common.HexToAddress(tokenAddressB.Hex()).Hex()
	}

	db.Create(event)

	return true
}

func checkTokens(db *gorm.DB, web3GolangHelper *web3helper.Web3GolangHelper) {
	events := make([]*models.EventsCatched, 0)
	db.Joins("INNER JOIN lp_pairs ON lp_pairs.events_catched_id = events_catcheds.id").Where("lp_pairs.has_liquidity = ?", 0).Preload("LPPairs").Find(&events)
	lo.ForEach(events, func(element *models.EventsCatched, _ int) {
		menuhelper.PrintTokenStatus(element)
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
		menuhelper.PrintMainMenu()
		mainMenuOption = readFromKeyBoard("Select any option: ")
		valid = mainMenuOption == "1" || mainMenuOption == "2" || mainMenuOption == "3" || mainMenuOption == "4"
		if !valid {
			fmt.Printf("\n%s\n", Red("Invalid option"))
		} else if mainMenuOption == "1" {
			menuhelper.PrintAccounts()
		} else if mainMenuOption == "2" {
			menuhelper.PrintLoginMenu()
		} else if mainMenuOption == "4" {
			licencehelper.ShowPaymentQr()
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
