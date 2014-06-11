package storage

import (
	"code.google.com/p/graphics-go/graphics"
	"github.com/dchest/uniuri"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path"
)

const (
	UploadsDir    = "./public/uploads"
	ThumbnailsDir = "./public/uploads/thumbnails"
)

func generateRandomFilename(contentType string) string {
	filename := uniuri.New()
	if contentType == "image/png" {
		return filename + ".png"
	}
	return filename + ".jpg"
}

type ImageProcessor interface {
	Process(src multipart.File, contentType string) (string, error)
}

type LocalImageProcessor struct {
	BaseDir       string
	ThumbnailsDir string
}

func (processor LocalImageProcessor) Process(src multipart.File, contentType string) (string, error) {

	filename := generateRandomFilename(contentType)

	if err := os.MkdirAll(processor.BaseDir, 0777); err != nil && !os.IsExist(err) {
		return filename, err
	}

	if err := os.MkdirAll(processor.ThumbnailsDir, 0777); err != nil && !os.IsExist(err) {
		return filename, err
	}

	// make thumbnail
	var (
		img image.Image
		err error
	)

	if contentType == "image/png" {
		img, err = png.Decode(src)
	} else {
		img, err = jpeg.Decode(src)
	}

	if err != nil {
		return filename, err
	}

	thumb := image.NewRGBA(image.Rect(0, 0, 300, 300))
	graphics.Thumbnail(thumb, img)

	dst, err := os.Create(path.Join(processor.ThumbnailsDir, filename))

	if err != nil {
		return filename, err
	}

	defer dst.Close()

	if contentType == "image/png" {
		png.Encode(dst, thumb)
	} else if contentType == "image/jpeg" {
		jpeg.Encode(dst, thumb, nil)
	}

	src.Seek(0, 0)

	dst, err = os.Create(path.Join(processor.BaseDir, filename))

	if err != nil {
		return filename, err
	}

	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return filename, err
	}

	return filename, nil

}

func NewImageProcessor() ImageProcessor {
	return LocalImageProcessor{UploadsDir, ThumbnailsDir}
}