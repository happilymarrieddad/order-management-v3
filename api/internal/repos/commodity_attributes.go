package repos

import (
	"context"

	"github.com/happilymarrieddad/order-management-v3/api/types"
	"xorm.io/xorm"
)

type CommodityAttributeFindOpts struct {
	IDs            []int64
	CommodityTypes []types.CommodityType // For IN query
	Limit          int
	Offset         int
}

//go:generate mockgen -source=./commodity_attributes.go -destination=./mocks/commodity_attributes.go -package=mock_repos CommodityAttributesRepo
type CommodityAttributesRepo interface {
	Create(ctx context.Context, ca *types.CommodityAttribute) error
	CreateTx(ctx context.Context, tx *xorm.Session, ca *types.CommodityAttribute) error
	Get(ctx context.Context, id int64) (*types.CommodityAttribute, bool, error)
	Update(ctx context.Context, ca *types.CommodityAttribute) error
	UpdateTx(ctx context.Context, tx *xorm.Session, ca *types.CommodityAttribute) error
	Find(ctx context.Context, opts *CommodityAttributeFindOpts) ([]*types.CommodityAttribute, int64, error)
}

type commodityAttributesRepo struct {
	db *xorm.Engine
}

func NewCommodityAttributesRepo(db *xorm.Engine) CommodityAttributesRepo {
	return &commodityAttributesRepo{db: db}
}

func (r *commodityAttributesRepo) Create(ctx context.Context, ca *types.CommodityAttribute) error {
	_, err := wrapInSession(r.db, func(tx *xorm.Session) (*struct{}, error) {
		return nil, r.CreateTx(ctx, tx, ca)
	})
	return err
}

func (r *commodityAttributesRepo) CreateTx(ctx context.Context, tx *xorm.Session, ca *types.CommodityAttribute) error {
	if err := types.Validate(ca); err != nil {
		return err
	}
	// CommodityAttribute does not have a 'Visible' field, so no soft delete logic here.
	_, err := tx.Context(ctx).Insert(ca)
	return err
}

func (r *commodityAttributesRepo) Get(ctx context.Context, id int64) (*types.CommodityAttribute, bool, error) {
	ca := new(types.CommodityAttribute)
	has, err := r.db.Context(ctx).ID(id).Get(ca)
	return ca, has, err
}

func (r *commodityAttributesRepo) Update(ctx context.Context, ca *types.CommodityAttribute) error {
	_, err := wrapInSession(r.db, func(tx *xorm.Session) (*struct{}, error) {
		return nil, r.UpdateTx(ctx, tx, ca)
	})
	return err
}

func (r *commodityAttributesRepo) UpdateTx(ctx context.Context, tx *xorm.Session, ca *types.CommodityAttribute) error {
	if err := types.Validate(ca); err != nil {
		return err
	}
	_, err := tx.Context(ctx).ID(ca.ID).Cols("name").Update(ca)
	return err
}

func (r *commodityAttributesRepo) Find(ctx context.Context, opts *CommodityAttributeFindOpts) ([]*types.CommodityAttribute, int64, error) {
	s := r.db.NewSession().Context(ctx)
	defer s.Close()

	applyCommodityAttributeFindOpts(s, opts)

	var commodityAttributes []*types.CommodityAttribute
	count, err := s.FindAndCount(&commodityAttributes)
	return commodityAttributes, count, err
}

func applyCommodityAttributeFindOpts(s *xorm.Session, opts *CommodityAttributeFindOpts) {
	if opts == nil {
		return
	}
	if len(opts.IDs) > 0 {
		s.In("id", opts.IDs)
	}

	if len(opts.CommodityTypes) > 0 {
		// Convert []types.CommodityType to []int for XORM In clause
		intCommodityTypes := make([]int, len(opts.CommodityTypes))
		for i, ct := range opts.CommodityTypes {
			intCommodityTypes[i] = int(ct)
		}
		s.In("commodity_type_name", intCommodityTypes)
	}

	if opts.Limit > 0 {
		s.Limit(opts.Limit, opts.Offset)
	}
}
