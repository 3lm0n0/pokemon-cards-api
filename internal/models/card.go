package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Card struct {
	gorm.Model // when constructing a query, gorm will use snake_case
	ID 			 uuid.UUID 			`json:"ID" 										gorm:"type:uuid;default:uuid_generate_v4();primary_key;unique"`
	Name 		 string				`json:"Name" validate:"required,min=1,max=32" 	gorm:"type:varchar(255);not null"`
	Hp		 	 int 				`json:"Hp" 										gorm:"type:int;not null"`
	FirstEdition bool 				`json:"FirstEdition" 							gorm:"default:false;not null""`
	Expansion	 string 			`json:"Expansion" 								gorm:"type:varchar(55);not null"`
	Kind 		 string 			`json:"Kind" 									gorm:"type:varchar(55);not null"`
	Rarity 		 string 			`json:"Rarity" 									gorm:"type:varchar(55);uniqueIndex;not null" validate:"required"`
	Price 		 decimal.Decimal 	`json:"Price" 									gorm:"type:numeric;not null" validate:"required"`
	Strangeness	 string 			`json:"Strangeness" 							gorm:"type:varchar(55);not null"`
	Image 		 string 			`json:"Image" 									gorm:"type:varchar(255);not null"`
	CreatedAt 	 time.Time 			`json:"CreatedAt" 								gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt 	 time.Time 			`json:"UpdatedAt" 								gorm:"default:NULL"`
	DeletedAt 	 time.Time 			`json:"DeletedAt"								gorm:"default:NULL"`
}
