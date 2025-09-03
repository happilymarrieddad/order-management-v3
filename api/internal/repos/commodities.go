package repos

import (
	"context"
	"strings"

	"github.com/happilymarrieddad/order-management-v3/api/types"
	"xorm.io/xorm"
)

// FindCommoditiesOpts defines the options for finding commodities.
type FindCommoditiesOpts struct {
	IDs           []int64
	Names         []string
	CommodityType types.CommodityType
	Limit         int
	Offset        int
}

//go:generate mockgen -source=./commodities.go -destination=./mocks/commodities.go -package=mock_repos
type CommoditiesRepo interface {
	Create(ctx context.Context, commodity *types.Commodity) error
	CreateTx(ctx context.Context, tx *xorm.Session, commodity *types.Commodity) error
	Get(ctx context.Context, id int64) (*types.Commodity, bool, error)
	Update(ctx context.Context, commodity *types.Commodity) error
	UpdateTx(ctx context.Context, tx *xorm.Session, commodity *types.Commodity) error
	Find(ctx context.Context, opts *FindCommoditiesOpts) ([]*types.Commodity, int64, error)
}

type commoditiesRepo struct {
	db *xorm.Engine
}

func NewCommoditiesRepo(db *xorm.Engine) CommoditiesRepo {
	return &commoditiesRepo{db: db}
}

func (r *commoditiesRepo) Create(ctx context.Context, commodity *types.Commodity) error {
	_, err := wrapInSession(r.db, func(tx *xorm.Session) (*struct{}, error) {
		return nil, r.CreateTx(ctx, tx, commodity)
	})
	return err
}

func (r *commoditiesRepo) CreateTx(ctx context.Context, tx *xorm.Session, commodity *types.Commodity) error {
	if err := types.Validate(commodity); err != nil {
		return err
	}
	_, err := tx.Context(ctx).Insert(commodity)
	return err
}

func (r *commoditiesRepo) Get(ctx context.Context, id int64) (*types.Commodity, bool, error) {
	commodity := new(types.Commodity)
	has, err := r.db.Context(ctx).ID(id).Get(commodity)
	return commodity, has, err
}

func (r *commoditiesRepo) Update(ctx context.Context, commodity *types.Commodity) error {
	_, err := wrapInSession(r.db, func(tx *xorm.Session) (*struct{}, error) {
		return nil, r.UpdateTx(ctx, tx, commodity)
	})
	return err
}

func (r *commoditiesRepo) UpdateTx(ctx context.Context, tx *xorm.Session, commodity *types.Commodity) error {
	if err := types.Validate(commodity); err != nil {
		return err
	}
	_, err := tx.Context(ctx).ID(commodity.ID).Update(commodity)
	return err
}

func (r *commoditiesRepo) Find(ctx context.Context, opts *FindCommoditiesOpts) ([]*types.Commodity, int64, error) {
	s := r.db.NewSession().Context(ctx)
	defer s.Close()

	s.Where("visible = ?", true)

	applyFindCommoditiesOpts(s, opts)

	var commodities []*types.Commodity
	count, err := s.FindAndCount(&commodities)
	return commodities, count, err
}

func applyFindCommoditiesOpts(s *xorm.Session, opts *FindCommoditiesOpts) {
	if opts == nil {
		return
	}
	if len(opts.IDs) > 0 {
		s.In("id", opts.IDs)
	}
	if len(opts.Names) > 0 {
		orConditions := make([]string, len(opts.Names))
		orArgs := make([]interface{}, len(opts.Names))
		for i, name := range opts.Names {
			orConditions[i] = "LOWER(name) LIKE LOWER(?)"
			orArgs[i] = "%" + name + "%"
		}
		s.And("("+strings.Join(orConditions, " OR ")+")", orArgs...)
	}
	if opts.CommodityType != types.CommodityTypeUnknown {
		s.Where("commodity_type = ?", opts.CommodityType)
	}
	if opts.Limit > 0 {
		s.Limit(opts.Limit, opts.Offset)
	}
}
