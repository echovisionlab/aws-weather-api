package query

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

type (
	Record struct {
		Minute    int `validate:"gte=0" form:"minute" binding:""`
		StationID int `validate:"required,gt=0" form:"station"`
		Page      int `form:"page"`
		PageSize  int `form:"page_size"`
	}
	Station struct {
		Name string `form:"name"`
		Addr string `form:"addr"`
		ID   int    `validate:"gte=0" form:"id"`
	}
)

func (r *Record) Scope() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Where("station_id = ?", r.StationID)

		if r.Minute == 0 {
			return db.Order("time DESC").Limit(1)
		}

		target := time.Now().Truncate(time.Minute).Add(time.Duration(-r.Minute) * time.Minute).UTC()

		page, size := r.Page, r.PageSize
		if page <= 0 {
			page = 1
		}
		if size <= 0 {
			size = 10
		} else if size > 100 {
			size = 100
		}
		offset := (page - 1) * size

		return db.
			Where("time > ?", target).
			Where("station_id = ?", r.StationID).
			Offset(offset).
			Limit(size)
	}
}

func (r *Station) Scope() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if r.ID > 0 { // exists
			db = db.Where("id = ?", r.ID)
		}
		if len(r.Addr) > 0 {
			addr := strings.ReplaceAll(r.Addr, "%", "")
			db = db.Where("address LIKE ?", fmt.Sprintf("%%%s%%", addr))
		}
		if len(r.Name) > 0 {
			name := strings.ReplaceAll(r.Name, "%", "")
			db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
		}
		return db
	}
}
