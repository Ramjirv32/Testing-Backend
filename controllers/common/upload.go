package common

import (
"context"
"fmt"
"io"
"net/url"
"path/filepath"
"time"

"cloud.google.com/go/storage"
"github.com/gofiber/fiber/v3"

"backend/config"
"backend/utils"
)

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

// Make the object publicly readable
if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
fmt.Printf("Warning: Failed to set public ACL for %s: %v\n", filename, err)
}

publicURL := fmt.Sprintf(
"https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media",
"ticpin-fa6d2.firebasestorage.app",
url.PathEscape(filename),
)

return utils.SuccessResponse(c, 200, "File uploaded successfully", fiber.Map{
"url":  publicURL,
"path": filename,
})
}
