package repository

import (
	stockopnamerepo "samrs-backend/internal/repository/stockopname"

	"gorm.io/gorm"
)

type StockOpnameFilter = stockopnamerepo.StockOpnameFilter

type StockOpnameRepository = stockopnamerepo.StockOpnameRepository

type StockOpnameItemRepository = stockopnamerepo.StockOpnameItemRepository

func NewStockOpnameRepository(db *gorm.DB) StockOpnameRepository {
	return stockopnamerepo.NewStockOpnameRepository(db)
}

func NewStockOpnameItemRepository(db *gorm.DB) StockOpnameItemRepository {
	return stockopnamerepo.NewStockOpnameItemRepository(db)
}
