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
	Other
	Complete
)

func statsQuery(db, base *gorm.DB, t ListType, date time.Time) *gorm.DB {
	subQuery := whereDate(
		db.Table("records").Select("max(id) as id").
			Where("deleted_at is NULL").Group("account_id"), date,
	).SubQuery()
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

	case Other:
		return base.Where("records.reason != ''")

	case Complete:
		return base
	}
	return nil
}

func statsBase(db *gorm.DB, acct Account, class, date string) (tx *gorm.DB, err error) {
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
	var err error
	base, err := statsBase(h.db, acct, class, r.FormValue("date"))
	if gorm.IsRecordNotFoundError(err) {
		h.message(w, r, "找不到班級", "找不到班級"+class)
		return
	}

	date, err := time.ParseInLocation("2006-01-02", r.FormValue("date"), time.Local)
	if err != nil {
		date = today()
	}

	page := struct {
		Class                                string
		Date                                 string
		Recorded, Unrecorded, Fevered, Other int
	}{
		Class: class,
		Date:  date.Format("2006-01-02"),
	}

	err = statsQuery(h.db, base, Recorded, date).Count(&page.Recorded).Error
	if err != nil {
		panic(err)
	}
	err = statsQuery(h.db, base, Unrecorded, date).Count(&page.Unrecorded).Error
	if err != nil {
		panic(err)
	}
	err = statsQuery(h.db, base, Fevered, date).Count(&page.Fevered).Error
	if err != nil {
		panic(err)
	}
	err = statsQuery(h.db, base, Other, date).Count(&page.Other).Error
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
	base, err := statsBase(h.db, acct, class, r.FormValue("date"))
	if gorm.IsRecordNotFoundError(err) {
		h.message(w, r, "找不到班級", "找不到班級"+class)
		return
	}

	date, err := time.ParseInLocation("2006-01-02", r.FormValue("date"), time.Local)
	if err != nil {
		date = today()
	}
	result := selectRecords(statsQuery(h.db, base, t, date))

	page := make(map[string]interface{})
	page["Records"] = result
	page["Type"] = t
	if err = h.tpls["stats.htm"].ExecuteTemplate(w, "account_list", page); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func parseListType(s string) ListType {
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 || n > 4 {
		return UnknownListType
	}
	return ListType(n)
}
