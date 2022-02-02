package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"

	logger "github.com/sirupsen/logrus"
)

type Wallet struct {
	PublicKey  string
	PrivateKey string
}

type Transaction struct {
	Sender   string
	Receiver string
	Amount   int
	Fee      int
	Message  string
}

func generateWallet() *Wallet {
	var privateKey *rsa.PrivateKey
	var derPkix []byte
	var err error
	if privateKey, err = rsa.GenerateKey(rand.Reader, 1024); err != nil {
		logger.Errorln(err)
	}
	derPkcs := x509.MarshalPKCS1PrivateKey(privateKey)
	publicKey := &privateKey.PublicKey
	if derPkix, err = x509.MarshalPKIXPublicKey(publicKey); err != nil {
		logger.Errorln(err)
	}
	wallet := new(Wallet)
	wallet.PrivateKey = base64.StdEncoding.EncodeToString(derPkcs)
	wallet.PublicKey = base64.StdEncoding.EncodeToString(derPkix)
	return wallet
}

func signTransaction(t *Transaction, privateKey string) ([]byte, error) {
	var bytes []byte
	var pvk *rsa.PrivateKey
	var signature []byte
	var err error = nil
	if bytes, err = base64.StdEncoding.DecodeString(privateKey); err != nil {
		logger.Errorln(err)
		return nil, err
	}
	if pvk, err = x509.ParsePKCS1PrivateKey(bytes); err != nil {
		logger.Errorln(err)
		return nil, err
	}
	tStr, _ := json.Marshal(t)
	hashed := sha256.Sum256(tStr)
	if signature, err = rsa.SignPKCS1v15(rand.Reader, pvk, crypto.SHA256, hashed[:]); err != nil {
		logger.Errorln(err)
		return nil, err
	}
	return signature, err
}
