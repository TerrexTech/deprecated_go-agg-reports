package report

import (
	"log"
	"strconv"

	"github.com/TerrexTech/go-mongoutils/mongo"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/pkg/errors"
)

// type ConfigSchema struct {
// 	Report    *report.Report
// 	Metric    *Metric
// 	Inventory *Inventory
// }

// DBIConfig is the configuration for the authDB.
type DBIConfig struct {
	Hosts               []string
	Username            string
	Password            string
	TimeoutMilliseconds uint32
	Database            string
	Collection          string
}

// DBI is the Database-interface for reporting.
// This fetches/writes data to/from database for generating reports
type DBI interface {
	Collection() *mongo.Collection
	GenReportData(report []Report) ([]Report, error)
	GenMetricData(metric []Metric) ([]Metric, error)
	GenInventoryData(inventory []Inventory) ([]Inventory, error)
	GenDeviceData(device []Device) ([]Device, error)
	SearchKeyVal(search []SearchByFieldVal) ([]interface{}, error)
	// InsertIntoReport(inv []Inventory, reportType string) (*mgo.InsertOneResult, error)
	InvAdvSearch(search map[string][]SearchParam) ([]Inventory, error)
	MetAdvSearch(searchInv []Inventory) ([]Metric, error)
	DevAdvSearch(searchInv []Inventory) ([]Device, error)
	DistributionInvFields() ([]InvenReport, error)
	// // SearchByTimestamp(search *SearchByDate) (*Report, error)
	// SearchByFieldVal(search []SearchByFieldVal) ([]interface{}, error)
}

// DB is the implementation for dbI.
// dbI is the Database-interface for generating reports.
type DB struct {
	collection *mongo.Collection
}

type InvenReport struct {
	ProdName    string  `bson:"prod_name,omitempty" json:"prod_name,omitempty"`
	ProdWeight  float64 `bson:"prod_weight,omitempty" json:"prod_weight,omitempty"`
	TotalWeight float64 `bson:"total_weight,omitempty" json:"total_weight,omitempty"`
	SoldWeight  float64 `bson:"sold_weight,omitempty" json:"sold_weight,omitempty"`
	WasteWeight float64 `bson:"waste_weight,omitempty" json:"waste_weight,omitempty"`
	ProductSold int64   `bson:"prod_sold,omitempty" json:"prod_sold,omitempty"`
}

// type SearchByDate struct {
// 	EndDate   int64 `bson:"end_date,omitempty" json:"end_date,omitempty"`
// 	StartDate int64 `bson:"start_date,omitempty" json:"start_date,omitempty"`
// }

// type SearchByFieldVal struct {
// 	SearchField string      `bson:"search_field,omitempty" json:"search_field,omitempty"`
// 	SearchVal   interface{} `bson:"search_val,omitempty" json:"search_val,omitempty"`
// }

// type AdvancedInvSearch struct {
// 	SKU         int64  `bson:"sku,omitempty" json:"sku,omitempty"`
// 	Name        string `bson:"name,omitempty" json:"name,omitempty"`
// 	Origin      string `bson:"origin,omitempty" json:"origin,omitempty"`
// 	DateArrived int64  `bson:"date_arrived,omitempty" json:"date_arrived,omitempty"`
// 	EndDate     int64  `bson:"end_date,omitempty" json:"end_date,omitempty"`
// 	StartDate   int64  `bson:"start_date,omitempty" json:"start_date,omitempty"`
// }

// type AdvancedMetSearch struct {
// 	// DeviceID uuuid.UUID `bson:"device_id,omitempty" json:"device_id,omitempty"`
// 	// TempIn   float64    `bson:"temp_in,omitempty" json:"temp_in,omitempty"`
// 	// Humidity float64    `bson:"humidity,omitempty" json:"humidity,omitempty"`
// 	// Ethylene float64    `bson:"ethylene,omitempty" json:"ethylene,omitempty"`
// 	// CarbonDi float64    `bson:"carbon_di,omitempty" json:"carbon_di,omitempty"`
// 	MetSearch Metric `bson:"met_search,omitempty" json:"met_search,omitempty"`
// 	EndDate   int64  `bson:"end_date,omitempty" json:"end_date,omitempty"`
// 	StartDate int64  `bson:"start_date,omitempty" json:"start_date,omitempty"`
// }

