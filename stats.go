package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

type ListType int

const (
	UnknownListType ListType = iota
	Recorded
	Unrecorded
	Fevered
)

func statsQuery(db, base *gorm.DB, t ListType, date time.Time) (tx *gorm.DB) {
	subQuery := db.Table("records").Select(
		"max(id) as id",
	).Where(
		"created_at > ? and created_at < ?", date, date.AddDate(0, 0, 1),
	).Where("deleted_at is NULL").Group("account_id").SubQuery()
	subQuery = db.Table("records").Joins("inner join ? as m on m.id = records.id", subQuery).SubQuery()
	base = base.Table("accounts").Joins(
		"left join ? as records on records.account_id = accounts.id", subQuery)

	switch t {
	case Recorded:
		return base.Where("records.id is not NULL")

	case Unrecorded:
		return base.Where("records.id is NULL")

	case Fevered:
		return base.Where(
			"(temperature >= 38 and type = 1) or (temperature >= 37.5 and type = 2)",
		)
	}
	return nil
}

func statsBase(db *gorm.DB, acct Account, class string) (tx *gorm.DB, err error) {
	tx = joinClasses(db).Set("gorm:auto_preload", true)
	if acct.Authority.Record == Group {
		class = acct.Class.Name
	}
	if class != "" {
		err = db.First(&Class{}, "name = ?", class).Error
		if err != nil {
			return
		}
		tx = tx.Where("classes.name = ?", class)
	}
	tx = tx.Where("role = ?", RoleStudent)
	tx = tx.Order("classes.name, number asc")
	return
}

// 今日總覽頁面
func (h Handler) stats(w http.ResponseWriter, r *http.Request) {
	acct := r.Context().Value(KeyAccount).(Account)
	class := r.FormValue("class")
	page := struct {
		Class                         string
		Recorded, Unrecorded, Fevered int
	}{
		Class: class,
	}
	var err error
	base, err := statsBase(h.db, acct, class)
	if gorm.IsRecordNotFoundError(err) {
		h.errorPage(w, r, "找不到班級", "找不到班級"+class)
		return
	}
	err = statsQuery(h.db, base, Recorded, today()).Count(&page.Recorded).Error
	if err != nil {
		panic(err)
	}
	err = statsQuery(h.db, base, Unrecorded, today()).Count(&page.Unrecorded).Error
	if err != nil {
		panic(err)
	}
	err = statsQuery(h.db, base, Fevered, today()).Count(&page.Fevered).Error
	if err != nil {
		panic(err)
	}
	h.HTML(w, r, "stats.htm", page)
}

func (h Handler) statsList(w http.ResponseWriter, r *http.Request) {
	acct := r.Context().Value(KeyAccount).(Account)
	class := r.FormValue("class")
	t := parseListType(r.FormValue("type"))
	if t == UnknownListType {
		http.Error(w, "Unknown list type", 415)
		return
	}

	var err error
	base, err := statsBase(h.db, acct, class)
	if gorm.IsRecordNotFoundError(err) {
		h.errorPage(w, r, "找不到班級", "找不到班級"+class)
		return
	}

	result := selectRecords(statsQuery(h.db, base, t, today()))
	if err = h.tpls["stats.htm"].ExecuteTemplate(w, "account_list", result); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func parseListType(s string) ListType {
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 || n > 3 {
		return UnknownListType
	}
	return ListType(n)
}
