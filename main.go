package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nikola43/snipergobot/dbutils"
	"github.com/nikola43/snipergobot/genericutils"
	"github.com/nikola43/snipergobot/licenceutils"
	"github.com/nikola43/snipergobot/menuutils"
	"github.com/nikola43/snipergobot/models"
	"github.com/nikola43/snipergobot/web3utils"

	//ierc20 "buytokenspancakegolang/contracts/IERC20"
	pancakeFactory "github.com/nikola43/snipergobot/contracts/IPancakeFactory"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/nikola43/web3golanghelper/web3helper"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

// Create SprintXxx functions to mix strings with other non-colorized strings:

func main() {
	web3GolangHelper := web3utils.InitWeb3()
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
	web3utils.ProccessContractEvents(db, web3GolangHelper, factoryAddress, factoryAbi)

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

func mainProgram() string {
	valid := false
	mainMenuOption := "1"

	for ok := true; ok; ok = !valid {
		menuutils.PrintMainMenu()
		mainMenuOption = genericutils.ReadFromKeyBoard("Select any option: ")
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
