package models

import "time"

type GSTINInfo struct {
	GSTIN  string `json:"gstin" firestore:"gstin"`
	Status string `json:"status" firestore:"status"`
	State  string `json:"state" firestore:"state"`
}

type PANToGSTINMapping struct {
	ReferenceID    int         `json:"reference_id" firestore:"reference_id"`
	VerificationID string      `json:"verification_id" firestore:"verification_id"`
	Status         string      `json:"status" firestore:"status"`
	PAN            string      `json:"pan" firestore:"pan"`
	GSTINList      []GSTINInfo `json:"gstin_list" firestore:"gstin_list"`
}

type PANVerification struct {
	Status                   string `json:"status" firestore:"status"`
	Message                  string `json:"message" firestore:"message"`
	ReferenceID              int    `json:"reference_id" firestore:"reference_id"`
	VerificationID           string `json:"verification_id" firestore:"verification_id"`
	RegisteredName           string `json:"registered_name" firestore:"registered_name"`
	NamePanCard              string `json:"name_pan_card" firestore:"name_pan_card"`
	NameProvided             string `json:"name_provided" firestore:"name_provided"`
	NameMatch                string `json:"name_match" firestore:"name_match"`
	DOB                      string `json:"dob" firestore:"dob"`
	DOBMatch                 string `json:"dob_match" firestore:"dob_match"`
	PanStatus                string `json:"pan_status" firestore:"pan_status"`
	AadhaarSeedingStatus     string `json:"aadhaar_seeding_status" firestore:"aadhaar_seeding_status"`
	AadhaarSeedingStatusDesc string `json:"aadhaar_seeding_status_desc" firestore:"aadhaar_seeding_status_desc"`
	FirstName                string `json:"first_name" firestore:"first_name"`
	LastName                 string `json:"last_name" firestore:"last_name"`
	Email                    string `json:"email" firestore:"email"`
	MobileNumber             string `json:"mobile_number" firestore:"mobile_number"`
	Type                     string `json:"type" firestore:"type"`
	Gender                   string `json:"gender" firestore:"gender"`
	MaskedAadhaar            string `json:"masked_aadhaar_number" firestore:"masked_aadhaar_number"`
	AadhaarLinked            bool   `json:"aadhaar_linked" firestore:"aadhaar_linked"`
	Address                  struct {
		FullAddress string `json:"full_address" firestore:"full_address"`
		Street      string `json:"street" firestore:"street"`
		City        string `json:"city" firestore:"city"`
		State       string `json:"state" firestore:"state"`
		Pincode     int    `json:"pincode" firestore:"pincode"`
		Country     string `json:"country" firestore:"country"`
	} `json:"address" firestore:"address"`
	VerifiedAt time.Time `json:"verified_at" firestore:"verified_at"`
}

type OrganizationDetails struct {
	Category        string            `json:"category" firestore:"category"`
	PAN             string            `json:"pan" firestore:"pan"`
	PANName         string            `json:"pan_name" firestore:"pan_name"`
	PANImage        string            `json:"pan_image" firestore:"pan_image"`
	PANVerification PANVerification   `json:"pan_verification" firestore:"pan_verification"`
	GSTINMapping    PANToGSTINMapping `json:"gstin_mapping" firestore:"gstin_mapping"`
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
	OrganizerEmail      string              `json:"organizer_email" firestore:"organizer_email"`
	City                string              `json:"city" firestore:"city"`
	State               string              `json:"state" firestore:"state"`
	OrganizationDetails OrganizationDetails `json:"organization_details" firestore:"organization_details"`
	GSTDetails          GSTDetails          `json:"gst_details" firestore:"gst_details"`
	BankDetails         BankDetails         `json:"bank_details" firestore:"bank_details"`
	BackupContact       BackupContact       `json:"backup_contact" firestore:"backup_contact"`
	Status              string              `json:"status" firestore:"status"`
	CreatedAt           time.Time           `json:"created_at" firestore:"created_at"`
	UpdatedAt           time.Time           `json:"updated_at" firestore:"updated_at"`
}
