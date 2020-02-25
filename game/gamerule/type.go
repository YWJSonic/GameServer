package gamerule

import (
	"githab.com/ServerUtility/attach"
	"githab.com/ServerUtility/game"
)

// JackPartBonusx2Index ...
const JackPartBonusx2Index int64 = 0

// JackPartBonusx3Index ...
const JackPartBonusx3Index int64 = 1

// JackPartBonusx5Index ...
const JackPartBonusx5Index int64 = 2

type jackPart struct {
	JackPartBonusx2 attach.Info
	JackPartBonusx3 attach.Info
	JackPartBonusx5 attach.Info
}

// Rule ...
type Rule struct {
	Version             string        `json:"Version"`             // game logic version
	GameTypeID          string        `json:"GameTypeID"`          // game unique id
	GameIndex           int           `json:"GameIndex"`           // game sort id
	WinScoreLimit       int64         `json:"WinScoreLimit"`       // game round win money limit
	WinBetRateLimit     int64         `json:"WinBetRateLimit"`     // game round win rate limit
	BetRate             []int64       `json:"BetRate"`             // bate money slice
	BetRateLinkIndex    []int64       `json:"BetRateLinkIndex"`    // player bet fest link index on BetRate slice
	BetRateDefaultIndex int64         `json:"BetRateDefaultIndex"` // default player bet index
	NormalReelSize      []int         `json:"NormalReelSize"`      // Normal reel row size ex:[3,4,5,4,3]
	NormalReelSymbol    [][]int       `json:"NormalReelSymbol"`    // Normal reel [[1,2,3],[4,5,6,7],[1,2,3,4,5],[4,5,6,7],[1,2,3]]
	FreeReelSize        []int         `json:"FreeReelSize"`        // Free reel row size ex:[3,4,5,4,3]
	FreeReelSymbol      [][]int       `json:"FreeReelSymbol"`      // Free reel
	RespinScroll        [][]int       `json:"RespinScroll"`
	RTPSetting          []int         `json:"RTPSetting"` //index 0: normal RTP. index 1:bonus RTP
	Space               int           `json:"Space"`
	WildsItemIndex      []int         `json:"WildsItemIndex"`
	ScotterItemIndex    []int         `json:"ScotterItemIndex"`
	ItemResults         [][]int       `json:"ItemResults"`
	JackPortResults     [][]int       `json:"JackPortResults"`
	RespinitemResults   [][]int       `json:"RespinitemResults"`
	SymbolGroup         map[int][]int `json:"SymbolGroup"`
	SpWhildWinRate      []int64       `json:"SpWhildWinRate"`
	JackPortTex         []float32     `json:"jackPortTex"`
	JackPartWinRate     []int         `json:"JackPartWinRate"`
}

// GetGameTypeID ...
func (r *Rule) GetGameTypeID() string {
	return r.GameTypeID
}

// GetBetMoney ...
func (r *Rule) GetBetMoney(index int64) int64 {
	return r.BetRate[index]
}

// GetInitScroll ...
func (r *Rule) GetInitScroll() interface{} {
	scrollmap := map[string][][]int{
		"normalreel": r.normalReel(),
		"respinreel": {r.respinReel()},
	}
	return scrollmap
}

func (r *Rule) normalReel() [][]int {
	return r.NormalReelSymbol
}
func (r *Rule) respinReel() []int {
	return r.RespinScroll[r.normalRTP()]
}
func (r *Rule) normalRTP() int {
	return r.RTPSetting[0]
}
func (r *Rule) respinRTP() int {
	return r.RTPSetting[1]
}

// Wild1 ...
func (r *Rule) Wild1() int {
	return r.WildsItemIndex[0]
}

// Wild2 ...
func (r *Rule) Wild2() int {
	return r.WildsItemIndex[1]
}

// Wild3 ...
func (r *Rule) Wild3() int {
	return r.WildsItemIndex[2]
}

// Wild4 ...
func (r *Rule) Wild4() int {
	return r.WildsItemIndex[3]
}

// func (r *GameLogic) LogicResult(betMoney int64, user *user.Info) map[string]interface{} {

// 	// return
// }

// func (r *GameLogic) OutputGame(betMoney int64, user *user.Info) (map[string]interface{}, map[string]interface{}, int64) {
// }

// func (r *GameLogic) outRespin(betMoney int64, user *user.Info) ([]interface{}, int64) {
// }

// GameRequest ...
func (r *Rule) GameRequest(config *game.RuleRequest) *game.RuleRespond {
	betMoney := r.GetBetMoney(config.BetIndex)
	jackPart := r.getJPFromAttach(config.Attach)
	result := make(map[string]interface{})
	otherData := make(map[string]interface{})
	var totalWin int64

	gameResult := r.newlogicResult(betMoney, &jackPart)

	result["normalresult"] = gameResult.Normalresult
	totalWin += gameResult.Normaltotalwin

	if gameResult.Respinresult != nil {
		result["respin"] = gameResult.Respinresult
		otherData["isrespin"] = 1
		totalWin += gameResult.Respintotalwin
	}

	result["totalwinscore"] = totalWin
	r.setJPFromAttach(betMoney, config.Attach, &jackPart)

	return &game.RuleRespond{
		Attach:        config.Attach,
		BetMoney:      betMoney,
		Totalwinscore: totalWin,
		GameResult:    result,
		OtherData:     otherData,
	}
}
func (r *Rule) setJPFromAttach(betMoney int64, att attach.IAttach, jP *jackPart) {
	att.SetValue(int64(r.GameIndex), JackPartBonusx2Index, "", jP.JackPartBonusx2.GetIValue()+int64(float32(betMoney)*r.JackPortTex[2]))
	att.SetValue(int64(r.GameIndex), JackPartBonusx3Index, "", jP.JackPartBonusx3.GetIValue()+int64(float32(betMoney)*r.JackPortTex[1]))
	att.SetValue(int64(r.GameIndex), JackPartBonusx5Index, "", jP.JackPartBonusx5.GetIValue()+int64(float32(betMoney)*r.JackPortTex[0]))
}
func (r *Rule) getJPFromAttach(att attach.IAttach) jackPart {
	value := jackPart{
		JackPartBonusx2: att.Get(int64(r.GameIndex), JackPartBonusx2Index),
		JackPartBonusx3: att.Get(int64(r.GameIndex), JackPartBonusx3Index),
		JackPartBonusx5: att.Get(int64(r.GameIndex), JackPartBonusx5Index),
	}
	return value
}

// // Result att 0: freecount
// func (r *Rule) logicResult(betMoney int64, JP *jackPart) map[string]interface{} {
// 	var result = make(map[string]interface{})
// 	var totalWin int64
// 	normalresult, otherdata, normaltotalwin := r.outputGame(betMoney, JP)
// 	result = foundation.AppendMap(result, otherdata)
// 	result["normalresult"] = normalresult
// 	totalWin += normaltotalwin
// 	if otherdata["isrespin"].(int) == 1 {
// 		respinresult, respintotalwin := r.outRespin(betMoney, JP)
// 		totalWin += respintotalwin
// 		result["respin"] = respinresult
// 		result["isrespin"] = 1
// 	}
// 	result["totalwinscore"] = totalWin
// 	return result
// }
