package report

import (
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/TerrexTech/uuuid"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/pkg/errors"
)

type Device struct {
	ID              objectid.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	DeviceID        uuuid.UUID        `bson:"device_id,omitempty" json:"device_id,omitempty"`
	RsCustomerID    uuuid.UUID        `bson:"rs_customer_id,omitempty" json:"rs_customer_id,omitempty"`
	InstallDate     int64             `bson:"install_date,omitempty" json:"install_date,omitempty"`
	MaintenanceDate int64             `bson:"maintenance_date,omitempty" json:"maintenance_date,omitempty"`
	Status          string            `bson:"status,omitempty" json:"status,omitempty"`
	NumReplacement  int64             `bson:"num_replacement,omitempty" json:"num_replacement,omitempty"`
	CostSaved       float64           `bson:"cost_saved,omitempty" json:"cost_saved,omitempty"`
	Version         int64             `bson:"version,omitempty" json:"version,omitempty"`
}

type marshalDevice struct {
	ID              objectid.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	DeviceID        string            `bson:"device_id,omitempty" json:"device_id,omitempty"`
	RsCustomerID    string            `bson:"rs_customer_id,omitempty" json:"rs_customer_id,omitempty"`
	InstallDate     int64             `bson:"install_date,omitempty" json:"install_date,omitempty"`
	MaintenanceDate int64             `bson:"maintenance_date,omitempty" json:"maintenance_date,omitempty"`
	Status          string            `bson:"status,omitempty" json:"status,omitempty"`
	CostSaved       float64           `bson:"cost_saved,omitempty" json:"cost_saved,omitempty"`
	Version         int64             `bson:"version,omitempty" json:"version,omitempty"`
}

func (d Device) MarshalBSON() ([]byte, error) {
	md := &marshalDevice{
		ID:              d.ID,
		InstallDate:     d.InstallDate,
		MaintenanceDate: d.MaintenanceDate,
		Version:         d.Version,
		Status:          d.Status,
		CostSaved:       d.CostSaved,
	}

	if d.DeviceID.String() != (uuuid.UUID{}).String() {
		md.DeviceID = d.DeviceID.String()
	}

	if d.RsCustomerID.String() != (uuuid.UUID{}).String() {
		md.RsCustomerID = d.RsCustomerID.String()
	}

	return bson.Marshal(md)
}

func (d Device) MarshalJSON() ([]byte, error) {
	md := &marshalDevice{
		ID:              d.ID,
		InstallDate:     d.InstallDate,
		MaintenanceDate: d.MaintenanceDate,
		Version:         d.Version,
		Status:          d.Status,
		CostSaved:       d.CostSaved,
	}

	if d.DeviceID.String() != (uuuid.UUID{}).String() {
		md.DeviceID = d.DeviceID.String()
	}

	if d.RsCustomerID.String() != (uuuid.UUID{}).String() {
		md.RsCustomerID = d.RsCustomerID.String()
	}

	return json.Marshal(md)
}

