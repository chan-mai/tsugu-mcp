package mcpserver

// MCPツールの型付き入力(jsonタグはinputjson/reginputのDTOと一致)
// jsonschemaタグでLLM向けの項目説明を付与

// --- 相続関係説明図 ---

type chartDecedent struct {
	Name               string `json:"name" jsonschema:"被相続人の氏名(必須)"`
	RegisteredDomicile string `json:"registeredDomicile,omitempty" jsonschema:"本籍"`
	LastAddress        string `json:"lastAddress,omitempty" jsonschema:"最後の住所"`
	RegistryAddress    string `json:"registryAddress,omitempty" jsonschema:"登記上の住所"`
	BirthDate          string `json:"birthDate,omitempty" jsonschema:"生年月日 YYYY-MM-DD"`
	DeathDate          string `json:"deathDate" jsonschema:"死亡日 YYYY-MM-DD(必須)"`
}

type chartPerson struct {
	Name         string `json:"name" jsonschema:"氏名"`
	Relationship string `json:"relationship,omitempty" jsonschema:"続柄(例:長男・妻)"`
	Address      string `json:"address,omitempty" jsonschema:"住所"`
	BirthDate    string `json:"birthDate,omitempty" jsonschema:"生年月日 YYYY-MM-DD"`
	DeathDate    string `json:"deathDate,omitempty" jsonschema:"死亡日 YYYY-MM-DD"`
	Outcome      string `json:"outcome,omitempty" jsonschema:"注記 inherit|renounce|division|by_representation"`
	Applicant    bool   `json:"applicant,omitempty" jsonschema:"申立人なら true"`
}

// 代襲の子孫(孫世代)
// jsonschema-go再帰非対応のため型は2世代まで(生成器GenerateFromJSONは任意深さ対応)
type chartGrandchild struct {
	Name         string       `json:"name" jsonschema:"氏名"`
	Relationship string       `json:"relationship,omitempty" jsonschema:"続柄(例:孫)"`
	Address      string       `json:"address,omitempty" jsonschema:"住所"`
	BirthDate    string       `json:"birthDate,omitempty" jsonschema:"生年月日 YYYY-MM-DD"`
	DeathDate    string       `json:"deathDate,omitempty" jsonschema:"死亡日 YYYY-MM-DD"`
	Outcome      string       `json:"outcome,omitempty" jsonschema:"注記 inherit|renounce|division|by_representation"`
	Applicant    bool         `json:"applicant,omitempty" jsonschema:"申立人なら true"`
	Spouse       *chartPerson `json:"spouse,omitempty" jsonschema:"配偶者"`
}

type chartChild struct {
	Name         string            `json:"name" jsonschema:"氏名"`
	Relationship string            `json:"relationship,omitempty" jsonschema:"続柄(例:長男)"`
	Address      string            `json:"address,omitempty" jsonschema:"住所"`
	BirthDate    string            `json:"birthDate,omitempty" jsonschema:"生年月日 YYYY-MM-DD"`
	DeathDate    string            `json:"deathDate,omitempty" jsonschema:"死亡日 YYYY-MM-DD"`
	Outcome      string            `json:"outcome,omitempty" jsonschema:"注記 inherit|renounce|division|by_representation"`
	Applicant    bool              `json:"applicant,omitempty" jsonschema:"申立人なら true"`
	Spouse       *chartPerson      `json:"spouse,omitempty" jsonschema:"配偶者"`
	Descendants  []chartGrandchild `json:"descendants,omitempty" jsonschema:"代襲の子孫(孫)"`
}

type chartPreparer struct {
	Address string `json:"address,omitempty" jsonschema:"作成者の住所"`
	Name    string `json:"name,omitempty" jsonschema:"作成者の氏名"`
}

type chartDoc struct {
	Decedent   chartDecedent `json:"decedent" jsonschema:"被相続人"`
	Spouse     *chartPerson  `json:"spouse,omitempty" jsonschema:"配偶者"`
	Ascendants []chartPerson `json:"ascendants,omitempty" jsonschema:"直系尊属(父母)"`
	Children   []chartChild  `json:"children,omitempty" jsonschema:"子(代襲は descendants)"`
	Siblings   []chartChild  `json:"siblings,omitempty" jsonschema:"兄弟姉妹"`
	Preparer   chartPreparer `json:"preparer,omitempty" jsonschema:"作成者"`
	PreparedAt string        `json:"preparedAt,omitempty" jsonschema:"作成日 YYYY-MM-DD"`
}

// --- 相続登記申請書 ---

type toukiCause struct {
	Date string `json:"date,omitempty" jsonschema:"原因日 YYYY-MM-DD"`
	Text string `json:"text" jsonschema:"文言(例:相続)"`
}

