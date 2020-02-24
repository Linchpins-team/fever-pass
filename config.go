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

func setupConfig(path string) {
	var ok string
	fmt.Print("This will overwrite your existing setting, continue? (y/n) ")
	fmt.Scanln(&ok)
	if ok != "y" && ok != "Y" {
		return
	}
	c := generateConfig()
	db, err := initDB(c)
	if err != nil {
		panic(err)
	}
	setupDB(c, db)
	writeConfig(c, path)
}

func generateConfig() (c Config) {
	fmt.Print("server base url (default localhost:8080): ")
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
			fmt.Print("database name: ")
			fmt.Scanln(&c.Database.Name)
			fmt.Print("database user: ")
			fmt.Scanln(&c.Database.User)
			fmt.Print("datbase password: ")
			fmt.Scanln(&c.Database.Password)
			if db, err := initDB(c); err != nil {
				fmt.Printf("Cannot connect to database: %s, initial it now? (y/n) ", err)
				var ans string
				fmt.Scanln(&ans)
				if ans == "y" {
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

	fmt.Print("admin password: ")
	fmt.Scanln(&c.Password)

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
	fmt.Println("configurations has generated at config.toml")
}

func loadConfig(path string) (c Config) {
	if _, err := toml.DecodeFile(path, &c); err != nil {
		log.Fatalln("No configuration file")
	}
	return
}

func createMySQLDatabase(c Config) {
	cmd := exec.Command("/bin/sh", "-c", "sudo mysql")
	cmd.Stdin = strings.NewReader(fmt.Sprintf(`
	CREATE DATABASE IF NOT EXISTS %s ;\n
	CREATE USER '%s'@'localhost' IDENTIFIED BY '%s'; \n
	GRANT ALL PRIVILEGES ON %s . * TO '%s'@'localhost'; \n
	FLUSH PRIVILEGES;
	`, c.Database.Name, c.Database.User, c.Database.Password, c.Database.Name, c.Database.User))
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
