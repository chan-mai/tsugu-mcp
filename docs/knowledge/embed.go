// Package knowledgeは相続登記ナレッジベースの埋め込みリソース
package knowledge

import "embed"

//go:embed *.md
var FS embed.FS

// Docは公開する1知識文書のメタ
type Doc struct {
	Slug        string // URIのslug(knowledge://<slug>)
	File        string // FS内のファイル名
	Name        string
	Description string
}

// Docsは公開する知識文書の一覧
var Docs = []Doc{
	{"index", "README.md", "索引", "目次・施行日年表・間違えやすい最新事項"},
	{"01-seido-houkaisei", "01-seido-houkaisei.md", "制度と法改正", "相続登記義務化・相続人申告登記・住所氏名変更義務化・国庫帰属・所有不動産記録証明"},
	{"02-tetsuzuki-flow", "02-tetsuzuki-flow.md", "手続フロー全体像", "相続パターン判定・申請方式・管轄・費用・オンライン/半ライン申請・完了"},
	{"03-toukishinseisho", "03-toukishinseisho.md", "登記申請書", "各項目の記載要件・ケース別A〜H・区分建物・物理的作成ルール"},
	{"04-tourokumenkyozei", "04-tourokumenkyozei.md", "登録免許税", "計算・端数処理・免税措置・税率表・納付・計算例"},
	{"05-hitsuyou-shorui-koseki", "05-hitsuyou-shorui-koseki.md", "必要書類と戸籍収集", "添付情報4分類・パターン別戸籍範囲・広域交付・法定相続情報証明・住所評価書類・原本還付"},
	{"06-fuzui-shorui", "06-fuzui-shorui.md", "付随書類", "遺産分割協議書・委任状・相続関係説明図・法定相続情報一覧図"},
	{"07-souzokunin-soutokubun", "07-souzokunin-soutokubun.md", "相続人の範囲と法定相続分", "判定アルゴリズム・民法条文・死亡日別年表"},
}
