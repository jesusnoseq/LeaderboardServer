package main

import (
	"log"

	"github.com/jesusnoseq/LeaderboardServer/functions/pkg/entry"
	"github.com/jesusnoseq/LeaderboardServer/functions/pkg/entry/persistence"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("app.env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Invalid config file", err)
	}

	viper.SetDefault("SALT", "ANY")
	viper.SetDefault("HTTP_PORT", "8080")
	salt := viper.Get("SALT").(string)
	port := viper.Get("HTTP_PORT").(string)
	dao := persistence.GetDao("")
	startServer(dao, port, salt)
}

func startServer(dao persistence.EntryDAO, port string, salt string) {
	router := entry.GetEntryServer()
	err := router.Run(":" + port)
	if err != err {
		log.Fatal("error runing http server", err)
	}
}
