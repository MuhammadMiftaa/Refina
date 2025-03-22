package dto

type AttachmentsResponse struct {
	ID            string `json:"id"`
	TransactionID string `json:"transaction_id"`
	Image         string `json:"image"`
}

type AttachmentsRequest struct {
	TransactionID string `json:"transaction_id"`
	Image         string `json:"image"`
}
