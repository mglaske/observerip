package observerip

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type StationResponse struct {
	Id                string
	Password          string
	TempF             int
	Humidity          int
	DewPointF         int
	WindChillF        int
	WindDirection     int
	WindSpeedMph      int
	WindGustMph       int
	RainInch          float64
	DailyRainInch     float64
	WeeklyRainInch    float64
	MonthlyRainInch   float64
	YearlyRainInch    float64
	SolarRadiation    int
	UV                int
	IndoorTempF       int
	IndoorHumidity    int
	BarometricInch    int
	LowBattery        int
	DateUTC           string
	SoftwareType      string
	Action            string
	RealTime          int
	RealTimeFrequency int
}

type EndpointResponse struct {
	Passkey     string
	StationType string
	DateUTC     string
	Endpoints   []Endpoint
}

type Endpoint struct {
	Id       int
	Name     string // For reference
	TempF    float64
	Humidity int
	Battery  int
}

func getValue(in *url.Values, key string) ([]string, error) {
	v, ok := (*in)[key]
	if ok {
		return v, nil
	}
	return v, errors.New("Value Not Found")
}

func getSingleValue(in *url.Values, key string, def string) string {
	vs, err := getValue(in, key)
	if err != nil {
		return def
	}
	return vs[0]
}

func getSingleIntValue(in *url.Values, key string, def int) int {
	sv := getSingleValue(in, key, "")
	if sv == "" {
		return def
	}
	i, e := strconv.Atoi(sv)
	if e != nil {
		panic(e)
	}
	return i
}

func getSingleFloatValue(in *url.Values, key string, def float64) float64 {
	sv := getSingleValue(in, key, "")
	if sv == "" {
		return def
	}
	f, e := strconv.ParseFloat(sv, 64)
	if e != nil {
		return def
	}
	return f
}

func (me *EndpointResponse) Parse(w http.ResponseWriter, req *http.Request) {
	// /endpoint?&PASSKEY=edd72cda33fc21b71a7c16603530f7c8&stationtype=WS-1501-IP&dateutc=2022-07-19+00:22:31&temp1f=74.30&humidity1=55&batt1=1&temp2f=74.66&humidity2=55&batt2=1&batt3=1&batt4=1&batt5=1&batt6=1&batt7=1&batt8=1
	rv := req.URL.Query()

	me.Passkey = getSingleValue(&rv, "PASSKEY", "")
	me.StationType = getSingleValue(&rv, "stationtype", "")
	me.DateUTC = getSingleValue(&rv, "dateutc", "")
	me.Endpoints = make([]Endpoint, 8)
	for x := 1; x <= 8; x++ {
		ep := Endpoint{
			Id:       x,
			TempF:    getSingleFloatValue(&rv, fmt.Sprintf("temp%df", x), -1),
			Humidity: getSingleIntValue(&rv, fmt.Sprintf("humidity%d", x), -1),
			Battery:  getSingleIntValue(&rv, fmt.Sprintf("batt%d", x), -1),
		}
		me.Endpoints[x-1] = ep
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "success")
}

func (me *StationResponse) Parse(w http.ResponseWriter, req *http.Request) {
	// GET /weatherstation/updateweatherstation.php?ID=4410&PASSWORD=&tempf=-9999&humidity=-9999&dewptf=-9999&windchillf=-9999&winddir=-9999&windspeedmph=-9999&windgustmph=-9999&rainin=0.00&dailyrainin=0.00&weeklyrainin=0.00&monthlyrainin=0.00&yearlyrainin=0.00&solarradiation=-9999&UV=-9999&indoortempf=-9999&indoorhumidity=-9999&baromin=-9999&lowbatt=0&dateutc=now&softwaretype=WH2602%20V4.6.2&action=updateraw&realtime=1&rtfreq=5
	rv := req.URL.Query()

	me.Id = getSingleValue(&rv, "ID", "")
	me.TempF = getSingleIntValue(&rv, "tempf", -1)
	me.Humidity = getSingleIntValue(&rv, "humidity", -1)
	me.DewPointF = getSingleIntValue(&rv, "dewptf", -1)
	me.WindChillF = getSingleIntValue(&rv, "windchillf", -1)
	me.WindDirection = getSingleIntValue(&rv, "winddir", -1)
	me.WindSpeedMph = getSingleIntValue(&rv, "windspeedmph", -1)
	me.WindGustMph = getSingleIntValue(&rv, "windgustmph", -1)
	me.RainInch = getSingleFloatValue(&rv, "rainin", -1)
	me.DailyRainInch = getSingleFloatValue(&rv, "dailyrainin", -1)
	me.WeeklyRainInch = getSingleFloatValue(&rv, "weeklyrainin", -1)
	me.MonthlyRainInch = getSingleFloatValue(&rv, "monthlyrainin", -1)
	me.YearlyRainInch = getSingleFloatValue(&rv, "yearlyrainin", -1)
	me.SolarRadiation = getSingleIntValue(&rv, "solarradiation", -1)
	me.UV = getSingleIntValue(&rv, "UV", -1)
	me.IndoorTempF = getSingleIntValue(&rv, "indoortempf", -1)
	me.IndoorHumidity = getSingleIntValue(&rv, "indoorhumidity", -1)
	me.BarometricInch = getSingleIntValue(&rv, "baromin", -1)
	me.LowBattery = getSingleIntValue(&rv, "lowbatt", -1)
	me.DateUTC = getSingleValue(&rv, "dateutc", "")
	me.SoftwareType = getSingleValue(&rv, "softwaretype", "")
	me.Action = getSingleValue(&rv, "action", "")
	me.RealTime = getSingleIntValue(&rv, "realtime", -1)
	me.RealTimeFrequency = getSingleIntValue(&rv, "rtfreq", -1)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "success")
}
