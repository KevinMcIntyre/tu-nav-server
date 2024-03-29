package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	//"tu-nav-server/controllers"

	"github.com/BurntSushi/toml"
	"github.com/KevinMcIntyre/tu-nav-server/controllers"
	"github.com/KevinMcIntyre/tu-nav-server/models"
	"github.com/KevinMcIntyre/tu-nav-server/utils"
	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
)

var db = setupDatabase()
var config = ReadConfig("./app-config.toml")

func init() {
	logFile, err := os.OpenFile("tu-nav-server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf(fmt.Sprintf("Error opening log file: %s \n", err))
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	main()
}

func main() {
	controllers.DB = db

	utils.WritePid()

	err := models.SeedBuildingsFile(db, config.BuildingsSeedFilePath)
	if err != nil {
		fmt.Println("Error seeding buildings", err)
	}

	router := httprouter.New()
	router.GET("/", controllers.HomeHandler)
	router.GET("/buildings", controllers.BuildingHandler)
	router.POST("/schedule", controllers.ScheduleHandler)
	router.GET("/versioncontrol", controllers.VerisonControlHandler)
	router.GET("/update", controllers.UpdateHandler)
	router.ServeFiles("/public/*filepath", http.Dir("/Users/Eric/go/src/tu-nav-server/public"))

	n := negroni.New(
		negroni.NewRecovery(),
	)

	// n.Use(negroni.NewStatic(http.Dir("/public")))
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
	ServerPort            string
	DBUser                string
	DBPassword            string
	DBName                string
	DBSSLMode             string
	BuildingsSeedFilePath string
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
