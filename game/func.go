package game

import (
	"encoding/json"
	"errors"

	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/ServerUtility/foundation/fileload"
	"gitlab.com/ServerUtility/restfult"
	"gitlab.com/ServerUtility/socket"
	"gitlab.com/Webserver/game/gamerule"
	"gitlab.com/baseserver/server"
)

// NewGameServer ...
func NewGameServer() {

	jsStr := fileload.Load("./file/config.json")
	config := foundation.StringToJSON(jsStr)
	baseSetting := server.NewSetting()
	baseSetting.SetData(config)

	gamejsStr := fileload.Load("./file/gameconfig.json")
	var gameRule = &gamerule.Rule{}
	if err := json.Unmarshal([]byte(gamejsStr), &gameRule); err != nil {
		panic(errors.New("gameconfig error: "))
	}

	var gameserver = server.NewService()
	var game = &Game{
		IGameRule: gameRule,
		Server:    gameserver,
		// ProtocolMap: protocol.NewProtocolMap(),
	}
	gameserver.Restfult = restfult.NewRestfultService()
	gameserver.Socket = socket.NewSocket()
	gameserver.IGame = game

	// start Server
	gameserver.Launch(baseSetting)

	// start DB service
	setting := gameserver.Setting.DBSetting()
	gameserver.LaunchDB("gameDB", setting)
	gameserver.LaunchDB("logDB", setting)
	gameserver.LaunchDB("payDB", setting)

	// start restful service
	go gameserver.LaunchRestfult(game.RESTfulURLs())
	go gameserver.LaunchSocket(game.SocketURLs())

	for {
		<-gameserver.ShotDown
	}
}
