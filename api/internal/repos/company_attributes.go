package repos

import (
	"context"

	"github.com/happilymarrieddad/order-management-v3/api/types"
	"xorm.io/xorm"
)

//go:generate mockgen -source=./company_attributes.go -destination=./mocks/company_attributes.go -package=mock_repos CompanyAttributesRepo
type CompanyAttributesRepo interface {
	Create(ctx context.Context, attr *types.CompanyAttribute) error
	CreateTx(ctx context.Context, tx *xorm.Session, attr *types.CompanyAttribute) error
	Find(ctx context.Context, opts *FindCompanyAttributesOpts) ([]*types.CompanyAttribute, int64, error)
	Delete(ctx context.Context, id int64) error
	DeleteTx(ctx context.Context, tx *xorm.Session, id int64) error
}

type FindCompanyAttributesOpts struct {
	Limit     int
	Offset    int
	CompanyID int64
	IDs       []int64
	Position  *int
}

func NewCompanyAttributesRepo(db *xorm.Engine) CompanyAttributesRepo {
	return &companyAttributesRepo{db: db}
}

type companyAttributesRepo struct {
	db *xorm.Engine
}

func (r *companyAttributesRepo) Create(ctx context.Context, attr *types.CompanyAttribute) error {
	_, err := wrapInSession(r.db, func(tx *xorm.Session) (*struct{}, error) {
		return nil, r.CreateTx(ctx, tx, attr)
	})
	return err
}

func (r *companyAttributesRepo) CreateTx(ctx context.Context, tx *xorm.Session, attr *types.CompanyAttribute) error {
	_, err := tx.Context(ctx).Insert(attr)
	return err
}

func (r *companyAttributesRepo) Find(ctx context.Context, opts *FindCompanyAttributesOpts) ([]*types.CompanyAttribute, int64, error) {
	sess := r.db.NewSession().Context(ctx)
	defer sess.Close()

	applyCompanyAttributeFindOpts(sess, opts)

	var attrs []*types.CompanyAttribute
	count, err := sess.FindAndCount(&attrs)
	return attrs, count, err
}

func (r *companyAttributesRepo) Delete(ctx context.Context, id int64) error {
	_, err := wrapInSession(r.db, func(tx *xorm.Session) (*struct{}, error) {
		return nil, r.DeleteTx(ctx, tx, id)
	})
	return err
}

func (r *companyAttributesRepo) DeleteTx(ctx context.Context, tx *xorm.Session, id int64) error {
	res, err := tx.Context(ctx).ID(id).Delete(&types.CompanyAttribute{})
	if err != nil {
		return err
	}
	if res == 0 {
		return types.NewNotFoundError("company attribute not found")
	}
	return nil
}

func applyCompanyAttributeFindOpts(s *xorm.Session, opts *FindCompanyAttributesOpts) {
	if opts == nil {
		return
	}
	if opts.Limit > 0 {
		s.Limit(opts.Limit, opts.Offset)
	}
	if opts.CompanyID > 0 {
		s.Where("company_id = ?", opts.CompanyID)
	}
	if len(opts.IDs) > 0 {
		s.In("id", opts.IDs)
	}
	if opts.Position != nil {
		s.Where("position = ?", *opts.Position)
	}
}
