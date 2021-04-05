package command

import (
	"context"
	"log"
	"strings"

	"github.com/spf13/cobra"

	cos_pb "github.com/BinacsLee/server/api/cos"
)

var (
	CosCmd = &cobra.Command{
		Use:   "cos",
		Short: "Cos Command:\t Just run `cli cos put/get sth.(file)`",
		Run: func(cmd *cobra.Command, args []string) {
			if !checkArgs(args, 2, 2) {
				return
			}

			op, arg := parseCosAgrs(args)

			switch op {
			case "put":
				file, data := processReadFile(arg)
				if len(file) != 0 {
					handleResp(node.Cos.CosPut(context.Background(), &cos_pb.CosPutReq{
						FileName:  file,
						FileBytes: data,
					}))
				} else {
					log.Printf(errorReadFile, file, data)
				}
			case "get":
				log.Printf("Error: Not support `get`.\n")
				// handleResp(node.Cos.CosGet(context.Background(), &cos_pb.CosGetReq{
				// 	CosURI: arg,
				// }))
			default:
				log.Printf(errorOpInvalid)
			}
		},
	}
)

func parseCosAgrs(args []string) (op, arg string) {
	op = strings.ToLower(args[0])
	arg = args[1]
	return
}
