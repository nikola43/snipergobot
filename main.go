package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
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

	tokenAddress := "0x1d61D38096f161410612b20D2c1080e5549dB3dd"
	lpAddress := "0x140a1b90945fCB4C453846e6B147da7e36f1eeC0"

	hasLiquidity := web3utils.GetTokenInfo(db, web3GolangHelper, tokenAddress, lpAddress)
	fmt.Println(hasLiquidity)

	/*
		data := []float64{3, 4, 9, 6, 2, 4, 5, 8, 5, 10, 2, 7, 2, 5, 6}
		graph := asciigraph.Plot(data)

		fmt.Println(graph)
	*/

	//web3GolangHelper.Buy(tokenAddress, 0.1)

	// check tokens on other goroutine each 5 seconds
	go func() {
		for {
			checkTokens(db, web3GolangHelper)
		}
	}()
	web3utils.ProccessContractEvents(db, web3GolangHelper, factoryAddress, factoryAbi)

}

func ParseDateTime(now time.Time) string {
	return strconv.Itoa(now.Year()) + "/" + now.Month().String() + "/" + strconv.Itoa(now.Day()) + " " + strconv.Itoa(now.Hour()) + ":" + strconv.Itoa(now.Minute()) + ":" + strconv.Itoa(now.Second()) + ":" + strconv.Itoa(now.Nanosecond())
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

func checkTokens(db *gorm.DB, web3GolangHelper *web3helper.Web3GolangHelper) {

	//fmt.Println("now", ParseDateTime(time.Now()))
	events := make([]*models.EventsCatched, 0)
	db.Joins("INNER JOIN lp_pairs ON lp_pairs.events_catched_id = events_catcheds.id").Where("lp_pairs.has_liquidity = ?", 0).Preload("LPPairs").Find(&events)
	lo.ForEach(events, func(element *models.EventsCatched, _ int) {
		hasLiquidity := web3utils.GetTokenInfo(db, web3GolangHelper, element.TokenAddress, element.LPPairs[0].LPAddress)
		//menuutils.PrintTokenPriceInfo(element)
		menuutils.PrintTokenStatus(element)
		dbutils.UpdateTokenStatus(db, web3GolangHelper, element)
		if hasLiquidity {
			//checkTradingActive(element.TokenAddress, web3GolangHelper)
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
