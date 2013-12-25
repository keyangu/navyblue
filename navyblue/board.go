package navyblue

const BOARD_SIZE = 5

// ボード
type Board struct {
	Cell	[BOARD_SIZE][BOARD_SIZE] string
}

// プレイヤー内のShip情報から
// 盤面のデータを作る
func (b *Board) setCellData(p Player) {

	for x := 0; x < BOARD_SIZE; x++ {
		for y := 0; y < BOARD_SIZE; y++ {
			if x == p.Warship.PX && y == p.Warship.PY && p.Warship.HP > 0 {
				b.Cell[x][y] = "W"
			} else if x == p.Cruiser.PX && y == p.Cruiser.PY && p.Cruiser.HP > 0 {
				b.Cell[x][y] = "C"
			} else if x == p.Submarine.PX && y == p.Submarine.PY && p.Submarine.HP > 0 {
				b.Cell[x][y] = "S"
			} else {
				b.Cell[x][y] = "　"
			}
		}
	}
}
