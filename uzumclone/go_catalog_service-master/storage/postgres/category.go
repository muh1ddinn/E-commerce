package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	ct "microservice/genproto/catalog_service"
	"microservice/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type categoryRepo struct {
	db *pgxpool.Pool
}

func NewCategoryRepo(db *pgxpool.Pool) storage.CategoryRepoI {
	return &categoryRepo{
		db: db,
	}
}

func (c *categoryRepo) Create(ctx context.Context, req *ct.CreateCategory) (resp *ct.GetCategory, err error) {

	resp = &ct.GetCategory{}

	id := uuid.NewString()

	if req.ParentId == "" {
		req.ParentId = id
	}

	_, err = c.db.Exec(ctx, `
		INSERT INTO category (
			id,
			slug,
			name_uz,
			name_ru,
			name_en,
			active,
			order_no,
			parent_id
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8
		) `, id, req.Slug, req.NameUz, req.NameRu, req.NameEn, req.Active, req.OrderNo, req.ParentId)

	if err != nil {
		log.Println("error while creating category")
		return nil, err
	}

	category, err := c.GetByID(ctx, &ct.CategoryPrimaryKey{Id: id})
	if err != nil {
		log.Println("error while getting category by id")
		return nil, err
	}

	return category, nil
}

func (c *categoryRepo) GetByID(ctx context.Context, req *ct.CategoryPrimaryKey) (resp *ct.GetCategory, err error) {

	resp = &ct.GetCategory{}

	var ParentId sql.NullString

	err = c.db.QueryRow(ctx, `
		SELECT
			id,
			slug,
			name_uz,
			name_ru,
			name_en,
			active,
			order_no,
			parent_id
		FROM category
		WHERE id = $1 AND deleted_at IS NULL
	`, req.Id).Scan(&resp.Id, &resp.Slug, &resp.NameUz, &resp.NameRu, &resp.NameEn, &resp.Active, &resp.OrderNo, &ParentId)

	if err != nil {
		log.Println("error while getting category by id")
		return nil, err
	}

	resp.ParentId = ParentId.String

	return resp, nil
}

func (c *categoryRepo) GetAll(ctx context.Context, req *ct.GetListCategoryRequest) (resp *ct.GetListCategoryResponse, err error) {

	resp = &ct.GetListCategoryResponse{}
	filter := ""

	offset := (req.Limit - 1) * req.Limit
	if offset < 0 {
		offset = 0
	}
	if req.Search != "" {
		filter += fmt.Sprintf(`AND slug  ILIKE  '%%%v%%' `, req.Search)

	}

	filter += fmt.Sprintf("OFFSET %v LIMIT %v", offset, req.Limit)
	fmt.Println("filter:", filter)
	rows, err := c.db.Query(context.Background(), `SELECT
	count(id) OVER(),
	id,
	slug,
	name_uz,
	name_ru,
	name_en,
	active,
	order_no,
	parent_id
	FROM category WHERE  deleted_at IS NULL `+filter+`
	`)
	if err != nil {
		return resp, err
	}
	for rows.Next() {
		var (
			category = &ct.GetCategory{}
		)
		if err := rows.Scan(
			&resp.Count,
			&category.Id,
			&category.Slug,
			&category.NameUz,
			&category.NameRu,
			&category.NameEn,
			&category.Active,
			&category.OrderNo,
			&category.ParentId,
		); err != nil {
			return resp, err
		}

		resp.Categories = append(resp.Categories, category)
	}

	return resp, nil
}

func (c *categoryRepo) Update(ctx context.Context, req *ct.UpdateCategory) (reso *ct.GetCategory, err error) {

	query := `UPDATE category SET
			slug = $2,
			name_uz = $3,
			name_ru = $4,
			name_en = $5,
			active = $6,
			order_no = $7,
			parent_id = $8,
			updated_at = NOW()
			WHERE id = $1`

	_, err = c.db.Exec(ctx, query,
		req.Id,
		req.Slug,
		req.NameUz,
		req.NameRu,
		req.NameEn,
		req.Active,
		req.OrderNo,
		req.ParentId)

	if err != nil {
		log.Println("error while updating category")
		return reso, err
	}
	category, err := c.GetByID(ctx, &ct.CategoryPrimaryKey{Id: req.Id})
	if err != nil {
		log.Println("error while getting category by id")
		return nil, err
	}

	return category, nil
}

func (c *categoryRepo) Delete(ctx context.Context, req *ct.CategoryPrimaryKey) (res *ct.CategoryPrimaryKey, err error) {

	query := `UPDATE category SET deleted_at=CURRENT_TIMESTAMP WHERE id=$1 AND deleted_at IS NULL`

	_, err = c.db.Exec(ctx, query, req.Id)

	if err != nil {
		log.Println("error while updating category")
		return res, err
	}
	return req, nil
}
