package command

import (
	"log"
	"net"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	cos_pb "github.com/BinacsLee/server/api/cos"
	crypto_pb "github.com/BinacsLee/server/api/crypto"
	pastebin_pb "github.com/BinacsLee/server/api/pastebin"
	tinyurl_pb "github.com/BinacsLee/server/api/tinyurl"
	user_pb "github.com/BinacsLee/server/api/user"

	"github.com/BinacsLee/cli/service"
	"github.com/BinacsLee/cli/util"
)

var instance, domain, port string

var (
	StartCmd = &cobra.Command{
		Use:   "start",
		Short: "Start Command",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			setInstance()
			if err := dailAndServe(); err != nil {
				return err
			}
			return nil
		},
	}
)

func init() {
	startCmdFlags(StartCmd)
}

func startCmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&instance, "instance", "", "Local instance name")
	cmd.PersistentFlags().StringVar(&domain, "domain", "api.binacs.cn", "API domain such as api.binacs.cn")
	cmd.PersistentFlags().StringVar(&port, "port", ":443", "API port such as :443")
}

func setInstance() {
	if len(instance) == 0 {
		hostname, err := os.Hostname()
		if err != nil {
			log.Printf("os.Hostname get err: %+v", err)
			instance = "defaultInstanceName"
		} else {
			instance = hostname
		}
	}
	log.Printf("instance = %s domain = %s port = %s", instance, domain, port)
}

func dailAndServe() error {
	// Dail API server
	conn, err := grpc.Dial(domain+port,
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(util.GetCertPool(), domain)),
		grpc.WithPerRPCCredentials(util.GetToken(instance)),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Listen to unix socket
	_ = os.Remove(util.GetSockPath())
	sockAddr, err := net.ResolveUnixAddr("unix", util.GetSockPath())
	if err != nil {
		return err
	}
	lis, err := net.ListenUnix("unix", sockAddr)
	if err != nil {
		return err
	}

	// Serve
	s := grpc.NewServer()
	node := service.InitService(conn)

	cos_pb.RegisterCosServer(s, node.Cos.(cos_pb.CosServer))
	crypto_pb.RegisterCryptoServer(s, node.Crypto.(crypto_pb.CryptoServer))
	pastebin_pb.RegisterPastebinServer(s, node.Pastebin.(pastebin_pb.PastebinServer))
	tinyurl_pb.RegisterTinyURLServer(s, node.TinyURL.(tinyurl_pb.TinyURLServer))
	user_pb.RegisterUserServer(s, node.User.(user_pb.UserServer))

	if err := s.Serve(lis); err != nil {
		return err
	}
	return nil
}
