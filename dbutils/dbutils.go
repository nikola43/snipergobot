package dbutils

import (
	"fmt"

	"github.com/nikola43/snipergobot/models"
	"github.com/nikola43/web3golanghelper/web3helper"
	"github.com/samber/lo"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ierc20 "github.com/nikola43/snipergobot/contracts/IERC20"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func Migrate(db *gorm.DB) {
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

func UpdateLiquidity(db *gorm.DB, eventID uint) bool {
	lpPair := new(models.LpPair)
	db.Model(&lpPair).Where("events_catched_id = ?", eventID).Update("has_liquidity", 1)
	return true
}

func UpdateTradingActive(db *gorm.DB, eventID uint) bool {
	lpPair := new(models.LpPair)
	db.Model(&lpPair).Where("events_catched_id = ?", eventID).Update("has_liquidity", 1)
	return true
}

func UpdateName(db *gorm.DB, token string, name string) bool {
	event := new(models.EventsCatched)
	db.Model(&event).Where("token_address = ?", token).Where("token_name != ?", name).Update("token_name", name)
	return true
}

func UpdateTokenStatus(db *gorm.DB, web3GolangHelper *web3helper.Web3GolangHelper, token *models.EventsCatched) {

	// create pancakeRouter pancakeRouterInstance
	tokenContractInstance, instanceErr := ierc20.NewPancake(common.HexToAddress(token.TokenAddress), web3GolangHelper.HttpClient())
	if instanceErr != nil {
		fmt.Println(instanceErr)
	}

	tokenName, getNameErr := tokenContractInstance.Name(nil)
	if getNameErr != nil {
		fmt.Println(getNameErr)
	}

	tokenSymbol, getSymbolErr := tokenContractInstance.Symbol(nil)
	if getSymbolErr != nil {
		fmt.Println(getSymbolErr)
	}
	UpdateName(db, token.TokenAddress, tokenName+" - $"+tokenSymbol)
	lo.ForEach(token.LPPairs, func(element *models.LpPair, _ int) {
		reserves := web3GolangHelper.GetReserves(element.LPAddress)
		if reserves.BlockTimestampLast != 0 {
			UpdateLiquidity(db, token.ID)
		}
	})

	UpdateTradingActive(db, token.ID)
}
