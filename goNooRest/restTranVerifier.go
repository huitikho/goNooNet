package restServer

import (
	"encoding/hex"
	"encoding/json"
	"goNooNet/helpers"
	"math"
	secp "secp256k1-go"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

func verifyTransaction(ctx *fasthttp.RequestCtx) (statusCode int,
	transactionForDB []byte,
	transactionTime int64) {
	//

	args := ctx.QueryArgs()
	for errNum, v := range helpers.MandatoryTransactionFields {
		//

		if !args.Has(v) {
			//

			return errNum, transactionForDB, transactionTime
		}
	}

	ttype := string(ctx.FormValue("TT"))
	//version := string(ctx.FormValue("VERSION"))
	sender := string(ctx.FormValue("SENDER"))
	receiver := string(ctx.FormValue("RECEIVER"))
	tst := string(ctx.FormValue("TST"))
	sign := string(ctx.FormValue("SIGNATURE"))

	if len(ttype) != 2 {
		//

		return helpers.StatusWrongAttr_TT, transactionForDB, transactionTime
	}

	// if len(version) != 2 {
	// 	//
	//
	// 	return helpers.StatusWrongAttr_VERSION, transactionForDB, transactionTime
	// }

	if len(sender) != 66 {
		//

		return helpers.StatusWrongAttr_SENDER, transactionForDB, transactionTime
	}

	if len(receiver) != 66 {
		//

		return helpers.StatusWrongAttr_RECEIVER, transactionForDB, transactionTime
	}

	if sender == receiver {
		//

		return helpers.StatusDontSendYourself, transactionForDB, transactionTime
	}

	if len(sign) != 130 {
		//

		return helpers.StatusWrongAttr_Signature, transactionForDB, transactionTime
	}

	if len(tst) != 10 {
		//

		return helpers.StatusWrongAttr_TST, transactionForDB, transactionTime
	}

	transactionTime, err := strconv.ParseInt(tst, 10, 64)
	if err != nil {
		//

		return helpers.StatusWrongAttr_TST, transactionForDB, transactionTime
	}

	timestamp := time.Unix(transactionTime, 0)
	if int64(math.Abs(float64(time.Since(timestamp)/time.Second))) > 10 {
		//

		return helpers.StatusWrongAttr_TST, transactionForDB, transactionTime
	}

	switch ttype {

	case "ST":

		args := ctx.QueryArgs()
		for errNum, v := range helpers.SimpleStructureFields {
			//

			if !args.Has(v) {
				//

				return errNum, transactionForDB, transactionTime
			}
		}

		ttoken := string(ctx.FormValue("TTOKEN"))
		ctoken := string(ctx.FormValue("CTOKEN"))

		_, err = strconv.ParseFloat(ctoken, 64)
		if err != nil {
			//

			return helpers.StatusWrongAttr_CTOKEN, transactionForDB, transactionTime
		}

		transcationForVerify := helpers.SimpleTransactionForVerify{ttype, sender, receiver, ttoken, ctoken, tst}
		js, err := json.Marshal(transcationForVerify)
		if err != nil {
			//

			return helpers.StatusInternalServerError, transactionForDB, transactionTime
		}

		decodedSignature, err := hex.DecodeString(sign)
		if err != nil {
			//

			return helpers.StatusWrongAttr_Signature, transactionForDB, transactionTime
		}

		publicKey, err := hex.DecodeString(sender)
		if err != nil {
			//

			return helpers.StatusWrongAttr_SENDER, transactionForDB, transactionTime
		}

		if secp.VerifySignature(js, decodedSignature, publicKey) != 1 {
			//

			return helpers.StatusSignVerifyError, transactionForDB, transactionTime
		}

		transcation := helpers.SimpleTransaction{ttype, sender, receiver, ttoken, ctoken, tst, sign}
		transactionForDB, err = json.Marshal(transcation)
		if err != nil {
			//

			return helpers.StatusInternalServerError, transactionForDB, transactionTime
		}

		return helpers.StatusOk, transactionForDB, transactionTime
	}

	return helpers.StatusUnknownTranType, transactionForDB, transactionTime
}

func verifyBalanceRequest(ctx *fasthttp.RequestCtx) (statusCode int, sender string, ttoken string) {
	//

	args := ctx.QueryArgs()
	for errNum, v := range helpers.RequestBalanceFields {
		//

		if !args.Has(v) {
			//

			return errNum, sender, ttoken
		}
	}

	ttoken = string(ctx.FormValue("TTOKEN"))
	sender = string(ctx.FormValue("SENDER"))

	if len(sender) != 66 {
		//

		return helpers.StatusWrongAttr_SENDER, sender, ttoken
	}

	return helpers.StatusOk, sender, ttoken
}

func verifyTranStatusRequest(ctx *fasthttp.RequestCtx) (statusCode int, key string) {
	//

	args := ctx.QueryArgs()
	for errNum, v := range helpers.RequestTranStatusFields {
		//

		if !args.Has(v) {
			//

			return errNum, key
		}
	}

	key = string(ctx.FormValue("KEY"))

	if len(key) != 64 {
		//

		return helpers.StatusWrongAttr_KEY, key
	}

	return helpers.StatusOk, key
}
