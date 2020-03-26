package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
)

func MockData(db *gorm.DB) {
	var accounts []Account
	var err error
	err = db.Find(&accounts).Error
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().Unix())
	for _, account := range accounts {
		record := Record{
			Account:    account,
			RecordedBy: account,
		}
		ok := record.random()
		if ok {
			db.Create(&record)
		} else {
			db.Where("account_id = ?", account.ID).Where("created_at > ?", today()).Delete(&Record{})
		}
	}
}

func (r *Record) random() bool {
	switch n := rand.Intn(1000); {
	case n < 10:
		r.Temperature = randomTemperature(37.5)

	case n < 60:
		return false

	default:
		r.Temperature = randomTemperature(35)
	}
	if n := rand.Intn(100); n < 2 {
		r.Reason = [7]string{"病假", "事假", "公假", "喪假", "自主健康管理", "居家檢疫", "居家隔離"}[rand.Intn(7)]
	}
	r.Type = TempType(rand.Intn(2) + 1)
	return true
}

func randomTemperature(offset float64) float64 {
	temp := rand.Float64() * 2.5
	return math.Floor(temp*10)/10 + offset
}
