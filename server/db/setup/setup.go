package setup

import (
	"fmt"
	"log"

	"server/helper"

	"gorm.io/gorm"
)

func CreateReportStatusEnum(db *gorm.DB) error {
	log.Println("[SETUP] Create report status enum.")

	var exists bool
	err := db.Raw(`
	SELECT EXISTS (
		SELECT 1 FROM pg_type WHERE typname = 'report_status'
	)
	`).Scan(&exists).Error
	if err != nil {
		return fmt.Errorf("error checking if report_status enum exists: %w", err)
	}

	if exists {
		fmt.Println("report_status enum already exists, skipping creation")
		return nil
	}

	sql := fmt.Sprintf(`CREATE TYPE report_status AS ENUM ('%s', '%s', '%s')`,
		helper.REPORT_STATUS_PROCESSING,
		helper.REPORT_STATUS_COMPLETED,
		helper.REPORT_STATUS_FAILED,
	)
	if err := db.Exec(sql).Error; err != nil {
		return err
	}
	return nil
}
