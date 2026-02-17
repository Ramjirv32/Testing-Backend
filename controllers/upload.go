package controllers

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v3"

	"backend/config"
	"backend/utils"
)

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

	// For simple mock storage, we return a public-ish URL
	// In production, you'd use SignedUrl or make object public
	// For this development env, we'll return the object path and assume the frontend knows how to handle it
	// or provide a placeholder for now
	url := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media", "ticpin-website.firebasestorage.app", url.PathEscape(filename))

	return utils.SuccessResponse(c, 200, "File uploaded successfully", fiber.Map{
		"url":  url,
		"path": filename,
	})
}
