package imageuploader

import (
	"context"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func ImageUploadHelper(input interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cld, err := cloudinary.NewFromParams(os.Getenv("CLOUDINARY_CLOUD_NAME"), os.Getenv("CLOUDINARY_API_KEY"), os.Getenv("CLOUDINARY_API_SECRET"))
	if err != nil {
		return "", err
	}
	uploadParam, err := cld.Upload.Upload(ctx, input, uploader.UploadParams{
		Transformation: "c_fill,g_face,h_250,w_250",
		Folder:         os.Getenv("CLOUDINARY_UPLOAD_FOLDER")})
	if err != nil {
		return "", err
	}
	return uploadParam.SecureURL, nil
}
