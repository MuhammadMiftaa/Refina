package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	ec "server/env/config"
	"server/helper"
	"server/internal/repository"
	"server/internal/types/dto"
	"server/internal/types/entity"
	"server/queue/config"

	"github.com/rabbitmq/amqp091-go"
)

type ReportsService interface {
	RequestReport(ctx context.Context, userID, fromDate, toDate string) error
	UpdateUserReport(ctx context.Context) error
}

type reportsService struct {
	reportsRepo repository.ReportsRepository
	usersRepo   repository.UsersRepository
}

func NewReportsService(reportsRepo repository.ReportsRepository, usersRepo repository.UsersRepository) ReportsService {
	return &reportsService{
		reportsRepo: reportsRepo,
		usersRepo:   usersRepo,
	}
}

func (s *reportsService) RequestReport(ctx context.Context, userID, fromDate, toDate string) error {
	// Check if user exists
	if _, err := s.usersRepo.GetUserByID(userID); err != nil {
		return err
	}

	// Check if there are any processed reports for the user
	userReports, err := s.reportsRepo.GetProcessedReportByUserID(ctx, nil, userID)
	if err != nil {
		return err
	}
	if len(userReports) > 0 {
		return fmt.Errorf("there are %d processed reports for user %s", len(userReports), userID)
	}

	// Create a new report
	report := entity.Reports{}
	report.UserID = userID
	fromDateTime, err := time.Parse(time.RFC3339, fromDate)
	if err != nil {
		return err
	}
	toDateTime, err := time.Parse(time.RFC3339, toDate)
	if err != nil {
		return err
	}
	report.FromDate = fromDateTime
	report.ToDate = toDateTime
	report.RequestAt = time.Now()
	report.NextRequestAt = time.Now().Add(24 * time.Hour)
	report.Status = helper.REPORT_STATUS_PROCESSING
	if report, err = s.reportsRepo.CreateReport(ctx, nil, report); err != nil {
		return err
	}

	// Publish a message to the RabbitMQ queue to process the report
	message := amqp091.Publishing{
		ContentType: "application/json",
		Body:        []byte(`{"id":"` + report.ID.String() + `", "user_id":"` + userID + `"}`),
	}
	prodRequestReports, err := config.GetChannel()
	if err != nil {
		return err
	}
	err = prodRequestReports.PublishWithContext(ctx, "user_report", "user_report", false, false, message)
	if err != nil {
		return err
	}

	return nil
}

func (s *reportsService) UpdateUserReport(ctx context.Context) error {
	// Consume messages from the RabbitMQ queue
	consRequestReports, err := config.GetChannel()
	if err != nil {
		return err
	}
	userReports, err := consRequestReports.ConsumeWithContext(ctx, "user_report", "consumer_user_report", true, false, false, false, nil)
	if err != nil {
		return err
	}

	// Process each user report
	for userReport := range userReports {
		// Unmarshal the message body to get the report ID and user ID
		var report struct {
			ID     string `json:"id"`
			UserID string `json:"user_id"`
		}
		if err := json.Unmarshal(userReport.Body, &report); err != nil {
			return fmt.Errorf("failed to unmarshal user report: %w", err)
		}
		if report.ID == "" {
			return fmt.Errorf("report ID is empty, skipping report processing")
		}

		// Get the user by ID
		user, err := s.usersRepo.GetUserByID(report.UserID)
		if err != nil {
			return fmt.Errorf("failed to get user by ID %s: %w", report.UserID, err)
		}

		// Check if the report exists and is in processing status
		existingReport, err := s.reportsRepo.GetReportByID(ctx, nil, report.ID)
		if err != nil {
			return fmt.Errorf("failed to get report by ID %s: %w", report.ID, err)
		}
		if existingReport.UserID != user.ID.String() {
			return fmt.Errorf("report user ID %s does not match user ID %s", existingReport.UserID, user.ID.String())
		}
		if existingReport.Status != helper.REPORT_STATUS_PROCESSING {
			return fmt.Errorf("report status is not processing, current status: %s", existingReport.Status)
		}

		// Simulate report generation and sending email
		var userReport dto.ReportResponse
		userReport.UserID = user.ID.String()
		userReport.UserName = user.Name
		userReport.FromDate = existingReport.FromDate
		userReport.ToDate = existingReport.ToDate
		userReport.GeneratedAt = time.Now()
		userReport.FileURL = fmt.Sprintf("https://example.com/reports/%s_report.pdf", user.ID.String())
		userReport.FileSize = 250000
		
		SMTPProvider := helper.NewZohoSMTP(ec.Cfg.ZSMTP)
		if err := helper.NewSMTPClient(SMTPProvider).SendSingleEmail(user.Email, "Report Generated", "financial-report-template.html", userReport); err != nil {
			// Handle error
			existingReport.Status = helper.REPORT_STATUS_FAILED
			if _, updateErr := s.reportsRepo.UpdateReport(ctx, nil, existingReport); updateErr != nil {
				return fmt.Errorf("failed to update report status to FAILED after email error: %w", updateErr)
			}
			return fmt.Errorf("failed to send email: %w", err)
		}

		// Update the report status and details
		existingReport.FileURL = &userReport.FileURL
		existingReport.FileSize = &userReport.FileSize
		existingReport.GeneratedAt = &userReport.GeneratedAt
		existingReport.Status = helper.REPORT_STATUS_COMPLETED
		if _, err := s.reportsRepo.UpdateReport(ctx, nil, existingReport); err != nil {
			return err
		}
	}

	return nil
}
