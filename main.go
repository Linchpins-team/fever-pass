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
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
)

var (
	hashKey, blockKey []byte
)

var (
	GenerateConfigure bool
	ConfPath          string
	SetAdminPassword  bool
	Mock              bool
	InitDB            bool
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

	flag.BoolVar(&GenerateConfigure, "g", false, "generate new configure")
	flag.BoolVar(&SetAdminPassword, "s", false, "set admin password")
	flag.BoolVar(&InitDB, "init", false, "use -init to generate the sql code to initialize database.")
	flag.StringVar(&ConfPath, "conf", "config.toml", "configuration file path")
	flag.BoolVar(&Mock, "mock", false, "use -mock to insert mock data for every accounts")
	flag.Parse()
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

func main() {
	if GenerateConfigure {
		genDefaultConf()
		return
	}
	c := loadConfig()
	if InitDB {
		createDatabaseCode(c)
		return
	}
	if SetAdminPassword {
		setupAdminPassword()
		return
	}

	db, err := initDB(c)
	if err != nil {
		fmt.Println(err)
		createDatabaseCode(c)
		return
	}

	if Mock {
		MockData(db)
		return
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

	log.Printf("Server is listening on %s", c.Server.Base)

	log.Fatal(srv.ListenAndServe())
}
