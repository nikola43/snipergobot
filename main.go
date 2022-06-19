package main

import (
	"bufio"
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
	pk := "3062979ebcda3efb3bae3919e003f8e3a3597d9244244a13e4a9ff7776221501"
	web3GolangHelper := web3utils.InitWeb3(pk)
	db := dbutils.InitDatabase()
	//dbutils.Migrate(db)
	factoryAddress := "0xB7926C0430Afb07AA7DEfDE6DA862aE0Bde767bc"
	factoryAbi, _ := abi.JSON(strings.NewReader(string(pancakeFactory.PancakeABI)))

	//checkLicense()
	createWalletFolder("wallets")

	//tokenAddress := readUserContractAddressInput()
	//lpAddress := readLpContractAddressInput()

	tokenAddress := "0x536f0A9fdC03eDcAF78720f6E3855F7bb6fEcA36"
	lpAddress := "0x692Cf4A84962E94Db915CD7dc890b05C6A196812"

	hasLiquidity := web3utils.GetTokenInfo(db, web3GolangHelper, tokenAddress, lpAddress)
	fmt.Println(hasLiquidity)

	web3GolangHelper.Buy(tokenAddress, 0.1)

	// check tokens on other goroutine each 5 seconds
	go func() {
		for {
			checkTokens(db, web3GolangHelper, tokenAddress, lpAddress)
			time.Sleep(time.Second * 5)
		}
	}()
	web3utils.ProccessContractEvents(db, web3GolangHelper, factoryAddress, factoryAbi, tokenAddress, lpAddress)

}

func readLpContractAddressInput() string {

	fmt.Print("Enter token LP: ")
	lpFilterAddress := ""

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lpFilterAddress = scanner.Text()
		fmt.Printf("%s: %s\n", menuutils.Cyan("LP Address"), menuutils.Yellow(lpFilterAddress))
		return lpFilterAddress
	}

	if scanner.Err() != nil {
		// Handle error.
		fmt.Println("scanner.Err()")
		fmt.Println(scanner.Err())
	}

	return lpFilterAddress
}

func readUserContractAddressInput() string {

	fmt.Print("Enter token contract Address: ")
	tokenFilterAddress := ""

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		tokenFilterAddress = scanner.Text()
		fmt.Printf("%s: %s\n", menuutils.Cyan("Selected Token"), menuutils.Yellow(tokenFilterAddress))
		return tokenFilterAddress
	}

	if scanner.Err() != nil {
		// Handle error.
		fmt.Println("scanner.Err()")
		fmt.Println(scanner.Err())
	}

	return tokenFilterAddress
}

func checkTokens(db *gorm.DB, web3GolangHelper *web3helper.Web3GolangHelper, tokenAddress string, lpAddress string) {
	events := make([]*models.EventsCatched, 0)
	db.Joins("INNER JOIN lp_pairs ON lp_pairs.events_catched_id = events_catcheds.id").Where("lp_pairs.has_liquidity = ?", 0).Preload("LPPairs").Find(&events)
	lo.ForEach(events, func(element *models.EventsCatched, _ int) {
		hasLiquidity := web3utils.GetTokenInfo(db, web3GolangHelper, tokenAddress, lpAddress)
		menuutils.PrintTokenStatus(element)
		dbutils.UpdateTokenStatus(db, web3GolangHelper, element)
		if hasLiquidity {
			checkTradingActive(tokenAddress, web3GolangHelper)
		}

	})
}

func checkTradingActive(tokenAddress string, web3GolangHelper *web3helper.Web3GolangHelper) bool {

	web3GolangHelper.Buy(tokenAddress, 0.1)

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
