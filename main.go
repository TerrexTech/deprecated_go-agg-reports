package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/TerrexTech/go-agg-reports/report"
	"github.com/TerrexTech/go-commonutils/commonutil"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type Env struct {
	Reportdb    report.DBI
	Metricdb    report.DBI
	Inventorydb report.DBI
	Devicedb    report.DBI
}

type ReportResponse struct {
	Inventory []report.Inventory
	Metric    []report.Metric
	Device    []report.Device
}

func main() {
	err := godotenv.Load()
	if err != nil {
		err = errors.Wrap(err,
			".env file not found, env-vars will be read as set in environment",
		)
		log.Println(err)
	}

	hosts := os.Getenv("MONGO_HOSTS")
	username := os.Getenv("MONGO_USERNAME")
	password := os.Getenv("MONGO_PASSWORD")
	database := os.Getenv("MONGO_DATABASE")
	collectionReport := os.Getenv("MONGO_REP_COLLECTION")
	collectionInv := os.Getenv("MONGO_INV_COLLECTION")
	collectionMet := os.Getenv("MONGO_METRIC_COLLECTION")
	collectionDev := os.Getenv("MONGO_DEVICE_COLLECTION")
	// collectionWarn := os.Getenv("MONGO_WARNING_COLLECTION")
	// collectionFlash := os.Getenv("MONGO_FLASHSALE_COLLECTION")

	timeoutMilliStr := os.Getenv("MONGO_TIMEOUT")
	parsedTimeoutMilli, err := strconv.Atoi(timeoutMilliStr)
	if err != nil {
		err = errors.Wrap(err, "Error converting Timeout value to int32")
		log.Println(err)
		log.Println("MONGO_TIMEOUT value will be set to 3000 as default value")
		parsedTimeoutMilli = 3000
	}
	timeoutMilli := uint32(parsedTimeoutMilli)

	log.Println(hosts)

	configReport := report.DBIConfig{
		Hosts:               *commonutil.ParseHosts(hosts),
		Username:            username,
		Password:            password,
		TimeoutMilliseconds: timeoutMilli,
		Database:            database,
		Collection:          collectionReport,
	}

	configMetric := report.DBIConfig{
		Hosts:               *commonutil.ParseHosts(hosts),
		Username:            username,
		Password:            password,
		TimeoutMilliseconds: timeoutMilli,
		Database:            database,
		Collection:          collectionMet,
	}

	configInv := report.DBIConfig{
		Hosts:               *commonutil.ParseHosts(hosts),
		Username:            username,
		Password:            password,
		TimeoutMilliseconds: timeoutMilli,
		Database:            database,
		Collection:          collectionInv,
	}
	configDev := report.DBIConfig{
		Hosts:               *commonutil.ParseHosts(hosts),
		Username:            username,
		Password:            password,
		TimeoutMilliseconds: timeoutMilli,
		Database:            database,
		Collection:          collectionDev,
	}

	// configWarn := report.DBIConfig{
	// 	Hosts:               *commonutil.ParseHosts(hosts),
	// 	Username:            username,
	// 	Password:            password,
	// 	TimeoutMilliseconds: timeoutMilli,
	// 	Database:            database,
	// 	Collection:          collectionWarn,
	// }

	// configFlash := report.DBIConfig{
	// 	Hosts:               *commonutil.ParseHosts(hosts),
	// 	Username:            username,
	// 	Password:            password,
	// 	TimeoutMilliseconds: timeoutMilli,
	// 	Database:            database,
	// 	Collection:          collectionFlash,
	// }

	dbReport, err := report.GenerateDB(configReport, &report.Report{})
	if err != nil {
		err = errors.Wrap(err, "Error connecting to Inventory DB")
		log.Println(err)
		return
	}

	dbMetric, err := report.GenerateDB(configMetric, &report.Metric{})
	if err != nil {
		err = errors.Wrap(err, "Error connecting to Inventory DB")
		log.Println(err)
		return
	}

	dbInventory, err := report.GenerateDB(configInv, &report.Inventory{})
	if err != nil {
		err = errors.Wrap(err, "Error connecting to Inventory DB")
		log.Println(err)
		return
	}

	dbDevice, err := report.GenerateDB(configDev, &report.Device{})
	if err != nil {
		err = errors.Wrap(err, "Error connecting to Inventory DB")
		log.Println(err)
		return
	}

	env := &Env{
		Reportdb:    dbReport,
		Metricdb:    dbMetric,
		Inventorydb: dbInventory,
		Devicedb:    dbDevice,
	}

	http.HandleFunc("/create-data", env.LoadDataInMongo)
	http.HandleFunc("/inv-report", env.InvReport)
	http.HandleFunc("/met-report", env.MetricReport)
	http.HandleFunc("/dev-report", env.DeviceReport)

	http.ListenAndServe(":8080", nil)

	// numValToGen := 5
	// generatedData := report.CreateAllData()
	// dbReport.GenReportData(numValToGen, generatedData)
	// dbMetric.GenMetricData(numValToGen, generatedData)
	// dbInventory.GenInventoryData(numValToGen, generatedData)

	//Searching Inventory
	// searchInv := []report.SearchByFieldVal{
	// 	report.SearchByFieldVal{
	// 		SearchField: "name",
	// 		SearchVal:   "Mango",
	// 	},
	// }

	// env.InvenReport(searchInv)

	//Search Metric
	// searchInv := []report.SearchByFieldVal{
	// 	report.SearchByFieldVal{
	// 		SearchField: "name",
	// 		SearchVal:   "Mango",
	// 	},
	// }

	// var test map[string]map[string]map[string]string
	// test = []byte(`{"inventory":{"fieldname"`)

}

