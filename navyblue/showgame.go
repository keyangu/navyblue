package navyblue

import (
	"appengine/user"
)

// 盤面表示で使用するゲームの全状態
type ShowGame struct {
	GameStored	Game
	Turn		int			// どちらのプレイヤーのターンか
	Friend		*Player
	Enemy		*Player
	GBoard		[2]*Board
}

// Friendが現在行動できるターンかどうか返す
func (sg *ShowGame) IsActiveTurn() bool {
	if sg.Turn == sg.Friend.Ptype {
		return true
	}
	return false
}

// FriendがPST_READY状態かどうか返す
func (sg *ShowGame) IsFriendReady() bool {
	return sg.Friend.IsReady()
}

// 盤面表示用のゲーム全データを取得
func (sg *ShowGame) Make(g *Game, u *user.User) {
	sg.GameStored = *g
	sg.Friend, sg.Enemy = g.getPlayer(u)
	if sg.Friend != nil {
		sg.GBoard[0] = new(Board)
		sg.GBoard[0].setCellData(*sg.Friend)
	}
	if sg.Enemy != nil {
		sg.GBoard[1] = new(Board)
		sg.GBoard[1].setCellData(*sg.Enemy)
	}
	if sg.GameStored.State == Turn1 {
		sg.Turn = Player1
	} else if sg.GameStored.State == Turn2 {
		sg.Turn = Player2
	}
}

