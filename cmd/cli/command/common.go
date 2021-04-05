package command

import (
	"io/ioutil"
	"log"
	"os"
)

const (
	errorCmdArgsEmpty         = "Error: Command args empty.\n"
	errorCmdArgsLengthInvalid = "Error: Command args length invalid, low = %+v, high = %+v, length = %+v\n"
	errorOpInvalid            = "Error: Operation invalid.\n"
	errorArgsLengthInvalid    = "Error: Args length invalid, theshould = %+v, length = %+v\n"
	errorReadFile             = "Error: os.Stat got: %+v, %+v\n"
	printResp                 = "Got: resp = %+v, err = %+v\n"
)

func checkArgs(args []string, low, high int) bool {
	length := len(args)
	if length == 0 {
		log.Printf(errorCmdArgsEmpty)
		return false
	}
	if low > 0 && length < low || high > 0 && length > high {
		log.Printf(errorCmdArgsLengthInvalid, low, high, length)
		return false
	}
	return true
}

func handleResp(resp interface{}, err error) {
	log.Printf(printResp, resp, err)
}

func processReadFile(file string) (string, []byte) {
	info, err := os.Stat(file)
	if err != nil || info.IsDir() {
		return err.Error(), nil
	}
	fileBytes, _ := ioutil.ReadFile(file)
	return info.Name(), fileBytes
}
