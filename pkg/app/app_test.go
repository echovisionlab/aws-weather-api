package app

import (
	"context"
	"github.com/echovisionlab/aws-weather-api/internal/testutil"
	"github.com/echovisionlab/aws-weather-api/pkg/model"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"io"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

func TestApp_Run(t *testing.T) {
	ctx := context.Background()
	container := testutil.SetupPostgres(ctx, t)
	defer testutil.ShutdownContainer(ctx, t, container)

	app, err := New()
	assert.NoError(t, err)

	db := app.Service.DB

	deleteAll := func() {
		db.Exec("DELETE FROM realtime_weather_record WHERE TRUE")
		db.Exec("DELETE FROM realtime_weather_station WHERE TRUE")
	}

	stop := app.Run()
	defer stop()

	t.Run("must return records by station", func(t *testing.T) {
		t.Cleanup(deleteAll)

		stations := getStations(100)
		records := getRecords(stations...)

		db.CreateInBatches(stations, 50)
		db.CreateInBatches(records, 50)

		var rec model.Record
		db.Where("station_id = 33").First(&rec)

		resp, err := http.Get("http://localhost:8080/record?station=33")
		assert.NoError(t, err)

		p := struct {
			Data []model.Record `json:"data"`
		}{}

		payload, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.NoError(t, resp.Body.Close())

		assert.NotNil(t, payload)

		assert.NoError(t, json.Unmarshal(payload, &p))
		assert.NotEmpty(t, p.Data)

		assert.Equal(t, rec, p.Data[0])
	})

	t.Run("must return stations by name or addr", func(t *testing.T) {
		t.Cleanup(deleteAll)

		stations := getStations(10)
		db.Create(stations)

		s := stations[rand.Intn(10)]
		addr := s.Address[:len(s.Address)/2]
		name := s.Name[:len(s.Name)/2]

		p := struct {
			Data []*model.Station `json:"data"`
		}{}

		resp, err := http.Get("http://localhost:8080/station?name=" + name)
		assert.NoError(t, err)
		payload, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.NoError(t, resp.Body.Close())
		assert.NoError(t, json.Unmarshal(payload, &p))
		assert.Equal(t, s, p.Data[0])

		resp, err = http.Get("http://localhost:8080/station?addr=" + addr)
		assert.NoError(t, err)
		payload, err = io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.NoError(t, resp.Body.Close())
		assert.NoError(t, json.Unmarshal(payload, &p))
		assert.Equal(t, s, p.Data[0])
	})

	t.Run("must support paged request", func(t *testing.T) {
		t.Cleanup(deleteAll)

		station := getStation(1)
		records := make([]*model.Record, 10)

		for i := 0; i < 10; i++ {
			records[i] = getRecord(station.Id)
			records[i].Time = records[i].Time.Add(time.Duration(-i) * time.Minute).UTC()
		}

		db.Create(station)
		db.Create(records)

		resp, err := http.Get("http://localhost:8080/record?station=1&minute=5&page_size=2")
		assert.NoError(t, err)
		bytes, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.NoError(t, resp.Body.Close())

		p := struct {
			Data []*model.Record `json:"data"`
		}{}

		assert.NoError(t, json.Unmarshal(bytes, &p))
		assert.Len(t, p.Data, 2)
	})
}

func getRecords(stations ...*model.Station) []*model.Record {
	stationCnt := len(stations)
	records := make([]*model.Record, stationCnt)
	for i := 0; i < stationCnt; i++ {
		records[i] = getRecord(stations[i].Id)
		records[i].Time = records[i].Time.In(time.UTC).Truncate(time.Second)
	}
	return records
}

// getStations returns stations with sequential ids starting from 1, not 0.
func getStations(size int) []*model.Station {
	stations := make([]*model.Station, size)
	for i := 0; i < size; i++ {
		stations[i] = getStation(i + 1)
	}
	return stations
}

func getStation(id int) *model.Station {
	return &model.Station{
		Id:            id,
		Name:          testutil.RandStringBytes(10),
		Altitude:      rand.Intn(100),
		HasRainSensor: rand.Int()%2 > 0,
		Address:       testutil.RandStringBytes(10),
	}
}

func getRecord(stationID int) *model.Record {
	return &model.Record{
		Id:                      uuid.New(),
		StationID:               stationID,
		RainAcc:                 rand.Float32(),
		RainFifteen:             rand.Float32(),
		RainHour:                rand.Float32(),
		RainThreeHour:           rand.Float32(),
		RainSixHour:             rand.Float32(),
		RainTwelveHour:          rand.Float32(),
		Temperature:             rand.Float32(),
		WindAverageMinute:       rand.Float32(),
		WindAverageMinuteDeg:    rand.Float32(),
		WindAverageTenMinute:    rand.Float32(),
		WindAverageTenMinuteDeg: rand.Float32(),
		Moisture:                rand.Intn(10),
		SeaLevelAirPressure:     rand.Float32(),
		Time:                    time.Now(),
	}
}
