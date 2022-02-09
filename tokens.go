package main

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	jwtExpiresDuration = time.Minute
	iamTokensUrl       = "https://iam.api.cloud.yandex.net/iam/v1/tokens"
)

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(data)
	if err != nil {
		return nil, fmt.Errorf("parse rsa private key from pem: %w", err)
	}

	return privateKey, nil
}

func genJwt(saId, keyId string, privateKey *rsa.PrivateKey) (string, error) {
	now := time.Now()

	claims := jwt.RegisteredClaims{
		Issuer:    saId,
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(jwtExpiresDuration)),
		Audience:  []string{iamTokensUrl},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodPS256, claims)
	token.Header["kid"] = keyId

	signed, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signed, nil
}

type iamV1Tokens struct {
	IamToken string `json:"iamToken"`
}

func genIam(jwToken string) (string, error) {
	req, err := http.NewRequest(http.MethodPost, iamTokensUrl, strings.NewReader(`{"jwt":"`+jwToken+`"}`))
	if err != nil {
		return "", fmt.Errorf("construct http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	body, err := httpDo(req)
	if err != nil {
		return "", err
	}

	var v iamV1Tokens
	err = json.Unmarshal(body, &v)
	if err != nil {
		return "", fmt.Errorf("unmarshal body: %w", err)
	}

	return v.IamToken, nil
}
