package util

import (
	"os/user"
)

const (
	SockDir  = "/var/run"
	SockName = "cli.sock"
)

func GetSockPath() string {
	user, err := user.Current()
	if err != nil {
		return SockDir + "/" + SockName
	}
	return user.HomeDir + "/" + SockName
}