// type AdvanceSearch struct {
// 	Inventory    Inventory `bson:"inventory,omitempty" json:"inventory,omitempty"`
// 	Metric       Metric    `bson:"inventory,omitempty" json:"inventory,omitempty"`
// 	TypeOfReport string
// }

func GenerateDB(dbConfig DBIConfig, schema interface{}) (*DB, error) {
	config := mongo.ClientConfig{
		Hosts:               dbConfig.Hosts,
		Username:            dbConfig.Username,
		Password:            dbConfig.Password,
		TimeoutMilliseconds: dbConfig.TimeoutMilliseconds,
	}

	client, err := mongo.NewClient(config)
	if err != nil {
		err = errors.Wrap(err, "Error creating DB-client")
		return nil, err
	}

	conn := &mongo.ConnectionConfig{
		Client:  client,
		Timeout: 5000,
	}

	// ====> Create New Collection
	collConfig := &mongo.Collection{
		Connection:   conn,
		Database:     dbConfig.Database,
		Name:         dbConfig.Collection,
		SchemaStruct: schema,
		// Indexes:      indexConfigs,
	}
	c, err := mongo.EnsureCollection(collConfig)
	if err != nil {
		err = errors.Wrap(err, "Error creating DB-client")
		return nil, err
	}
	return &DB{
		collection: c,
	}, nil
}

func (d *DB) Collection() *mongo.Collection {
	return d.collection
}

func (db *DB) GenReportData(report []Report) ([]Report, error) {
	// report := []Report{}
	// for i := 0; i < numOfVal; i++ {
	// 	genData := CreateAllData()
	// 	report = append(report, genData.RType)
	// }

	// log.Println(report)

	for _, v := range report {
		insertResult, err := db.collection.InsertOne(v)
		if err != nil {
			err = errors.Wrap(err, "Unable to insert data")
			log.Println(err)
			return nil, err
		}
		log.Println(insertResult)
	}
	return report, nil
}

func (db *DB) GenMetricData(metric []Metric) ([]Metric, error) {
	// metric := []Metric{}
	// for i := 0; i < numOfVal; i++ {
	// 	genData := CreateAllData()
	// 	metric = append(metric, genData.MType)
	// }

	for _, v := range metric {
		insertResult, err := db.collection.InsertOne(v)
		if err != nil {
			err = errors.Wrap(err, "Unable to insert data")
			log.Println(err)
			return nil, err
		}
		log.Println(insertResult)
	}
	return metric, nil
}

func (db *DB) GenInventoryData(inventory []Inventory) ([]Inventory, error) {
	// inventory := []Inventory{}
	// for i := 0; i < numOfVal; i++ {
	// 	genData := CreateAllData()
	// 	inventory = append(inventory, genData.IType)
	// }

	for _, v := range inventory {
		insertResult, err := db.collection.InsertOne(v)
		if err != nil {
			err = errors.Wrap(err, "Unable to insert data")
			log.Println(err)
			return nil, err
		}
		log.Println(insertResult)
	}
	return inventory, nil
}

func (db *DB) GenDeviceData(device []Device) ([]Device, error) {
	// device := []Device{}
	// for i := 0; i < numOfVal; i++ {
	// 	genData := CreateAllData()
	// 	device = append(device, genData.DType)
	// }

	for _, v := range device {
		insertResult, err := db.collection.InsertOne(v)
		if err != nil {
			err = errors.Wrap(err, "Unable to insert data")
			log.Println(err)
			return nil, err
		}
		log.Println(insertResult)
	}
	return device, nil
}

type SearchByFieldVal struct {
	SearchField string      `bson:"search_field,omitempty" json:"search_field,omitempty"`
	SearchVal   interface{} `bson:"search_val,omitempty" json:"search_val,omitempty"`
}

func (db *DB) SearchKeyVal(search []SearchByFieldVal) ([]interface{}, error) {

	var findResults []interface{}
	var err error

	for _, v := range search {
		if v.SearchField != "" && v.SearchVal != "" {
			findResults, err = db.collection.Find(map[string]interface{}{
				v.SearchField: map[string]interface{}{
					"$eq": &v.SearchVal,
				},
			})
		}
	}

	if err != nil {
		err = errors.Wrap(err, "Error while fetching product.")
		log.Println(err)
		return nil, err
	}

	//length
	if len(findResults) == 0 {
		msg := "No results found - SearchByDate"
		return nil, errors.New(msg)
	}
	return findResults, nil
}

