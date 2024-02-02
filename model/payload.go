package model

import "encoding/json"

type Payload struct {
	Table        string                     `json:"table" binding:"required"`
	Records      map[string]json.RawMessage `json:"records" binding:"required"`
	Result       map[string]json.RawMessage `json:"-"`
	Database     string                     `json:"database"`
	Distribution string                     `json:"distribution"` //children distribution:even (default),random
	RecordCount  uint                       `json:"record_count" binding:"required"`
	Children     []Payload                  `json:"children"`
}
