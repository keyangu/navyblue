package navyblue

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"errors"
	"strconv"
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

//==========================================================================
// イベントハンドラから呼ばれ、ゲームデータの更新を行うAPI
//==========================================================================

// registerハンドラから呼ばれ、プレイヤーの登録処理を行う
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

// doDeployハンドラから呼ばれ、戦艦の配置を行う
// すべてのプレイヤーの準備が整ったらゲームを開始する
func (g *Game) setDeployData(c *appengine.Context, u *user.User,
	wsx, wsy, crx, cry, smx, smy int) error {
	err := datastore.RunInTransaction(*c, func(c appengine.Context) error {

		g.getFromStore(c)
		p, _ := g.getPlayer(u)
		if p == nil {
			return errors.New("ゲームに参加してないよ!")
		}

		// 指定された座標をセット
		p.Warship.PX = wsx
		p.Warship.PY = wsy
		p.Cruiser.PX = crx
		p.Cruiser.PY = cry
		p.Submarine.PX = smx
		p.Submarine.PY = smy

		// 準備ができたことにする
		p.State = PST_READY

		// プレイヤーのどちらもREADYになったら、ゲーム開始ステートへ
		if g.Player1.State == PST_READY && g.Player2.State == PST_READY {
			g.State = Turn1
			g.GMessage = "開戦！！"
		}

		if err := g.putToStore(c); err != nil {
			return err
		}

		return nil

	}, nil)
	return err
}

// doAtackハンドラから呼ばれ、攻撃の処理を行う
func (g *Game) doAttack(c *appengine.Context, u *user.User, atx, aty int) error {
	err := datastore.RunInTransaction(*c, func(c appengine.Context) error {
		g.getFromStore(c)
		// ログインユーザから自軍、敵軍データを取得
		fri, enm := g.getPlayer(u)
		if fri == nil {
			return errors.New("ゲームに参加してないよ！")
		}

		// 指定した箇所が攻撃可能かどうか判定
		if fri.checkAttackable(atx, aty) == 0 {
			// 攻撃不可
			fri.Message = "そこは攻撃できる地点ではありません"
		} else {

			fri.Message = ""

			ret, _ := enm.Attacked(atx, aty)
			pstr := conv_point2str(atx, aty)
			switch ret {
			case ATTACK_SUCCESS_SUNK:
				g.GMessage = fri.Name + "が" + pstr + "を攻撃！" + ">>>>" + "撃沈！！！"
			case ATTACK_SUCCESS_DAMAGE:
				g.GMessage = fri.Name + "が" + pstr + "を攻撃！" + ">>>>" + "命中！！"
			case ATTACK_FAIL_NEAR:
				g.GMessage = fri.Name + "が" + pstr + "を攻撃！" + ">>>>" + "波高し！"
			case ATTACK_FAIL_FAR:
				g.GMessage = fri.Name + "が" + pstr + "を攻撃！" + ">>>>" + "ミス"
			default:
				return errors.New("player.Attacked()の内部エラー")
			}

			// すべての艦が撃沈したら終了
			if enm.checkSunked() {
				g.Winner = *fri
				g.State = Finish
			} else {
				if g.State == Turn1 {
					g.State = Turn2
				} else if g.State == Turn2 {
					g.State = Turn1
				} else {
					return errors.New("ステート不正")
				}
			}
		}
		if err := g.putToStore(c); err != nil {
			return err
		}
		return nil
	}, nil)
	return err
}

// doMoveハンドラから呼ばれ、移動の処理を行う
func (g *Game) doMove(c *appengine.Context, u *user.User, stype, dir, dist int) error {
	err := datastore.RunInTransaction(*c, func(c appengine.Context) error {
		g.getFromStore(c)

		// ログインユーザから自軍、敵軍データを取得
		fri, _ := g.getPlayer(u)
		if fri == nil {
			return errors.New("ゲームに参加してないよ！")
		}

		// 移動する船のタイプ、方向、マスを取得
		var s *Ship
		var s_str string
		switch stype {
		case 0:
			s = &fri.Warship
			s_str = "戦艦"
		case 1:
			s = &fri.Cruiser
			s_str = "駆逐艦"
		case 2:
			s = &fri.Submarine
			s_str = "潜水艦"
		default:
			return errors.New("CGIのデータをうまく取得できてないよ")
		}
		// 船の移動
		switch fri.moveShip(s, dir, dist) {
		case MOVE_FAIL_SUNKEN:
			fri.Message = "移動しようとした船はすでに沈没しています"
		case MOVE_FAIL_AREA_OVER:
			fri.Message = "エリアをはみ出してしまうので移動できません"
		case MOVE_FAIL_DUPLICATE:
			fri.Message = "他の味方艦のいる場所へは移動できません"
		case MOVE_FAIL_INTERNAL:
			return errors.New("p.moveShip()の内部エラー")
		case MOVE_SUCCESS:
			fri.Message = ""
			g.GMessage = fri.Name + "の" + s_str + "が" + conv_way2str(dir) + "へ" + strconv.Itoa(dist) + "移動！"
			if g.State == Turn1 {
				g.State = Turn2
			} else if g.State == Turn2 {
				g.State = Turn1
			} else {
				return errors.New("ステート不正")
			}
		}

		if err := g.putToStore(c); err != nil {
			return err
		}
		return nil

	}, nil)
	return err
}

func conv_way2str(way int) string {
	switch way {
	case 0:
		return "北"
	case 1:
		return "東"
	case 2:
		return "南"
	case 3:
		return "西"
	default:
		return ""
	}
}

func conv_point2str(x, y int) string {
	var ret string
	switch x {
	case 0:
		ret = "A-"
	case 1:
		ret = "B-"
	case 2:
		ret = "C-"
	case 3:
		ret = "D-"
	case 4:
		ret = "E-"
	}
	return ret + strconv.Itoa(y+1)
}
