package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"regexp"
	"time"

	"github.com/tpc3/DuckPolice/lib/common"
	"github.com/tpc3/DuckPolice/lib/config"

	"github.com/bwmarrin/discordgo"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

var db *bun.DB

func InitDB() error {
	sqldb, err := sql.Open(sqliteshim.ShimName, config.CurrentConfig.Db.Parameter)
	if err != nil {
		return err
	}
	db = bun.NewDB(sqldb, sqlitedialect.New())
	_, err = db.NewCreateTable().Model((*common.Config)(nil)).IfNotExists().Exec(context.Background())
	if err != nil {
		return err
	}
	_, err = db.NewCreateTable().Model((*common.Log)(nil)).IfNotExists().Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func Close() {
	err := db.Close()
	if err != nil {
		log.Fatal("DB Close error", err)
	}
}

var guildCache = map[string]*config.Guild{}

func LoadGuild(id string) *config.Guild {
	val, exists := guildCache[id]
	if exists {
		return val
	}
	var rawData string
	guild := config.CurrentConfig.Guild
	row := db.QueryRow("SELECT data FROM config WHERE guild = ?", id)
	err := row.Scan(&rawData)

	if err == nil {
		json.Unmarshal([]byte(rawData), &guild)
	} else if err == sql.ErrNoRows {
		// skip
	} else {
		log.Fatal("LoadGuild scan error: ", err)
	}

	compileGuild(&guild)

	guildCache[id] = &guild
	return &guild
}

func compileGuild(guild *config.Guild) {
	guild.ParsedIgnore = make([]*regexp.Regexp, 0, len(guild.Ignore))
	for _, v := range guild.Ignore {
		regex, err := regexp.Compile(v)
		if err != nil {
			log.Print("Failed to compile regex: ", err)
			continue
		}
		guild.ParsedIgnore = append(guild.ParsedIgnore, regex)
	}
}

func SaveGuild(id string) error {
	data, err := json.Marshal(guildCache[id])
	if err != nil {
		log.Print("WARN: Marshal guild error: ", err)
		return err
	}
	src := common.Config{
		Guild: id,
		Data:  data,
	}
	_, err = db.NewInsert().Model(&src).On("Conflict (guild) DO UPDATE").Set("data = EXCLUDED.data").Exec(context.Background())
	if err != nil {
		log.Print("WARN: SaveGuild error: ", err)
	}

	compileGuild(guildCache[id])

	return err
}

func AddLog(orgMsg *discordgo.MessageCreate, guildId, groupId string, content *string, channelid, messageid string) {
	src := common.Log{
		Guild:     guildId,
		GroupID:   groupId,
		Content:   *content,
		ChannelID: channelid,
		MessageID: messageid,
	}
	_, err := db.NewInsert().Model(&src).Exec(context.Background())
	if err != nil {
		log.Print("WARN: adding log error: ", err)
	}
}

func SearchLog(session *discordgo.Session, guildId, groupId string, content *string) (found bool, dst common.Log) {
	src := common.Log{
		Guild:   guildId,
		GroupID: groupId,
		Content: *content,
	}
	_, err := db.NewSelect().Model(&src).Where("guild = ?", guildId).Where("groupid = ?", groupId).Where("content = ?", content).Exec(context.Background(), &dst)
	if err != nil {
		if err == sql.ErrNoRows {
			found = false
		} else {
			log.Fatal("Search Log error: ", err)
		}
	} else {
		_, err := session.State.Message(dst.ChannelID, dst.MessageID)
		if err == discordgo.ErrStateNotFound {
			_, err = session.ChannelMessage(dst.ChannelID, dst.MessageID)
		}
		if err == nil {
			found = true
			return
		} else {
			DeleteLog(guildId, groupId, content)
		}
	}
	return
}

func DeleteLog(guildId, groupId string, content *string) {
	src := common.Log{
		Guild:   guildId,
		GroupID: groupId,
		Content: *content,
	}
	_, err := db.NewDelete().Model(&src).Where("guild = ?", guildId).Where("groupid = ?").Where("content = ?").Exec(context.Background())
	if err != nil {
		log.Print("WARN: deleting log error: ", err)
	}
}

func timeToSnowflake(t time.Time) int64 {
	return (t.UnixMilli() - 1420070400000) << 22
}

func CleanOldLog() (int64, error) {
	res, err := db.NewDelete().Table("log").Where("messageid < ?", timeToSnowflake(time.Now().Add(-time.Duration(config.CurrentConfig.LogPeriod)*time.Second))).Exec(context.Background())
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func AutoLogCleaner() {
	for {
		_, err := CleanOldLog()
		if err != nil {
			log.Print("Failed to auto clean log: ", err)
		}
		nextTime := time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour)
		time.Sleep(time.Until(nextTime))
	}
}
