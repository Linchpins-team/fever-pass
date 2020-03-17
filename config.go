package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/BurntSushi/toml"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	Debug   = "debug"
	Release = "release"
)

type Config struct {
	Mode   string
	Server struct {
		Base string
		Port int
	}

	Database struct {
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

func setupConfig() {
	var ok string
	fmt.Print("Do you want to create a new setting? (y/n) ")
	fmt.Scanln(&ok)
	var c Config
	if parseYN(ok) {
		c = generateConfig()
	} else {
		c = loadConfig()
	}

	fmt.Print("admin password: ")
	fmt.Scanln(&c.Password)

	db, err := initDB(c)
	if err != nil {
		panic(err)
	}
	setupDB(c, db)
	if parseYN(ok) {
		writeConfig(c, ConfPath)
	}
}

func generateConfig() (c Config) {
	fmt.Print("server base url (default http://localhost:8080): ")
	fmt.Scanln(&c.Server.Base)
	if c.Server.Base == "" {
		c.Server.Base = "http://localhost:8080"
	}

	fmt.Print("server listen port (default: 8080): ")
	_, err := fmt.Scanln(&c.Server.Port)
	if err != nil {
		c.Server.Port = 8080
	}

	fmt.Print("database mode debug/release (default: debug): ")
	fmt.Scanln(&c.Mode)
	if c.Mode != Release && c.Mode != Debug {
		c.Mode = Debug
	}

	if c.Mode == Release {
		fmt.Println("MySQL database setup. Please create database before configuration")
		for {
			fmt.Print("Database name: ")
			fmt.Scanln(&c.Database.Name)
			fmt.Print("Database user: ")
			fmt.Scanln(&c.Database.User)
			fmt.Print("Database password: ")
			fmt.Scanln(&c.Database.Password)
			if db, err := initDB(c); err != nil {
				fmt.Printf("Cannot connect to database: %s, initial it now? (y/n) ", err)
				var ans string
				fmt.Scanln(&ans)
				if parseYN(ans) {
					createMySQLDatabase(c)
					break
				} else {
					fmt.Println("Please enter again")
				}
			} else {
				db.Close()
				break
			}
		}
	}

	return
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
	fmt.Println("Configurations has been generated at config.toml")
}

func loadConfig() (c Config) {
	if _, err := toml.DecodeFile(ConfPath, &c); err != nil {
		log.Fatalln("No configuration file")
	}
	return
}

func createMySQLDatabase(c Config) {
	cmd := exec.Command("/bin/sh", "-c", "sudo mysql")
	cmd.Stdin = strings.NewReader(fmt.Sprintf(`
	CREATE DATABASE IF NOT EXISTS %s ;
	DROP USER IF EXISTS '%s'@'localhost';
	FLUSH PRIVILEGES;
	CREATE USER '%s'@'localhost' IDENTIFIED BY '%s'; 
	GRANT ALL PRIVILEGES ON %s . * TO '%s'@'localhost'; 
	FLUSH PRIVILEGES;
	`, c.Database.Name, c.Database.User, c.Database.User, c.Database.Password, c.Database.Name, c.Database.User))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
