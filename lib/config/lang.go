package config

type Strings struct {
	Lang     string
	Help     string
	CurrConf string
	Usage    usagestr
	Clean    cleanstr
	Error    errorstr
}

type errorstr struct {
	Title        string
	UnknownTitle string
	UnknownDesc  string
	NoCmd        string
	SubCmd       string
	Syntax       string
	SyntaxDesc   string
}

type usagestr struct {
	Title  string
	Config configusagestr
}

type cleanstr struct {
	Title      string
	DeletedLog string
}

type configusagestr struct {
	Desc   string
	Prefix string
	Lang   string
}

var (
	Lang map[string]Strings
)

func loadLang() {
	Lang = map[string]Strings{}
	Lang["japanese"] = Strings{
		Lang:     "japanese",
		Help:     "Botの使い方に関しては、下記Wikiをご参照ください。",
		CurrConf: "現在の設定",
		Usage: usagestr{
			Title: "使い方: ",
			Config: configusagestr{
				Desc:   "各種設定を行います。\n設定項目と内容は以下をご確認ください。",
				Prefix: "コマンドの接頭詞を指定します。\nデフォルトは`" + CurrentConfig.Guild.Prefix + "`です。",
				Lang:   "言語を指定します。デフォルトは`" + CurrentConfig.Guild.Lang + "`です。",
			},
		},
		Clean: cleanstr{
			Title:      "実行結果",
			DeletedLog: "件の古いログを消去しました",
		},
		Error: errorstr{
			UnknownTitle: "予期せぬエラーが発生しました。",
			UnknownDesc:  "この問題は管理者に報告されます。",
			NoCmd:        "そのようなコマンドはありません。",
			SubCmd:       "引数が不正です。",
			Syntax:       "構文エラー",
			SyntaxDesc:   "パラメータの解析に失敗しました。\nコマンドの構文が正しいかお確かめください。",
		},
	}
	Lang["english"] = Strings{
		Lang:     "english",
		Help:     "Usage is available on the Wiki.",
		CurrConf: "Current config",
		Usage: usagestr{
			Title: "Usage: ",
			Config: configusagestr{
				Desc:   "Do configuration.\nItem list is below.",
				Prefix: "Specify command prefix.\nDefaults to `" + CurrentConfig.Guild.Prefix + "`",
				Lang:   "Specify language.\nDefaults to `" + CurrentConfig.Guild.Lang + "`",
			},
		},
		Clean: cleanstr{
			Title:      "Delete old log",
			DeletedLog: "old log",
		},
		Error: errorstr{
			UnknownTitle: "Unexpected error is occured.",
			UnknownDesc:  "This issue will be reported",
			NoCmd:        "No such command.",
			SubCmd:       "Invalid argument.",
			Syntax:       "Syntax error",
			SyntaxDesc:   "Failed to parse parameter.\nPlease check your command syntax.",
		},
	}
}
