package database

import (
	"context"
	"errors"
	"log"
)

func (d *PAT_db) Add_PAT(
	ctx context.Context,
	owner_id uint64,
	hash string,
	pat_name string,
) (uint64, error) {
	pat_id, err := d.db.Add_PAT(ctx, owner_id, hash, pat_name)
	if !errors.Is(err, Safe_Error) && err != nil {
		log.Printf("Error while adding PAT: %v\n", err.Error())
		return 0, err
	}

	return pat_id, err
}

func (d *PAT_db) Validate_PAT(
	ctx context.Context,
	owner_id uint64,
	hash string,
) (uint64, error) {
	pat_id, err := d.db.Validate_PAT(ctx, owner_id, hash)
	if !errors.Is(err, Error_Invalid) && err != nil {
		log.Printf("Error while adding PAT: %v\n", err.Error())
		return 0, err
	}

	return pat_id, err
}

func (d *PAT_db) Get_Id(
	ctx context.Context,
	owner_id uint64,
	pat_name string,
) (uint64, error) {
	pat_id, err := d.db.Get_Id(ctx, owner_id, pat_name)
	if !errors.Is(err, Error_Invalid) && err != nil {
		log.Printf("Error while getting pat id: %v\n", err.Error())
	}

	return pat_id, err
}

func (d *PAT_db) Delete_PAT(
	ctx context.Context,
	pat_id uint64,
) error {
	err := d.db.Delete_PAT(ctx, pat_id)
	if !errors.Is(err, Error_Invalid) && err != nil {
		log.Printf("Error while deleteing pat: %v\n", err.Error())
	}

	return err
}

func (d *PAT_db) Get_PATs(
	ctx context.Context,
	user_id uint64,
) ([]string, error) {
	pats, err := d.db.Get_PATs(ctx, user_id)
	if !errors.Is(err, Error_Invalid) && err != nil {
		log.Printf("Error while deleteing pat: %v\n", err.Error())
	}

	return pats, err
}

func (d *PAT_db) Info(
	ctx context.Context,
	pat_id uint64,
) (*PAT_Info, error) {
	info, err := d.db.Info(ctx, pat_id)
	if !errors.Is(err, Error_Invalid) && err != nil {
		log.Printf("Error while getting pat info: %v\n", err.Error())
	}

	return info, err
}

func (d *PAT_db) Update_Use(
	ctx context.Context,
	pat_id uint64,
) error {
	err := d.db.Update_Use(ctx, pat_id)
	if !errors.Is(err, Error_Invalid) && err != nil {
		log.Printf("Error while updateing pat use: %v\n", err.Error())
	}

	return err
}
