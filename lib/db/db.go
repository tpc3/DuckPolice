package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/tpc3/DuckPolice/lib/config"

	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var getDBverStmt *sql.Stmt
var loadGuildStmt *sql.Stmt
var insertGuildStmt *sql.Stmt
var updateGuildStmt *sql.Stmt

var addLogStmt map[string]*sql.Stmt
var cleanOldLogStmt map[string]*sql.Stmt

var guild_loaded bool

const db_version = 1

func init() {
	var err error
	db, err = sql.Open(config.CurrentConfig.Db.Kind, config.CurrentConfig.Db.Parameter)
	if err != nil {
		log.Fatal("DB load error: ", err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS " + config.CurrentConfig.Db.Tableprefix + "guilds (" +
		"id BIGINT NOT NULL PRIMARY KEY," +
		"db_version INT NOT NULL," +
		"prefix VARCHAR," +
		"lang VARCHAR)")
	if err != nil {
		log.Fatal("Create guild table error: ", err)
	}
	_, err = db.Exec("ALTER TABLE " + config.CurrentConfig.Db.Tableprefix + "guilds DROP COLUMN bots")
	if err == nil {
		log.Print("WARN: Database update successfully")
	} else {
		log.Print("FINE: Guilds DB up-to-date!")
	}
	_, err = db.Exec("VACUUM")
	if err != nil {
		log.Fatal("DB VACUUM error: ", err)
	}
	getDBverStmt, err = db.Prepare("SELECT db_version FROM " + config.CurrentConfig.Db.Tableprefix + "guilds WHERE " + "id = ?")
	if err != nil {
		log.Fatal("Prepare getDBverStmt error: ", err)
	}
	loadGuildStmt, err = db.Prepare("SELECT * FROM " + config.CurrentConfig.Db.Tableprefix + "guilds WHERE " + "id = ?")
	if err != nil {
		log.Fatal("Prepare loadGuildStmt error: ", err)
	}
	insertGuildStmt, err = db.Prepare("INSERT INTO " + config.CurrentConfig.Db.Tableprefix + "guilds(id,db_version,prefix,lang) VALUES(?,?,?,?)")
	if err != nil {
		log.Fatal("Prepare insertGuildStmt error: ", err)
	}
	updateGuildStmt, err = db.Prepare("UPDATE " + config.CurrentConfig.Db.Tableprefix + "guilds " + "SET db_version = ?, prefix = ?, lang = ? " + "WHERE id = ?")
	if err != nil {
		log.Fatal("Prepare updateGuildStmt error", err)
	}
	addLogStmt = map[string]*sql.Stmt{}
	cleanOldLogStmt = map[string]*sql.Stmt{}
}

func Close() {
	err := db.Close()
	if err != nil {
		log.Fatal("DB Close error", err)
	}
}

var guildCache = map[string]*config.Guild{}

func LoadGuild(id *string) *config.Guild {
	val, exists := guildCache[*id]
	if exists {
		return val
	}
	rows, err := getDBverStmt.Query(id)
	if err != nil {
		log.Fatal("GetDBver query error", err)
	}
	URLsTable := config.CurrentConfig.Db.Tableprefix + *id + "_URLs"
	var dbVersion int
	if rows.Next() {
		err := rows.Scan(&dbVersion)
		rows.Close()
		if err != nil {
			log.Fatal("GetDBver Scan error: ", err)
		}
	} else {
		rows.Close()
		log.Print("WARN: Guild not found, making row.")
		guild := config.CurrentConfig.Guild
		_, err = insertGuildStmt.Exec(id, db_version, guild.Prefix, guild.Lang)
		if err != nil {
			log.Fatal("LoadGuild insert error: ", err)
		}
		_, err = db.Exec("CREATE TABLE IF NOT EXISTS " + URLsTable + " (" +
			"content VARCHAR NOT NULL," +
			"timeat BIGINT NOT NULL," +
			"channelid BIGINT NOT NULL," +
			"messageid BIGINT NOT NULL)")
		if err != nil {
			log.Fatal("Create URL table error:", err)
		}
	}
	var guild config.Guild
	var guildID int64
	rows, err = loadGuildStmt.Query(id)
	if err != nil {
		log.Fatal("LoadGuild query error: ", err)
	}
	defer rows.Close()
	if !rows.Next() {
		log.Fatal("LoadGuild next returned false")
	}
	err = rows.Scan(&guildID, &dbVersion, &guild.Prefix, &guild.Lang)
	if err != nil {
		log.Fatal("LoadGuild scan error: ", err)
	}
	defer rows.Close()

	addLogStmt[*id], err = db.Prepare("INSERT INTO " + URLsTable + "(content,timeat,channelid,messageid) VALUES(?,?,?,?)")
	if err != nil {
		log.Fatal("Prepare addLogStmt error: ", err)
	}

	cleanOldLogStmt[*id], err = db.Prepare("DELETE FROM " + URLsTable + " " + "WHERE timeat < ?")
	if err != nil {
		log.Fatal("Prepare cleanOldLogStmt error: ", err)
	}

	guildCache[*id] = &guild
	guild_loaded = true
	return &guild
}

func SaveGuild(id *string, guild *config.Guild) error {
	_, err := updateGuildStmt.Exec(db_version, guild.Prefix, guild.Lang, *id)
	if err != nil {
		log.Print("WARN: SaveGuild error: ", err)
	} else {
		delete(guildCache, *id)
	}
	return err
}

func AddLog(orgMsg *discordgo.MessageCreate, guildId *string, content *string, channelid, messageid *string) {
	addLogStmt[*guildId].Exec(content, time.Now().Unix(), channelid, messageid)
}

func SearchLog(orgMsg *discordgo.MessageCreate, guildId *string, content *string) (found bool, channelid string, messageid string) {
	if !(guild_loaded) {
		LoadGuild(guildId)
	}
	URLsTable := config.CurrentConfig.Db.Tableprefix + *guildId + "_URLs"
	err := db.QueryRow("SELECT channelid, messageid FROM "+URLsTable+" "+"WHERE content = ?", content).Scan(&channelid, &messageid)
	if err != nil {
		if err == sql.ErrNoRows {
			found = false
		} else {
			log.Fatal("Search Log error: ", err)
		}
	} else {
		found = true
	}
	return
}

func CleanOldLog(guildId *string) (*int64, error) {
	res, err := cleanOldLogStmt[*guildId].Exec(time.Now().Unix() - config.CurrentConfig.LogPeriod)
	if err != nil {
		return nil, err
	}
	num, err := res.RowsAffected()
	return &num, err
}

func AutoLogCleaner() {
	for {
		for guildId := range guildCache {
			_, err := CleanOldLog(&guildId)
			if err != nil {
				log.Print("Failed to auto clean log: ", err)
			}
		}
		nextTime := time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour)
		time.Sleep(time.Until(nextTime))
	}
}
