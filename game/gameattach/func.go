package gameattach

import "githab.com/ServerUtility/attach"

// NewAttach ...
func NewAttach(userID int64) attach.IAttach {
	attach := &UserAttach{
		userID:  userID,
		dataMap: make(map[int64]map[int64]attach.Info),
	}
	// attach.InitData(userID)
	return attach
}
