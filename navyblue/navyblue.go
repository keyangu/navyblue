package navyblue

// TODO Trunsaction と ancestor クエリを使わないとまずそう

import (
	"appengine"
	//    "appengine/datastore"
	"appengine/user"
	"html/template"
	"net/http"
	//  "time"
	"fmt"
	"strconv"
)

var RegisterHTMLTemplate = template.Must(template.ParseFiles("register.html"))
var DeployHTMLTemplate = template.Must(template.ParseFiles("deploy.html"))
var BattleHTMLTemplate = template.Must(template.ParseFiles("battle.html"))
var FinishHTMLTemplate = template.Must(template.ParseFiles("finish.html"))

func init() {
	http.HandleFunc("/doRegister", doRegister)
	http.HandleFunc("/doDeploy", doDeploy)
	http.HandleFunc("/reset", reset)
	http.HandleFunc("/doAttack", doAttack)
	http.HandleFunc("/doMove", doMove)
	http.HandleFunc("/debug", debug)
	http.HandleFunc("/", navyhandler)
}

// ハンドラ
// ステートに応じて各処理にリダイレクトする
func navyhandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)

	/*
		   app.yaml に login:required を書いたので、自分でログイン処理を
		   書く必要がない
		// ユーザがnil(ログインしていない)
		if u == nil {
			url, err := user.LoginURL(c, r.URL.String())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Location", url)
			w.WriteHeader(http.StatusFound)
			return
		}
	*/

	// ゲームの状態をデータストアから取得
	g := new(Game)
	g.getFromStore(c)

	switch g.State {
	case Init:
		register(w, r, c, u, g)
	case Deploy:
		deploy(w, r, c, u, g)
	case Turn1:
		fallthrough
	case Turn2:
		battle1(w, r, c, u, g)
	case Finish:
		finish(w, r, c, u, g)
	}
}

// プレイヤー登録画面表示
func register(w http.ResponseWriter, r *http.Request, c appengine.Context, u *user.User, g *Game) {

	sg := new(ShowGame)
	sg.Make(g, u)
	if err := RegisterHTMLTemplate.Execute(w, sg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// プレイヤー登録CGI
func doRegister(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	u := user.Current(c)

	if r.FormValue("name") == "" {
		fmt.Fprintf(w, "なんか入力してね！")
		return
	}

	// 名前からプレイヤーオブジェクトを生成
	p := Player{
		User:  *u,
		Name:  r.FormValue("name"),
		Ptype: Player1,
	}
	p.Warship = *new(Ship)
	p.Warship.SetDefaultWarship()
	p.Cruiser = *new(Ship)
	p.Cruiser.SetDefaultCruiser()
	p.Submarine = *new(Ship)
	p.Submarine.SetDefaultSubmarine()

	g := new(Game)
	if err := g.setNewPlayer(&c, &p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	/*
		g.getFromStore(c)

		if g.Player1.User == *u {
			// 同じアカウントなので重複登録させない
			fmt.Fprintf(w, "すでに登録してるよ！")
			return
		}
		if g.Player1.Name != "" {
			p.Ptype = Player2
			g.Player2 = p
			g.State = Deploy
		} else {
			g.Player1 = p
		}

		if err := g.putToStore(c); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	*/

	http.Redirect(w, r, "/", http.StatusFound)
}

/**
 * 配置画面
 */
func deploy(w http.ResponseWriter, r *http.Request, c appengine.Context, u *user.User, g *Game) {

	//fmt.Fprintf(w, "%v", g)
	//return

	// ログイン情報からプレイヤー情報を取得
	p, _ := g.getPlayer(u)

	if p == nil {
		fmt.Fprintf(w, "ゲームに参加してないよ!")
		return
	}

	sg := new(ShowGame)
	sg.Make(g, u)

	// ユーザに応じた盤面の表示
	if err := DeployHTMLTemplate.Execute(w, sg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/**
 * 配置処理画面
 */
func doDeploy(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)

	// フォーム入力内容の取得
	wsx, _ := strconv.Atoi(r.FormValue("warship_x"))
	wsy, _ := strconv.Atoi(r.FormValue("warship_y"))
	crx, _ := strconv.Atoi(r.FormValue("cruiser_x"))
	cry, _ := strconv.Atoi(r.FormValue("cruiser_y"))
	smx, _ := strconv.Atoi(r.FormValue("submarine_x"))
	smy, _ := strconv.Atoi(r.FormValue("submarine_y"))

	// 指定された座標がかぶっていないかチェック
	if check_duplicate(wsx, wsy, crx, cry, smx, smy) {
		fmt.Fprintf(w, "戦艦の位置が重複しています。戻って修正してください。")
		return
	}

	g := new(Game)
	err := g.setDeployData(&c, u, wsx, wsy, crx, cry, smx, smy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
	return

}

// プレイヤーの攻撃ターン
func battle1(w http.ResponseWriter, r *http.Request, c appengine.Context, u *user.User, g *Game) {

	p, _ := g.getPlayer(u)
	if p == nil {
		fmt.Fprintf(w, "ゲームに参加してないよ!")
		return
	}

	// 表示に使用するゲーム情報を生成
	sg := new(ShowGame)
	sg.Make(g, u)

	// ユーザに応じた盤面の表示
	if err := BattleHTMLTemplate.Execute(w, sg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// 攻撃実行
func doAttack(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)

	// フォームの情報を取得
	atx, _ := strconv.Atoi(r.FormValue("attack_x"))
	aty, _ := strconv.Atoi(r.FormValue("attack_y"))

	g := new(Game)
	err := g.doAttack(&c, u, atx, aty)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

// 移動
func doMove(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)

	// フォーム入力の取得
	stype, _ := strconv.Atoi(r.FormValue("shiptype"))
	way, _ := strconv.Atoi(r.FormValue("way"))
	cell, _ := strconv.Atoi(r.FormValue("cell"))

	g := new(Game)
	err := g.doMove(&c, u, stype, way, cell)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

// ゲーム終了
func finish(w http.ResponseWriter, r *http.Request, c appengine.Context, u *user.User, g *Game) {

	sg := new(ShowGame)
	sg.Make(g, u)

	if err := FinishHTMLTemplate.Execute(w, sg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// デバッグ
func debug(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	g := new(Game)
	g.getFromStore(c)
	fmt.Fprintf(w, "%v", g)
}

// リセット
func reset(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	g := new(Game)
	g.deleteStore(c)

	http.Redirect(w, r, "/", http.StatusFound)
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

// 指定した座標が重複しているかどうかチェック
// @return true 重複あり
// @return false 重複なし
func check_duplicate(wsx, wsy, crx, cry, smx, smy int) bool {
	if wsx == crx && wsy == cry {
		return true
	}
	if wsx == smx && wsy == smy {
		return true
	}
	if crx == smx && cry == smy {
		return true
	}

	return false
}
