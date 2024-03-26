package query

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
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
