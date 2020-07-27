package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/matt-doug-davidson/timestamps"
	"gopkg.in/yaml.v2"
)

// SensorYaml holds the configuration for a single sensor
type SensorYaml struct {
	Sensor  string  `yaml:"sensor"`
	Minimum float64 `yaml:"min"`
	Maximum float64 `yaml:"max"`
}

type Config struct {
	Sensors []SensorYaml
}

type Instrument struct {
	SerialNumber  string
	TcpPort       string
	UserName      string
	Password      string
	SensorsConfig Config
	PreviousData  map[string]interface{}
}

// The sensor data creatation attributes
var sensors = `sensors:
-
   sensor: NO2
   min: 0.0
   max: 100.0
-
   sensor: Ox
   min: 0.0
   max: 100.0
-
   sensor: O3
   min: 0.0
   max: 100.0
-
   sensor: "O3 raw"
   min: 0.0
   max: 100.0
-
   sensor: "PM2.5"
   min: 0.0
   max: 100.0
-
   sensor: "PM2.5 raw"
   min: 0.0
   max: 100.0
-
   sensor: "PM10 raw"
   min: 0.0
   max: 100.0
-
   sensor: "PM10"
   min: 0.0
   max: 100.0
-
   sensor: "TEMP"
   min: 0.0
   max: 100.0
-
   sensor: "RH"
   min: 0.0
   max: 100.0
-
   sensor: "DP"
   min: 0.0
   max: 100.0`

func (i *Instrument) init() {
	for i, v := range i.SensorsConfig.Sensors {
		fmt.Println(i)
		fmt.Println(v)
		fmt.Println(v.Sensor)

	}
	now, _, _, _ := createIntervalTimestrings()
	i.PreviousData = i.createDataMap(now)
}

func (i *Instrument) handleLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleLogin")
	fmt.Println(i)
	if r.Method != "POST" {
		http.Error(w, "405 Not FOund", http.StatusNotFound)
	}
	contentType := r.Header.Get("content-type")
	fmt.Println("content-type ", contentType)
	userName := ""
	password := ""
	if contentType == "application/x-www-form-urlencoded" {
		err := r.ParseForm()
		if err != nil {
			fmt.Println("Error parsing form. Cause: ", err.Error())
		}
		fmt.Println(r.Body)
		v := r.Form
		fmt.Println(v)
		userName = r.Form.Get("UserName")
		fmt.Println("userName: ", userName)
		password = r.Form.Get("Password")
		fmt.Println("password: ", password)
	}
	if userName == "" {
		fmt.Println("Error: No UserName in parameters")
		http.Error(w, "Username not found", http.StatusNotFound)
		return
	}
	if password == "" {
		fmt.Println("Error: No Password in parameters")
		http.Error(w, "Password not found", http.StatusNotFound)
		return
	}
	if userName != i.UserName || password != i.Password {
		fmt.Println("Error: Username/Password mismatch")
		http.Error(w, "Error Unathorized", http.StatusUnauthorized)
		return
	}

	// Expire after 5 years
	expires := time.Now().AddDate(5, 0, 0)

	loginCookie := http.Cookie{
		Name:    ".MONOAUTH",
		Value:   "90lkjsdfsdfkjsdfksfkjsdfksjdflskjdfjsdfhsfzxcvbmlAzECDE12as345323r52222",
		Path:    "/",
		Expires: expires,
	}
	http.SetCookie(w, &loginCookie)
	w.Header().Set("Content-type", "text/html")
	w.WriteHeader(200)
	w.Write([]byte("Success"))
}

func (i *Instrument) getInstruments() [1]string {
	instruments := [1]string{i.SerialNumber}
	fmt.Println(instruments)
	jsonString, err := json.Marshal(instruments)
	if err != nil {
		fmt.Println("Error marshal'ing instrument. Cause: ", err.Error())
	}
	fmt.Println(string(jsonString))
	fmt.Println(jsonString)
	return instruments
}

func (i *Instrument) handleInstruments(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleInstruments")
	if r.Method != "GET" {
		fmt.Println("Not GET API")
		http.Error(w, "Error method not allowed", http.StatusMethodNotAllowed)
		return
	}
	instruments := i.getInstruments()
	jsonBytes, err := json.Marshal(instruments)
	if err != nil {
		fmt.Println("Error marshal'ing instrument. Cause: ", err.Error())
	}
	fmt.Println(string(jsonBytes))
	fmt.Println(jsonBytes)
	w.Header().Set("Content-type", "json")
	w.WriteHeader(200)
	w.Write(jsonBytes)
	fmt.Println(instruments)
}

func (i *Instrument) getInstrument() map[string]interface{} {

	sensors := [11]map[string]interface{}{
		{"name": "NO2", "units": "ppb", "decimalPlaces": 1},
		{"name": "Ox", "units": "ppb", "decimalPlaces": 1},
		{"name": "O3", "units": "ppb", "decimalPlaces": 1},
		{"name": "O3 raw", "units": "ppb", "decimalPlaces": 1},
		{"name": "PM2.5 raw", "units": "ug/m3", "decimalPlaces": 1},
		{"name": "PM2.5", "units": "ug/m3", "decimalPlaces": 1},
		{"name": "PM10 raw", "units": "ug/m3", "decimalPlaces": 1},
		{"name": "PM10", "units": "ug/m3", "decimalPlaces": 1},
		{"name": "TEMP", "units": "C", "decimalPlaces": 2},
		{"name": "RH", "units": "%", "decimalPlaces": 1},
		{"name": "DP", "units": "C", "decimalPlaces": 1},
	}

	fmt.Println(sensors)
	instrument := map[string]interface{}{
		"serial":       i.SerialNumber,
		"name":         i.SerialNumber,
		"organisation": "Aeroqual",
		"network":      "Unassigned instruments",
		"timeZone":     "(UTC+02:00) Athens, Bucharest",
		"sensors":      sensors,
	}
	fmt.Println(instrument)
	jsonString, err := json.Marshal(instrument)
	if err != nil {
		fmt.Println("Error marshal'ing instrument. Cause: ", err.Error())
	}
	fmt.Println(string(jsonString))
	return instrument
}

