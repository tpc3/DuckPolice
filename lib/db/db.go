package db

import (
	"database/sql"
	"encoding/json"
	"log"
	"regexp"
	"time"

	"github.com/tpc3/DuckPolice/lib/config"

	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var loadGuildConfigStmt *sql.Stmt
var setGuildConfigStmt *sql.Stmt

var searchLogStmt *sql.Stmt
var addLogStmt *sql.Stmt
var deleteLogStmt *sql.Stmt
var cleanOldLogStmt *sql.Stmt

func init() {
	var err error
	db, err = sql.Open(config.CurrentConfig.Db.Kind, config.CurrentConfig.Db.Parameter)
	if err != nil {
		log.Fatal("DB load error: ", err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS " + config.CurrentConfig.Db.Tableprefix + "config (" +
		"guild BIGINT NOT NULL PRIMARY KEY," +
		"data VARCHAR NOT NULL)")
	if err != nil {
		log.Fatal("Create guild config table error: ", err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS " + config.CurrentConfig.Db.Tableprefix + "log (" +
		"guild VARCHAR NOT NULL, " +
		"groupid VARCHAR NOT NULL, " +
		"content VARCHAR NOT NULL, " +
		"channelid BIGINT NOT NULL, " +
		"messageid BIGINT NOT NULL)")
	if err != nil {
		log.Fatal("Create log table error: ", err)
	}
	_, err = db.Exec("VACUUM")
	if err != nil {
		log.Fatal("DB VACUUM error: ", err)
	}
	loadGuildConfigStmt, err = db.Prepare("SELECT data FROM " + config.CurrentConfig.Db.Tableprefix + "config WHERE guild = ?")
	if err != nil {
		log.Fatal("Prepare loadGuildStmt error: ", err)
	}
	setGuildConfigStmt, err = db.Prepare("INSERT INTO " + config.CurrentConfig.Db.Tableprefix + "config" +
		"(guild,data) " +
		"VALUES (?,?) " +
		"ON CONFLICT (guild) " +
		"DO UPDATE set data = ?")
	if err != nil {
		log.Fatal("Prepare insertGuildStmt error: ", err)
	}

	searchLogStmt, err = db.Prepare("SELECT channelid, messageid FROM " + config.CurrentConfig.Db.Tableprefix + "log WHERE guild = ? AND groupid = ? AND content = ?")
	if err != nil {
		log.Fatal("Prepare searchLogStmt error: ", err)
	}

	addLogStmt, err = db.Prepare("INSERT INTO " + config.CurrentConfig.Db.Tableprefix + "log (guild,groupid,content,channelid,messageid) VALUES(?,?,?,?,?)")
	if err != nil {
		log.Fatal("Prepare addLogStmt error: ", err)
	}

	deleteLogStmt, err = db.Prepare("DELETE FROM " + config.CurrentConfig.Db.Tableprefix + "log WHERE guild = ? AND groupid = ? AND content = ?")
	if err != nil {
		log.Fatal("Prepare updateLogStmt error: ", err)
	}

	cleanOldLogStmt, err = db.Prepare("DELETE FROM " + config.CurrentConfig.Db.Tableprefix + "log WHERE messageid < ?")
	if err != nil {
		log.Fatal("Prepare cleanOldLogStmt error: ", err)
	}
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
	row := loadGuildConfigStmt.QueryRow(id)
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
	_, err = setGuildConfigStmt.Exec(id, data, data)
	if err != nil {
		log.Print("WARN: SaveGuild error: ", err)
	}

	compileGuild(guildCache[id])

	return err
}

func AddLog(orgMsg *discordgo.MessageCreate, guildId, groupId string, content *string, channelid, messageid string) {
	addLogStmt.Exec(guildId, groupId, content, channelid, messageid)
}

func SearchLog(session *discordgo.Session, guildId, groupId string, content *string) (found bool, channelid string, messageid string) {
	err := searchLogStmt.QueryRow(guildId, groupId, content).Scan(&channelid, &messageid)
	if err != nil {
		if err == sql.ErrNoRows {
			found = false
		} else {
			log.Fatal("Search Log error: ", err)
		}
	} else {
		_, err := session.State.Message(channelid, messageid)
		if err == discordgo.ErrStateNotFound {
			_, err = session.ChannelMessage(channelid, messageid)
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
	deleteLogStmt.Exec(guildId, groupId, content)
}

func timeToSnowflake(t time.Time) int64 {
	return (t.UnixMilli() - 1420070400000) << 22
}

func CleanOldLog() (int64, error) {
	res, err := cleanOldLogStmt.Exec(timeToSnowflake(time.Now().Add(-time.Duration(config.CurrentConfig.LogPeriod) * time.Second)))
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
