package repos

import (
	"context"

	"github.com/happilymarrieddad/order-management-v3/api/types"
	"xorm.io/xorm"
)

type CompanyFindOpts struct {
	IDs    []int64 `json:"-"`
	Limit  int     `json:"limit"`
	Offset int     `json:"offset"`
}

//go:generate mockgen -source=./companies.go -destination=./mocks/companies.go -package=mock_repos CompaniesRepo
type CompaniesRepo interface {
	Create(ctx context.Context, company *types.Company) error
	CreateTx(ctx context.Context, tx *xorm.Session, company *types.Company) error
	Get(ctx context.Context, id int64) (*types.Company, bool, error)
	Update(ctx context.Context, company *types.Company) error
	UpdateTx(ctx context.Context, tx *xorm.Session, company *types.Company) error
	Delete(ctx context.Context, id int64) error
	DeleteTx(ctx context.Context, tx *xorm.Session, id int64) error
	Find(ctx context.Context, opts *CompanyFindOpts) ([]*types.Company, int64, error)
	GetIncludeInvisible(ctx context.Context, id int64) (*types.Company, bool, error)
}

type companiesRepo struct {
	db *xorm.Engine
}

func NewCompaniesRepo(db *xorm.Engine) CompaniesRepo {
	return &companiesRepo{db: db}
}

func (r *companiesRepo) Create(ctx context.Context, company *types.Company) error {
	_, err := wrapInSession(r.db, func(tx *xorm.Session) (*struct{}, error) {
		return nil, r.CreateTx(ctx, tx, company)
	})
	return err
}

func (r *companiesRepo) CreateTx(ctx context.Context, tx *xorm.Session, company *types.Company) error {
	if err := types.Validate(company); err != nil {
		return err
	}
	company.Visible = true
	_, err := tx.Context(ctx).Insert(company)
	return err
}

func (r *companiesRepo) Get(ctx context.Context, id int64) (*types.Company, bool, error) {
	company := new(types.Company)
	has, err := r.db.Context(ctx).ID(id).Where("visible = ?", true).Get(company)
	return company, has, err
}

func (r *companiesRepo) GetIncludeInvisible(ctx context.Context, id int64) (*types.Company, bool, error) {
	company := new(types.Company)
	has, err := r.db.Context(ctx).ID(id).Get(company)
	return company, has, err
}

func (r *companiesRepo) Update(ctx context.Context, company *types.Company) error {
	_, err := wrapInSession(r.db, func(tx *xorm.Session) (*struct{}, error) {
		return nil, r.UpdateTx(ctx, tx, company)
	})
	return err
}

func (r *companiesRepo) UpdateTx(ctx context.Context, tx *xorm.Session, company *types.Company) error {
	if err := types.Validate(company); err != nil {
		return err
	}
	_, err := tx.Context(ctx).ID(company.ID).Update(company)
	return err
}

func (r *companiesRepo) Delete(ctx context.Context, id int64) error {
	_, err := wrapInSession(r.db, func(tx *xorm.Session) (*struct{}, error) {
		return nil, r.DeleteTx(ctx, tx, id)
	})
	return err
}

func (r *companiesRepo) DeleteTx(ctx context.Context, tx *xorm.Session, id int64) error {
	_, err := tx.Context(ctx).ID(id).Cols("visible").Update(&types.Company{Visible: false})
	return err
}

func (r *companiesRepo) Find(ctx context.Context, opts *CompanyFindOpts) ([]*types.Company, int64, error) {
	s := r.db.NewSession().Context(ctx)
	defer s.Close()
	s.Where("visible = ?", true)
	applyCompanyFindOpts(s, opts)
	var companies []*types.Company
	count, err := s.FindAndCount(&companies)
	return companies, count, err
}

func applyCompanyFindOpts(s *xorm.Session, opts *CompanyFindOpts) {
	if opts == nil {
		return
	}
	if len(opts.IDs) > 0 {
		s.In("id", opts.IDs)
	}
	if opts.Limit > 0 {
		s.Limit(opts.Limit, opts.Offset)
	}
}
