package navyblue

import (
	"appengine"
	"appengine/datastore"
	"net/http"
	"appengine/user"
)

// アプリステート
const (
	Init = iota
	Deploy
	Turn1
	Turn2
	Finish
)

const gameKeyStrID = "12345678"

// ゲームの状態(DataStore保存用)
type Game struct {
	State		int
	Player1		Player
	Player2		Player
	Winner		Player
	GMessage	string
}

// ゲームデータを取得
func (g *Game) getFromStore(c appengine.Context) *Game {
	if err := datastore.Get(c, datastore.NewKey(c, "Game", gameKeyStrID, 0, nil), g); err != nil {
		return nil
	}
	return g
}

// ゲームデータをDataStoreに保存
func (g *Game) putToStore(c appengine.Context, w http.ResponseWriter) int {
	g.deleteStore(c)
	_, err := datastore.Put(c, datastore.NewKey(c, "Game", gameKeyStrID, 0, nil), g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return -1
	}
	return 0
}

// ゲームデータをDataStoreから削除
func (g *Game) deleteStore(c appengine.Context) {
	datastore.Delete(c, datastore.NewKey(c, "Game", gameKeyStrID, 0, nil))
}

// ユーザ情報から味方と敵のプレイヤー情報を返す
// 最初に自分のユーザと一致するプレイヤー
// 次に相手のプレイヤー情報を返す
func (g *Game) getPlayer(u *user.User) (*Player, *Player) {
	switch *u {
	case g.Player1.User:
		return &g.Player1, &g.Player2
	case g.Player2.User:
		return &g.Player2, &g.Player1
	default:
		return nil, nil
	}
}

