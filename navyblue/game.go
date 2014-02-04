package navyblue

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"errors"
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
	State    int
	Player1  Player
	Player2  Player
	Winner   Player
	GMessage string
}

// ゲームデータを取得
func (g *Game) getFromStore(c appengine.Context) *Game {
	if err := datastore.Get(c, datastore.NewKey(c, "Game", gameKeyStrID, 0, nil), g); err != nil {
		return nil
	}
	return g
}

// ゲームデータをDataStoreに保存
func (g *Game) putToStore(c appengine.Context) error {
	g.deleteStore(c)
	_, err := datastore.Put(c, datastore.NewKey(c, "Game", gameKeyStrID, 0, nil), g)
	return err
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

func (g *Game) setNewPlayer(c *appengine.Context, p *Player) error {
	err := datastore.RunInTransaction(*c, func(c appengine.Context) error {
		g.getFromStore(c)
		if g.Player1.User == p.User {
			return errors.New("すでに登録してるよ！")
		}
		if g.Player1.Name != "" {
			p.Ptype = Player2
			g.Player2 = *p
			g.State = Deploy
		} else {
			g.Player1 = *p
		}
		err := g.putToStore(c)
		return err
	}, nil)

	return err
}
