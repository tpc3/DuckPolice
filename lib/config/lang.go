package config

type Strings struct {
	Lang     string
	Help     string
	CurrConf string
	Usage    usagestr
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

type configusagestr struct {
	Desc    string
	Prefix  string
	Lang    string
	Alert   string
	Domain  string
	Channel string
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
				Alert:  "通知の設定を行います。\ntype: reply, message, dmのいずれかを指定できます。\nmessage: 通知文を設定します。\nreact: 通知先メッセージにつけるリアクションを設定します。\nreject: 通知メッセージを削除するリアクションを指定します。",
				Domain: "検出するドメインリストの設定を行います。\nmode: white, blackのいずれかを指定できます。\nadd,delete: リストの内容を追加/削除します。",
				Channel: "チャンネルグループの設定を行います。\n" +
					"同じグループに含まれるチャンネル内でのみ検出を行います。\n" +
					"[selector]: 設定対象を指定します。チャンネルID,`channel`,`thread`のいずれかを設定できます。省略した場合コマンドを実行したチャンネルのIDが使用されます。\n" +
					"[value]: グループ名を指定します。selectorが`channel`,`thread`の場合は`categoryId`,`channelId`,`threadId`を指定できます。\n" +
					"メッセージが送信されたチャンネルID→チャンネル種類→上位チャンネル(カテゴリなど)のIDの順に検索が行われ、最初に一致したグループIDが使用されます。\n" +
					"一致がなければ\"default\"になります",
			},
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
		Error: errorstr{
			UnknownTitle: "Unexpected error is occurred.",
			UnknownDesc:  "This issue will be reported",
			NoCmd:        "No such command.",
			SubCmd:       "Invalid argument.",
			Syntax:       "Syntax error",
			SyntaxDesc:   "Failed to parse parameter.\nPlease check your command syntax.",
		},
	}
}
