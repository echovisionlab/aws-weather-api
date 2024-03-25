package model

import (
	"github.com/google/uuid"
	"time"
)

type (
	Station struct {
		Id            int    `gorm:"column:id" json:"id"`
		Name          string `gorm:"column:name" json:"name"`
		Altitude      int    `gorm:"column:altitude" json:"altitude"`
		HasRainSensor bool   `gorm:"column:has_rain_sensor" json:"has_rain_sensor"`
		Address       string `gorm:"column:address" json:"address"`
	}

	Record struct {
		Id                      uuid.UUID `gorm:"column:id" json:"id"`
		StationID               int       `gorm:"column:station_id" json:"station_id"`
		RainAcc                 float32   `gorm:"column:rain_acc" json:"rain_acc"`
		RainFifteen             float32   `gorm:"column:rain_fifteen" json:"rain_fifteen"`
		RainHour                float32   `gorm:"column:rain_hour" json:"rain_hour"`
		RainThreeHour           float32   `gorm:"column:rain_three_hour" json:"rain_three_hour"`
		RainSixHour             float32   `gorm:"column:rain_six_hour" json:"rain_six_hour"`
		RainTwelveHour          float32   `gorm:"column:rain_twelve_hour" json:"rain_twelve_hour"`
		Temperature             float32   `gorm:"column:temperature" json:"temperature"`
		WindAverageMinute       float32   `gorm:"column:wind_avg_minute" json:"wind_average_minute"`
		WindAverageMinuteDeg    float32   `gorm:"column:wind_avg_minute_deg" json:"wind_average_minute_deg"`
		WindAverageTenMinute    float32   `gorm:"column:wind_avg_ten_minute" json:"wind_average_ten_minute"`
		WindAverageTenMinuteDeg float32   `gorm:"column:wind_avg_ten_minute_deg" json:"wind_average_ten_minute_deg"`
		Moisture                int       `gorm:"column:moisture" json:"moisture"`
		SeaLevelAirPressure     float32   `gorm:"column:sea_level_air_pressure" json:"sea_level_air_pressure"`
		Time                    time.Time `gorm:"column:time" json:"time"`
	}
)

func (s *Station) TableName() string {
	return "realtime_weather_station"
}

func (r *Record) TableName() string {
	return "realtime_weather_record"
}
