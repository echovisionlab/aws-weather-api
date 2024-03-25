package service

import (
	"github.com/echovisionlab/aws-weather-api/pkg/model"
	"github.com/echovisionlab/aws-weather-api/pkg/query"
	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *Service {
	return &Service{db}
}

func (s *Service) FindRecordBy(q *query.Record) []*model.Record {
	var records []*model.Record
	s.DB.Scopes(q.Scope()).Find(&records)
	return records
}

func (s *Service) FindStationBy(q *query.Station) []*model.Station {
	var stations []*model.Station
	s.DB.Scopes(q.Scope()).Find(&stations)
	return stations
}
