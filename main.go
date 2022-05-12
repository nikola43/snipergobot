package main

import (
	"buytokenspancakegolang/models"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	ierc20 "buytokenspancakegolang/contracts/IERC20"
	pancakeFactory "buytokenspancakegolang/contracts/IPancakeFactory"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	ccolor "github.com/fatih/color"
	"github.com/go-cmd/cmd"
	"github.com/kyokomi/emoji"
	"github.com/mattn/go-colorable"
	"github.com/mdp/qrterminal/v3"
	"github.com/nikola43/web3golanghelper/web3helper"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Create SprintXxx functions to mix strings with other non-colorized strings:
var yellow = ccolor.New(ccolor.FgYellow).SprintFunc()
var red = ccolor.New(ccolor.FgRed).SprintFunc()
var cyan = ccolor.New(ccolor.FgCyan).SprintFunc()
var green = ccolor.New(ccolor.FgGreen).SprintFunc()

type Wallet struct {
	PublicKey  string `json:"PublicKey"`
	PrivateKey string `json:"PrivateKey"`
}
type Reserve struct {
	Reserve0           *big.Int
	Reserve1           *big.Int
	BlockTimestampLast uint32
}

func main() {
	qrr()
	printWelcome()

	//GenerateWallet()
	fmt.Println(parseDateTime())
	getWallets()

	// Declarations
	web3GolangHelper := initWeb3()
	db := InitDatabase()
	migrate(db)
	factoryAddress := "0xB7926C0430Afb07AA7DEfDE6DA862aE0Bde767bc"
	factoryAbi, _ := abi.JSON(strings.NewReader(string(pancakeFactory.PancakeABI)))

	// LOGIC -----------------------------------------------------------
	checkTokens(db, web3GolangHelper)
	proccessEvents(db, web3GolangHelper, factoryAddress, factoryAbi)
}

func proccessEvents(db *gorm.DB, web3GolangHelper *web3helper.Web3GolangHelper, contractAddress string, contractAbi abi.ABI) {

	logs := make(chan types.Log)
	sub := web3GolangHelper.BuildContractEventSubscription(web3GolangHelper, contractAddress, logs)

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
	pk := "b366406bc0b4883b9b4b3b41117d6c62839174b7d21ec32a5ad0cc76cb3496bd"
	rpcUrl := "https://speedy-nodes-nyc.moralis.io/84a2745d907034e6d388f8d6/bsc/testnet"
	wsUrl := "wss://speedy-nodes-nyc.moralis.io/84a2745d907034e6d388f8d6/bsc/testnet/ws"
	web3GolangHelper := web3helper.NewWeb3GolangHelper(rpcUrl, wsUrl, pk)

	chainID, err := web3GolangHelper.HttpClient().NetworkID(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Chain Id: " + chainID.String())
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
	lpPairs = append(lpPairs, &models.LpPair{
		LPAddress:    newEvent[0].(common.Address).Hex(),
		LPPairA:      common.HexToAddress(tokenAddressA.Hex()).Hex(),
		LPPairB:      common.HexToAddress(tokenAddressB.Hex()).Hex(),
		HasLiquidity: false,
	})

	event.TxHash = vLog.TxHash.Hex()
	event.LPPairs = lpPairs
	if common.HexToAddress(tokenAddressA.Hex()).Hex() != wBnbContractAddress {
		event.TokenAddress = common.HexToAddress(tokenAddressA.Hex()).Hex()
	} else {
		event.TokenAddress = common.HexToAddress(tokenAddressB.Hex()).Hex()
	}

	db.Create(event)

	return true
}

func UpdateLiquidity(db *gorm.DB, eventID uint) bool {
	lpPair := new(models.LpPair)
	db.Where(&lpPair, "events_catched_id = ?", eventID).Update("has_liquidity", 1)

	return true
}

func UpdateName(db *gorm.DB, token string, name string) bool {
	event := new(models.EventsCatched)
	db.Where(&event, "token_address = ?", token).Update("token_name", 1)

	return true
}

func checkTokens(db *gorm.DB, web3GolangHelper *web3helper.Web3GolangHelper) {
	events := make([]*models.EventsCatched, 0)
	db.Find(&events)
	lo.ForEach(events, func(element *models.EventsCatched, _ int) {
		printTokenStatus(element)
		updateTokenStatus(db, web3GolangHelper, element)
	})

}

func updateTokenStatus(db *gorm.DB, web3GolangHelper *web3helper.Web3GolangHelper, token *models.EventsCatched) {

	// create pancakeRouter pancakeRouterInstance
	tokenContractInstance, instanceErr := ierc20.NewPancake(common.HexToAddress(token.TokenAddress), web3GolangHelper.HttpClient())
	if instanceErr != nil {
		fmt.Println(instanceErr)
	}

	tokenName, getNameErr := tokenContractInstance.Name(nil)
	if getNameErr != nil {
		UpdateName(db, token.TokenAddress, tokenName)
		fmt.Println(getNameErr)
	}

	reserves := web3GolangHelper.getReserves(token.TokenAddress)
	if reserves.Reserve0.Uint64() > web3helper.EtherToWei(big.NewFloat(0)).Uint64() {
		UpdateLiquidity(db, token.ID)
	}

}

func printTokenStatus(token *models.EventsCatched) {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetOutput(colorable.NewColorableStdout())
	logrus.Info("TOKEN INFO")

	fmt.Printf("%s: %s\n", cyan("Token Address"), yellow(token.TokenAddress))
	fmt.Printf("%s:\n", cyan("LP Pairs"))
	lo.ForEach(token.LPPairs, func(element *models.LpPair, _ int) {
		fmt.Printf("\t%s: %s\n", cyan("LP Address"), yellow(element.LPAddress))
		fmt.Printf("\t%s: %s\n", cyan("LP TokenA Address"), yellow(element.LPPairA))
		fmt.Printf("\t%s: %s\n", cyan("LP TokenB Address"), yellow(element.LPPairB))
		fmt.Printf("\t%s: %s\n\n", cyan("LP Has Liquidity"), getPairLiquidityIcon(element))
		fmt.Printf("\t%s: %s\n\n", cyan("Trading Enabled"), getPairTradingIcon(element))
	})
}

func getPairTradingIcon(pair *models.LpPair) string {
	icon := "🔴"
	if pair.TradingEnabled {
		icon = "🟢"
	}
	return icon
}

func getPairLiquidityIcon(pair *models.LpPair) string {
	icon := "🔴"
	if pair.HasLiquidity {
		icon = "🟢"
	}
	return icon
}

func GenerateWallet() {

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println(hexutil.Encode(privateKeyBytes)[2:]) // fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println(hexutil.Encode(publicKeyBytes)[4:]) // 9a7df67f79246283fdc93af76d4f8cdd62c4886e8cd870944e817dd0b97934fdd7719d0810951e03418205868a5c1b40b192451367f28e0088dd75e15de40c05

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println(address) // 0x96216849c49358B10257cb55b28eA603c874b05E

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	fmt.Println(hexutil.Encode(hash.Sum(nil)[12:])) // 0x96216849c49358b10257cb55b28ea603c874b05e

	wallet := Wallet{
		PublicKey:  address,
		PrivateKey: hexutil.Encode(privateKeyBytes)[2:],
	}

	file, _ := json.MarshalIndent(wallet, "", " ")
	_ = ioutil.WriteFile("wallets/"+address+".json", file, 0644)
}

func getWallets() {
	wallets := make([]Wallet, 0)

	wPath := "./wallets"
	files, err := ioutil.ReadDir(wPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fileName := file.Name()
		fmt.Println("fileName", fileName)

		wallet := Wallet{
			PublicKey:  "",
			PrivateKey: "",
		}

		// Open our jsonFile
		jsonFile, _ := os.Open(wPath + "/" + fileName)
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &wallet)
		fmt.Println(wallet)
		wallets = append(wallets, wallet)
	}

	fmt.Println(wallets)
}

func parseDateTime() string {
	now := time.Now()
	return strconv.Itoa(now.Year()) + "/" + now.Month().String() + "/" + strconv.Itoa(now.Day()) + " " + strconv.Itoa(now.Hour()) + ":" + strconv.Itoa(now.Minute()) + ":" + strconv.Itoa(now.Second()) + ":" + strconv.Itoa(now.Nanosecond())
}

const DefaultTimeoutTime = "1m"

func RunCMD(name string, args ...string) (err error, stdout, stderr []string) {
	c := cmd.NewCmd(name, args...)
	s := <-c.Start()
	stdout = s.Stdout
	stderr = s.Stderr
	return
}

func printWelcome() {
	clearScreen()

	fmt.Println()
	fmt.Println()
	fmt.Printf("%s\n", red("██████╗░███████╗███╗░░░███╗███████╗███╗░░██╗████████╗░█████╗░██████╗░  ██╗░░░██╗░█████╗░░░░███████╗░░░██████╗░"))
	fmt.Printf("%s\n", red("██╔══██╗██╔════╝████╗░████║██╔════╝████╗░██║╚══██╔══╝██╔══██╗██╔══██╗  ██║░░░██║██╔══██╗░░░██╔════╝░░░╚════██╗"))
	fmt.Printf("%s\n", red("██║░░██║█████╗░░██╔████╔██║█████╗░░██╔██╗██║░░░██║░░░██║░░██║██████╔╝  ╚██╗░██╔╝╚██████║░░░██████╗░░░░░░███╔═╝"))
	fmt.Printf("%s\n", red("██║░░██║██╔══╝░░██║╚██╔╝██║██╔══╝░░██║╚████║░░░██║░░░██║░░██║██╔══██╗  ░╚████╔╝░░╚═══██║░░░╚════██╗░░░██╔══╝░░"))
	fmt.Printf("%s\n", red("██████╔╝███████╗██║░╚═╝░██║███████╗██║░╚███║░░░██║░░░╚█████╔╝██║░░██║  ░░╚██╔╝░░░█████╔╝██╗██████╔╝██╗███████╗"))
	fmt.Printf("%s\n", red("╚═════╝░╚══════╝╚═╝░░░░░╚═╝╚══════╝╚═╝░░╚══╝░░░╚═╝░░░░╚════╝░╚═╝░░╚═╝  ░░░╚═╝░░░░╚════╝░╚═╝╚═════╝░╚═╝╚══════╝"))
	fmt.Println()
	fmt.Println()

	devMessage := emoji.Sprint("Developed with :smiling_imp:By Mr. Nobody")
	fmt.Println(devMessage)
}

func clearScreen() {
	os := runtime.GOOS
	cmdString := "clear"
	fmt.Println(os)
	switch os {
	case "windows":
		cmdString = "cls.exe"
	}

	cmd := exec.Command(cmdString)
	err := cmd.Run()
	if err != nil {
		fmt.Println(cmd)
	}
}

func qrr() {

	config := qrterminal.Config{
		Level:     qrterminal.M,
		Writer:    os.Stdout,
		BlackChar: qrterminal.WHITE,
		WhiteChar: qrterminal.BLACK,
		QuietZone: 1,
	}
	qrterminal.GenerateWithConfig("https://github.com/mdp/qrterminal", config)

	/*
		err := qrcode.WriteColorFile("singana", qrcode.Medium, 256, color.Black, color.White, "secondfile.png")
		if err != nil {
			fmt.Printf("Sorry couldn't create qrcode:,%v", err)

		}
	*/
}
