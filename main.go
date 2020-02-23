package main

import (
	"encoding/base64"
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

func initDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", "/tmp/gorm.sqlite")
	if err != nil {
		panic(err)
	}
	migrate(db)
	setupDB(db)
	return db
}

func migrate(db *gorm.DB) {
	if err := db.AutoMigrate(&Record{}, &Account{}).Error; err != nil {
		panic(err)
	}
}

func setupDB(db *gorm.DB) {
	var admin Account
	err := db.First(&admin, 1).Error
	if gorm.IsRecordNotFoundError(err) {
		admin = Account{
			Name: "admin",
			Role: Admin,
		}
		password := os.Getenv("ADMIN_PASSWORD")
		admin.Password, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		err = db.Create(&admin).Error
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}
}

func main() {
	h := NewHandler(initDB())
	defer h.db.Close()

	srv := &http.Server{
		Handler: h.Router(),
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Server is listen on http://%s", srv.Addr)

	log.Fatal(srv.ListenAndServe())
}
