package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var (
	hashKey, blockKey []byte
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("Cannot load .env")
	}
	loadKeys()
}

func loadKeys() {
	if k := os.Getenv("HASH_KEY"); k != "" {
		hashKey = decodeKey(k)
	} else {
		hashKey = securecookie.GenerateRandomKey(32)
		fmt.Println("HASH_KEY=" + encodeKey(hashKey))
	}
	if k := os.Getenv("BLOCK_KEY"); k != "" {
		blockKey = decodeKey(k)
	} else {
		blockKey = securecookie.GenerateRandomKey(32)
		fmt.Println("BLOCK_KEY=" + encodeKey(blockKey))
	}
}

func encodeKey(key []byte) string {
	return base64.StdEncoding.EncodeToString(key)
}

func decodeKey(key string) []byte {
	if dst, err := base64.StdEncoding.DecodeString(key); err == nil {
		return dst
	} else {
		panic(err)
	}
}

func initDB(c Config) (db *gorm.DB, err error) {
	switch c.Mode {
	case Debug:
		db, err = gorm.Open("sqlite3", "/tmp/gorm.sqlite")

	case Release:
		connection := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", c.Database.User, c.Database.Password, c.Database.Name)
		db, err = gorm.Open("mysql", connection)
	}
	if err != nil {
		return
	}
	migrate(db)
	setupDB(c, db)
	return
}

func migrate(db *gorm.DB) {
	if err := db.AutoMigrate(&Record{}, &Account{}).Error; err != nil {
		panic(err)
	}
}

func setupDB(c Config, db *gorm.DB) {
	var admin Account
	err := db.First(&admin, 1).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		panic(err)
	}
	admin.Name = "admin"
	admin.Role = Admin
	password := c.Password
	admin.Password, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	err = db.Save(&admin).Error
	if err != nil {
		panic(err)
	}
}

func main() {
	var init bool
	var confPath string
	flag.BoolVar(&init, "init", false, "init configuration")
	flag.StringVar(&confPath, "conf", "config.toml", "configuration file path")
	flag.Parse()

	if init {
		setupConfig()
		return
	}

	c := loadConfig(confPath)
	db, err := initDB(c)
	if err != nil {
		panic(err)
	}
	h := NewHandler(db)
	defer h.db.Close()

	srv := &http.Server{
		Handler: h.Router(),
		Addr:    fmt.Sprintf("127.0.0.1:%d", c.Server.Port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Server is listen on http://%s", c.Server.Base)

	log.Fatal(srv.ListenAndServe())
}
