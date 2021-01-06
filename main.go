//    Copyright 2021 Stefen Sharkey
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
package main

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
)

var (
	db     *sql.DB
	logger = NewLogger()
)

func main() {
	logger.SetLevel(logrus.TraceLevel)

	err := InitializeDatabase()

	if err != nil {
		return
	}

	discord, err := InitializeDiscord()

	if err != nil {
		return
	}

	// Run bot until it is interrupted.
	logger.Started()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Stop the bot once it is interrupted.
	logger.Stopping()
	db.Close()
	discord.Close()
}

func InitializeDatabase() error {
	// Obtain SQL credentials
	// Set the config file
	viper.SetConfigFile("sql.yml")

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		logger.SQLError(err)
	}

	// Create a configurations struct
	var config Configurations
	err := viper.Unmarshal(&config)

	if err != nil {
		logger.ConfigError(err)
		return err
	}

	// Get the Data Source Name
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s",
		config.Database.DBUser,
		config.Database.DBPassword,
		config.Server.Protocol,
		config.Server.IP,
		config.Server.Port,
		config.Database.DBName)

	logger.DSNDebug(dsn)

	// Initialize the database connection
	logger.SQLOpeningDebug(dsn)
	db, err = sql.Open(config.Database.DBDriver, dsn)

	if err != nil {
		logger.SQLError(err)
		return err
	}

	logger.SQLOpenedDebug(dsn)

	// TODO: Log creating table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS channel_assignments (
    guild_id                BIGINT UNSIGNED,
    open_ticket_category    BIGINT UNSIGNED,
    closed_ticket_category  BIGINT UNSIGNED,
    PRIMARY KEY (guild_id)
)`)
	// TODO: Log created table if not exists

	return nil
}

func InitializeDiscord() (*discordgo.Session, error) {
	token, _ := ioutil.ReadFile("token")
	discord, err := discordgo.New("Bot " + string(token))

	if err != nil {
		logger.DiscordSessionError(err)
		return nil, err
	}

	discord.AddHandler(GuildCreate)
	discord.AddHandler(GuildDelete)

	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuilds)

	err = discord.Open()

	if err != nil {
		logger.DiscordConnectionError(err)
		return nil, err
	}

	return discord, nil
}

func GuildCreate(session *discordgo.Session, guildCreate *discordgo.GuildCreate) {
	logger.JoinedGuild(guildCreate.Name)

	// Perform the following operations upon joining a server.
	// (1) Check if the category exists in the DB for this server.
	// (1)(a) If exists, go to step (3).
	// (1)(b) If not exists, go to step (2).
	// (2) Ask the user to specify a category.
	// (2)(a) If category does not exist, make it.
	// (2)(b) If category exists, update DB.
	// (3) TBD
	//_, err := db.Exec("SELECT * FROM channel_assignments;")
	row := db.QueryRow("SELECT * FROM channel_assignments WHERE guild_id = ?;", guildCreate.ID)

	var (
		guild_id          uint64
		opened_channel_id uint64
		closed_channel_id uint64
	)

	// If there is an error, check if it is an SQL error or if there were just no rows.
	if err := row.Scan(&guild_id, &opened_channel_id, &closed_channel_id); err != nil {
		if err == sql.ErrNoRows {
			// TODO: Didn't find rows.
		} else {
			logger.SQLError(err)
		}
	}
}

func GuildDelete(session *discordgo.Session, guildDelete *discordgo.GuildDelete) {
	logger.JoinedGuild(guildDelete.Name)
}
