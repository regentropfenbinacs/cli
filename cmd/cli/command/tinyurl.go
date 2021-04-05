package command

import (
	"context"
	"log"
	"strings"

	"github.com/spf13/cobra"

	tinyurl_pb "github.com/BinacsLee/server/api/tinyurl"
)

var (
	TinyurlCmd = &cobra.Command{
		Use:   "tinyurl",
		Short: "TinyURL Command:\t Just run `cli tinyurl encode/decode sth.`",
		Run: func(cmd *cobra.Command, args []string) {
			if !checkArgs(args, 2, 2) {
				return
			}

			op, url := parseTinyurlAgrs(args)

			switch op {
			case "encode":
				handleResp(node.TinyURL.TinyURLEncode(context.Background(), &tinyurl_pb.TinyURLEncodeReq{
					Url: url,
				}))
			case "decode":
				handleResp(node.TinyURL.TinyURLDecode(context.Background(), &tinyurl_pb.TinyURLDecodeReq{
					Turl: url,
				}))
			default:
				log.Printf(errorOpInvalid)
			}
		},
	}
)

func parseTinyurlAgrs(args []string) (op, url string) {
	op = strings.ToLower(args[0])
	url = args[1]
	return
}
