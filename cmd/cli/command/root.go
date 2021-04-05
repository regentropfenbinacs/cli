package command

import (
	"net"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/BinacsLee/cli/service"
	"github.com/BinacsLee/cli/util"
)

var node *service.NodeServiceImpl

func unixConnect(addr string, t time.Duration) (net.Conn, error) {
	unix_addr, err := net.ResolveUnixAddr("unix", util.GetSockPath())
	conn, err := net.DialUnix("unix", nil, unix_addr)
	return conn, err
}

var (
	RootCmd = &cobra.Command{
		Use:   "root",
		Short: "Terminal client for https://binacs.cn\nMore at https://github.com/BinacsLee/cli",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			conn, err := grpc.Dial(util.GetSockPath(),
				grpc.WithBlock(),
				grpc.WithInsecure(),
				grpc.WithDialer(unixConnect),
				grpc.FailOnNonTempDialError(true),
			)
			if err != nil {
				return err
			}
			node = service.InitService(conn)
			return nil
		},
	}
)
