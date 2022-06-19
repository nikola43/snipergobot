package menuutils

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/nikola43/snipergobot/models"

	ccolor "github.com/fatih/color"
	"github.com/kyokomi/emoji"
	"github.com/samber/lo"
)

// Create SprintXxx functions to mix strings with other non-colorized strings:
var Yellow = ccolor.New(ccolor.FgYellow).SprintFunc()
var Red = ccolor.New(ccolor.FgRed).SprintFunc()
var Cyan = ccolor.New(ccolor.FgCyan).SprintFunc()
var Green = ccolor.New(ccolor.FgGreen).SprintFunc()

func PrintTokenStatus(token *models.EventsCatched) {
	//logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	//logrus.SetOutput(colorable.NewColorableStdout())
	//logrus.Info("TOKEN INFO")

	fmt.Printf("%s\n", time.Now())
	fmt.Printf("%s: %s\n", Cyan("Token Address"), Yellow(token.TokenAddress))
	fmt.Printf("%s: %s\n", Cyan("Token Name"), Yellow(token.TokenName))
	fmt.Printf("%s:\n", Cyan("LP Pairs"))
	lo.ForEach(token.LPPairs, func(element *models.LpPair, _ int) {
		fmt.Printf("\t%s: %s\n", Cyan("LP Address"), Yellow(element.LPAddress))
		fmt.Printf("\t%s: %s\n", Cyan("LP TokenA Address"), Yellow(element.LPPairA))
		fmt.Printf("\t%s: %s\n", Cyan("LP TokenB Address"), Yellow(element.LPPairB))
		fmt.Printf("\t%s: %s\n", Cyan("LP Has Liquidity"), GetPairLiquidityIcon(true))
		fmt.Printf("\t%s:  %s\n\n", Cyan("Trading Enabled"), GetPairTradingIcon(false))
	})
	fmt.Println("------------------------------------------------------------------------")
}

func GetPairTradingIcon(tradingEnabled bool) string {
	icon := "🔴"
	if tradingEnabled {
		icon = "🟢"
	}
	return icon
}

func GetPairLiquidityIcon(hasLiquidity bool) string {
	icon := "🔴"
	if hasLiquidity {
		icon = "🟢"
	}
	return icon
}

func PrintLoginMenu() {
	fmt.Printf("\t%s\n", Cyan("Login"))
	fmt.Printf("\t%s: %s\n", Cyan("1. "), Yellow("Import wallet using private key"))
	fmt.Printf("\t%s: %s\n", Cyan("2. "), Yellow("Import wallet using wallet file"))
	fmt.Printf("\t%s: %s\n", Cyan("3. "), Yellow("Import wallet using Wallet connnect"))
}

func PrintAccounts() {
	fmt.Printf("\t%s\n", Cyan("Accounts"))
	fmt.Printf("\t%s: %s\n", Cyan("1. "), Yellow("Generate new wallet file"))
	fmt.Printf("\t%s: %s\n", Cyan("2. "), Yellow("Show wallets files"))
	fmt.Printf("\t%s: %s\n", Cyan("3. "), Yellow("Return back"))

	/*
		valid := false
		accountMenuOption := "1"

		for ok := true; ok; ok = !valid {
			accountMenuOption = readFromKeyBoard("Select any option: ")
			valid = accountMenuOption == "1" || accountMenuOption == "2" || accountMenuOption == "3"
			if !valid {
				fmt.Printf("\n%s\n", Red("Invalid option"))
			} else if accountMenuOption == "1" {
				PrintAccounts()
			} else if accountMenuOption == "2" {
				PrintLoginMenu()
			} else if accountMenuOption == "3" {
				mainProgram()
			}
		}
		fmt.Println("You select " + accountMenuOption)
	*/
}

func PrintMainMenu() {
	fmt.Printf("%s\n", Red("\tMAIN MENU"))
	fmt.Printf("\t%s: %s\n", Cyan("1"), Yellow("Init Sniper"))
	fmt.Printf("\t%s: %s\n", Cyan("3"), Yellow("Manage accounts"))
	fmt.Printf("\t%s: %s\n", Cyan("4"), Yellow("Show Payment QR"))
}

func ClearScreen() {
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

func PrintWelcome() {
	ClearScreen()

	fmt.Println()
	fmt.Println()
	fmt.Printf("%s\n", Red("\t██████╗░███████╗███╗░░░███╗███████╗███╗░░██╗████████╗░█████╗░██████╗░  ██╗░░░██╗░█████╗░░░░███████╗░░░██████╗░"))
	fmt.Printf("%s\n", Red("\t██╔══██╗██╔════╝████╗░████║██╔════╝████╗░██║╚══██╔══╝██╔══██╗██╔══██╗  ██║░░░██║██╔══██╗░░░██╔════╝░░░╚════██╗"))
	fmt.Printf("%s\n", Red("\t██║░░██║█████╗░░██╔████╔██║█████╗░░██╔██╗██║░░░██║░░░██║░░██║██████╔╝  ╚██╗░██╔╝╚██████║░░░██████╗░░░░░░███╔═╝"))
	fmt.Printf("%s\n", Red("\t██║░░██║██╔══╝░░██║╚██╔╝██║██╔══╝░░██║╚████║░░░██║░░░██║░░██║██╔══██╗  ░╚████╔╝░░╚═══██║░░░╚════██╗░░░██╔══╝░░"))
	fmt.Printf("%s\n", Red("\t██████╔╝███████╗██║░╚═╝░██║███████╗██║░╚███║░░░██║░░░╚█████╔╝██║░░██║  ░░╚██╔╝░░░█████╔╝██╗██████╔╝██╗███████╗"))
	fmt.Printf("%s\n", Red("\t╚═════╝░╚══════╝╚═╝░░░░░╚═╝╚══════╝╚═╝░░╚══╝░░░╚═╝░░░░╚════╝░╚═╝░░╚═╝  ░░░╚═╝░░░░╚════╝░╚═╝╚═════╝░╚═╝╚══════╝"))
	fmt.Println()
	fmt.Println()

	devMessage := emoji.Sprint("\t\t\t\t\t\tDeveloped with :smiling_imp:By Mr. Nobody")
	fmt.Printf("%s\n", Red(devMessage))
	fmt.Println()
	fmt.Println()

	// valid := false
	// mainMenuOption := "1"

	// for ok := true; ok; ok = !valid {
	// 	printMainMenu()
	// 	mainMenuOption = readFromKeyBoard("Select any option: ")
	// 	valid = mainMenuOption == "1"
	// 	if !valid {
	// 		fmt.Printf("\n%s\n", Red("Invalid option"))
	// 	}
	// }
	// fmt.Println("You select " + mainMenuOption)
}
