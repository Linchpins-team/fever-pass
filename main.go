package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
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
	godotenv.Load()
	file, err := os.OpenFile(".env", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	hashKey = loadKey("HASH_KEY", file)
	blockKey = loadKey("BLOCK_KEY", file)
}

func loadKey(name string, w io.Writer) (key []byte) {
	if k := os.Getenv(name); k != "" {
		key = decodeKey(k)
	} else {
		key = securecookie.GenerateRandomKey(32)
		fmt.Fprintln(w, name+"="+encodeKey(key))
	}
	return
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
	return
}

func migrate(db *gorm.DB) {
	if err := db.AutoMigrate(&Record{}, &Account{}, &URL{}).Error; err != nil {
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
		setupConfig(confPath)
		return
	}

	c := loadConfig(confPath)
	db, err := initDB(c)
	if err != nil {
		panic(err)
	}
	h := NewHandler(db, c)
	defer h.db.Close()

	srv := &http.Server{
		Handler: h.Router(),
		Addr:    fmt.Sprintf("0.0.0.0:%d", c.Server.Port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Server is listen on %s", c.Server.Base)

	log.Fatal(srv.ListenAndServe())
}
