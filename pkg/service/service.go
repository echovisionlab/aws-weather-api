package service

import (
	"fmt"
	"github.com/echovisionlab/aws-weather-api/pkg/model"
	"github.com/echovisionlab/aws-weather-api/pkg/query"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Config struct {
	MaxPageSize int `validate:"gt=0,lte=100"`
}

type Service struct {
	DB          *gorm.DB
	MaxPageSize int
}

func New(db *gorm.DB, config *Config) *Service {
	return &Service{db, config.MaxPageSize}
}

func (s *Service) FindStationBy(q *query.Station) []*model.Station {
	var stations []*model.Station
	s.DB.Scopes(getQueryStationScopes(q)).Find(&stations)
	return stations
}

func (s *Service) FindRecordBy(q *query.Record) []*model.Record {
	var records []*model.Record
	s.DB.Scopes(getQueryRecordScopes(q, s.MaxPageSize)).Find(&records)
	return records
}

func getQueryStationScopes(q *query.Station) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if q.ID > 0 { // exists
			db = db.Where("id = ?", q.ID)
		}
		if len(q.Addr) > 0 {
			addr := strings.ReplaceAll(q.Addr, "%", "")
			db = db.Where("address LIKE ?", fmt.Sprintf("%%%s%%", addr))
		}
		if len(q.Name) > 0 {
			name := strings.ReplaceAll(q.Name, "%", "")
			db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
		}
		return db
	}
}

func getQueryRecordScopes(q *query.Record, maxPageSize int) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Where("station_id = ?", q.StationID)

		if q.Minute == 0 {
			return db.Order("time DESC").Limit(1)
		}

		target := time.Now().Truncate(time.Minute).Add(time.Duration(-q.Minute) * time.Minute).UTC()
		offset, size := getOffsetLimit(q, maxPageSize)

		return db.
			Where("time > ?", target).
			Where("station_id = ?", q.StationID).
			Offset(offset).
			Limit(size)
	}
}

func getOffsetLimit(q *query.Record, maxPageSize int) (int, int) {
	page, size := q.Page, q.PageSize
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	} else if size > 100 {
		size = maxPageSize
	}
	offset := (page - 1) * size
	return offset, size
}
