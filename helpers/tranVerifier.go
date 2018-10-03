package helpers

import (
	"encoding/hex"
	"encoding/json"
	"math"
	"net"
	secp "secp256k1-go"
	"strconv"
	"time"
)

func PubkeyFromSeckey(privateKey []byte) []byte {
	//

	return secp.PubkeyFromSeckey([]byte(privateKey))
}

func GetRawTransactionType(rawData string) string {
	//

	var tranType TranType
	err := json.Unmarshal([]byte(rawData), &tranType)
	if err != nil {
		//

		return ""

	} else {
		//

		return tranType.TT
	}
}

func ParseHelloTransaction(rawData string) (HelloTransaction, error) {
	//

	var helloTransaction HelloTransaction
	err := json.Unmarshal([]byte(rawData), &helloTransaction)
	return helloTransaction, err
}

func VerifyHelloTransaction(tran HelloTransaction) (ok bool) {
	//

	ok = false

	if len(tran.SENDER) != 66 {
		//

		return
	}

	if net.ParseIP(tran.ADDRESS) == nil {
		//

		return
	}

	if len(tran.TST) != 10 {
		//

		return
	}

	transactionTime, err := strconv.ParseInt(tran.TST, 10, 64)
	if err != nil {
		//

		return
	}

	timestamp := time.Unix(transactionTime, 0)
	if int64(math.Abs(float64(time.Since(timestamp)/time.Second))) > 10 {
		//

		return
	}

	if len(tran.SIGNATURE) != 130 {
		//

		return
	}

	transcationForVerify := HelloTransactionForVerify{tran.TT, tran.SENDER, tran.ADDRESS, tran.TST}
	js, err := json.Marshal(transcationForVerify)
	if err != nil {
		//

		return
	}

	decodedSignature, err := hex.DecodeString(tran.SIGNATURE)
	if err != nil {
		//

		return
	}

	publicKey, err := hex.DecodeString(tran.SENDER)
	if err != nil {
		//

		return
	}

	if secp.VerifySignature(js, decodedSignature, publicKey) != 1 {
		//

		return
	}

	ok = true

	return
}

func ParseSimpleTransaction(rawData string) (SimpleTransaction, error) {
	//

	var simpleTransaction SimpleTransaction
	err := json.Unmarshal([]byte(rawData), &simpleTransaction)
	return simpleTransaction, err
}

func VerifySimpleTransaction(tran SimpleTransaction) (transactionTime int64, ok bool) {
	//

	transactionTime = 0
	ok = false

	if len(tran.SENDER) != 66 {
		//

		return
	}

	if len(tran.RECEIVER) != 66 {
		//

		return
	}

	if tran.SENDER == tran.RECEIVER {
		//

		return
	}

	if len(tran.TST) != 10 {
		//

		return
	}

	transactionTime, err := strconv.ParseInt(tran.TST, 10, 64)
	if err != nil {
		//

		return
	}

	timestamp := time.Unix(transactionTime, 0)
	if int64(math.Abs(float64(time.Since(timestamp)/time.Second))) > 10 {
		//

		return
	}

	if len(tran.TTOKEN) == 0 {
		//

		return
	}

	if len(tran.CTOKEN) == 0 {
		//

		return

	} else {
		//

		_, err = strconv.ParseFloat(tran.CTOKEN, 64)
		if err != nil {
			//

			return
		}
	}

	if len(tran.SIGNATURE) != 130 {
		//

		return
	}

	transcationForVerify := SimpleTransactionForVerify{tran.TT, tran.SENDER, tran.RECEIVER, tran.TTOKEN, tran.CTOKEN, tran.TST}
	js, err := json.Marshal(transcationForVerify)
	if err != nil {
		//

		return
	}

	decodedSignature, err := hex.DecodeString(tran.SIGNATURE)
	if err != nil {
		//

		return
	}

	publicKey, err := hex.DecodeString(tran.SENDER)
	if err != nil {
		//

		return
	}

	if secp.VerifySignature(js, decodedSignature, publicKey) != 1 {
		//

		return
	}

	ok = true

	return
}

func CreateHelloTransaction(publicKey string, privateKey []byte, ip string) (tran string, ok bool) {
	//

	tran = ""
	ok = false

	helloTransactionForVerify := HelloTransactionForVerify{"HL", publicKey, ip, strconv.FormatInt(time.Now().Unix(), 10)}
	js, err := json.Marshal(helloTransactionForVerify)
	if err != nil {
		//

		return
	}

	signature := secp.Sign(js, []byte(privateKey))
	if len(signature) == 0 {
		//

		return
	}

	helloTransaction := HelloTransaction{"HL", publicKey, ip, helloTransactionForVerify.TST, string(signature)}
	js, err = json.Marshal(helloTransaction)
	if err != nil {
		//

		return
	}

	tran = string(js)
	ok = true

	return
}
