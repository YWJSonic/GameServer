package protocol

import (
	"net/http"

	"github.com/YWJSonic/ServerUtility/foundation"
	"github.com/YWJSonic/ServerUtility/myhttp"
)

// // NewProtocolMap ...
// func NewProtocolMap() map[string]func(r *http.Request) IProtocol {
// 	Map := map[string]func(r *http.Request) IProtocol{
// 		"createNewSocketProtocol": func(r *http.Request) IProtocol {
// 			prot := &CreateNewSocketProtocol{}
// 			IProtocol.InitData(prot, r)
// 			return prot
// 		},
// 		"gameresultProtocol": func(r *http.Request) IProtocol {
// 			prot := &GameResultProtocol{}
// 			IProtocol.InitData(prot, r)
// 			return prot
// 		},
// 	}
// 	return Map
// }
// // IProtocol ...
// type IProtocol interface {
// 	InitData(r *http.Request) // Load http Data GET,POST,Socket Head
// }

// CreateNewSocketProtocol ...
type CreateNewSocketProtocol struct {
}

// InitData ...
func (c *CreateNewSocketProtocol) InitData(r *http.Request) {

}

// GameResultProtocol ...
type GameResultProtocol struct {
	Token      string
	BetIndex   int64
	GameTypeID string
	PlayerID   int64
}

// InitData ...
func (c *GameResultProtocol) InitData(r *http.Request) {
	postData := myhttp.PostData(r)
	c.Token = foundation.InterfaceToString(postData["token"])
	c.BetIndex = foundation.InterfaceToInt64(postData["bet"])
	c.GameTypeID = foundation.InterfaceToString(postData["gametypeid"])
	c.PlayerID = foundation.InterfaceToInt64(postData["playerid"])
}
