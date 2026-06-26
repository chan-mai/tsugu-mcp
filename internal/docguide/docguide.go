// Package docguideは相続登記の必要書類を案内する(根拠: docs/knowledge/05)
// 添付情報の4分類を、相続方法と相続人パターンに応じて組み立てる
package docguide

// Inputは案内の条件
type Input struct {
	Method                            string // legal(法定相続) | agreement(遺産分割協議) | will(遺言)
	HeirPattern                       string // children | substitution | ascendants | siblings | siblings_substitution
	RegistryAddressDiffersFromHonseki bool   // 登記上の住所≠本籍(同一性証明の要否)
	UseLegalInfoNumber                bool   // 法定相続情報番号/一覧図で戸籍束を省略
	ApplicantAtWindow                 bool   // 本人が窓口で広域交付を使う(注意喚起の要否)
}

// Categoryは添付情報の1分類
type Category struct {
	Name  string
	Items []string
}

// Resultは必要書類の案内
type Result struct {
	Categories []Category
	Notes      []string
	Fees       []string
}

// RequiredDocumentsは条件から必要書類リストを組み立てる
func RequiredDocuments(in Input) Result {
	var res Result

	cause := []string{
		"被相続人の出生から死亡までの連続した戸籍(除籍・改製原戸籍を含む)",
		"相続人全員の現在の戸籍(被相続人の死亡日以降に取得)",
	}
	switch in.HeirPattern {
	case "substitution":
		cause = append(cause, "先に死亡した子の出生から死亡までの連続戸籍(孫=代襲相続人の確定)")
	case "ascendants":
		cause = append(cause, "子・孫が存在しないことの確認(被相続人の連続戸籍による。直系尊属に代襲なし)")
	case "siblings":
		cause = append(cause, "被相続人の父・母それぞれの出生から死亡までの連続戸籍(兄弟姉妹の確定・直系尊属全員死亡の確認)")
	case "siblings_substitution":
		cause = append(cause,
			"被相続人の父・母それぞれの出生から死亡までの連続戸籍(兄弟姉妹の確定・直系尊属全員死亡の確認)",
			"先に死亡した兄弟姉妹の出生から死亡までの連続戸籍(甥姪の確定。一代限り)")
	}
	switch in.Method {
	case "agreement":
		cause = append(cause,
			"遺産分割協議書(相続人全員が実印で押印)",
			"相続人全員の印鑑証明書(各1通。有効期限なし)")
	case "will":
		cause = append(cause, "遺言書(自筆証書は家庭裁判所の検認済み。公正証書・法務局保管の自筆証書は検認不要)")
	}
	res.Categories = append(res.Categories, Category{"登記原因証明情報(相続人の確定)", cause})

	res.Categories = append(res.Categories, Category{
		"住所証明情報", []string{"相続する人(新名義人)全員の住民票の写し(マイナンバー記載なし・原本)"},
	})

	if in.RegistryAddressDiffersFromHonseki {
		res.Categories = append(res.Categories, Category{
			"同一性証明", []string{"被相続人の住民票の除票または戸籍の附票(登記上の住所と本籍が異なる場合)"},
		})
	}

	res.Categories = append(res.Categories, Category{
		"課税価格の根拠", []string{"固定資産評価証明書・固定資産課税明細書・名寄帳のいずれか(最新年度の『価格』『評価額』欄)"},
	})

	if (in.HeirPattern == "siblings" || in.HeirPattern == "siblings_substitution") && in.ApplicantAtWindow {
		res.Notes = append(res.Notes, "戸籍の広域交付は傍系(兄弟姉妹)の戸籍・電子化前の改製原戸籍・戸籍の附票が対象外。これらは本籍地へ請求(郵送可)")
	}
	if in.Method == "agreement" {
		res.Notes = append(res.Notes, "遺産分割協議書に添付する印鑑証明書に有効期限はない(売買登記や金融機関手続と異なる)")
	}
	if in.UseLegalInfoNumber {
		res.Notes = append(res.Notes, "法定相続情報一覧図の写し(または法定相続情報番号)で戸籍束の添付を省略でき、住所記載があれば住所証明情報も兼ねられる")
	}

	res.Fees = []string{
		"戸籍全部/個人事項証明書(謄抄本) 450円",
		"除籍全部事項証明書(除籍謄本) 750円",
		"改製原戸籍謄本 750円",
		"戸籍の附票の写し 約300円",
		"住民票の写し 約300円",
		"印鑑登録証明書 約300円",
	}
	return res
}
