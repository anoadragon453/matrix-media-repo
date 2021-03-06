package util_exif

import (
	"context"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/sirupsen/logrus"
	"github.com/turt2live/matrix-media-repo/storage"
	"github.com/turt2live/matrix-media-repo/types"
)

type ExifOrientation struct {
	RotateDegrees  int // should be 0, 90, 180, or 270
	FlipVertical   bool
	FlipHorizontal bool
}

func GetExifOrientation(media *types.Media) (*ExifOrientation, error) {
	if media.ContentType != "image/jpeg" && media.ContentType != "image/jpg" {
		return nil, errors.New("image is not a jpeg")
	}

	filePath, err := storage.ResolveMediaLocation(context.TODO(), &logrus.Entry{}, media.DatastoreId, media.Location)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	exifData, err := exif.Decode(file)
	if err != nil {
		return nil, err
	}

	rawValue, err := exifData.Get(exif.Orientation)
	if err != nil {
		return nil, err
	}

	orientation, err := rawValue.Int(0)
	if err != nil {
		return nil, err
	}

	if orientation < 1 || orientation > 8 {
		return nil, errors.New(fmt.Sprintf("orientation out of range: %d", orientation))
	}

	flipHorizontal := orientation < 5 && (orientation%2) == 0
	flipVertical := orientation > 4 && (orientation%2) != 0
	degrees := 0

	// TODO: There's probably a better way to represent this
	if orientation == 1 || orientation == 2 {
		degrees = 0
	} else if orientation == 3 || orientation == 4 {
		degrees = 180
	} else if orientation == 5 || orientation == 6 {
		degrees = 270
	} else if orientation == 7 || degrees == 8 {
		degrees = 90
	}

	return &ExifOrientation{degrees, flipVertical, flipHorizontal}, nil
}
