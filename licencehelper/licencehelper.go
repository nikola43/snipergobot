package licencehelper

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/nikola43/snipergobot/menuhelper"
	"github.com/nikola43/snipergobot/sysinfo"

	"github.com/mdp/qrterminal"
)

func SaveLicense(plainData string, filename string) {
	//rand.Read(encryptionKey)
	//fmt.Println(encryptionKey)
	//encryption(plainData, filename)
}

func CheckLicense() {
	fileName := "licence.dat"
	if !FileExist(fileName) {
		fmt.Printf("%s:", menuhelper.Cyan("Lincense file not found"))

		now := time.Now()
		licenceEndDate := now.AddDate(0, 0, 5)

		fmt.Println("now", ParseDateTime(now))
		fmt.Println("licenceEndDate", ParseDateTime(licenceEndDate))

		info := sysinfo.NewSysInfo()
		fmt.Printf("%+s\n", info.ToHash())
		SaveLicense(info.ToHash(), fileName)

	} else {
		fmt.Printf("%v file exist\n", fileName)
		//mt.Printf("Decrypted Msg : %s", decryption(fileName))

		isLicenceValid := true
		if isLicenceValid {
			fmt.Println()
			ShowPaymentQr()
			fmt.Println()
			fmt.Printf("%s: %s\n", menuhelper.Cyan("Send 1 ETH To"), menuhelper.Yellow("0x6d5F00aE01F715D3082Ad40dfB5c18A1a35d3A17"))
			fmt.Printf("%s\n", menuhelper.Cyan("You will receive email with API KEY after our system process payment"))
			fmt.Println()
		}
	}
}

func ShowPaymentQr() {
	const Red = "\033[44m  \033[0m"
	const BLUE = "\033[43m  \033[0m"

	config := qrterminal.Config{
		Level:     qrterminal.M,
		Writer:    os.Stdout,
		BlackChar: Red,
		WhiteChar: BLUE,
		QuietZone: 1,
	}
	qrterminal.GenerateWithConfig("0x6d5F00aE01F715D3082Ad40dfB5c18A1a35d3A17", config)
}

func FileExist(filename string) bool {
	_, error := os.Stat(filename)
	if os.IsNotExist(error) {
		return false
	} else {
		return true
	}
}

func ParseDateTime(now time.Time) string {
	return strconv.Itoa(now.Year()) + "/" + now.Month().String() + "/" + strconv.Itoa(now.Day()) + " " + strconv.Itoa(now.Hour()) + ":" + strconv.Itoa(now.Minute()) + ":" + strconv.Itoa(now.Second()) + ":" + strconv.Itoa(now.Nanosecond())
}