func (i *Instrument) handleInstrument(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleInstrument")
	// Check that the instrument configured matches the one in the message

	if r.Method != "GET" {
		fmt.Println("Not GET API")
	}
	instrument := i.getInstrument()
	fmt.Println(instrument)
	jsonBytes, err := json.Marshal(instrument)
	if err != nil {
		fmt.Println("Error marshal'ing instrument. Cause: ", err.Error())
	}
	fmt.Println(string(jsonBytes))
	w.Header().Set("Content-type", "json")
	w.WriteHeader(200)
	w.Write(jsonBytes)
	fmt.Println(instrument)
}

func (i *Instrument) handleData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleData")
	if r.Method != "GET" {
		fmt.Println("Not GET API")
	}
	data := i.createData()
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshal'ing data. Cause: ", err.Error())
	}
	fmt.Println(string(jsonBytes))
	w.Header().Set("Content-type", "json")
	w.WriteHeader(200)
	w.Write(jsonBytes)
}

func (i *Instrument) run() {
	http.HandleFunc("/api/account/login", i.handleLogin)
	http.HandleFunc("/api/instrument", i.handleInstruments)
	http.HandleFunc("/api/instrument/"+i.SerialNumber, i.handleInstrument)
	http.HandleFunc("/api/data/"+i.SerialNumber, i.handleData)
	log.Fatal(http.ListenAndServe(":"+i.TcpPort, nil))
	for {
		time.Sleep(1 * time.Second)
		fmt.Println("Instrument ", i.SerialNumber)
	}
}
func (i *Instrument) createData() map[string]interface{} {
	now, twenty, twelve, two := createIntervalTimestrings()
	fmt.Println(twenty, twelve, two)
	currentData := i.createDataMap(two)
	fmt.Println("previousData:\n", i.PreviousData)
	fmt.Println("currentData:\n", currentData)
	data := [2]interface{}{
		i.PreviousData, currentData,
	}
	sdData := map[string]interface{}{
		"serial":             i.SerialNumber,
		"name":               i.SerialNumber,
		"organisation":       "Aeroqual",
		"network":            "Unassigned instruments",
		"from":               twenty,
		"to":                 now,
		"averagingPeriod":    10,
		"timeZone":           "(UTC+02:00) Athens, Bucharest",
		"summerTimeAdjusted": true,
		"data":               data,
	}
	i.PreviousData = currentData
	fmt.Println("sdData:\n", sdData)
	jsonString, err := json.Marshal(sdData)
	if err != nil {
		fmt.Println("Error marshal'ing sd data. Cause: ", err.Error())
	}
	fmt.Println(string(jsonString))
	return sdData
}

func (i *Instrument) createDataMap(timestring string) map[string]interface{} {
	data := map[string]interface{}{}
	data["Time"] = timestring
	for _, v := range i.SensorsConfig.Sensors {
		value := createSensorValue(v.Minimum, v.Maximum)
		data[v.Sensor] = value
	}

	fmt.Println(data)
	return data
}

func createSensorValue(min float64, max float64) float64 {
	valueRange := max - min
	value := rand.Float64() * valueRange
	value += min
	// Round to 1 decimal place
	return math.Round(value*10) / 10
}

func createIntervalTimestrings() (string, string, string, string) {
	timeStamp := timestamps.RoundDownMinutes(timestamps.LocalTimestamp())
	return timestamps.TimestampToTimestringNoMilli(timeStamp),
		timestamps.TimestampToTimestringNoMilli(timestamps.MinutesEarlier(timeStamp, 20)),
		timestamps.TimestampToTimestringNoMilli(timestamps.MinutesEarlier(timeStamp, 12)),
		timestamps.TimestampToTimestringNoMilli(timestamps.MinutesEarlier(timeStamp, 2))
}

func main() {

	sn := os.Getenv("SN")
	tcpPort := os.Getenv("PORT")
	userName := os.Getenv("UN")
	password := os.Getenv("PW")

	if sn == "" {
		panic("Enviroment variable SN not defined")
	}
	if tcpPort == "" {
		panic("Enviroment variable PORT not defined")
	}
	if userName == "" {
		panic("Enviroment variable UN not defined")
	}
	if password == "" {
		panic("Enviroment variable PW not defined")
	}

	//Create config structure
	config := Config{}
	err := yaml.Unmarshal([]byte(sensors), &config)
	if err != nil {
		fmt.Println("Error unmarshal'ing sensors. Cause: ", err.Error())
	}

	var inst = Instrument{
		SerialNumber: sn, TcpPort: tcpPort,
		UserName: userName, Password: password,
		SensorsConfig: config}
	inst.init()
	inst.run()

	for {
		time.Sleep(1 * time.Second)
	}

}
