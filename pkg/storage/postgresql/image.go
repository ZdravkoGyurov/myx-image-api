package postgresql

import (
	"context"
	"fmt"

	"github.com/ZdravkoGyurov/myx-image-api/pkg/errors"
	"github.com/ZdravkoGyurov/myx-image-api/pkg/types"
)

var (
	imageTable = "images"

	insertImageQuery = fmt.Sprintf(`INSERT INTO %s 
	(name, location, longitude, latitude)
	VALUES ($1, $2, $3, $4)
	RETURNING *`, imageTable)

	readImagesQuery = fmt.Sprintf(`SELECT * FROM %s WHERE isInsideBBox($1, $2, $3, $4, longitude, latitude)`, imageTable)

	deleteImageQuery = fmt.Sprintf(`DELETE FROM %s WHERE name=$1`, imageTable)
)

func (s *Storage) CreateImage(ctx context.Context, image types.Image) error {
	dbCtx, cancel := context.WithTimeout(ctx, s.cfg.RequestTimeout)
	defer cancel()

	if _, err := s.pool.Exec(dbCtx, insertImageQuery, image.Fields()...); err != nil {
		return dbError(errors.Newf("failed to insert image: %w", err))
	}

	return nil
}

func (s *Storage) GetImages(ctx context.Context, bbox types.BBox) ([]*types.Image, error) {
	dbCtx, cancel := context.WithTimeout(ctx, s.cfg.RequestTimeout)
	defer cancel()

	rows, err := s.pool.Query(dbCtx, readImagesQuery,
		bbox.SouthWestLongitude, bbox.SouthWestLatitude, bbox.NorthEastLongitude, bbox.NorthEastLatitude)
	if err != nil {
		return nil, dbError(errors.Newf("failed to get images: %w", err))
	}

	images := make([]*types.Image, 0)
	for rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, dbError(errors.Newf("failed to read image row: %w", err))
		}
		image, err := readImageRecord(rows)
		if err != nil {
			return nil, dbError(errors.Newf("failed to deserialize image row: %w", err))
		}
		images = append(images, image)
	}

	return images, nil
}

func (s *Storage) DeleteImage(ctx context.Context, name string) error {
	dbCtx, cancel := context.WithTimeout(ctx, s.cfg.RequestTimeout)
	defer cancel()

	if _, err := s.pool.Exec(dbCtx, deleteImageQuery, name); err != nil {
		return dbError(errors.Newf("failed to delete image with name %s: %w", name, err))
	}

	return nil
}

func readImageRecord(row dbRecord) (*types.Image, error) {
	var image types.Image

	err := row.Scan(
		&image.Name,
		&image.Location,
		&image.Longitude,
		&image.Latitude,
	)
	if err != nil {
		return nil, err
	}

	return &image, nil
}
