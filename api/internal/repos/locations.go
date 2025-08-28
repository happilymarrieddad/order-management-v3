package repos

import (
	"context"
	"errors"

	"github.com/happilymarrieddad/order-management-v3/api/types"
	"xorm.io/xorm"
)

var (
	// ErrLocationNameExists is returned when a location with the same name already exists for a company.
	ErrLocationNameExists = errors.New("a location with this name already exists for the company")
)

// LocationFindOpts defines the options for finding locations.
type LocationFindOpts struct {
	IDs        []int64
	CompanyIDs []int64
	AddressIDs []int64
	Names      []string
	Limit      int
	Offset     int
}

//go:generate mockgen -source=./locations.go -destination=./mocks/locations.go -package=mock_repos LocationsRepo
type LocationsRepo interface {
	Create(ctx context.Context, location *types.Location) error
	CreateTx(ctx context.Context, tx *xorm.Session, location *types.Location) error
	Get(ctx context.Context, id int64) (*types.Location, bool, error)
	Update(ctx context.Context, location *types.Location) error
	UpdateTx(ctx context.Context, tx *xorm.Session, location *types.Location) error
	Delete(ctx context.Context, id int64) error
	DeleteTx(ctx context.Context, tx *xorm.Session, id int64) error
	Find(ctx context.Context, opts *LocationFindOpts) ([]*types.Location, int64, error)
	CountByCompanyID(ctx context.Context, companyID int64) (int64, error)
}

type locationsRepo struct {
	db *xorm.Engine
}

func NewLocationsRepo(db *xorm.Engine) LocationsRepo {
	return &locationsRepo{db: db}
}

func (r *locationsRepo) Create(ctx context.Context, location *types.Location) error {
	_, err := wrapInSession(r.db, func(tx *xorm.Session) (*struct{}, error) {
		return nil, r.CreateTx(ctx, tx, location)
	})
	return err
}

func (r *locationsRepo) CreateTx(ctx context.Context, tx *xorm.Session, location *types.Location) error {
	exists, err := tx.Context(ctx).Where("company_id = ? AND name = ?", location.CompanyID, location.Name).Exist(&types.Location{})
	if err != nil {
		return err
	}
	if exists {
		return ErrLocationNameExists
	}
	_, err = tx.Context(ctx).Insert(location)
	return err
}

func (r *locationsRepo) Get(ctx context.Context, id int64) (*types.Location, bool, error) {
	location := new(types.Location)
	has, err := r.db.Context(ctx).ID(id).Get(location)
	return location, has, err
}

func (r *locationsRepo) Update(ctx context.Context, location *types.Location) error {
	_, err := wrapInSession(r.db, func(tx *xorm.Session) (*struct{}, error) {
		return nil, r.UpdateTx(ctx, tx, location)
	})
	return err
}

func (r *locationsRepo) UpdateTx(ctx context.Context, tx *xorm.Session, location *types.Location) error {
	exists, err := tx.Context(ctx).
		Where("company_id = ? AND name = ? AND id != ?", location.CompanyID, location.Name, location.ID).
		Exist(&types.Location{})
	if err != nil {
		return err
	}
	if exists {
		return ErrLocationNameExists
	}
	_, err = tx.Context(ctx).ID(location.ID).AllCols().Update(location)
	return err
}

func (r *locationsRepo) Delete(ctx context.Context, id int64) error {
	_, err := wrapInSession(r.db, func(tx *xorm.Session) (*struct{}, error) {
		return nil, r.DeleteTx(ctx, tx, id)
	})
	return err
}

func (r *locationsRepo) DeleteTx(ctx context.Context, tx *xorm.Session, id int64) error {
	_, err := tx.Context(ctx).ID(id).Delete(new(types.Location))
	return err
}

func (r *locationsRepo) Find(ctx context.Context, opts *LocationFindOpts) ([]*types.Location, int64, error) {
	s := r.db.NewSession().Context(ctx)
	defer s.Close()
	applyLocationFindOpts(s, opts)
	var locations []*types.Location
	count, err := s.FindAndCount(&locations)
	return locations, count, err
}

func (r *locationsRepo) CountByCompanyID(ctx context.Context, companyID int64) (int64, error) {
	count, err := r.db.Context(ctx).Where("company_id = ?", companyID).Count(new(types.Location))
	return count, err
}

func applyLocationFindOpts(s *xorm.Session, opts *LocationFindOpts) {
	if opts == nil {
		return
	}
	if len(opts.IDs) > 0 {
		s.In("id", opts.IDs)
	}
	if len(opts.CompanyIDs) > 0 {
		s.In("company_id", opts.CompanyIDs)
	}
	if len(opts.AddressIDs) > 0 {
		s.In("address_id", opts.AddressIDs)
	}
	if len(opts.Names) > 0 {
		s.In("name", opts.Names)
	}
	if opts.Limit > 0 {
		s.Limit(opts.Limit, opts.Offset)
	}
}
