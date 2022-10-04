package controller

import (
	"bytes"
	"context"
	"io"
	"strconv"
	"strings"

	"github.com/ZdravkoGyurov/myx-image-api/pkg/errors"
	"github.com/ZdravkoGyurov/myx-image-api/pkg/types"
	"github.com/rwcarlsen/goexif/exif"
)

func (c *Controller) UploadImage(ctx context.Context, fileName string, file io.Reader) (types.Image, error) {
	var fileBuffer bytes.Buffer
	fileReader := io.TeeReader(file, &fileBuffer)

	imageExif, err := exif.Decode(fileReader)
	if err != nil {
		return types.Image{}, errors.Newf("failed to extract exif data: %w", err)
	}

	latitude, longitude, err := imageExif.LatLong()
	if err != nil {
		return types.Image{}, errors.Newf("failed to extract latitude and longitude: %w", err)
	}

	location, err := c.fileStorage.StoreFile(ctx, fileName, &fileBuffer)
	if err != nil {
		return types.Image{}, errors.Newf("failed to store file: %w", err)
	}

	image := types.Image{
		Name:      fileName,
		Location:  location,
		Longitude: longitude,
		Latitude:  latitude,
	}

	if err = c.storage.CreateImage(ctx, image); err != nil {
		return types.Image{}, errors.Newf("failed to store file metadata: %w", err)
	}

	return image, nil
}

func (c *Controller) GetImages(ctx context.Context, bboxS string) ([]*types.Image, error) {
	bbox, err := parseBBox(bboxS)
	if err != nil {
		return nil, errors.Newf("failed to parse bbox: %s: %w", err, errors.ErrInvalidEntity)
	}

	return c.storage.GetImages(ctx, bbox)
}

func (c *Controller) DeleteImage(ctx context.Context, fileName string) error {
	if err := c.storage.DeleteImage(ctx, fileName); err != nil {
		return errors.Newf("failed to delete file metadata: %w", err)
	}

	if err := c.fileStorage.DeleteFile(ctx, fileName); err != nil {
		return errors.Newf("failed to delete file: %w", err)
	}

	return nil
}

func parseBBox(bboxString string) (types.BBox, error) {
	bboxCoords := strings.Split(bboxString, ",")

	if len(bboxCoords) != 4 {
		return types.BBox{}, errors.New("failed to parse bbox, 4 parameters are required")
	}

	southWestLongitude, err := strconv.ParseFloat(bboxCoords[0], 64)
	if err != nil {
		return types.BBox{}, errors.Newf("failed to parse south west longitude: %w", err)
	}
	southWestLatitude, err := strconv.ParseFloat(bboxCoords[1], 64)
	if err != nil {
		return types.BBox{}, errors.Newf("failed to parse south west latitude: %w", err)
	}
	northEastLongitude, err := strconv.ParseFloat(bboxCoords[2], 64)
	if err != nil {
		return types.BBox{}, errors.Newf("failed to parse north east longitude: %w", err)
	}
	northEastLatitude, err := strconv.ParseFloat(bboxCoords[3], 64)
	if err != nil {
		return types.BBox{}, errors.Newf("failed to parse north east latitude: %w", err)
	}

	return types.BBox{
		SouthWestLongitude: southWestLongitude,
		SouthWestLatitude:  southWestLatitude,
		NorthEastLongitude: northEastLongitude,
		NorthEastLatitude:  northEastLatitude,
	}, nil
}