func (env *Env) LoadDataInMongo(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}

	numValues := 10

	// genData := report.CreateAllData()

	var rep []report.Report
	var metric []report.Metric
	var inventory []report.Inventory
	var device []report.Device
	for i := 0; i < numValues; i++ {
		genData := report.CreateAllData()
		rep = append(rep, genData.RType)
		metric = append(metric, genData.MType)
		inventory = append(inventory, genData.IType)
		device = append(device, genData.DType)
	}

	// DB connection
	reportData, err := env.Reportdb.GenReportData(rep)
	if err != nil {
		err = errors.Wrap(err, "Unable to create new data in mongo")
		log.Println(err)
		return
	}
	metricData, err := env.Metricdb.GenMetricData(metric)
	if err != nil {
		err = errors.Wrap(err, "Unable to create new data in mongo")
		log.Println(err)
		return
	}
	inventoryData, err := env.Inventorydb.GenInventoryData(inventory)

	if err != nil {
		err = errors.Wrap(err, "Unable to create new data in mongo")
		log.Println(err)
		return
	}

	deviceData, err := env.Devicedb.GenDeviceData(device)

	if err != nil {
		err = errors.Wrap(err, "Unable to create new data in mongo")
		log.Println(err)
		return
	}

	// log.Println(reportData)
	rData, err := json.Marshal(&reportData)
	if err != nil {
		err = errors.Wrap(err, "Unable to create response body")
		log.Println(err)
	}
	mData, err := json.Marshal(&metricData)
	if err != nil {
		err = errors.Wrap(err, "Unable to create response body")
		log.Println(err)
	}
	iData, err := json.Marshal(&inventoryData)
	if err != nil {
		err = errors.Wrap(err, "Unable to create response body")
		log.Println(err)
		return
	}
	dData, err := json.Marshal(&deviceData)
	if err != nil {
		err = errors.Wrap(err, "Unable to create response body")
		log.Println(err)
		return
	}
	w.Write(rData)
	w.Write(mData)
	w.Write(iData)
	w.Write(dData)
}

