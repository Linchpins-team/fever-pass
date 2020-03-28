package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	Debug   = "debug"
	Release = "release"
)

type Config struct {
	Mode string

	Site struct {
		Name string
		Icon string
	}

	Server struct {
		Base string
		Port int
	}

	Database struct {
		Host     string
		Name     string
		User     string
		Password string
	}

	Password string `toml:"-"`
}

func parseYN(str string) bool {
	switch str {
	case "y", "Y":
		return true

	default:
		return false
	}
}

func setupAdminPassword() {
	var c Config
	c = loadConfig()

	fmt.Print("admin password: ")
	fmt.Scanln(&c.Password)

	db, err := initDB(c)
	if err != nil {
		panic(err)
	}
	setupDB(c, db)
}

// generate default configure
func genDefaultConf() {
	var c Config
	c.Server.Base = "http://localhost:8080"
	c.Server.Port = 8080

	c.Mode = Release

	c.Site.Name = "Fever Pass"
	c.Site.Icon = "/static/img/icon.png"

	c.Database.Host = "localhost"
	c.Database.Name = "fever_pass"
	c.Database.User = "fever_pass_user"
	writeConfig(c, ConfPath)
}

func writeConfig(c Config, path string) {
	os.Remove(path)
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	enc := toml.NewEncoder(file)
	err = enc.Encode(c)
	if err != nil {
		panic(err)
	}
	fmt.Println("Configurations has been generated at", path)
}

func loadConfig() (c Config) {
	if _, err := toml.DecodeFile(ConfPath, &c); err != nil {
		log.Fatalln("No configuration file, please use -g to generate configure.")
	}
	return
}

func createDatabaseCode(c Config) {
	fmt.Println("Copy the following code to sql.")
	fmt.Printf(`
CREATE DATABASE IF NOT EXISTS %[1]s ;
DROP USER IF EXISTS '%[2]s'@'localhost';
FLUSH PRIVILEGES;
CREATE USER '%[2]s'@'localhost' IDENTIFIED BY '%[3]s'; 
GRANT ALL PRIVILEGES ON %[1]s . * TO '%[2]s'@'localhost'; 
FLUSH PRIVILEGES;
`, c.Database.Name, c.Database.User, c.Database.Password)
}
