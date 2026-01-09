package models

import (
	"encoding/json"
	"strconv"
	"time"
)

// Ensure Timestamp implements json.Unmarshaler
var _ json.Unmarshaler = (*Timestamp)(nil)
