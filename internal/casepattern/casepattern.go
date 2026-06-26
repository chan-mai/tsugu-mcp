// Package casepatternは相続登記申請書のケース別様式を選択する
// 相続方法と特殊事情から、登記の目的・原因・申請構造・税率・添付の要点を返す
package casepattern

// 選択条件
type Input struct {
	Method           string // legal | agreement | will_specified | bequest
	BequestToHeir    bool   // 遺贈で受遺者が相続人(D-1)か第三者(D-2)
	Multilevel       bool   // 数次相続(E)
	Substitution     bool   // 代襲相続(F)
	Renunciation     bool   // 相続放棄あり(G)
	SpousalResidence bool   // 配偶者居住権の設定(H)
}

// ケース別の様式要点
type Pattern struct {
	Key        string
	Name       string
	Purpose    string // 登記の目的
	Cause      string // 原因
	Structure  string // 申請構造
	TaxRate    string // 税率
	OriginInfo string // 登記原因証明情報の要点
	Caveat     string // 確実性の注意(空可)
}

// 選択結果
type Result struct {
	Primary   Pattern
	Modifiers []string // 数次・代襲・放棄・居住権の追加注意
	Note      string
}

var primary = map[string]Pattern{
	"A": {"A", "法定相続分による相続", "所有権移転", "相続(死亡日)",
		"相続人全員に法定相続分の持分を記載。相続人1人が保存行為として全員のため申請可(委任状不要、申請人にならない者に登記識別情報は通知されない)",
		"4/1000", "被相続人の出生〜死亡の戸除籍+相続人全員の戸籍。遺産分割協議書・印鑑証明書は不要", ""},
	"B": {"B", "遺産分割協議による相続", "所有権移転", "相続(死亡日。協議成立日ではない)",
		"取得する相続人のみ記載。単独取得は持分記載不要、共有取得は各人に持分。取得しない相続人は記載しない",
		"4/1000", "戸除籍一式+遺産分割協議書+協議した相続人全員の印鑑証明書。住所証明は取得者のみ", ""},
	"C": {"C", "特定財産承継遺言(相続させる遺言)", "所有権移転", "相続(遺贈ではなく相続。死亡日)",
		"当該相続人の単独申請",
		"4/1000", "遺言書+被相続人の死亡記載の戸除籍+相続人の戸籍。自筆証書は家裁の検認済(法務局保管は検認不要)", ""},
	"D1": {"D-1", "相続人に対する遺贈", "所有権移転", "遺贈(遺贈者の死亡日)",
		"受遺者である相続人の単独申請(令和5年改正)。義務者の押印・印鑑証明書・登記識別情報は不要",
		"4/1000", "遺言書+遺贈者の死亡記載の戸除籍+受遺者が相続人と分かる戸籍+受遺者の住所証明", ""},
	"D2": {"D-2", "第三者(相続人以外)への遺贈", "所有権移転", "遺贈(遺贈者の死亡日)",
		"共同申請。権利者=受遺者(第三者)、義務者=遺言執行者がいればその者、いなければ相続人全員",
		"20/1000", "遺言書+遺贈者の死亡記載の戸除籍+義務者の登記識別情報(権利証)+義務者の印鑑証明書(3か月以内)+遺言執行者の資格証明+受遺者の住所証明",
		"公式専用様式がなく実務上確立。要確認"},
}

// 条件からケース別様式を選択
func Select(in Input) Result {
	res := Result{Note: "ケースの目的・原因・申請構造・添付の要点。最終確認は司法書士へ"}

	switch in.Method {
	case "legal":
		res.Primary = primary["A"]
	case "agreement":
		res.Primary = primary["B"]
	case "will_specified":
		res.Primary = primary["C"]
	case "bequest":
		if in.BequestToHeir {
			res.Primary = primary["D1"]
		} else {
			res.Primary = primary["D2"]
		}
	default:
		res.Primary = Pattern{Key: "?", Name: "unknown method (legal|agreement|will_specified|bequest)"}
	}

	if in.Multilevel {
		res.Modifiers = append(res.Modifiers, "数次相続(E): 原因を連記。中間が単独相続なら1件で登記名義人から最終相続人へ直接移転可、複数共有なら相続ごとに別件申請。第1被相続人の出生〜死亡+中間・最終相続人全員の戸籍が必要")
	}
	if in.Substitution {
		res.Modifiers = append(res.Modifiers, "代襲相続(F): 原因は単一日付(被相続人の死亡日)。被代襲者(先に死亡した子)の出生〜死亡の連続戸除籍が必須。申請書に代襲の特別注記は不要")
	}
	if in.Renunciation {
		res.Modifiers = append(res.Modifiers, "相続放棄(G): 放棄者は申請書に記載せず相続分を再計算。家裁の相続放棄申述受理証明書を登記原因証明情報に追加")
	}
	if in.SpousalResidence {
		res.Modifiers = append(res.Modifiers, "配偶者居住権(H): 別途『配偶者居住権設定』の共同申請(権利者=配偶者、義務者=建物取得者)。建物のみ・存続期間必須・税率2/1000。前提として建物の相続登記が必要")
	}
	return res
}
