package storage

import (
	"context"
	ct "microservice/genproto/catalog_service"
)

type StorageI interface {
	CloseDB()
	Category() CategoryRepoI
}

type CategoryRepoI interface {
	Create(ctx context.Context, req *ct.CreateCategory) (resp *ct.GetCategory, err error)
	GetByID(ctx context.Context, req *ct.CategoryPrimaryKey) (resp *ct.GetCategory, err error)
	GetAll(ctx context.Context, req *ct.GetListCategoryRequest) (resp *ct.GetListCategoryResponse, err error)
	Update(ctx context.Context, req *ct.UpdateCategory) (resp *ct.GetCategory, err error)
	Delete(ctx context.Context, req *ct.CategoryPrimaryKey) (resp *ct.CategoryPrimaryKey, err error)
}
