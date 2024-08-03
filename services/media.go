package services

// Import Cloudinary and other necessary libraries
//===================
import (
	"context"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func credentials() (*cloudinary.Cloudinary, context.Context) {
	cld, _ := cloudinary.New()
	cld.Config.URL.Secure = true
	ctx := context.Background()
	return cld, ctx
}

func UploadFile(file interface{}, resourceType string) (string, error) {
	cld, ctx := credentials()

	uploadedFile, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		ResourceType: resourceType,
	})
	if err != nil {
		return "", err
	}

	return uploadedFile.SecureURL, nil
}
