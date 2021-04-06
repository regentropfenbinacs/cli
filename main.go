package main

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"os"
	// "bytes"
	// "io"
	// "io/ioutil"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	cos_pb "github.com/BinacsLee/server/api/cos"
	crypto_pb "github.com/BinacsLee/server/api/crypto"
	pastebin_pb "github.com/BinacsLee/server/api/pastebin"
	tinyurl_pb "github.com/BinacsLee/server/api/tinyurl"
	user_pb "github.com/BinacsLee/server/api/user"
)

var instance, endpoint string

var crtBytes = []byte(`-----BEGIN CERTIFICATE-----
MIIFijCCBHKgAwIBAgIQA+NpxZV2SO5BH3dBMxSFyDANBgkqhkiG9w0BAQsFADBy
MQswCQYDVQQGEwJDTjElMCMGA1UEChMcVHJ1c3RBc2lhIFRlY2hub2xvZ2llcywg
SW5jLjEdMBsGA1UECxMURG9tYWluIFZhbGlkYXRlZCBTU0wxHTAbBgNVBAMTFFRy
dXN0QXNpYSBUTFMgUlNBIENBMB4XDTIxMDEzMTAwMDAwMFoXDTIyMDEzMDIzNTk1
OVowGDEWMBQGA1UEAxMNYXBpLmJpbmFjcy5jbjCCASIwDQYJKoZIhvcNAQEBBQAD
ggEPADCCAQoCggEBAMdeFuoHEyKhbNMwKCRlnAxWQdOiWeGl/RG4HE3wWsmq4ePC
BohTaupWMTkWJaL+gCy+6X3mDkUlE4tJlFMIj6dfUmjAnf0RmrDH4DOb9XgLcQC6
ss5czpWj7EEAW7yDoHuLDaEaw5+V4uiyRzR49VyWUc1d1UY/cZyschwzI5/PxG01
b8uTtZ1s/BhSN2ruRggaiLi+cUC+W4Mh6w0hfMQlAt+DITS1HOUmVCi+EQ5ThKDi
o3svZ+LCkypPaGaxzX+WrjvqowRF8gn8F+jN1pk+o0GKFROxT7k58Kr5iPclbrYO
VCiKQtkengH7V17iiiP0HoVXLD+yzRfAb39OjjECAwEAAaOCAnQwggJwMB8GA1Ud
IwQYMBaAFH/TmfOgRw4xAFZWIo63zJ7dygGKMB0GA1UdDgQWBBQ0o5bv3sSGKxth
pjmusSnsdf3t9TAYBgNVHREEETAPgg1hcGkuYmluYWNzLmNuMA4GA1UdDwEB/wQE
AwIFoDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwPgYDVR0gBDcwNTAz
BgZngQwBAgEwKTAnBggrBgEFBQcCARYbaHR0cDovL3d3dy5kaWdpY2VydC5jb20v
Q1BTMIGSBggrBgEFBQcBAQSBhTCBgjA0BggrBgEFBQcwAYYoaHR0cDovL3N0YXR1
c2UuZGlnaXRhbGNlcnR2YWxpZGF0aW9uLmNvbTBKBggrBgEFBQcwAoY+aHR0cDov
L2NhY2VydHMuZGlnaXRhbGNlcnR2YWxpZGF0aW9uLmNvbS9UcnVzdEFzaWFUTFNS
U0FDQS5jcnQwCQYDVR0TBAIwADCCAQMGCisGAQQB1nkCBAIEgfQEgfEA7wB1ACl5
vvCeOTkh8FZzn2Old+W+V32cYAr4+U1dJlwlXceEAAABd1aP/wEAAAQDAEYwRAIg
NZzRICxxxd9+XwqlTyM9RLvRYX80IqZY5leYJv0SyO4CIBPeYuhSv4UbukqFtfGw
LOSkyL8/n92GJVAHvIpOuM0WAHYAIkVFB1lVJFaWP6Ev8fdthuAjJmOtwEt/XcaD
XG7iDwIAAAF3Vo//QwAABAMARzBFAiAiLzgSww0dpks+pE/XLjxQQYwXEusBQaK0
g8nV4IADHgIhAMVd1co3pMxNEKmXuHe8Dpim85eHbgVkMj75pqb2HsXjMA0GCSqG
SIb3DQEBCwUAA4IBAQCD+nBU8+V3760IdAJ9fTF0PJxCVICtU7eBcwxE1qFn4LyS
74DoK/1RQ80/6eicjPTF94motbZ9OjmnAwhPinQK8ej0yhyWfLjIvZAjGQr3ESFX
uJPs1ZK5R4own0fIA1hLTpd4e1Bx42mXCsFgdPnTNGfaOciqOeAx9IwneXPFYwOE
pSythIBcBggMDh+JP9blpqssQhvBSbUNDHu9pHiEFxyvSQA1hhoLLmUXSU0dJwNz
A7Q6lLP/DxRDK14O2T3YCnEig/QuvUhF3gixR5CN2dy0ZVNTtZ05mDkGiRr0FDCL
NBPuw6GVAFUXNmGWNO1XCGh/McyfD7jy+vpUQk+8
-----END CERTIFICATE-----
`)

// var serverName = "server.grpc.io"
var serverName = "api.binacs.cn"

type Token struct {
	Value string
}

const headerAuthorize string = "authorization"

func (t *Token) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{headerAuthorize: t.Value}, nil
}

func (t *Token) RequireTransportSecurity() bool {
	return true
}

func init() {
	flag.StringVar(&instance, "instance", "", "Instance name")
	flag.StringVar(&endpoint, "endpoint", "api.binacs.cn:443", "API endpoints")
	// flag.StringVar(&endpoint, "endpoint", "server.grpc.io:9500", "API endpoints")
}

