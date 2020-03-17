package main

import (
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
)

type ListType int

const (
	UnknownListType ListType = iota
	Recorded
	Unrecorded
	Fevered
)

func statsQuery(db, base *gorm.DB, t ListType) (tx *gorm.DB) {
	subQuery := db.Table("records").Select("max(id) as id, account_id").Where("created_at > ?", today()).Where("deleted_at is NULL").Group("account_id").SubQuery()

	switch t {
	case Recorded:
		return base.Table("accounts").Joins("inner join ? as r on accounts.id = r.account_id", subQuery)

	case Unrecorded:
		return base.Table("accounts").Joins("left join ? as r on accounts.id = r.account_id", subQuery).Where("r.id is NULL")

	case Fevered:
		return base.Table("accounts").Joins(
			"inner join records as r on r.account_id = accounts.id",
		).Joins("inner join ? as m on m.id = r.id", subQuery).Where(
			"(temperature > 38 and type = 1) or (temperature > 37.5 and type = 2)",
		)

	}
	return nil
}

func statsBase(db *gorm.DB, acct Account, class string) (tx *gorm.DB, err error) {
	tx = joinClasses(db).Set("gorm:auto_preload", true)
	if acct.Role == Teacher {
		class = acct.Class.Name
	}
	if class != "" {
		err = db.First(&Class{}, "name = ?", class).Error
		if err != nil {
			return
		}
		tx = tx.Where("classes.name = ?", class)
	}
	tx = tx.Where("role = ? or role = ?", Student, Teacher)
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
	err = statsQuery(h.db, base, Recorded).Count(&page.Recorded).Error
	if err != nil {
		panic(err)
	}
	err = statsQuery(h.db, base, Unrecorded).Count(&page.Unrecorded).Error
	if err != nil {
		panic(err)
	}
	err = statsQuery(h.db, base, Fevered).Count(&page.Fevered).Error
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

	var accounts []Account
	var err error
	base, err := statsBase(h.db, acct, class)
	if gorm.IsRecordNotFoundError(err) {
		h.errorPage(w, r, "找不到班級", "找不到班級"+class)
		return
	}
	err = statsQuery(h.db, base, t).Find(&accounts).Error
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err = h.tpls["stats.htm"].ExecuteTemplate(w, "account_list", accounts); err != nil {
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
