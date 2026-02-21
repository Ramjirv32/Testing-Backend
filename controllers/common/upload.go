package common

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gofiber/fiber/v3"

	"backend/config"
	"backend/utils"
)

// secureFolders are uploaded without a public ACL; access requires a signed URL
// or an authenticated Firebase Storage request.
var secureFolders = map[string]bool{
	"pan_cards":       true,
	"kyc":             true,
	"bank_documents":  true,
	"gstin_documents": true,
}

func isSecureFolder(folder string) bool {
	return secureFolders[strings.ToLower(strings.TrimPrefix(folder, "/"))]
}

// UploadFile handles multipart file uploads to Firebase Storage.
// Accepts a "file" field and optional "folder" field (default: "others").
func UploadFile(c fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return utils.ErrorResponse(c, 400, "No file uploaded")
	}

	folder := c.FormValue("folder", "others")

	src, err := file.Open()
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to open file")
	}
	defer src.Close()

	ctx := context.Background()
	bucket, err := config.FirebaseStorage.DefaultBucket()
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to get storage bucket")
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s/%d_%s", folder, time.Now().Unix(), utils.GenerateUUIDv7()+ext)

	obj := bucket.Object(filename)
	wc := obj.NewWriter(ctx)
	if _, err = io.Copy(wc, src); err != nil {
		return utils.ErrorResponse(c, 500, fmt.Sprintf("Failed to write to storage: %v", err))
	}
	if err := wc.Close(); err != nil {
		return utils.ErrorResponse(c, 500, fmt.Sprintf("Failed to close storage writer: %v", err))
	}

	// Make the object publicly readable only for non-secure folders
	if !isSecureFolder(folder) {
		if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
			fmt.Printf("Warning: Failed to set public ACL for %s: %v\n", filename, err)
		}
	}

	publicURL := fmt.Sprintf(
		"https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media",
		"ticpin-fa6d2.firebasestorage.app",
		url.PathEscape(filename),
	)

	return utils.SuccessResponse(c, 200, "File uploaded successfully", fiber.Map{
		"url":    publicURL,
		"path":   filename,
		"secure": isSecureFolder(folder),
	})
}

// GetSignedURL generates a 15-minute signed download URL for a secure-folder file.
// Query param: path  (e.g. "pan_cards/1234_abc.jpg")
// This endpoint must be protected by middleware.Auth + middleware.AdminOnly.
func GetSignedURL(c fiber.Ctx) error {
	objectPath := c.Query("path")
	if objectPath == "" {
		return utils.ErrorResponse(c, 400, "query param 'path' is required")
	}

	// Only allow secure folders to be accessed via this endpoint
	topFolder := strings.SplitN(objectPath, "/", 2)[0]
	if !isSecureFolder(topFolder) {
		return utils.ErrorResponse(c, 403, "only secure-folder paths can be fetched via this endpoint")
	}

	// Load service-account credentials (client_email + private_key)
	cfg := config.LoadConfig()
	var credJSON string
	if cfg.FirebaseCredentials != "" {
		credJSON = cfg.FirebaseCredentials
	} else {
		keyPath := cfg.FirebaseKeyPath
		data, err := os.ReadFile(keyPath)
		if err != nil {
			renderPath := "/etc/secrets/" + keyPath
			data, err = os.ReadFile(renderPath)
			if err != nil {
				return utils.ErrorResponse(c, 500, "Failed to load service-account credentials")
			}
		}
		credJSON = string(data)
	}

	var sa struct {
		ClientEmail string `json:"client_email"`
		PrivateKey  string `json:"private_key"`
	}
	if err := json.Unmarshal([]byte(credJSON), &sa); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to parse service-account credentials")
	}

	opts := &storage.SignedURLOptions{
		GoogleAccessID: sa.ClientEmail,
		PrivateKey:     []byte(sa.PrivateKey),
		Method:         "GET",
		Expires:        time.Now().Add(15 * time.Minute),
	}

	const bucketName = "ticpin-fa6d2.firebasestorage.app"
	signedURL, err := storage.SignedURL(bucketName, objectPath, opts)
	if err != nil {
		return utils.ErrorResponse(c, 500, fmt.Sprintf("Failed to generate signed URL: %v", err))
	}

	return utils.SuccessResponse(c, 200, "Signed URL generated", fiber.Map{
		"url":        signedURL,
		"expires_in": "15 minutes",
	})
}
