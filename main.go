package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/BurntSushi/toml"
	"github.com/KevinMcIntyre/tu-nav-server/controllers"
	"github.com/KevinMcIntyre/tu-nav-server/utils"
	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
)

var db = setupDatabase()
var config = ReadConfig("./app-config.toml")

func main() {
	utils.WritePid()

	router := httprouter.New()
	router.POST("/schedule", controllers.ScheduleHandler)

	n := negroni.New(
		negroni.NewRecovery(),
	)

	n.UseHandler(router)
	n.Run(":" + config.ServerPort)
}

func setupDatabase() *sql.DB {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s  sslmode=%s",
		config.DBUser, config.DBPassword, config.DBName, config.DBSSLMode),
	)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return db
}

type Config struct {
	ServerPort string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func ReadConfig(configfile string) Config {
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file is missing: ", configfile)
	}

	var config Config
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}
	return config
}
