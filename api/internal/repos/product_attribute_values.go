package repos

import (
	"context"

	"github.com/happilymarrieddad/order-management-v3/api/types"
	"xorm.io/xorm"
)

//go:generate mockgen -source=./product_attribute_values.go -destination=./mocks/product_attribute_values.go -package=mock_repos ProductAttributeValuesRepo
type ProductAttributeValuesRepo interface {
	Get(ctx context.Context, id int64) (*types.ProductAttributeValue, bool, error)
	Find(ctx context.Context, opts *ProductAttributeValueFindOpts) ([]*types.ProductAttributeValue, int64, error)
	Create(ctx context.Context, attr *types.ProductAttributeValue) error
	CreateTx(ctx context.Context, tx *xorm.Session, attr *types.ProductAttributeValue) error
	Update(ctx context.Context, attr *types.ProductAttributeValue) error
	UpdateTx(ctx context.Context, tx *xorm.Session, attr *types.ProductAttributeValue) error
	Delete(ctx context.Context, ids ...int64) error
	DeleteTx(ctx context.Context, tx *xorm.Session, ids ...int64) error
}

type productAttributeValuesRepo struct {
	db *xorm.Engine
}

func NewProductAttributeValuesRepo(db *xorm.Engine) ProductAttributeValuesRepo {
	return &productAttributeValuesRepo{db: db}
}

// ProductAttributeValueFindOpts provides options for finding product attribute values.
type ProductAttributeValueFindOpts struct {
	CompanyIDs []int64
	ProductIDs []int64
	Limit      int
	Offset     int
}

func (r *productAttributeValuesRepo) Get(ctx context.Context, id int64) (*types.ProductAttributeValue, bool, error) {
	attr := new(types.ProductAttributeValue)
	has, err := r.db.Context(ctx).ID(id).Get(attr)
	return attr, has, err
}

func (r *productAttributeValuesRepo) Create(ctx context.Context, attr *types.ProductAttributeValue) error {
	_, err := wrapInSession(r.db, func(tx *xorm.Session) (*struct{}, error) {
		return nil, r.CreateTx(ctx, tx, attr)
	})
	return err
}

func (r *productAttributeValuesRepo) CreateTx(ctx context.Context, tx *xorm.Session, attr *types.ProductAttributeValue) error {
	if err := types.Validate(attr); err != nil {
		return err
	}
	_, err := tx.Context(ctx).Insert(attr)
	return err
}

func (r *productAttributeValuesRepo) Update(ctx context.Context, attr *types.ProductAttributeValue) error {
	_, err := wrapInSession(r.db, func(tx *xorm.Session) (*struct{}, error) {
		return nil, r.UpdateTx(ctx, tx, attr)
	})
	return err
}

func (r *productAttributeValuesRepo) UpdateTx(ctx context.Context, tx *xorm.Session, attr *types.ProductAttributeValue) error {
	if err := types.Validate(attr); err != nil {
		return err
	}
	_, err := tx.Context(ctx).ID(attr.ID).Cols("value").Update(attr)
	return err
}

func (r *productAttributeValuesRepo) Delete(ctx context.Context, ids ...int64) error {
	_, err := wrapInSession(r.db, func(tx *xorm.Session) (*struct{}, error) {
		return nil, r.DeleteTx(ctx, tx, ids...)
	})
	return err
}

func (r *productAttributeValuesRepo) DeleteTx(ctx context.Context, tx *xorm.Session, ids ...int64) error {
	if len(ids) == 0 {
		return nil
	}
	_, err := tx.Context(ctx).In("id", ids).Delete(&types.ProductAttributeValue{})
	return err
}

func (r *productAttributeValuesRepo) Find(ctx context.Context, opts *ProductAttributeValueFindOpts) ([]*types.ProductAttributeValue, int64, error) {
	s := r.db.NewSession().Context(ctx)
	defer s.Close()

	applyProductAttributeValueFindOpts(s, opts)

	var values []*types.ProductAttributeValue
	count, err := s.FindAndCount(&values)
	return values, count, err
}

func applyProductAttributeValueFindOpts(s *xorm.Session, opts *ProductAttributeValueFindOpts) {
	if opts == nil {
		return
	}

	if len(opts.CompanyIDs) > 0 {
		s.In("company_id", opts.CompanyIDs)
	}

	if len(opts.ProductIDs) > 0 {
		s.In("product_id", opts.ProductIDs)
	}

	if opts.Limit > 0 {
		s.Limit(opts.Limit, opts.Offset)
	}
}
