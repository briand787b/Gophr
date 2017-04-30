package main

import (
	"time"
	"net/http"
	"mime"
	"path/filepath"
	"os"
	"io"
	"mime/multipart"
	"github.com/disintegration/imaging"
	"image"
	"github.com/briand787b/validation"
)

type Image struct {
	ID 		string
	UserID 		string
	Name 		string
	Location 	string
	Size 		int64
	CreatedAt 	time.Time
	Description 	string
}

type ImageStore interface {
	Save(image *Image) error
	Find(id string) (*Image, error)
	FindAll(offset int) ([]Image, error)
	FindAllByUser(user *User, offset int) ([]Image, error)
}

const imageIDLength = 10

var thumbnailWidth = 400
var widthPreview = 800

// A map of accepted mime types and their file extension
var mimeExtensions = map[string]string {
	"image/png": 	".png",
	"image/jpeg":	".jpg",
	"image/gif": 	".gif",
}

func NewImage(user *User) *Image {
	return &Image{
		ID:		GenerateID("img", imageIDLength),
		UserID:		user.ID,
		CreatedAt: 	time.Now(),
	}
}

func (image *Image) CreateFromURL(imageURL string) error {
	// Get the response from the URL
	response, err := http.Get(imageURL)
	if err != nil {
		return err
	}

	// Make sure we have a valid response
	if response.StatusCode != http.StatusOK {
		return validation.ErrImageURLInvalid
	}

	defer response.Body.Close()

	// Ascertain the type of the file we downloaded
	mimeType, _, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if err != nil {
		return validation.ErrInvalidImageType
	}

	// Get an extension for the file
	ext, valid := mimeExtensions[mimeType]
	if !valid {
		return validation.ErrInvalidImageType
	}

	// Get a name from the URL
	image.Name = filepath.Base(imageURL)
	image.Location = image.ID + ext

	// Open a file at target location
	savedFile, err := os.Create("./data/images/" + image.Location)
	if err != nil {
		return err
	}
	defer savedFile.Close()

	// Copy the entire response to teh output file
	size, err := io.Copy(savedFile, response.Body)
	if err != nil {
		return err
	}
	// The returned value from io.Copy is the number of bytes copied
	image.Size = size

	// Create the various resizes of the image
	err = image.CreateResizedImages()
	if err != nil {
		return err
	}

	// Save our image to the store
	return globalImageStore.Save(image)
}

func (image *Image) CreateFromFile(file multipart.File, headers *multipart.FileHeader) error {
	// Move our file to an appropriate place, with an appropriate name
	image.Name = headers.Filename
	image.Location = image.ID + filepath.Ext(image.Name)

	// Open a file at the target location
	savedFile, err := os.Create("./data/images/" + image.Location)
	if err != nil {
		return err
	}

	defer savedFile.Close()

	// Copy the uploaded file to the target location
	size, err := io.Copy(savedFile, file)
	if err != nil {
		return err
	}
	image.Size = size

	// Create the various resizes of the image
	err = image.CreateResizedImages()
	if err != nil {
		return err
	}

	// Save the image to the database
	return globalImageStore.Save(image)
}

func (image *Image) CreateResizedImages() error {
	// Generate an image from a file
	srcImage, err := imaging.Open("./data/images/" + image.Location)
	if err != nil {
		return err
	}

	// Create a channel to send errors on
	errChan := make(chan error)

	// Process each size
	go image.resizePreview(errChan, srcImage)
	go image.resizeThumbnail(errChan, srcImage)

	// Wait for images to finish resizing
	for i := 0; i < 2; i++ {
		err = <- errChan
		if err != nil {
			return err
		}
	}

	return nil
}

func (image *Image) resizeThumbnail(errChan chan error, srcImage image.Image) {
	dstImage := imaging.Thumbnail(srcImage, thumbnailWidth, thumbnailWidth, imaging.Lanczos)
	destination := "./data/images/thumbnail/" + image.Location
	errChan <- imaging.Save(dstImage, destination)
}

func (image *Image) resizePreview(errChan chan error, srcImage image.Image) {
	size := srcImage.Bounds().Size()
	ratio := float64(size.Y) / float64(size.X)
	targetHeight := int(float64(widthPreview) * ratio)

	dstImage := imaging.Resize(srcImage, widthPreview, targetHeight, imaging.Lanczos)
	destination := "./data/images/preview/" + image.Location

	errChan <- imaging.Save(dstImage, destination)
}

func (image *Image) StaticRoute() string {
	return "/im/" + image.Location
}

func (image *Image) ShowRoute() string {
	return "/image/" + image.ID
}

func (image *Image) StaticThumbnailRoute() string {
	return "/im/thumbnail/" + image.Location
}

func (image *Image) StaticPreviewRoute() string {
	return "/im/preview/" + image.Location
}
