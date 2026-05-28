package repository

import (
	assetrepo "samrs-backend/internal/repository/asset"

	"gorm.io/gorm"
)

type AssetFilter = assetrepo.AssetFilter

type AssetRepository = assetrepo.AssetRepository

type AssetEventFilter = assetrepo.AssetEventFilter

type AssetEventRepository = assetrepo.AssetEventRepository

type AssetBrandFilter = assetrepo.AssetBrandFilter

type AssetBrandRepository = assetrepo.AssetBrandRepository

type AssetModelFilter = assetrepo.AssetModelFilter

type AssetModelRepository = assetrepo.AssetModelRepository

type AssetStatusFilter = assetrepo.AssetStatusFilter

type AssetStatusRepository = assetrepo.AssetStatusRepository

type AssetMutationFilter = assetrepo.AssetMutationFilter

type AssetMutationRepository = assetrepo.AssetMutationRepository

func NewAssetRepository(db *gorm.DB) AssetRepository {
	return assetrepo.NewAssetRepository(db)
}

func NewAssetEventRepository(db *gorm.DB) AssetEventRepository {
	return assetrepo.NewAssetEventRepository(db)
}

func NewAssetBrandRepository(db *gorm.DB) AssetBrandRepository {
	return assetrepo.NewAssetBrandRepository(db)
}

func NewAssetModelRepository(db *gorm.DB) AssetModelRepository {
	return assetrepo.NewAssetModelRepository(db)
}

func NewAssetStatusRepository(db *gorm.DB) AssetStatusRepository {
	return assetrepo.NewAssetStatusRepository(db)
}

func NewAssetMutationRepository(db *gorm.DB) AssetMutationRepository {
	return assetrepo.NewAssetMutationRepository(db)
}
