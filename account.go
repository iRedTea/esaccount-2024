package esaccount

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type ConfidentialType string
type ConfidentialSettings map[string]ConfidentialType
type Int64Slice []int64



const (
	All     ConfidentialType = "ALL"
	Friends ConfidentialType = "FRIENDS"
	Nobody  ConfidentialType = "NOBODY"
)

type Account struct {
	Id            int64                `json:"id" gorm:"unique;primaryKey"`
	Description   string               `json:"description"`
	DateOfBirth   string               `json:"date_of_birth"`
	Follows       Int64Slice           `json:"follows" gorm:"type:text"`
	Followers     Int64Slice            `json:"followers" gorm:"type:text"`
	Confidentials ConfidentialSettings `json:"confidentials" gorm:"type:text"`
}

// Implementing the driver.Valuer interface for Int64Slice
func (s Int64Slice) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Implementing the sql.Scanner interface for Int64Slice
func (s *Int64Slice) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), s); err != nil {
		return fmt.Errorf("failed to unmarshal Int64Slice: %w", err)
	}
	return nil
}

// Implementing the driver.Valuer interface for ConfidentialSettings
func (c ConfidentialSettings) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// Implementing the sql.Scanner interface for ConfidentialSettings
func (c *ConfidentialSettings) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), c); err != nil {
		return fmt.Errorf("failed to unmarshal ConfidentialSettings: %w", err)
	}
	return nil
}