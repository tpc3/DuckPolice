package common

import "github.com/uptrace/bun"

type Config struct {
	bun.BaseModel `bun:"table:config"`

	Guild string `bun:"guild,pk,notnull"`
	Data  []byte `bun:"data,notnull"`
}

type Log struct {
	bun.BaseModel `bun:"table:log"`

	Guild     string `bun:"guild,notnull"`
	GroupID   string `bun:"groupid,notnull"`
	Content   string `bun:"content,notnull"`
	ChannelID string `bun:"channelid,notnull"`
	MessageID string `bun:"messageid,notnull"`
}
