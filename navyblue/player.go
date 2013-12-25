package navyblue

import (
	"appengine/user"
)

// プレイヤータイプ
// const と iota を使って enum のように定義できる
const (
	Player1 = iota		// Player1 == 0
	Player2				// Player2 == 1
	Guest				// Guest   == 2
)

const (
	PST_INIT = iota
	PST_READY
	PST_TURN
	PST_WAIT
	PST_FIN
)

// プレイヤー
type Player struct {
	User		user.User
	Name		string
	Ptype		int
	State		int
	Warship		Ship
	Cruiser		Ship
	Submarine	Ship
	Message		string
}

// StateがPST_READYかどうかを返す
func (p *Player) IsReady() bool {
	if p.State == PST_READY {
		return true
	}
	return false
}

// 指定した座標が攻撃可能かどうか判定
func (p *Player) checkAttackable(x, y int) int {
	if p.Warship.checkAttackable(x, y) == 1 {
		return 1
	}
	if p.Cruiser.checkAttackable(x, y) == 1 {
		return 1
	}
	if p.Submarine.checkAttackable(x, y) == 1 {
		return 1
	}
	return 0
}

const (
	ATTACK_SUCCESS_DAMAGE = iota
	ATTACK_SUCCESS_SUNK
	ATTACK_FAIL_NEAR
	ATTACK_FAIL_FAR
	ATTACK_ERR
)

// 攻撃実行処理
// 攻撃対象のPlayerを主体とすること
// 攻撃可能かどうかは別途判定済み
func (enm *Player) Attacked(x, y int) (int, int) {
	var s *Ship
	var ret1, ret2, ret3 int
	// TODO 沈没、ダメージ、近く、MISSで戻る優先度が変わるはず
	s = &enm.Warship
	ret1 = s.Attacked(x, y)
	if ret1 == SHIP_ATTACK_HIT_DAMAGE {
		return ATTACK_SUCCESS_DAMAGE, (-1)
	} else if ret1 == SHIP_ATTACK_HIT_SUNK {
		return ATTACK_SUCCESS_SUNK, s.ShipType
	}

	s = &enm.Cruiser
	ret2 = s.Attacked(x, y)
	if ret2 == SHIP_ATTACK_HIT_DAMAGE {
		return ATTACK_SUCCESS_DAMAGE, (-1)
	} else if ret2 == SHIP_ATTACK_HIT_SUNK {
		return ATTACK_SUCCESS_SUNK, s.ShipType
	}

	s = &enm.Submarine
	ret3 = s.Attacked(x, y)
	if ret3 == SHIP_ATTACK_HIT_DAMAGE {
		return ATTACK_SUCCESS_DAMAGE, (-1)
	} else if ret3 == SHIP_ATTACK_HIT_SUNK {
		return ATTACK_SUCCESS_SUNK, s.ShipType
	}

	// どれもミスならFAR
	if ret1 == SHIP_ATTACK_MISS && ret2 == SHIP_ATTACK_MISS && ret3 == SHIP_ATTACK_MISS {
		return ATTACK_FAIL_FAR, (-1)
	}

	// どれか一つでも近ければNEAR
	return ATTACK_FAIL_NEAR, (-1)
}

const (
	MOVE_SUCCESS = iota
	MOVE_FAIL_AREA_OVER		// 範囲外
	MOVE_FAIL_SUNKEN		// 沈没してる
	MOVE_FAIL_DUPLICATE		// 他の船が移動先にいる
	MOVE_FAIL_INTERNAL
)
// 船の移動
func (p *Player) moveShip(s *Ship, way, cell int) int {
	// 他の船の位置チェックも必要なので、
	// Shipのメソッドではなく、Playerのメソッドにしている

	// 移動しようとしている船が沈没していないか
	if s.HP <= 0 {
		return MOVE_FAIL_SUNKEN
	}

	x := s.PX
	y := s.PY

	// 移動量を出す
	if way == 0 || way == 3 {
		cell *= -1
	}
	// 仮移動
	if way == 0 || way == 2 {
		y += cell
	} else if way == 1 || way == 3 {
		x += cell
	} else {
		return MOVE_FAIL_INTERNAL
	}

	// 移動可能かどうかチェック
	if (x < 0 || x > BOARD_SIZE) || (y < 0 || y > BOARD_SIZE) {
		// 移動先がボードをはみ出した
		return MOVE_FAIL_AREA_OVER
	}
	// 移動先に別の船がいないかどうか
	if (s.ShipType != p.Warship.ShipType) && (s.PX == p.Warship.PX && s.PY == p.Warship.PY) &&
	(p.Warship.HP > 0) {
		return MOVE_FAIL_DUPLICATE
	}
	if (s.ShipType != p.Cruiser.ShipType) && (s.PX == p.Cruiser.PX && s.PY == p.Cruiser.PY) &&
	(p.Cruiser.HP > 0) {
		return MOVE_FAIL_DUPLICATE
	}
	if (s.ShipType != p.Submarine.ShipType) && (s.PX == p.Submarine.PX && s.PY == p.Submarine.PY) &&
	(p.Submarine.HP > 0) {
		return MOVE_FAIL_DUPLICATE
	}

	// 移動可能であれば移動する
	s.PX = x
	s.PY = y
	return MOVE_SUCCESS
}

// プレイヤーが所持している船が全て撃沈されたかどうか
func (p *Player) checkSunked() bool {
	if (p.Warship.HP <= 0) && (p.Cruiser.HP <= 0) && (p.Submarine.HP <= 0) {
		return true
	}
	return false
}
