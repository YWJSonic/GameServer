package game

import (
	"github.com/YWJSonic/GameServer/game/cache"
	"github.com/YWJSonic/GameServer/game/gameattach"

	"github.com/YWJSonic/BaseServer/server"
	"github.com/YWJSonic/ServerUtility/igame"
	"github.com/YWJSonic/ServerUtility/playerinfo"
	"github.com/YWJSonic/ServerUtility/restfult"
	"github.com/YWJSonic/ServerUtility/socket"
	"github.com/YWJSonic/ServerUtility/user"
)

// Game ...
type Game struct {
	Server    *server.Service
	Cache     *cache.GameCache
	IGameRule igame.IGameRule
	// ProtocolMap map[string]func(r *http.Request) protocol.IProtocol
}

// RESTfulURLs ...
func (g *Game) RESTfulURLs() []restfult.Setting {
	return []restfult.Setting{
		restfult.Setting{
			RequestType: "POST",
			URL:         "lobby/login",
			Fun:         g.login,
			ConnType:    restfult.Client,
		},
		restfult.Setting{
			RequestType: "POST",
			URL:         "lobby/gameinit",
			Fun:         g.gameinit,
			ConnType:    restfult.Client,
		},
		restfult.Setting{
			RequestType: "POST",
			URL:         "lobby/refresh",
			Fun:         g.refresh,
			ConnType:    restfult.Client,
		},
		restfult.Setting{
			RequestType: "POST",
			URL:         "lobby/exchange",
			Fun:         g.exchange,
			ConnType:    restfult.Client,
		},
		restfult.Setting{
			RequestType: "POST",
			URL:         "lobby/gameresult",
			Fun:         g.gameresult,
			ConnType:    restfult.Client,
		},
	}
}

// SocketURLs ...
func (g *Game) SocketURLs() []socket.Setting {
	return []socket.Setting{
		socket.Setting{
			URL: "lobby/createNewSocket",
			Fun: g.createNewSocket,
		},
	}
}

// NewUser ...
func (g *Game) NewUser(token, gameAccount string) *user.Info {
	return &user.Info{}
}

// GetUser ...
func (g *Game) GetUser(userToken string) (*user.Info, error) {
	return &user.Info{}, nil
}

// GetUserByGameID ...
func (g *Game) GetUserByGameID(userToken string, UserID int64) (*user.Info, error) {
	return &user.Info{
		UserServerInfo: &playerinfo.AccountInfo{},
		UserGameInfo: &playerinfo.Info{
			ID: UserID,
		},
		IAttach: gameattach.NewAttach(UserID),
	}, nil
}

// CheckGameType ...
func (g *Game) CheckGameType(clientGameTypeID string) bool {
	return true
}

// CheckToken ...
func (g *Game) CheckToken(token string) bool {
	return true
}
