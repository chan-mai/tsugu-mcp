# tsugu-mcp

相続手続き書類をGoで生成するツール。JSONから次のPDFを出力する。日付は和暦表示(`-era`で切替)。

- **相続関係説明図**(1ページ): 左→右の横型家系図、人物は枠で囲わず`(続柄)　氏名`で表記、婚姻は縦の二重線、親子はブラケット線。被相続人の左に詳細(最後の住所・本籍・出生・死亡)、左下に作成日・作成者、右下に「以下余白」。
- **相続登記申請書**(所有権移転・複数ページ): 1枚目上部に受付番号表用の破線枠、登記の目的・原因・相続人・申請人(複数+持分)・添付情報・申請日・課税価格・登録免許税、続いて不動産の表示(土地/建物)を流し込み。件数が多い場合は自動でページ送り。

## 使い方

### CLI

```sh
go build -o tsugu ./cmd/tsugu

./tsugu chart -in testdata/sample_full.json  -out chart.pdf   # 相続関係説明図
./tsugu touki -in testdata/touki_sample.json -out touki.pdf   # 相続登記申請書
./tsugu chart -era both < family.json > chart.pdf             # 標準入出力 / 西暦併記
```

サブコマンド `chart`(相続関係説明図) / `touki`(相続登記申請書)。各々共通フラグ:

| フラグ | 既定 | 内容 |
|--------|------|------|
| `-in`  | 標準入力 | 入力JSONのパス |
| `-out` | 標準出力 | 出力PDFのパス |
| `-era` | `wareki` | 日付表記: `wareki`(和暦) / `both`(和暦+西暦) / `seireki` |

### ライブラリ

```go
// 相続関係説明図
pdf, err := relationchart.GenerateFromJSON(chartJSON, relationchart.DefaultOptions())
// 相続登記申請書
pdf, err := registration.GenerateFromJSON(toukiJSON, registration.DefaultOptions())
```

## 入力JSON(相続関係説明図)

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

## 入力JSON(相続登記申請書)

登記の目的は「所有権移転」固定のためJSONに無い。`applicants`は複数+持分、`properties`は土地/建物を混在可。件数が多いと自動でページ送り。

```jsonc
{
  "causes": [                                  // 原因(数次相続で複数併記)
    { "date": "2020-03-10", "text": "山田太郎相続" }, { "date": "2024-06-15", "text": "相続" }
  ],
  "decedent": { "name": "山田 一郎", "address": "東京都千代田区一番町1番地" }, // 住所は死亡時の本籍
  "applicants": [
    { "name": "山田 花子", "address": "東京都新宿区西新宿二丁目2番2号", "share": "2分の1",
      "nameKana": "やまだ はなこ", "birthDate": "1950-05-05",
      "email": "hanako@example.com", "phone": "090-0000-0000", "contact": true },
    { "name": "山田 次郎", "address": "神奈川県横浜市中区本町三丁目3番3号", "share": "2分の1" }
  ],
  "attachments": ["登記原因証明情報", "住所証明情報"],   // 改行+インデントで列挙
  "declineIdInfo": false,                       // 登記識別情報の通知欄(常に表示): false=□ / true=☑
  "applicationDate": "2024-12-10", "registry": "東京法務局",
  "taxValue": "0", "registrationTax": "0",
  "properties": [
    { "kind": "land",     "number": "...", "location": "...", "lotNumber": "1番", "landCategory": "宅地", "area": "123.45" },// numberは 土地家屋償却資産（補充）課税台帳（名寄帳） の場合、物件番号として記載される場合もあるみたい
    { "kind": "building", "number": "...", "location": "...", "houseNumber": "1番1",
      "buildingType": "居宅", "structure": "木造2階建", "floorArea": "1階 60.00平方メートル" }
  ]
}
```

- `kind`: `land`(土地: 地番/地目/地積) / `building`(建物: 家屋番号/種類/構造/床面積)。
- `contact: true` の申請人に氏名ふりがな・生年月日(西暦)・メールの枠表と連絡先電話を付す。
- 課税価格・登録免許税・地積等は表示文字列。本ツールは算定せず忠実描画。

## 設計

データを一方向に流す層構造とし、各層を描画ライブラリ非依存で単体テスト可能にしている。2書類は中間表現と描画を共有し、モデル・JSON境界・レイアウトを書類ごとに分離する。

```
相続関係説明図:  JSON ─[inputjson]→ family.Document  ─[layout]───→ scene.Scene    ─┐
相続登記申請書:  JSON ─[reginput]──→ touki.Application ─[reglayout]→ []scene.Scene  ─┤─[render]→ PDF
共有: ymd / wareki / scene / render
```

| パッケージ | 責務 |
|------------|------|
| `ymd` | 年月日Dateと解析(全書類で共有) |
| `internal/wareki` | 西暦→和暦変換(元号テーブル駆動の純粋関数) |
| `internal/scene` | 描画指示の中間表現(枠・線・ラベル, mm座標) |
| `internal/render` | Scene→PDF。`Canvas`でgopdf依存を隔離。`ToPDFMulti`で複数ページ |
| `family` / `internal/inputjson` / `internal/layout` / `relationchart` | 相続関係説明図: モデル / JSON境界 / 横型ツリー配置 / 公開API |
| `touki` / `internal/reginput` / `internal/reglayout` / `registration` | 相続登記申請書: モデル / JSON境界 / 流し込み+ページ送り / 公開API |

レイアウトは Document を1本の家系ツリー(`tree.go`)へ変換し、左→右へ世代を列に割り当てて配置する。人物は枠なしカード(`card.go`)、関係線は縦の婚姻二重線と親子ブラケット。純粋な幾何計算のため `HeuristicMeasurer` で実フォント無しに座標を検証でき、描画は `Canvas` 抽象越しのためフェイクで呼び出しを検証できる。

## フォント

日本語表示にIPAexゴシック(TrueType)を`go:embed`で同梱し、PDFへサブセット埋め込みする。ライセンスは[IPAフォントライセンスv1.0](assets/fonts/IPA_Font_License_Agreement_v1.0.txt)(埋め込み・再配布可)。

## 開発

```sh
go test ./...   # 全層の単体・統合テスト
go vet ./...
```
