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
	RegexCompile string
}

type usagestr struct {
	Title  string
	Config configusagestr
}

type configusagestr struct {
	Desc    string
	Prefix  string
	Lang    string
	Ignore  string
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
				Ignore: "無視するメッセージ内容を正規表現で指定します。メッセージ全体と指定された正規表現のいずれかが一致した場合は、処理をスキップします。",
				Alert:  "通知の設定を行います。\ntype: reply, message, dmのいずれかを指定できます。\nmessage: 通知文を設定します。\nreact: 通知先メッセージにつけるリアクションを設定します。\nreject: 通知メッセージを削除するリアクションを指定します。",
				Domain: "検出するドメインリストの設定を行います。\nmode: white, blackのいずれかを指定できます。\nadd,del: リストの内容を追加/削除します。",
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
			RegexCompile: "正規表現の処理に失敗しました。",
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
				Ignore: "Specify regexp that ignore message. If entire message matched, ignore message.",
				Alert:  "Setting about alert.\ntype: `reply`, `message`, `dm`\nmessage: message to send\nreact: emoji to reaction message alerted\nreject: reaction emoji to reject(delete) alert",
				Domain: "domain list to used to detection\nmode: `white`, `black`\nadd,del: add or delete list item",
				Channel: "Setting channel group\n" +
					"detection only work in same channel group\n" +
					"[selector]: specify channel to be set. channel ID,`channel`,`thread` can be used. if omit, channel id which command sent is used.\n" +
					"[value]: Specify group name. If selector is `channel` or `thread`, `categoryId`,`channelId`,`threadId` can be set.\n" +
					"Check channel id which message sent → type of channel → parent channels id. Which matched first to be used.\n" +
					"If no match, \"default\" to be used.",
			},
		},
		Error: errorstr{
			UnknownTitle: "Unexpected error is occurred.",
			UnknownDesc:  "This issue will be reported",
			NoCmd:        "No such command.",
			SubCmd:       "Invalid argument.",
			Syntax:       "Syntax error",
			SyntaxDesc:   "Failed to parse parameter.\nPlease check your command syntax.",
			RegexCompile: "Failed to parse regexp",
		},
	}
}
