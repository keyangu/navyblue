package navyblue

// 船のタイプ
const (
	SHIP_TYPE_NONE = iota
	SHIP_TYPE_WARSHIP
	SHIP_TYPE_CRUISER
	SHIP_TYPE_SUBMARINE
)

// Cellに該当する船の情報
type Ship struct {
	ShipType	int
	HP			int
	PX			int
	PY			int
}

// デフォルトの戦艦を生成する
func (s *Ship) SetDefaultWarship() {
	s.ShipType = SHIP_TYPE_WARSHIP
	s.HP = 3
	s.PX = 0
	s.PY = 0
}

// デフォルトの戦艦を生成する
func (s *Ship) SetDefaultCruiser() {
	s.ShipType = SHIP_TYPE_CRUISER
	s.HP = 2
	s.PX = 0
	s.PY = 1
}

// デフォルトの戦艦を生成する
func (s *Ship) SetDefaultSubmarine() {
	s.ShipType = SHIP_TYPE_SUBMARINE
	s.HP = 1
	s.PX = 0
	s.PY = 2
}

// 指定した座標が攻撃可能かどうか判定
func (s *Ship) checkAttackable(x, y int) int {
	if s.HP > 0 {
		if (s.PX - 1 <= x) && (x <= s.PX + 1) && (s.PY - 1 <= y) && (y <= s.PY + 1) {
			return 1
		}
	}
	return 0
}

const (
	SHIP_ATTACK_HIT_DAMAGE = iota
	SHIP_ATTACK_HIT_SUNK
	SHIP_ATTACK_NEAR
	SHIP_ATTACK_MISS
)
// 指定した座標に攻撃される
func (s *Ship) Attacked(x, y int) int {
	// すでに沈没している場合はMISS
	if s.HP == 0 {
		return SHIP_ATTACK_MISS
	}
	if (s.PX == x) && (s.PY == y) {
		// HIT
		if s.HP -= 1; s.HP == 0 {
			// 沈没
			return SHIP_ATTACK_HIT_SUNK
		} else {
			// ダメージ
			return SHIP_ATTACK_HIT_DAMAGE
		}
	}
	// 近い
	if (s.PX - 1 <= x && x <= s.PX + 1) && (s.PY - 1 <= y && y <= s.PY + 1) {
		return SHIP_ATTACK_NEAR
	}

	return SHIP_ATTACK_MISS
}