type toukiApplicant struct {
	Name      string `json:"name" jsonschema:"氏名(必須)"`
	Address   string `json:"address,omitempty" jsonschema:"住所"`
	Share     string `json:"share,omitempty" jsonschema:"持分(例:2分の1。単独なら空)"`
	NameKana  string `json:"nameKana,omitempty" jsonschema:"氏名ふりがな"`
	BirthDate string `json:"birthDate,omitempty" jsonschema:"生年月日 YYYY-MM-DD"`
	Email     string `json:"email,omitempty" jsonschema:"メールアドレス"`
	Phone     string `json:"phone,omitempty" jsonschema:"連絡先電話番号"`
	Contact   bool   `json:"contact,omitempty" jsonschema:"連絡先表を付す代表者なら true"`
}

type landRight struct {
	Symbol      string `json:"symbol,omitempty" jsonschema:"符号"`
	LocationLot string `json:"locationLot,omitempty" jsonschema:"所在及び地番"`
	Category    string `json:"category,omitempty" jsonschema:"地目"`
	Area        string `json:"area,omitempty" jsonschema:"地積(平方メートルは自動付与)"`
	RightType   string `json:"rightType,omitempty" jsonschema:"敷地権の種類(例 所有権)"`
	RightShare  string `json:"rightShare,omitempty" jsonschema:"敷地権の割合(例 1000分の35)"`
}

type toukiProperty struct {
	Kind         string      `json:"kind" jsonschema:"land(土地) / building(建物) / condominium(区分建物)"`
	Number       string      `json:"number,omitempty" jsonschema:"不動産番号"`
	Location     string      `json:"location" jsonschema:"所在(区分建物では一棟の建物の所在)"`
	LotNumber    string      `json:"lotNumber,omitempty" jsonschema:"地番(土地)"`
	LandCategory string      `json:"landCategory,omitempty" jsonschema:"地目(土地)"`
	Area         string      `json:"area,omitempty" jsonschema:"地積(土地。平方メートルは自動付与)"`
	HouseNumber  string      `json:"houseNumber,omitempty" jsonschema:"家屋番号(建物・区分建物の専有部分)"`
	BuildingType string      `json:"buildingType,omitempty" jsonschema:"種類(建物・区分建物)"`
	Structure    string      `json:"structure,omitempty" jsonschema:"構造(建物・区分建物)"`
	FloorArea    string      `json:"floorArea,omitempty" jsonschema:"床面積(建物・区分建物)"`
	BuildingName string      `json:"buildingName,omitempty" jsonschema:"一棟の建物の名称(区分建物)"`
	UnitName     string      `json:"unitName,omitempty" jsonschema:"専有部分の建物の名称(区分建物。例 301号)"`
	LandRights   []landRight `json:"landRights,omitempty" jsonschema:"敷地権の表示(区分建物)"`
	Value        int         `json:"value,omitempty" jsonschema:"固定資産評価額(円)。指定すると登録免許税を自動計算"`
	Exemption    string      `json:"exemption,omitempty" jsonschema:"免税 none|small_value|intermediate(登録免許税の自動計算時)"`
}

type toukiDecedent struct {
	Name string `json:"name" jsonschema:"被相続人の氏名"`
}

type toukiDoc struct {
	Causes          []toukiCause     `json:"causes,omitempty" jsonschema:"原因(数次相続で複数併記)"`
	Decedent        toukiDecedent    `json:"decedent" jsonschema:"被相続人"`
	Applicants      []toukiApplicant `json:"applicants" jsonschema:"申請人(1名以上)"`
	Attachments     []string         `json:"attachments,omitempty" jsonschema:"添付情報"`
	DeclineIDInfo   bool             `json:"declineIdInfo,omitempty" jsonschema:"登記識別情報の通知を希望しない欄。trueでチェック"`
	ApplicationDate string           `json:"applicationDate,omitempty" jsonschema:"申請日 YYYY-MM-DD"`
	Registry        string           `json:"registry,omitempty" jsonschema:"法務局"`
	TaxValue        string           `json:"taxValue,omitempty" jsonschema:"課税価格(表示文字列。算定はしない)"`
	RegistrationTax string           `json:"registrationTax,omitempty" jsonschema:"登録免許税(表示文字列。算定はしない)"`
	Properties      []toukiProperty  `json:"properties" jsonschema:"不動産の表示(1件以上)"`
}

// --- ツール入出力 ---

type chartToolInput struct {
	Document   chartDoc `json:"document" jsonschema:"相続関係説明図の内容"`
	OutputPath string   `json:"outputPath,omitempty" jsonschema:"出力PDFのパス(省略時は一時ファイル)"`
	Era        string   `json:"era,omitempty" jsonschema:"日付表記 wareki|both|seireki(既定 wareki)"`
}

type toukiToolInput struct {
	Document   toukiDoc `json:"document" jsonschema:"相続登記申請書の内容"`
	OutputPath string   `json:"outputPath,omitempty" jsonschema:"出力PDFのパス(省略時は一時ファイル)"`
	Era        string   `json:"era,omitempty" jsonschema:"日付表記 wareki|both|seireki(既定 wareki)"`
}

type toolResult struct {
	Path  string `json:"path" jsonschema:"生成したPDFの絶対パス"`
	Bytes int    `json:"bytes" jsonschema:"PDFのバイト数"`
}
