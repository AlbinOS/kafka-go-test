package models

import "github.com/elodina/go-avro"

type Device struct {
	Type            *avro.GenericEnum
	Ip              string
	Password        string
	Scene           string
	ProjectDatabase string
}

func NewDevice() *Device {
	return &Device{
		Type: avro.NewGenericEnum([]string{"CYCLOP", "PHIDI"}),
	}
}

func (o *Device) Schema() avro.Schema {
	if _Device_schema_err != nil {
		panic(_Device_schema_err)
	}
	return _Device_schema
}

// Enum values for DEVICE_TYPE
const (
	DEVICE_TYPE_CYCLOP int32 = 0
	DEVICE_TYPE_PHIDI  int32 = 1
)

// Generated by codegen. Please do not modify.
var _Device_schema, _Device_schema_err = avro.ParseSchema(`{
    "type": "record",
    "namespace": "amoobi.monitoring",
    "name": "Device",
    "fields": [
        {
            "name": "type",
            "type": {
                "type": "enum",
                "name": "DEVICE_TYPE",
                "symbols": [
                    "CYCLOP",
                    "PHIDI"
                ]
            }
        },
        {
            "name": "ip",
            "type": "string"
        },
        {
            "name": "password",
            "type": "string"
        },
        {
            "name": "scene",
            "type": "string"
        },
        {
            "name": "projectDatabase",
            "type": "string"
        }
    ]
}`)
