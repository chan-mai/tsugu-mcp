# tsugu-mcp

相続手続き書類をGoで生成するツール。A4縦1ページの**相続関係説明図**を、JSONからPDF出力する。

体裁は実務の伝統的様式に準拠する: 左→右の横型家系図、人物は枠で囲わず`(続柄)　氏名`で表記、婚姻は縦の二重線、親子はブラケット線。日付は和暦表示。被相続人の左に詳細(最後の住所・本籍・出生・死亡)、左下に作成日・作成者、右下に「以下余白」。

## 使い方

### CLI

```sh
go build -o tsugu ./cmd/tsugu

./tsugu -in testdata/sample_full.json -out chart.pdf      # ファイル
./tsugu -era both < family.json > chart.pdf                # 標準入出力 / 西暦併記
```

| フラグ | 既定 | 内容 |
|--------|------|------|
| `-in`  | 標準入力 | 入力JSONのパス |
| `-out` | 標準出力 | 出力PDFのパス |
| `-era` | `wareki` | 日付表記: `wareki`(和暦) / `both`(和暦+西暦) / `seireki` |

### ライブラリ

```go
pdf, err := relationchart.GenerateFromJSON(data, relationchart.DefaultOptions())
// もしくはモデルから直接
pdf, err := relationchart.Generate(doc, relationchart.Options{Era: relationchart.EraWareki})
```

## 入力JSON

人物は所属(spouse / children / descendants / ascendants / siblings)とネストで関係を表す。各 children/siblings ノードは自分の `spouse` と `descendants`(代襲)を持てる。

```jsonc
{
  "decedent": {                          // 被相続人(必須: name, deathDate)
    "name": "山田 太郎",
    "registeredDomicile": "東京都千代田区一番町1番地",  // 本籍
    "lastAddress": "東京都千代田区一番町1番1号",        // 最後の住所
    "birthDate": "1940-01-01",           // YYYY-MM-DD
    "deathDate": "2024-06-15"
  },
  "spouse": { "name": "山田 花子", "relationship": "妻" },
  "children": [
    { "name": "山田 一郎", "relationship": "長男" },
    { "name": "山田 次郎", "relationship": "二男", "deathDate": "2020-03-10",
      "spouse": { "name": "山田 梅子", "relationship": "妻" },
      "descendants": [                                   // 代襲(任意世代)
        { "name": "山田 三郎", "relationship": "長男" },
        { "name": "佐藤 良子", "relationship": "二女", "applicant": true,
          "address": "東京都新宿区西新宿2番2号", "birthDate": "1975-05-05",
          "outcome": "inherit" }
      ] }
  ],
  "ascendants": [ { "name": "山田 祖一", "relationship": "父" } ],  // 直系尊属(父母)
  "siblings":   [ { "name": "山田 春夫", "relationship": "弟" } ],  // 兄弟姉妹
  "preparer": { "address": "東京都新宿区西新宿2番2号", "name": "佐藤 良子" },
  "preparedAt": "2024-12-10"
}
```

- `outcome`(注記): `inherit`→`(相)` / `renounce`→`(相続放棄)` / `division`→`(分割)` / `by_representation`→`(代襲)`。和文(`相続`等)も可。
- `applicant: true` で `(申立人)` を付す。
- `relationship`(続柄)は表示用の自由文字列。役割は配置(children/ascendants/siblings, descendants)とネストで表す。
- `ascendants` があると尊属を根とし、被相続人と `siblings` をその子として描く。
- 本ツールは与えられた家族構成を作図するレンダラに徹し、**法定相続人や相続分の算定は行わない**。

## 設計

データを一方向に流す層構造とし、各層を描画ライブラリ非依存で単体テスト可能にしている。

```
JSON ─[inputjson]→ family.Document ─[layout]→ scene.Scene ─[render]→ PDF
      DTO境界          ドメインモデル     純粋な幾何      Canvas越し描画
```

| パッケージ | 責務 |
|------------|------|
| `family` | 入力ドメインモデルと意味的検証。描画・JSON非依存 |
| `internal/inputjson` | JSONとモデルの境界。日付書式・列挙語彙の解釈 |
| `internal/wareki` | 西暦→和暦変換(元号テーブル駆動の純粋関数) |
| `internal/scene` | 描画指示の中間表現(枠・線・ラベル, mm座標) |
| `internal/layout` | モデル→Scene。横型ツリー配置とAutoFit。`Measurer`で文字幅を注入 |
| `internal/render` | Scene→PDF。`Canvas`インターフェースでgopdf依存を隔離 |
| `relationchart` | 公開API。CLI・将来のMCPサーバーが共通利用 |

レイアウトは Document を1本の家系ツリー(`tree.go`)へ変換し、左→右へ世代を列に割り当てて配置する。人物は枠なしカード(`card.go`)、関係線は縦の婚姻二重線と親子ブラケット。純粋な幾何計算のため `HeuristicMeasurer` で実フォント無しに座標を検証でき、描画は `Canvas` 抽象越しのためフェイクで呼び出しを検証できる。

## フォント

日本語表示にIPAexゴシック(TrueType)を`go:embed`で同梱し、PDFへサブセット埋め込みする。ライセンスは[IPAフォントライセンスv1.0](assets/fonts/IPA_Font_License_Agreement_v1.0.txt)(埋め込み・再配布可)。

## 開発

```sh
go test ./...   # 全層の単体・統合テスト
go vet ./...
```
