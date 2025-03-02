package prepopulate

import (
	"fmt"
	"log"
	"sync"

	"server/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func InvestmentTypesSeeder(db *gorm.DB) {
	// Data investasi dan satuannya
	investmentData := []struct {
		Name string
		Unit string
	}{
		{"Gold", "Gram / Troy Ounce"},
		{"Stocks", "Lembar"},
		{"Mutual Funds", "Unit Penyertaan (UP)"},
		{"Bonds", "Nominal / Lot"},
		{"Government Securities", "Nominal / Unit"},
		{"Deposits", "Nominal"},
		{"Others", "-"},
	}

	// Kosongkan tabel investment_types
	db.Exec("TRUNCATE TABLE investment_types RESTART IDENTITY CASCADE")

	var wg sync.WaitGroup

	for _, inv := range investmentData {
		wg.Add(1)
		go func(name, unit string) {
			defer wg.Done()
			investment := entity.InvestmentTypes{
				Base: entity.Base{ID: uuid.New()},
				Name: name,
				Unit: unit,
			}

			// Insert ke database
			if err := db.Create(&investment).Error; err != nil {
				log.Printf("Error inserting investment type %s: %v\n", name, err)
			} else {
				fmt.Printf("Inserted investment type: %s\n", name)
			}
		}(inv.Name, inv.Unit)
	}

	// Tunggu semua goroutine selesai
	wg.Wait()
}