func (d *Device) UnmarshalBSON(in []byte) error {
	var ok bool

	m := make(map[string]interface{})
	err := bson.Unmarshal(in, m)
	if err != nil {
		err = errors.Wrap(err, "Unmarshal Error")
		return err
	}

	if m["_id"] != nil {
		d.ID = m["_id"].(objectid.ObjectID)
	}

	if m["device_id"] != nil {
		d.DeviceID, err = uuuid.FromString(m["device_id"].(string))
		if err != nil {
			err = errors.Wrap(err, "Error parsing DeviceID for inventory")
			return err
		}
	}

	if m["rs_customer_id"] != nil {
		d.RsCustomerID, err = uuuid.FromString(m["rs_customer_id"].(string))
		if err != nil {
			err = errors.Wrap(err, "Error parsing DeviceID for inventory")
			return err
		}
	}

	if m["install_date"] != nil {
		installDateType := reflect.TypeOf(m["install_date"]).Kind()
		d.InstallDate, ok = m["install_date"].(int64)
		if !ok {
			if installDateType == reflect.Float64 {
				d.InstallDate = int64(m["install_date"].(float64))
			} else {
				val, _ := strconv.Atoi((m["install_date"]).(string))
				d.InstallDate = int64(val)
			}
		}
	}

	if m["maintenance_date"] != nil {
		maintenanceDateType := reflect.TypeOf(m["maintenance_date"]).Kind()
		d.MaintenanceDate, ok = m["maintenance_date"].(int64)
		if !ok {
			if maintenanceDateType == reflect.Float64 {
				d.MaintenanceDate = int64(m["maintenance_date"].(float64))
			} else {
				val, _ := strconv.Atoi((m["maintenance_date"]).(string))
				d.MaintenanceDate = int64(val)
			}
		}
	}

	if m["status"] != nil {
		d.Status = m["status"].(string)
	}

	if m["version"] != nil {
		versionType := reflect.TypeOf(m["version"]).Kind()
		d.Version, ok = m["version"].(int64)
		if !ok {
			if versionType == reflect.Float64 {
				d.Version = int64(m["version"].(float64))
			} else {
				val, _ := strconv.Atoi((m["version"]).(string))
				d.Version = int64(val)
			}
		}
	}

	if m["cost_saved"] != nil {
		costSavedType := reflect.TypeOf(m["cost_saved"]).Kind()
		d.CostSaved, ok = m["cost_saved"].(float64)
		if !ok {
			if costSavedType != reflect.Float64 {
				val, _ := strconv.Atoi((m["cost_saved"]).(string))
				d.CostSaved = float64(val)
			}
		}
	}

	return nil
}

func (d *Device) UnmarshalJSON(in []byte) error {
	var ok bool

	m := make(map[string]interface{})
	err := bson.Unmarshal(in, m)
	if err != nil {
		err = errors.Wrap(err, "Unmarshal Error")
		return err
	}

	if m["_id"] != nil {
		d.ID = m["_id"].(objectid.ObjectID)
	}

	if m["device_id"] != nil {
		d.DeviceID, err = uuuid.FromString(m["device_id"].(string))
		if err != nil {
			err = errors.Wrap(err, "Error parsing DeviceID for inventory")
			return err
		}
	}

	if m["rs_customer_id"] != nil {
		d.RsCustomerID, err = uuuid.FromString(m["rs_customer_id"].(string))
		if err != nil {
			err = errors.Wrap(err, "Error parsing DeviceID for inventory")
			return err
		}
	}

	if m["install_date"] != nil {
		installDateType := reflect.TypeOf(m["install_date"]).Kind()
		d.InstallDate, ok = m["install_date"].(int64)
		if !ok {
			if installDateType == reflect.Float64 {
				d.InstallDate = int64(m["install_date"].(float64))
			} else {
				val, _ := strconv.Atoi((m["install_date"]).(string))
				d.InstallDate = int64(val)
			}
		}
	}

	if m["maintenance_date"] != nil {
		maintenanceDateType := reflect.TypeOf(m["maintenance_date"]).Kind()
		d.MaintenanceDate, ok = m["maintenance_date"].(int64)
		if !ok {
			if maintenanceDateType == reflect.Float64 {
				d.MaintenanceDate = int64(m["maintenance_date"].(float64))
			} else {
				val, _ := strconv.Atoi((m["maintenance_date"]).(string))
				d.MaintenanceDate = int64(val)
			}
		}
	}

	if m["status"] != nil {
		d.Status = m["status"].(string)
	}

	if m["version"] != nil {
		versionType := reflect.TypeOf(m["version"]).Kind()
		d.Version, ok = m["version"].(int64)
		if !ok {
			if versionType == reflect.Float64 {
				d.Version = int64(m["version"].(float64))
			} else {
				val, _ := strconv.Atoi((m["version"]).(string))
				d.Version = int64(val)
			}
		}
	}

	if m["cost_saved"] != nil {
		costSavedType := reflect.TypeOf(m["cost_saved"]).Kind()
		d.CostSaved, ok = m["cost_saved"].(float64)
		if !ok {
			if costSavedType != reflect.Float64 {
				val, _ := strconv.Atoi((m["cost_saved"]).(string))
				d.CostSaved = float64(val)
			}
		}
	}

	return nil
}
