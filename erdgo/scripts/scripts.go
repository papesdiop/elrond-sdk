package main

import (
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	_ "go/token"
	"io/ioutil"
	"strings"

	elrondGoCore "github.com/ElrondNetwork/elrond-go/core"
	_ "github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/elrond-sdk/erdgo"
	"github.com/ElrondNetwork/elrond-sdk/erdgo/blockchain"
	"github.com/ElrondNetwork/elrond-sdk/erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk/erdgo/interactors"
)

const SC_DELEGATION = ""

func main() {

	const URL string = "http://158.175.191.253:8080"
	const delegationSC string = "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqhllllsajxzat"
	//sendFromMinterToWallets(URL,"/home/andrei/testing/owner.pem", "/home/andrei/100Kdelegators.pem")
	sendFromWalletsToDelegationSCRustAddNodes(URL, "/home/andrei/ownerTestnet.pem", "/home/andrei/allInOneValidatorsKeys.pem", delegationSC)

}

func sendFromWalletsToDelegationSCRustAddNodes(URL string, ownerDelegSCPemPath string, blsKeysPemPath string, delegationSC string) {
	proxy := blockchain.NewElrondProxy(URL)
	txSigner := blockchain.NewTxSigner()
	txInteractor, err := interactors.NewTransactionInteractor(proxy, txSigner)
	if err != nil {
		fmt.Println(err)
		return
	}
	ownerSk, err := erdgo.LoadPrivateKeyFromPemFile(ownerDelegSCPemPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	blsKeysSK, blsKeysPk, err := readAllSkPkFromPem(blsKeysPemPath, 2831)
	if err != nil {
		fmt.Println(err)
		return
	}

	publicKey, err := erdgo.GetAddressFromPrivateKey(ownerSk)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(publicKey)

	addressHandler, err := data.NewAddressFromBech32String(publicKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	account, err := proxy.GetAccount(addressHandler)
	if err != nil {
		fmt.Println(err)
		return
	}

	txArgs := interactors.ArgCreateTransaction{
		Nonce:     account.Nonce,
		Value:     "2500000000000000000000",
		RcvAddr:   delegationSC,
		SndAddr:   account.Address,
		GasPrice:  1000000000,
		GasLimit:  50000000,
		Data:      []byte(""),
		Signature: "",
		ChainID:   "1",
		Version:   1,
		Options:   0,
	}

	for index := 0; index < len(blsKeysPk); index++ {
		blsSignature, err := txSigner.SignMessage(addressHandler.AddressBytes(), blsKeysSK[index])
		if err != nil {
			fmt.Println(err)
			return
		}
		inputData := "addNodes@" + string(blsKeysPk[index]) + "@" + hex.EncodeToString(blsSignature)
		txArgs.Data = []byte(inputData)
		txArgsSigner, err := txInteractor.ApplySignatureAndSender(ownerSk, txArgs)
		if err != nil {
			fmt.Println(err)
			return
		}
		tx := txInteractor.CreateTransaction(txArgsSigner)
		txInteractor.AddTransaction(tx)
		txArgs.Nonce = txArgs.Nonce + 1

	}

	msgs, err := txInteractor.SendTransactionsAsBunch(100)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, msg := range msgs {
		fmt.Println(msg)
	}

}

func sendFromMinterToWallets(URL string, minterPemPath string, walletsPemPath string) {
	proxy := blockchain.NewElrondProxy(URL)
	txSigner := blockchain.NewTxSigner()
	txInteractor, err := interactors.NewTransactionInteractor(proxy, txSigner)
	if err != nil {
		fmt.Println(err)
		return
	}
	ownerSK, err := erdgo.LoadPrivateKeyFromPemFile(minterPemPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	publicKey, err := erdgo.GetAddressFromPrivateKey(ownerSK)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(publicKey)

	addressHandler, err := data.NewAddressFromBech32String(publicKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	account, err := proxy.GetAccount(addressHandler)
	if err != nil {
		fmt.Println(err)
		return
	}

	txArgs := interactors.ArgCreateTransaction{
		Nonce:     account.Nonce,
		Value:     "5000000000000000000",
		RcvAddr:   "",
		SndAddr:   "",
		GasPrice:  1000000000,
		GasLimit:  55000,
		Data:      []byte(""),
		Signature: "",
		ChainID:   "1",
		Version:   1,
		Options:   0,
	}

	_, receiversKeys, err := readAllSkPkFromPem(walletsPemPath, 100000)
	if err != nil {
		return
	}

	for index := 0; index < len(receiversKeys); index++ {
		txArgs.RcvAddr = receiversKeys[index]
		txArgsSigned, err := txInteractor.ApplySignatureAndSender(ownerSK, txArgs)
		if err != nil {
			fmt.Println(err)
			return
		}

		txInteractor.AddTransaction(txInteractor.CreateTransaction(txArgsSigned))
		txArgs.Nonce = txArgs.Nonce + 1
	}

	msgs, err := txInteractor.SendTransactionsAsBunch(50)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, msg := range msgs {
		fmt.Println(msg)
	}
}

func sendFromWalletsToSystemDelegationDelegate(URL string, walletsPemPath string, delegationSC string) {
	proxy := blockchain.NewElrondProxy(URL)
	txSigner := blockchain.NewTxSigner()
	txInteractor, err := interactors.NewTransactionInteractor(proxy, txSigner)
	if err != nil {
		fmt.Println(err)
		return
	}

	txArgs := interactors.ArgCreateTransaction{
		Nonce:     1,
		Value:     "0", //"1000000000000000000",
		RcvAddr:   delegationSC,
		SndAddr:   "",
		GasPrice:  1000000000,
		GasLimit:  12000000,
		Data:      []byte("unDelegate@0DE0B6B3A7640000"),
		Signature: "",
		ChainID:   "1",
		Version:   1,
		Options:   0,
	}

	senderSk, senderPk, err := readAllSkPkFromPem(walletsPemPath, 100000)
	if err != nil {
		return
	}

	for index := 0; index < len(senderPk); index++ {
		txArgs.SndAddr = senderPk[index]
		txArgsSigned, err := txInteractor.ApplySignatureAndSender(senderSk[index], txArgs)
		if err != nil {
			fmt.Println(err)
			return
		}

		txInteractor.AddTransaction(txInteractor.CreateTransaction(txArgsSigned))

	}

	msgs, err := txInteractor.SendTransactionsAsBunch(50)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, msg := range msgs {
		fmt.Println(msg)
	}
}

func readAllSkPkFromPem(pemPath string, keysToRead int) ([][]byte, []string, error) {
	file, err := elrondGoCore.OpenFile(pemPath)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	buff, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	if len(buff) == 0 {
		fmt.Println("Empty File")
		return nil, nil, errors.New("empty file error")
	}
	var blkRecovered *pem.Block

	pkKeys := []string{}
	skKeys := [][]byte{}

	for i := 0; i < keysToRead; i++ {
		if len(buff) == 0 {
			//less private pkKeys present in the file than required
			return nil, nil, errors.New("Nil")
		}

		blkRecovered, buff = pem.Decode(buff)

		if blkRecovered == nil {
			return nil, nil, errors.New("Nil block recivered")
		}
		blockType := blkRecovered.Type
		header := "PRIVATE KEY for "
		if strings.Index(blockType, header) != 0 {
			return nil, nil, errors.New("missing header")
		}
		blockTypeString := blockType[len(header):]
		pkKeys = append(pkKeys, blockTypeString)
		skKey, err := hex.DecodeString(string(blkRecovered.Bytes))
		if err != nil {
			fmt.Println(err)
			return nil, nil, err
		}
		skKeys = append(skKeys, skKey)
	}
	return skKeys, pkKeys, nil

}