// func (db *DB) InvAdvSearch(search Inventory) ([]Inventory, error) {

// 	var findResults []interface{}
// 	var err error

// 	// searchFields := map[string]interface{}

// 	// y
// 	// y["field"]
// 	// y["type"]
// 	// switch
// 	// string
// 	// int
// 	// y["equal"].(int)

// 	findResults, err = db.collection.Find(map[string]interface{}{
// 		"sku": map[string]int64{
// 			"$eq": search.SKU,
// 		},
// 		"name": map[string]string{
// 			"$eq": search.Name,
// 		},
// 		// "origin": map[string]string{
// 		// 	"$eq": search.Origin,
// 		// },
// 		// map["inventory"].(map[string]interface{})
// 		// "timestamp": map[string]int64{
// 		// 	"$lte": search.EndDate,
// 		// 	"$gte": search.StartDate,
// 		// },
// 	})

// 	if err != nil {
// 		err = errors.Wrap(err, "Error while fetching results from inventory.")
// 		log.Println(err)
// 		return nil, err
// 	}

// 	//length
// 	if len(findResults) == 0 {
// 		msg := "No results found - InvAdvSearch"
// 		return nil, errors.New(msg)
// 	}

// 	inventory := []Inventory{}

// 	for _, v := range findResults {
// 		result := v.(*Inventory)
// 		inventory = append(inventory, *result)
// 	}
// 	return inventory, nil
// }

//InvAdvSearch - uses AdvancedInvSearch struct
func (db *DB) InvAdvSearch(search map[string][]SearchParam) ([]Inventory, error) {

	var findResults []interface{}
	var err error
	// var type string

	inv := search["inventory"]

	findParams := map[string]interface{}{}

	if inv != nil {
		for _, v := range inv {
			if v.Type == "" {
				err := errors.Wrap(err, "Type required - InvAdvSearch.")
				log.Println(err)
				return nil, err
			}
			if v.Field == "" {
				err := errors.Wrap(err, "Field is required - InvAdvSearch")
				log.Println(err)
				return nil, err
			}
			if v.Equal == "" && v.LowerLimit == 0 && v.UpperLimit == 0 {
				err = errors.Wrap(err, "Missing value in equal. No lowerlimit and upperlimit set - InvAdvSearch.")
				log.Println(err)
				return nil, err
			}

			switch v.Type {
			case "string":
				findParams[v.Field] = map[string]string{
					"$eq": v.Equal,
				}

			case "float":
				if v.Equal != "" {
					floatValue, err := strconv.ParseFloat(v.Equal, 64)
					if err != nil {
						err = errors.Wrap(err, "Error converting value of equal to float - InvAdvSearch")
						log.Println(err)
						return nil, err
					}
					findParams[v.Field] = map[string]float64{
						"$eq": floatValue,
					}
				} else {
					limitMap := map[string]map[string]float64{}
					if v.LowerLimit != 0 {
						limitMap[v.Field]["$gt"] = v.LowerLimit
					}
					if v.UpperLimit != 0 {
						limitMap[v.Field]["$lt"] = v.UpperLimit
					}
					findParams[v.Field] = limitMap[v.Field]
				}

			case "int":
				if v.Equal != "" {
					intValue, err := strconv.ParseInt(v.Equal, 10, 64)
					if err != nil {
						err = errors.Wrap(err, "Error converting equal to int - InvAdvSearch")
						log.Println(err)
						return nil, err
					}
					findParams[v.Field] = map[string]int64{
						"$eq": intValue,
					}
				} else {
					limitMap := map[string]map[string]int64{}
					if v.LowerLimit != 0 {
						limitMap[v.Field]["$gt"] = int64(v.LowerLimit)
					}
					if v.UpperLimit != 0 {
						limitMap[v.Field]["$lt"] = int64(v.UpperLimit)
					}

					findParams[v.Field] = limitMap[v.Field]
				}
			}
			if v.Type == "string" {
				findParams[v.Field] = map[string]interface{}{
					"$eq": v.Equal,
				}
			}

			// else if v.Equal {
			// 	findParams[v.Field] = map[string]interface{}{
			// 		"$eq": v.Equal,
			// 	}
			// 	log.Println(findParams[v.Field], "dsklfjdskfsl")
			// }
		}
	}

	log.Println(findParams, "###123################")

	findResults, err = db.collection.Find(findParams)
	log.Println("FFFFFFFFFFFFFFF")
	for _, r := range findResults {
		log.Printf("%+v", r)
	}

	if err != nil {
		err = errors.Wrap(err, "Error while fetching results from inventory.")
		log.Println(err)
		return nil, err
	}

	//length
	if len(findResults) == 0 {
		msg := "No results found - InvAdvSearch"
		return nil, errors.New(msg)
	}

	inventory := []Inventory{}

	for _, v := range findResults {
		result := v.(*Inventory)
		inventory = append(inventory, *result)
	}
	return inventory, nil
}

