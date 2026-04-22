package dto

import (
	"hrsync-backend/internal/db"
)

type ReimburseResponse struct {
	db.InnerReimburse
	EmployeeName string  `json:"employeeName"`
	Department   string  `json:"department"`
}

type CreateReimburseRequest struct {
	Email        string  `json:"email"`
	Amount       int     `json:"amount"`
	Description  string  `json:"description"`
	AttachBill   string  `json:"attachBill"`
	PaymentProof string  `json:"paymentProof"`
	CreatorRole  *string `json:"creatorRole"`
}

type UpdateReimburseRequest struct {
	Amount        *int    `json:"amount"`
	Description   *string `json:"description"`
	AttachBill    *string `json:"attachBill"`
	PaymentProof  *string `json:"paymentProof"`
	Status        *string `json:"status"`
	UpdatedByRole *string `json:"updatedByRole"`
}