func main() {
	flag.Parse()

	if len(instance) == 0 {
		hostname, err := os.Hostname()
		if err != nil {
			fmt.Println("os.Hostname err:", err)
			instance = "defaultInstanceName"
		} else {
			instance = hostname
		}
	}

	log.Println("instance =", instance, "endpoint =", endpoint)

	certDERBlock, _ := pem.Decode(crtBytes)
	if certDERBlock == nil {
		fmt.Println("certDERBlock nil")
	}

	crt, err := x509.ParseCertificate(certDERBlock.Bytes)
	if err != nil {
		fmt.Println("x509.ParseCertificate certDERBlock err:", err)
	}

	cp := x509.NewCertPool()
	cp.AddCert(crt)

	conn, err := grpc.Dial(endpoint,
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(cp, serverName)),
		grpc.WithPerRPCCredentials(&Token{
			Value: "bearer token.server.grpc.io",
		}),
	)
	if err != nil {
		fmt.Println("grpc.Dial get error:", err)
	}
	defer conn.Close()

	// userTest(conn)
	// cryptoTest(conn)
	// pastebinTest(conn)
	// tinyurlTest(conn)
	cosTest(conn)
}

func userTest(conn *grpc.ClientConn) {
	userCli := user_pb.NewUserClient(conn)

	{
		resp, err := userCli.UserTest(context.Background(), &user_pb.UserTestReq{})
		if err != nil {
			fmt.Println("userCli.UserTest get error:", err)
		}
		fmt.Println("resp =", resp)
	}

	{
		resp, err := userCli.UserRegister(context.Background(), &user_pb.UserRegisterReq{
			Id:  "id",
			Pwd: "pwd",
		})
		if err != nil {
			fmt.Println("userCli.UserTest get error:", err)
		}
		fmt.Println("resp =", resp)
	}

	{
		resp, err := userCli.UserAuth(context.Background(), &user_pb.UserAuthReq{
			Id:  "id",
			Pwd: "pwd",
		})
		if err != nil {
			fmt.Println("userCli.UserTest get error:", err)
		}
		fmt.Println("resp =", resp)
	}
}

func cryptoTest(conn *grpc.ClientConn) {
	cryptoCli := crypto_pb.NewCryptoClient(conn)

	var t string
	{
		resp, err := cryptoCli.CryptoEncrypt(context.Background(), &crypto_pb.CryptoEncryptReq{
			Algorithm: "BASE64",
			PlainText: "PlainText",
		})
		if err != nil {
			fmt.Println("userCli.UserTest get error:", err)
		}
		fmt.Println("resp =", resp)
		t = resp.Data.EncryptText
	}

	fmt.Println("t = ", t)
	{
		resp, err := cryptoCli.CryptoDecrypt(context.Background(), &crypto_pb.CryptoDecryptReq{
			Algorithm:   "BASE64",
			EncryptText: t,
		})
		if err != nil {
			fmt.Println("userCli.UserTest get error:", err)
		}
		fmt.Println("resp =", resp)
	}
}

func pastebinTest(conn *grpc.ClientConn) {
	pastebinCli := pastebin_pb.NewPastebinClient(conn)

	{
		resp, err := pastebinCli.PastebinSubmit(context.Background(), &pastebin_pb.PastebinSubmitReq{
			Author: "author",
			Syntax: "syntax",
			Text:   "text",
		})
		if err != nil {
			fmt.Println("userCli.UserTest get error:", err)
		}
		fmt.Println("resp =", resp)
		// Only purl such as: Y7ZVbu
	}
}

func tinyurlTest(conn *grpc.ClientConn) {
	tinyurlCli := tinyurl_pb.NewTinyURLClient(conn)

	var t string
	{
		resp, err := tinyurlCli.TinyURLEncode(context.Background(), &tinyurl_pb.TinyURLEncodeReq{
			Url: "https://binacs.cn/toys",
		})
		if err != nil {
			fmt.Println("userCli.UserTest get error:", err)
		}
		fmt.Println("resp =", resp)
		t = resp.Data.Turl
		// 包含binacs.cn前缀
	}

	{
		resp, err := tinyurlCli.TinyURLDecode(context.Background(), &tinyurl_pb.TinyURLDecodeReq{
			Turl: t,
		})
		if err != nil {
			fmt.Println("userCli.UserTest get error:", err)
		}
		fmt.Println("resp =", resp)
	}
}

func cosTest(conn *grpc.ClientConn) {
	cosCli := cos_pb.NewCosClient(conn)

	{
		resp, err := cosCli.CosBucketURL(context.Background(), &cos_pb.CosBucketURLReq{})
		if err != nil {
			fmt.Println("userCli.UserTest get error:", err)
		}
		fmt.Println("resp =", resp)
	}

	
	// fileBytes, err := ioutil.ReadFile("aaa.md")
	// if err != nil {
	// 	fmt.Println("ioutil.ReadFile get error:", err)
	// }
	// {
	// 	resp, err := cosCli.CosPut(context.Background(), &cos_pb.CosPutReq{
	// 		FileName:  "aaaasdfasasddfasfasd.md",
	// 		FileBytes: fileBytes,
	// 		// FileBytes: buf.Bytes(),
	// 	})
	// 	if err != nil {
	// 		fmt.Println("userCli.UserTest get error:", err)
	// 	}
	// 	fmt.Println("resp =", resp)
	// }

	{
		resp, err := cosCli.CosGet(context.Background(), &cos_pb.CosGetReq{
			CosURI: "aaaasdfasasddfasfasd_2021-04-02%2009:59:40.md",
		})
		if err != nil {
			fmt.Println("userCli.UserTest get error:", err)
		}
		fmt.Println("resp =", resp)
	}
}
