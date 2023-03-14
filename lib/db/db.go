package db

import (
	"database/sql"
	"log"
	"sync"
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
		log.Fatalf("DB load error: %s", err)
	}

	if _, err := db.Exec(`
    CREATE TABLE IF NOT EXISTS ` + config.CurrentConfig.Db.Tableprefix + `guilds (
      id BIGINT NOT NULL PRIMARY KEY,
      db_version INT NOT NULL,
      prefix VARCHAR,
      lang VARCHAR
    );
  `); err != nil {
		log.Fatalf("Create guild table error: %s", err)
	}

	if _, err := db.Exec(`
		ALTER TABLE ` + config.CurrentConfig.Db.Tableprefix + `guilds DROP COLUMN bots;
	`); err == nil {
		log.Println("WARN: Database update successfully")
	} else {
		log.Println("FINE: Guilds DB up-to-date!")
	}

	if _, err := db.Exec("VACUUM"); err != nil {
		log.Fatalf("DB VACUUM error: %s", err)
	}

	getDBverStmt, err = db.Prepare("SELECT db_version FROM " + config.CurrentConfig.Db.Tableprefix + "guilds WHERE " + "id = ?")
	if err != nil {
		log.Fatalf("Prepare getDBverStmt error: %s", err)
	}

	loadGuildStmt, err = db.Prepare("SELECT * FROM " + config.CurrentConfig.Db.Tableprefix + "guilds WHERE " + "id = ?")
	if err != nil {
		log.Fatalf("Prepare loadGuildStmt error: %s", err)
	}

	insertGuildStmt, err = db.Prepare("INSERT INTO " + config.CurrentConfig.Db.Tableprefix + "guilds(id,db_version,prefix,lang) VALUES(?,?,?,?)")
	if err != nil {
		log.Fatalf("Prepare insertGuildStmt error: %s", err)
	}

	updateGuildStmt, err = db.Prepare("UPDATE " + config.CurrentConfig.Db.Tableprefix + "guilds " + "SET db_version = ?, prefix = ?, lang = ? " + "WHERE id = ?")
	if err != nil {
		log.Fatalf("Prepare updateGuildStmt error: %s", err)
	}

	addLogStmt = map[string]*sql.Stmt{}
	cleanOldLogStmt = map[string]*sql.Stmt{}
}

func Close() {
	if err := db.Close(); err != nil {
		log.Fatalf("DB Close error: %v", err)
	}
}

var guildCache sync.Map

func LoadGuild(id *string) *config.Guild {
	val, exists := guildCache.Load(*id)
	if exists {
		return val.(*config.Guild)
	}
	rows, err := getDBverStmt.Query(id)
	if err != nil {
		log.Fatalf("GetDBver query error: %s", err)
	}
	defer rows.Close()

	URLsTable := config.CurrentConfig.Db.Tableprefix + *id + "_URLs"
	var dbVersion int
	foundGuild := false

	if rows.Next() {
		err := rows.Scan(&dbVersion)
		if err != nil {
			log.Fatalf("GetDBver Scan error: %s", err)
		}
		foundGuild = true
	}

	if !foundGuild {
		rows.Close()
		log.Println("WARN: Guild not found, making row.")
		guild := config.CurrentConfig.Guild
		_, err = insertGuildStmt.Exec(id, db_version, guild.Prefix, guild.Lang)
		if err != nil {
			log.Fatalf("LoadGuild insert error: %s", err)
		}

		_, err = db.Exec("CREATE TABLE IF NOT EXISTS " + URLsTable + " (" +
			"content VARCHAR NOT NULL," +
			"timeat BIGINT NOT NULL," +
			"channelid BIGINT NOT NULL," +
			"messageid BIGINT NOT NULL)")
		if err != nil {
			log.Fatalf("Create URL table error: %s", err)
		}
	}

	var guild config.Guild
	var guildID int64

	rows, err = loadGuildStmt.Query(id)
	if err != nil {
		log.Fatalf("LoadGuild query error: %s", err)
	}
	defer rows.Close()

	foundInDb := false

	for rows.Next() {
		err = rows.Scan(&guildID, &dbVersion, &guild.Prefix, &guild.Lang)
		if err != nil {
			log.Fatalf("LoadGuild scan error: %s", err)
		}
		foundInDb = true
	}

	if !foundInDb {
		log.Fatal("LoadGuild next returned false")
	}

	addLogStmt[*id], err = db.Prepare("INSERT INTO " + URLsTable + "(content,timeat,channelid,messageid) VALUES(?,?,?,?)")
	if err != nil {
		log.Fatalf("Prepare addLogStmt error: %s", err)
	}

	cleanOldLogStmt[*id], err = db.Prepare("DELETE FROM " + URLsTable + " WHERE timeat < ?")
	if err != nil {
		log.Fatalf("Prepare cleanOldLogStmt error: %s", err)
	}

	guildCache.Store(*id, &guild)
	guild_loaded = true
	return &guild
}

func SaveGuild(id *string, guild *config.Guild) error {
	_, err := updateGuildStmt.Exec(db_version, guild.Prefix, guild.Lang, *id)
	if err != nil {
		log.Println("ERROR: SaveGuild error: ", err)
		return err
	}
	return nil
}

func AddLog(orgMsg *discordgo.MessageCreate, guildId *string, content *string, channelid, messageid *string) {
	addLogStmt[*guildId].Exec(content, time.Now().Unix(), channelid, messageid)
}

func SearchLog(orgMsg *discordgo.MessageCreate, guildId *string, content *string) (found bool, channelid string, messageid string) {
	if !guild_loaded {
		LoadGuild(guildId)
	}
	URLsTable := config.CurrentConfig.Db.Tableprefix + *guildId + "_URLs"
	err := db.QueryRow("SELECT channelid, messageid FROM "+URLsTable+" "+"WHERE content = ?", content).Scan(&channelid, &messageid)
	switch {
	case err == sql.ErrNoRows:
		found = false
	case err != nil:
		log.Fatalf("Search Log error: %v", err)
	default:
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
		guildCache.Range(func(guildID, value interface{}) bool {
			_, err := CleanOldLog(guildID.(*string))
			if err != nil {
				log.Println("Failed to auto clean log: ", err)
			}
			return true
		})
		nextTime := time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour)
		time.Sleep(time.Until(nextTime))
	}
}
