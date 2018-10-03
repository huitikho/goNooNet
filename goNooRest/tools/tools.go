package tools

import (
	"encoding/json"
	"goNooNet/helpers"

	"github.com/valyala/fasthttp"
)

func MakeResponse(statusCode int,
	ctx *fasthttp.RequestCtx) {
	//

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(statusCode)
}

func MakeBalanceResponse(balance string,
	statusCode int,
	ctx *fasthttp.RequestCtx) {
	//

	response := helpers.BalanceResponse{balance}
	jsResponse, _ := json.Marshal(response)
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(statusCode)
	ctx.SetBody(jsResponse)
}