func (env *Env) InvReport(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrap(err, "Unable to read the request body")
		log.Println(err)
		return
	}

	var query map[string][]report.SearchParam

	// advSearch := report.AdvancedInvSearch{}
	err = json.Unmarshal(body, &query)
	if err != nil {
		err = errors.Wrap(err, "Unable to unmarshal - advSearch")
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	invSearchResult, err := env.Inventorydb.InvAdvSearch(query)
	if err != nil {
		err = errors.Wrap(err, "Unable to read the search Inventory - invResult")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	invByte, err := json.Marshal(&invSearchResult)
	if err != nil {
		err = errors.Wrap(err, "Unable to marshal Inventory results - InvReport")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(invByte)
}

func (env *Env) MetricReport(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrap(err, "Unable to read the request body")
		log.Println(err)
		return
	}

	log.Println(string(body))

	var query map[string][]report.SearchParam

	// advSearch := report.AdvancedInvSearch{}
	err = json.Unmarshal(body, &query)
	if err != nil {
		err = errors.Wrap(err, "Unable to unmarshal - advSearch")
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	invSearchResult, err := env.Inventorydb.InvAdvSearch(query)
	if err != nil {
		err = errors.Wrap(err, "Unable to read the search Inventory - invResult")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println(invSearchResult)

	metricResult, err := env.Metricdb.MetAdvSearch(invSearchResult)
	if err != nil {
		err = errors.Wrap(err, "Did not get metric query result - MetricReport")
		log.Println(err)
		return
	}

	log.Println(metricResult, "&&&&&&&&&&&&")

	respObject := ReportResponse{
		Inventory: invSearchResult,
		Metric:    metricResult,
	}

	metricByte, err := json.Marshal(&respObject)
	if err != nil {
		err = errors.Wrap(err, "Unable to marshal Metric results - MetricReport")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(metricByte)
}

func (env *Env) DeviceReport(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrap(err, "Unable to read the request body")
		log.Println(err)
		return
	}

	log.Println(string(body))

	var query map[string][]report.SearchParam

	// advSearch := report.AdvancedInvSearch{}
	err = json.Unmarshal(body, &query)
	if err != nil {
		err = errors.Wrap(err, "Unable to unmarshal - advSearch")
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	invSearchResult, err := env.Inventorydb.InvAdvSearch(query)
	if err != nil {
		err = errors.Wrap(err, "Unable to read the search Inventory - invResult")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	deviceResult, err := env.Devicedb.DevAdvSearch(invSearchResult)
	if err != nil {
		err = errors.Wrap(err, "Did not get metric query result - MetricReport")
		log.Println(err)
		return
	}

	respObject := ReportResponse{
		Inventory: invSearchResult,
		Device:    deviceResult,
	}

	deviceByte, err := json.Marshal(&respObject)
	if err != nil {
		err = errors.Wrap(err, "Unable to marshal Metric results - MetricReport")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(deviceByte)
}

func (env *Env) DistInvReport(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}

	// body, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	err = errors.Wrap(err, "Unable to read the request body")
	// 	log.Println(err)
	// 	return
	// }

	// log.Println(string(body))

	// var query map[string][]report.SearchParam

	// // advSearch := report.AdvancedInvSearch{}
	// err = json.Unmarshal(body, &query)
	// if err != nil {
	// 	err = errors.Wrap(err, "Unable to unmarshal - advSearch")
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	log.Println(err)
	// 	return
	// }

	// invSearchResult, err := env.Inventorydb.InvAdvSearch(query)
	// if err != nil {
	// 	err = errors.Wrap(err, "Unable to read the search Inventory - invResult")
	// 	log.Println(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// invResults, err := env.Inventorydb.DistributeInvFields()
	// if err != nil {
	// 	err = errors.Wrap(err, "Did not get metric query result - MetricReport")
	// 	log.Println(err)
	// 	return
	// }

	// respObject := ReportResponse{
	// 	Inventory: invSearchResult,
	// 	Device:    deviceResult,
	// }

	// deviceByte, err := json.Marshal(&respObject)
	// if err != nil {
	// 	err = errors.Wrap(err, "Unable to marshal Metric results - MetricReport")
	// 	log.Println(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// w.Write(deviceByte)
}
