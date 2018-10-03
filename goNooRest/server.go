package restServer

import (
	"goNooNet/goNooRest/tools"
	"goNooNet/helpers"
	"log"
	"net"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/valyala/fasthttp"
)

var redisDB *redis.Client
var tranChannel *chan string

func fastHTTPRawHandler(ctx *fasthttp.RequestCtx) {
	if string(ctx.Method()) == "GET" {
		//

		switch string(ctx.Path()) {

		case "/wallet/transaction":
			statusCode, transactionForDb, transactionTime := verifyTransaction(ctx)
			if statusCode >= 600 {
				//

				tools.MakeResponse(statusCode, ctx)
				return

			} else {
				//

				errRedis := redisDB.ZAdd("RAW TRANSACTIONS", redis.Z{
					Score:  float64(transactionTime),
					Member: transactionForDb,
				})

				if helpers.IsRedisError(errRedis) {
					//

					tools.MakeResponse(helpers.StatusInternalServerError, ctx)
					return
				}

				*tranChannel <- string(transactionForDb)

				tools.MakeResponse(helpers.StatusOk, ctx)
			}

		case "/wallet/getBalance":
			statusCode, sender, ttoken := verifyBalanceRequest(ctx)
			if statusCode >= 600 {
				//

				tools.MakeResponse(statusCode, ctx)
				return

			} else {
				//

				log.Printf("REQUESTING: %s tokens", ttoken)
				zScore := redisDB.ZScore("BALANCE", sender)
				if helpers.IsRedisError(zScore) {
					//

					tools.MakeBalanceResponse("0", helpers.StatusOk, ctx)
					return
				}

				tools.MakeBalanceResponse(strconv.FormatFloat(zScore.Val(), 'f', -1, 64), helpers.StatusOk, ctx)
			}

		case "/wallet/tranStatus":
			statusCode, key := verifyTranStatusRequest(ctx)
			if statusCode >= 600 {
				//

				tools.MakeResponse(statusCode, ctx)
				return

			} else {
				//

				zScore := redisDB.ZScore("COMPLETE TRANSACTIONS", key)
				if !helpers.IsRedisError(zScore) {
					//

					tools.MakeResponse(helpers.StatusOk, ctx)
					return
				}

				zScore = redisDB.ZScore("FAILED TRANSACTIONS", key)
				if !helpers.IsRedisError(zScore) {
					//

					tools.MakeResponse(helpers.StatusTranFailed, ctx)
					return
				}

				tools.MakeResponse(helpers.StatusTranNotFound, ctx)
				return
			}

		case "/blockchain/getHeight":
			tools.MakeResponse(helpers.StatusOk, ctx)

		case "/blockchain/getTran":
			tools.MakeResponse(helpers.StatusOk, ctx)

		case "/blockchain/getBlock":
			tools.MakeResponse(helpers.StatusOk, ctx)

		case "/blockchain/getVersion":
			tools.MakeResponse(helpers.StatusOk, ctx)

		default:
			//

			ctx.Error("Unsupported path", fasthttp.StatusNotFound)
		}

		return
	}

	ctx.Error("Unsupported method", fasthttp.StatusMethodNotAllowed)
}

func Start(r *redis.Client, c *chan string, ip string) {
	//

	redisDB = r
	tranChannel = c

	server := &fasthttp.Server{
		Handler: fastHTTPRawHandler,
		//	DisableKeepalive: true, // We have some strange things with Android clients when True
		GetOnly: true,
	}

	panic(server.ListenAndServe(net.JoinHostPort(ip, "5000")))
}
