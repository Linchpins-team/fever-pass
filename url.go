package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/skip2/go-qrcode"
)

type URL struct {
	ID        string `gorm:"primary_key;type:varchar(32)"`
	CreatedAt time.Time
	ExpireAt  time.Time
	Last      int
}

func (h Handler) newURL(w http.ResponseWriter, r *http.Request) {
	var url URL
	var err error
	url.Last, err = strconv.Atoi(r.FormValue("max_usage"))
	if err != nil {
		http.Error(w, err.Error(), 415)
		return
	}

	valid, err := strconv.Atoi(r.FormValue("valid_time"))
	if err != nil {
		http.Error(w, err.Error(), 415)
		return
	}

	url.CreatedAt = time.Now()
	url.ExpireAt = url.CreatedAt.AddDate(0, 0, valid)
	url.ID = generateURL()
	if err = h.db.Create(&url).Error; err != nil {
		panic(err)
	}

	data := make(map[string]string)
	data["url"] = h.inviteURL(url.ID)
	data["qrcode"] = h.createQRCode(url.ID)

	enc := json.NewEncoder(w)
	if err = enc.Encode(&data); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func generateURL() string {
	random := make([]byte, 9)
	if _, err := rand.Read(random); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(random)
}

func (h Handler) inviteURL(key string) string {
	return fmt.Sprintf("%s/register?key=%s", h.config.Server.Base, key)
}

func (h Handler) createQRCode(key string) string {
	path := fmt.Sprintf("/qrcodes/%s.png", key)
	err := qrcode.WriteFile(h.inviteURL(key), qrcode.Medium, 256, "static"+path)
	if err != nil {
		panic(err)
	}
	return h.config.Server.Base + path
}

func (h Handler) registerPage(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	var url URL
	err := h.db.Where("id = ? and expire_at > ? and last > 0", key, time.Now()).First(&url).Error
	if gorm.IsRecordNotFoundError(err) {
		// not found
		http.Error(w, "invalid key", 404)
		return
	} else if err != nil {
		panic(err)
	}

	page := struct {
		Key string
	}{key}
	h.HTML(w, "register.htm", page)
}
