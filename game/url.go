package game

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"gitlab.com/ServerUtility/code"
	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/ServerUtility/game"
	"gitlab.com/ServerUtility/httprouter"
	"gitlab.com/ServerUtility/messagehandle"
	"gitlab.com/ServerUtility/myhttp"
	"gitlab.com/ServerUtility/socket"
	"gitlab.com/Webserver/game/protocol"
)

func (g *Game) createNewSocket(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	if !g.CheckToken(token) {
		log.Fatal("createNewSocket: not this token\n")
		return
	}

	_, err := g.GetUser(token)
	if err != nil {
		log.Fatal("createNewSocket: not this user\n")
		return
	}

	c, err := g.Server.Socket.Upgrade(w, r, r.Header)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	g.Server.Socket.AddNewConn("f", c, g.SocketMessageHandle)
	// g.Server.Socket.AddNewConn(user.GetGameInfo().GameAccount, c, g.SocketMessageHandle)

	time.Sleep(time.Second * 3)
	g.Server.Socket.ConnMap["f"].Send(websocket.CloseMessage, []byte{})
}

// SocketMessageHandle ...
func (g *Game) SocketMessageHandle(msg socket.Message) error {
	fmt.Println("#-- socket --#", msg)
	return nil
}

func (g *Game) login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var result = make(map[string]interface{})

	postData := myhttp.PostData(r)
	// logintype := foundation.InterfaceToInt(postData["logintype"])
	gametypeid := foundation.InterfaceToString(postData["gametypeid"])
	token := foundation.InterfaceToString(postData["accounttoken"])

	if g.CheckToken(gametypeid) {
		return
	}

	user, err := g.GetUser(token)
	if err != nil {
		return
	}

	result["gameaccount"] = user.UserServerInfo.GameAccount
	result["token"] = user.UserGameInfo.GameToken
	result["serversetting"] = g.Server.Setting.ToClient()

	g.Server.Log("login Log")
	// myhttp.HTTPResponse(w, result, err)
}

func (g *Game) gameinit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
}

func (g *Game) refresh(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
}

func (g *Game) exchange(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
}

func (g *Game) checkout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
}

func (g *Game) gameresult(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var proto protocol.GameResultProtocol
	proto.InitData(r)

	// check token
	// if err = foundation.CheckToken(mycache.GetToken(playerInfo.GameAccount), token); err.ErrorCode != code.OK {
	// 	messagehandle.ErrorLogPrintln("GetPlayerInfoByPlayerID-3", err, token, betIndex, betMoney)
	// 	myhttp.HTTPResponse(w, "", err)
	// 	return
	// }

	// gametype check
	// if proto.GameTypeID != g.IGameRule.GetGameTypeID() {
	// 	err := messagehandle.New()
	// 	err.ErrorCode = code.GameTypeError
	// 	err.Msg = "GameTypeError"
	// 	// messagehandle.ErrorLogPrintln("GetPlayerInfoByPlayerID-1", err, token, betIndex, betMoney)
	// 	myhttp.HTTPResponse(w, "", err)
	// 	return
	// }

	// // get player
	player, err := g.GetUserByGameID(proto.Token, proto.PlayerID)
	if err != nil {
		err := messagehandle.New()
		err.Msg = "GameTypeError"
		err.ErrorCode = code.GameTypeError
		// messagehandle.ErrorLogPrintln("GetPlayerInfoByPlayerID-2", err, token, betIndex, betMoney)
		g.Server.HTTPResponse(w, "", err)
		return
	}

	// // money check
	// if player.UserGameInfo.Money < g.IGameRule.GetBetMoney(proto.BetIndex) {
	// 	err := messagehandle.New()
	// 	err.Msg = "NoMoneyToBet"
	// 	err.ErrorCode = code.NoMoneyToBet
	// 	// messagehandle.ErrorLogPrintln("GetPlayerInfoByPlayerID-4", err, token, betIndex, betMoney, playerInfo)
	// 	myhttp.HTTPResponse(w, "", err)
	// 	return
	// }

	// get attach
	player.LoadData()

	// get game result
	RuleRequest := &game.RuleRequest{
		BetIndex: proto.BetIndex,
		Attach:   player.IAttach,
	}
	result := g.IGameRule.GameRequest(RuleRequest)
	player.UserGameInfo.SumMoney(result.Totalwinscore - result.BetMoney)

	resultMap := make(map[string]interface{})
	resultMap["totalwinscore"] = result.Totalwinscore
	resultMap["playermoney"] = player.UserGameInfo.GetMoney()
	resultMap["normalresult"] = result.GameResult["normalresult"]
	resultMap["attach"] = result.Attach

	respin, ok := result.OtherData["isrespin"]
	if ok && respin == 1 {
		resultMap["isrespin"] = 1
		resultMap["respin"] = result.GameResult["respin"]
	} else {
		resultMap["isrespin"] = 0
		resultMap["respin"] = []interface{}{}
	}
	g.Server.HTTPResponse(w, resultMap, messagehandle.New())
}
