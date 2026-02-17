package models

import "time"

type OrganizationDetails struct {
	Category string `json:"category" firestore:"category"`
	PAN      string `json:"pan" firestore:"pan"`
	PANName  string `json:"pan_name" firestore:"pan_name"`
	PANImage string `json:"pan_image" firestore:"pan_image"` // URL from Storage
}

type GSTDetails struct {
	HasGST bool   `json:"has_gst" firestore:"has_gst"`
	GSTIN  string `json:"gstin" firestore:"gstin"`
}

type BankDetails struct {
	AccountHolderName string `json:"account_holder_name" firestore:"account_holder_name"`
	AccountNumber     string `json:"account_number" firestore:"account_number"`
	IFSCCode          string `json:"ifsc_code" firestore:"ifsc_code"`
	BankName          string `json:"bank_name" firestore:"bank_name"`
	BranchName        string `json:"branch_name" firestore:"branch_name"`
}

type BackupContact struct {
	Name  string `json:"name" firestore:"name"`
	Email string `json:"email" firestore:"email"`
	Phone string `json:"phone" firestore:"phone"`
}

type PartnerProfile struct {
	ID                  string              `json:"id" firestore:"id"`
	UserID              string              `json:"user_id" firestore:"user_id"`
	OrganizationDetails OrganizationDetails `json:"organization_details" firestore:"organization_details"`
	GSTDetails          GSTDetails          `json:"gst_details" firestore:"gst_details"`
	BankDetails         BankDetails         `json:"bank_details" firestore:"bank_details"`
	BackupContact       BackupContact       `json:"backup_contact" firestore:"backup_contact"`
	Status              string              `json:"status" firestore:"status"` // pending, approved, rejected
	CreatedAt           time.Time           `json:"created_at" firestore:"created_at"`
	UpdatedAt           time.Time           `json:"updated_at" firestore:"updated_at"`
}
