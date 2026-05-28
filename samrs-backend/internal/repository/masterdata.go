package repository

import (
	masterdatarepo "samrs-backend/internal/repository/masterdata"

	"gorm.io/gorm"
)

type RoomFilter = masterdatarepo.RoomFilter

type RoomRepository = masterdatarepo.RoomRepository

type BedFilter = masterdatarepo.BedFilter

type BedRepository = masterdatarepo.BedRepository

type CategoryFilter = masterdatarepo.CategoryFilter

type CategoryRepository = masterdatarepo.CategoryRepository

type VendorFilter = masterdatarepo.VendorFilter

type VendorRepository = masterdatarepo.VendorRepository

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return masterdatarepo.NewRoomRepository(db)
}

func NewBedRepository(db *gorm.DB) BedRepository {
	return masterdatarepo.NewBedRepository(db)
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return masterdatarepo.NewCategoryRepository(db)
}

func NewVendorRepository(db *gorm.DB) VendorRepository {
	return masterdatarepo.NewVendorRepository(db)
}
