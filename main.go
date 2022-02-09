package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var (
	debugModeFlag bool

	serviceAccountIdFlag string
	keyIdFlag            string
	privateKeyPathFlag   string

	certIdFlag string
)

func init() {
	flag.BoolVar(&debugModeFlag, "debug", false, "debug mode")

	flag.StringVar(&serviceAccountIdFlag, "sa-id", "", "service account ID")
	flag.StringVar(&keyIdFlag, "key-id", "", "key ID")
	flag.StringVar(&privateKeyPathFlag, "privkey", "./private.pem", "private key file path (https://cloud.yandex.ru/docs/iam/operations/authorized-key/create)")

	flag.StringVar(&certIdFlag, "cert-id", "", "certificate id")
}

func main() {
	flag.Parse()

	privateKey, err := loadPrivateKey(privateKeyPathFlag)
	if err != nil {
		fmt.Printf("ERR load private key: %s\n", err)
		os.Exit(1)
		return
	}
	if debugModeFlag {
		fmt.Printf("DBG private key loaded")
	}

	jwToken, err := genJwt(serviceAccountIdFlag, keyIdFlag, privateKey)
	if err != nil {
		fmt.Printf("ERR generate JWT: %s\n", err)
		os.Exit(1)
		return
	}
	if debugModeFlag {
		fmt.Printf("DBG jwt generated: %s", jwToken)
	}

	iamToken, err := genIam(jwToken)
	if err != nil {
		fmt.Printf("ERR generate IAM: %s\n", err)
		os.Exit(1)
		return
	}
	if debugModeFlag {
		fmt.Printf("DBG iam token generated: %s", iamToken)
	}

	cert, err := getCertificate(iamToken, certIdFlag)
	if err != nil {
		fmt.Printf("ERR generate IAM: %s\n", err)
		os.Exit(1)
		return
	}

	fullchainFilePath := fmt.Sprintf("./%s_fullchain.pem", certIdFlag)
	err = ioutil.WriteFile(fullchainFilePath, []byte(strings.Join(cert.CertificateChain, "")), 0600)
	if err != nil {
		fmt.Printf("ERR write fullchain: %s\n", err)
		os.Exit(1)
		return
	}
	fmt.Printf("full chain file was created in %s\n", fullchainFilePath)

	privkeyFilePath := fmt.Sprintf("./%s_privkey.pem", certIdFlag)
	err = ioutil.WriteFile(privkeyFilePath, []byte(cert.PrivateKey), 0600)
	if err != nil {
		fmt.Printf("ERR write fullchain: %s\n", err)
		os.Exit(1)
		return
	}
	fmt.Printf("private key file was created in %s\n", privkeyFilePath)
}

type Certificate struct {
	CertificateId    string   `json:"certificateId"`
	CertificateChain []string `json:"certificateChain"`
	PrivateKey       string   `json:"privateKey"`
}

func getCertificate(token, certificateId string) (*Certificate, error) {
	url := fmt.Sprintf("https://data.certificate-manager.api.cloud.yandex.net/certificate-manager/v1/certificates/%s:getContent", certificateId)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("construct http request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	body, err := httpDo(req)
	if err != nil {
		return nil, err
	}

	var v Certificate
	err = json.Unmarshal(body, &v)
	if err != nil {
		return nil, fmt.Errorf("unmarshal body: %w", err)
	}

	return &v, nil
}
