package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/utils"
	"gorm.io/gorm"
)

type JSONStrings []string

func (j *JSONStrings) Scan(value interface{}) error {
	if value == nil {
		*j = []string{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("cannot scan type %T into JSONStrings", value)
	}

	var result []string
	if err := json.Unmarshal(bytes, &result); err != nil {
		*j = []string{}
		return nil
	}
	*j = result
	return nil
}

func (j JSONStrings) Value() (driver.Value, error) {
	if j == nil {
		return "[]", nil
	}
	bytes, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}
	return string(bytes), nil
}

type ExperienceClient struct {
	gorm.Model
	ExperienceID     uint        `json:"experience_id" gorm:"not null;index"`
	Name             string      `json:"name" gorm:"type:varchar(255);not null"`
	URL              *string     `json:"url,omitempty" gorm:"type:varchar(500)"`
	StartDate        utils.Date  `json:"start_date" gorm:"type:date;not null"`
	EndDate          *utils.Date `json:"end_date,omitempty"`
	Description      string      `json:"description" gorm:"type:text"`
	Achievements     JSONStrings `json:"achievements" gorm:"type:text;default:'[]'"`
	Responsibilities JSONStrings `json:"responsibilities" gorm:"type:text;default:'[]'"`
	Technologies     JSONStrings `json:"technologies" gorm:"type:text;default:'[]'"`
}
