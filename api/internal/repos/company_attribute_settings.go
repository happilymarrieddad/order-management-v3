package repos

import (
	"context"

	"github.com/happilymarrieddad/order-management-v3/api/types"
	"xorm.io/xorm"
)

// CompanyAttributeSettingFindOpts provides options for finding company attribute settings.
type CompanyAttributeSettingFindOpts struct {
	CompanyIDs []int64
}

//go:generate mockgen -source=./company_attribute_settings.go -destination=./mocks/company_attribute_settings.go -package=mock_repos CompanyAttributeSettingsRepo
type CompanyAttributeSettingsRepo interface {
	Get(ctx context.Context, id int64) (*types.CompanyAttributeSetting, bool, error)
	Find(ctx context.Context, opts *CompanyAttributeSettingFindOpts) ([]*types.CompanyAttributeSetting, int64, error)
	Create(ctx context.Context, setting *types.CompanyAttributeSetting) error
	CreateTx(ctx context.Context, tx *xorm.Session, setting *types.CompanyAttributeSetting) error
	Update(ctx context.Context, setting *types.CompanyAttributeSetting) error
	UpdateTx(ctx context.Context, tx *xorm.Session, setting *types.CompanyAttributeSetting) error
	Delete(ctx context.Context, ids ...int64) error
	DeleteTx(ctx context.Context, tx *xorm.Session, ids ...int64) error
}

type companyAttributeSettingsRepo struct {
	db *xorm.Engine
}

func NewCompanyAttributeSettingsRepo(db *xorm.Engine) CompanyAttributeSettingsRepo {
	return &companyAttributeSettingsRepo{db: db}
}

func (r *companyAttributeSettingsRepo) Get(ctx context.Context, id int64) (*types.CompanyAttributeSetting, bool, error) {
	setting := new(types.CompanyAttributeSetting)
	has, err := r.db.Context(ctx).ID(id).Get(setting)
	return setting, has, err
}

func (r *companyAttributeSettingsRepo) Find(ctx context.Context, opts *CompanyAttributeSettingFindOpts) ([]*types.CompanyAttributeSetting, int64, error) {
	s := r.db.NewSession().Context(ctx)
	defer s.Close()

	if opts != nil {
		if len(opts.CompanyIDs) > 0 {
			s.In("company_id", opts.CompanyIDs)
		}
	}

	var settings []*types.CompanyAttributeSetting
	count, err := s.FindAndCount(&settings)
	return settings, count, err
}

func (r *companyAttributeSettingsRepo) Create(ctx context.Context, setting *types.CompanyAttributeSetting) error {
	_, err := wrapInSession(r.db, func(tx *xorm.Session) (*struct{}, error) {
		return nil, r.CreateTx(ctx, tx, setting)
	})
	return err
}

func (r *companyAttributeSettingsRepo) CreateTx(ctx context.Context, tx *xorm.Session, setting *types.CompanyAttributeSetting) error {
	if err := types.Validate(setting); err != nil {
		return err
	}
	_, err := tx.Context(ctx).Insert(setting)
	return err
}

func (r *companyAttributeSettingsRepo) Update(ctx context.Context, setting *types.CompanyAttributeSetting) error {
	_, err := wrapInSession(r.db, func(tx *xorm.Session) (*struct{}, error) {
		return nil, r.UpdateTx(ctx, tx, setting)
	})
	return err
}

func (r *companyAttributeSettingsRepo) UpdateTx(ctx context.Context, tx *xorm.Session, setting *types.CompanyAttributeSetting) error {
	if err := types.Validate(setting); err != nil {
		return err
	}
	_, err := tx.Context(ctx).ID(setting.ID).Cols("display_order", "commodity_attribute_id").Update(setting)
	return err
}

func (r *companyAttributeSettingsRepo) Delete(ctx context.Context, ids ...int64) error {
	_, err := wrapInSession(r.db, func(tx *xorm.Session) (*struct{}, error) {
		return nil, r.DeleteTx(ctx, tx, ids...)
	})
	return err
}

func (r *companyAttributeSettingsRepo) DeleteTx(ctx context.Context, tx *xorm.Session, ids ...int64) error {
	if len(ids) == 0 {
		return nil
	}
	_, err := tx.Context(ctx).In("id", ids).Delete(&types.CompanyAttributeSetting{})
	return err
}