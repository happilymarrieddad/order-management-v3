package repos

import (
	"context"
	"fmt"

	"github.com/happilymarrieddad/order-management-v3/api/types"
	"googlemaps.github.io/maps"
	"xorm.io/xorm"
)

// AddressFindOpts defines the options for finding addresses.
type AddressFindOpts struct {
	IDs    []int64
	Limit  int
	Offset int
}

//go:generate mockgen -source=./addresses.go -destination=./mocks/addresses.go -package=mock_repos GoogleAPIClient
type GoogleAPIClient interface {
	Geocode(ctx context.Context, r *maps.GeocodingRequest) ([]maps.GeocodingResult, error)
}

//go:generate mockgen -source=./addresses.go -destination=./mocks/addresses.go -package=mock_repos AddressesRepo
type AddressesRepo interface {
	Create(ctx context.Context, address *types.Address) (*types.Address, error)
	CreateTx(ctx context.Context, tx *xorm.Session, address *types.Address) (*types.Address, error)
	Get(ctx context.Context, id int64) (*types.Address, bool, error)
	Update(ctx context.Context, address *types.Address) error
	UpdateTx(ctx context.Context, tx *xorm.Session, address *types.Address) error
	Delete(ctx context.Context, id int64) error
	DeleteTx(ctx context.Context, tx *xorm.Session, id int64) error
	Find(ctx context.Context, opts *AddressFindOpts) ([]*types.Address, int64, error)
	GeocodeAddress(ctx context.Context, addr *types.Address) (string, error)
}

type addressesRepo struct {
	db      *xorm.Engine
	gclient GoogleAPIClient
}

func NewAddressesRepo(db *xorm.Engine, gclient GoogleAPIClient) AddressesRepo {
	return &addressesRepo{db: db, gclient: gclient}
}

func (r *addressesRepo) Create(ctx context.Context, address *types.Address) (*types.Address, error) {
	return wrapInSession(r.db, func(tx *xorm.Session) (*types.Address, error) {
		return r.CreateTx(ctx, tx, address)
	})
}

func (r *addressesRepo) CreateTx(ctx context.Context, tx *xorm.Session, address *types.Address) (*types.Address, error) {
	if err := types.Validate(address); err != nil {
		return nil, err
	}
	if _, err := tx.Context(ctx).Insert(address); err != nil {
		return nil, err
	}
	return address, nil
}

func (r *addressesRepo) Get(ctx context.Context, id int64) (*types.Address, bool, error) {
	address := new(types.Address)
	has, err := r.db.Context(ctx).ID(id).Get(address)
	return address, has, err
}

func (r *addressesRepo) Update(ctx context.Context, address *types.Address) error {
	_, err := wrapInSession(r.db, func(tx *xorm.Session) (*struct{}, error) {
		return nil, r.UpdateTx(ctx, tx, address)
	})
	return err
}

func (r *addressesRepo) UpdateTx(ctx context.Context, tx *xorm.Session, address *types.Address) error {
	if err := types.Validate(address); err != nil {
		return err
	}
	_, err := tx.Context(ctx).ID(address.ID).Update(address)
	return err
}

func (r *addressesRepo) Delete(ctx context.Context, id int64) error {
	_, err := wrapInSession(r.db, func(tx *xorm.Session) (*struct{}, error) {
		return nil, r.DeleteTx(ctx, tx, id)
	})
	return err
}

func (r *addressesRepo) DeleteTx(ctx context.Context, tx *xorm.Session, id int64) error {
	_, err := tx.Context(ctx).ID(id).Delete(&types.Address{})
	return err
}

func (r *addressesRepo) Find(ctx context.Context, opts *AddressFindOpts) ([]*types.Address, int64, error) {
	s := r.db.NewSession().Context(ctx)
	defer s.Close()
	applyAddressFindOpts(s, opts)
	var addresses []*types.Address
	count, err := s.FindAndCount(&addresses)
	return addresses, count, err
}

func (r *addressesRepo) GeocodeAddress(ctx context.Context, addr *types.Address) (string, error) {
	res, err := r.gclient.Geocode(ctx, &maps.GeocodingRequest{
		Address: fmt.Sprintf("%s, %s, %s %s", addr.Line1, addr.City, addr.State, addr.PostalCode),
	})
	if err != nil || len(res) == 0 {
		return "", err
	}
	return res[0].PlaceID, nil
}

func applyAddressFindOpts(s *xorm.Session, opts *AddressFindOpts) {
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
