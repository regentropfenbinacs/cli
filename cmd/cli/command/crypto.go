package command

import (
	"context"
	"log"
	"strings"

	"github.com/spf13/cobra"

	crypto_pb "github.com/BinacsLee/server/api/crypto"
)

var (
	CryptoCmd = &cobra.Command{
		Use:   "crypto",
		Short: "Crypto Command:\t Just run `cli crypto encrypt/decrypt BASE64/AES/DES sth.(string)`",
		Run: func(cmd *cobra.Command, args []string) {
			if !checkArgs(args, 3, 0) {
				return
			}

			op, algo, text := parseCryptoAgrs(args)

			switch op {
			case "encrypt":
				handleResp(node.Crypto.CryptoEncrypt(context.Background(), &crypto_pb.CryptoEncryptReq{
					Algorithm: algo,
					PlainText: text,
				}))
			case "decrypt":
				handleResp(node.Crypto.CryptoDecrypt(context.Background(), &crypto_pb.CryptoDecryptReq{
					Algorithm:   algo,
					EncryptText: text,
				}))
			default:
				log.Printf(errorOpInvalid)
			}
		},
	}
)

func parseCryptoAgrs(args []string) (op, algo, text string) {
	op = strings.ToLower(args[0])
	algo = strings.ToUpper(args[1])
	text = args[2]
	args = args[3:]
	for _, s := range args {
		text += " " + s
	}
	return
}