func (db *DB) MetAdvSearch(searchInv []Inventory) ([]Metric, error) {

	var findResults []interface{}
	var err error

	findParams := map[string]interface{}{}

	for _, v := range searchInv {
		log.Println("%+v", v.ItemID)
		findParams["item_id"] = map[string]interface{}{
			"$eq": v.ItemID.String(),
		}
	}

	findResults, err = db.collection.Find(findParams)

	if err != nil {
		err = errors.Wrap(err, "Error while fetching results from inventory.")
		log.Println(err)
		return nil, err
	}

	//length
	if len(findResults) == 0 {
		msg := "No results found - InvMetSearch"
		return nil, errors.New(msg)
	}

	metric := []Metric{}

	for _, v := range findResults {
		result := v.(*Metric)
		metric = append(metric, *result)
	}
	return metric, nil
}

func (db *DB) DevAdvSearch(searchInv []Inventory) ([]Device, error) {

	var findResults []interface{}
	var err error

	findParams := map[string]interface{}{}

	for _, v := range searchInv {
		log.Println("%+v", v.DeviceID)
		findParams["device_id"] = map[string]interface{}{
			"$eq": v.DeviceID.String(),
		}
	}

	findResults, err = db.collection.Find(findParams)

	if err != nil {
		err = errors.Wrap(err, "Error while fetching results from device - DevAdvSearch.")
		log.Println(err)
		return nil, err
	}

	//length
	if len(findResults) == 0 {
		msg := "No results found - DevAdvSearch"
		return nil, errors.New(msg)
	}

	device := []Device{}

	for _, v := range findResults {
		result := v.(*Device)
		device = append(device, *result)
	}
	return device, nil
}

func (db *DB) DistributionInvFields() ([]InvenReport, error) {
	pipeline := bson.NewArray(
		bson.VC.Document(
			bson.NewDocument(
				bson.EC.SubDocumentFromElements(
					"$group",
					bson.EC.String("_id", "$name"),
					bson.EC.SubDocumentFromElements(
						"total_weight",
						bson.EC.String("$sum", "$total_weight"),
					),
					bson.EC.SubDocumentFromElements(
						"waste_weight",
						bson.EC.String("$sum", "$waste_weight"),
					),
					bson.EC.SubDocumentFromElements(
						"sold_weight",
						bson.EC.String("$sum", "$sold_weight"),
					),
				),
			),
		),
	)
	aggResults, err := db.collection.Aggregate(pipeline)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	// log.Println(aggResults)

	dist := []InvenReport{}

	for _, v := range aggResults {
		value := v.(map[string]interface{})
		strValue := value["_id"].(string)
		twValue := value["total_weight"].(float64)
		wwValue := value["waste_weight"].(float64)
		swValue := value["sold_weight"].(float64)
		dist = append(dist, InvenReport{
			ProdName:    strValue,
			TotalWeight: twValue,
			WasteWeight: wwValue,
			SoldWeight:  swValue,
		})
	}

	// distWeight, err := json.Marshal(&dist)
	// if err != nil {
	// 	err = errors.Wrap(err, "Unable to marshal distribution")
	// 	log.Println(err)
	// 	return nil, err
	// }

	return dist, nil
}
