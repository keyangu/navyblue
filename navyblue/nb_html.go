package navyblue

import (
	"html/template"
)

var RegisterHTMLTemplate = template.Must(template.New("Register").Parse(RegisterHTML))
var DebugHTMLTemplate = template.Must(template.New("Debug").Parse(DebugHTML))
var DeployHTMLTemplate = template.Must(template.New("Deploy").Parse(DeployHTML))
var BattleHTMLTemplate = template.Must(template.New("Battle").Parse(BattleHTML))
var FinishHTMLTemplate = template.Must(template.New("Finish").Parse(FinishHTML))

const RegisterHTML = `
<html>
	<body>
	<h1>Player登録</h1>
	{{if .Friend}}
	<p>ようこそ、{{.Friend.Name}}さん。</p>
	<p>対戦相手が登録するまでお待ちください。</p>
	<p><a href=".">更新</a></p>
	{{else}}
	<form action="/doRegister" method="post">
		<div>名前：<input name="name" type="text" /><input type="submit" value="参加" /></div>
	</form>
	{{end}}
<!--
	<hr />
	<p>debug</p>
	{{.}}
	<hr />
-->
	</body>
</html>
`

const DeployHTML = `
<html>
	<body>
	<h1>配備</h1>
	<p>配置を決めてください。</p>
	<table bordercolor="black" border="1">
		<tr>
			<td>{{with index .GBoard 0}}{{index .Cell 0 0}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 1 0}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 2 0}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 3 0}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 4 0}}{{end}}</td>
			<th>1</th>
		</tr>
		<tr>
			<td>{{with index .GBoard 0}}{{index .Cell 0 1}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 1 1}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 2 1}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 3 1}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 4 1}}{{end}}</td>
			<th>2</th>
		</tr>
		<tr>
			<td>{{with index .GBoard 0}}{{index .Cell 0 2}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 1 2}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 2 2}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 3 2}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 4 2}}{{end}}</td>
			<th>3</th>
		</tr>
		<tr>
			<td>{{with index .GBoard 0}}{{index .Cell 0 3}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 1 3}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 2 3}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 3 3}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 4 3}}{{end}}</td>
			<th>4</th>
		</tr>
		<tr>
			<td>{{with index .GBoard 0}}{{index .Cell 0 4}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 1 4}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 2 4}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 3 4}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 4 4}}{{end}}</td>
			<th>5</th>
		</tr>
		<tr><th>A</th><th>B</th><th>C</th><th>D</th><th>E</th><th>　</th></tr>
	</table>
	<form action="/doDeploy", method="post">
		<li>W(戦艦):<select name="warship_x">
						<option value="0">A</option>
						<option value="1">B</option>
						<option value="2">C</option>
						<option value="3">D</option>
						<option value="4">E</option>
					</select> x
					<select name="warship_y">
						<option value="0">1</option>
						<option value="1">2</option>
						<option value="2">3</option>
						<option value="3">4</option>
						<option value="4">5</option>
					</select>
		</li>
		<li>C(巡洋艦):<select name="cruiser_x">
						<option value="0">A</option>
						<option value="1">B</option>
						<option value="2">C</option>
						<option value="3">D</option>
						<option value="4">E</option>
					</select> x
					<select name="cruiser_y">
						<option value="0">1</option>
						<option value="1">2</option>
						<option value="2">3</option>
						<option value="3">4</option>
						<option value="4">5</option>
					</select>
		</li>
		<li>S(潜水艦):<select name="submarine_x">
						<option value="0">A</option>
						<option value="1">B</option>
						<option value="2">C</option>
						<option value="3">D</option>
						<option value="4">E</option>
					</select> x
					<select name="submarine_y">
						<option value="0">1</option>
						<option value="1">2</option>
						<option value="2">3</option>
						<option value="3">4</option>
						<option value="4">5</option>
					</select>
		</li>
		<input type="submit" value="配置" />
	</form>
	{{if .IsFriendReady}}
	<p>
	相手が配置を終えるのをお待ちください。
	相手が配置を終えるまでの間は再配置ができます。
	</p>
	<p><a href=".">更新</a></p>
	{{end}}
<!--
	<hr />
	<p>debug</p>
	{{.}}
	<hr />
-->
	</body>
</html>
`

const BattleHTML = `
<html>
	<body>
	<h1>戦闘</h1>
	<h3>{{.GameStored.GMessage}}</h3>
	<p color="red">{{.Friend.Message}}</p>
	<table bordercolor="black" border="1">
		<tr>
			<td>{{with index .GBoard 0}}{{index .Cell 0 0}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 1 0}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 2 0}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 3 0}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 4 0}}{{end}}</td>
			<th>1</th>
		</tr>
		<tr>
			<td>{{with index .GBoard 0}}{{index .Cell 0 1}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 1 1}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 2 1}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 3 1}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 4 1}}{{end}}</td>
			<th>2</th>
		</tr>
		<tr>
			<td>{{with index .GBoard 0}}{{index .Cell 0 2}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 1 2}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 2 2}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 3 2}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 4 2}}{{end}}</td>
			<th>3</th>
		</tr>
		<tr>
			<td>{{with index .GBoard 0}}{{index .Cell 0 3}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 1 3}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 2 3}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 3 3}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 4 3}}{{end}}</td>
			<th>4</th>
		</tr>
		<tr>
			<td>{{with index .GBoard 0}}{{index .Cell 0 4}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 1 4}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 2 4}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 3 4}}{{end}}</td>
			<td>{{with index .GBoard 0}}{{index .Cell 4 4}}{{end}}</td>
			<th>5</th>
		</tr>
		<tr><th>A</th><th>B</th><th>C</th><th>D</th><th>E</th><th>　</th></tr>
	</table>
	<br />
	{{if .IsActiveTurn}}
	<p>今回のターンの行動を決めてください。</p>
	<form action="/doAttack", method="post">
		<select name="attack_x">
			<option value="0">A</option>
			<option value="1">B</option>
			<option value="2">C</option>
			<option value="3">D</option>
			<option value="4">E</option>
		</select>
		x
		<select name="attack_y">
			<option value="0">1</option>
			<option value="1">2</option>
			<option value="2">3</option>
			<option value="3">4</option>
			<option value="4">5</option>
		</select>
		を魚雷で<input type="submit" value="攻撃！" />
	</form>
	<form action="/doMove", method="post">
		<select name="shiptype">
			<option value="0">W(戦艦)</option>
			<option value="1">C(巡洋艦)</option>
			<option value="2">S(潜水艦)</option>
		</select>
		を
		<select name="way">
			<option value="0">北</option>
			<option value="1">東</option>
			<option value="2">南</option>
			<option value="3">西</option>
		</select>
		へ
		<select name="cell">
			<option value="1">1</option>
			<option value="2">2</option>
			<option value="3">3</option>
			<option value="4">4</option>
		</select>
		マス
		<input type="submit" value="移動する" />
	</form>
	{{else}}
	<h4>相手のターンです。</h4>
	<p>相手が行動を選択するまで少しお待ちください。</p>
	<p><a href=".">更新</a></p>
	{{end}}

<!--
	<hr />
	<p>debug</p>
	{{.}}
	<hr />
-->
	</body>
</html>
`

const FinishHTML = `
<html>
	<body>
	<h1>勝負あり</h1>
	{{.GameStored.Winner.Name}}の勝ちです。
	</body>
</html>
`

const DebugHTML = `
<html>
	<body>
	{{.Winner.Name}}の勝利です。
	</body>
</html>
`
