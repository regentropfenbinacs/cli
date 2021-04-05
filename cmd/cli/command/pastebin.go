package command

import (
	"context"
	"log"
	"strings"

	"github.com/spf13/cobra"

	pastebin_pb "github.com/BinacsLee/server/api/pastebin"
)

var (
	PastebinCmd = &cobra.Command{
		Use:   "pastebin",
		Short: "PasteBin Command:\t Just run `cli pastebin submit sth.(file)`",
		Run: func(cmd *cobra.Command, args []string) {
			if !checkArgs(args, 2, 2) {
				return
			}

			op, file := parsePastebinAgrs(args)

			switch op {
			case "submit":
				file, data := processReadFile(file)
				if len(file) != 0 {
					handleResp(node.Pastebin.PastebinSubmit(context.Background(), &pastebin_pb.PastebinSubmitReq{
						Text: string(data),
					}))
				} else {
					log.Printf(errorReadFile, file, data)
				}
			default:
				log.Printf(errorOpInvalid)
			}
		},
	}
)

func parsePastebinAgrs(args []string) (op, file string) {
	op = strings.ToLower(args[0])
	file = args[1]
	return
}
