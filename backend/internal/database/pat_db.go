package database

import (
	"context"
	"errors"
	"log"
)

func (d *PATDb) AddPAT(
	ctx context.Context,
	ownerId uint64,
	hash string,
	patName string,
) (uint64, error) {
	patId, err := d.db.AddPAT(ctx, ownerId, hash, patName)
	if !errors.Is(err, ErrSafe) && err != nil {
		log.Printf("Error while adding PAT: %v\n", err.Error())
		return 0, err
	}

	return patId, err
}

func (d *PATDb) ValidatePAT(
	ctx context.Context,
	ownerId uint64,
	hash string,
) (uint64, error) {
	patId, err := d.db.ValidatePAT(ctx, ownerId, hash)
	if !errors.Is(err, ErrInvalid) && err != nil {
		log.Printf("Error while adding PAT: %v\n", err.Error())
		return 0, err
	}

	return patId, err
}

func (d *PATDb) GetId(
	ctx context.Context,
	ownerId uint64,
	patName string,
) (uint64, error) {
	patId, err := d.db.GetId(ctx, ownerId, patName)
	if !errors.Is(err, ErrInvalid) && err != nil {
		log.Printf("Error while getting pat id: %v\n", err.Error())
	}

	return patId, err
}

func (d *PATDb) DeletePAT(
	ctx context.Context,
	patId uint64,
) error {
	err := d.db.DeletePAT(ctx, patId)
	if !errors.Is(err, ErrInvalid) && err != nil {
		log.Printf("Error while deleting pat: %v\n", err.Error())
	}

	return err
}

func (d *PATDb) GetPATs(
	ctx context.Context,
	userId uint64,
) ([]string, error) {
	pats, err := d.db.GetPATs(ctx, userId)
	if !errors.Is(err, ErrInvalid) && err != nil {
		log.Printf("Error while deleting pat: %v\n", err.Error())
	}

	return pats, err
}

func (d *PATDb) Info(
	ctx context.Context,
	patId uint64,
) (*PATInfo, error) {
	info, err := d.db.Info(ctx, patId)
	if !errors.Is(err, ErrInvalid) && err != nil {
		log.Printf("Error while getting pat info: %v\n", err.Error())
	}

	return info, err
}

func (d *PATDb) UpdateUse(
	ctx context.Context,
	patId uint64,
) error {
	err := d.db.UpdateUse(ctx, patId)
	if !errors.Is(err, ErrInvalid) && err != nil {
		log.Printf("Error while updating pat use: %v\n", err.Error())
	}

	return err
}
