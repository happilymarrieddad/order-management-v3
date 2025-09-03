package repos

import (
	"sync"

	"xorm.io/xorm"
)

//go:generate mockgen -source=./global_repo.go -destination=./mocks/global_repo.go -package=mock_repos GlobalRepo
type GlobalRepo interface {
	Addresses() AddressesRepo
	Users() UsersRepo
	Companies() CompaniesRepo
	Locations() LocationsRepo
	CommodityAttributes() CommodityAttributesRepo
	Commodities() CommoditiesRepo
	CompanyAttributes() CompanyAttributesRepo
}

func NewGlobalRepo(db *xorm.Engine, gclient GoogleAPIClient) GlobalRepo {
	return &globalRepo{
		db:      db,
		mutex:   &sync.RWMutex{},
		repos:   make(map[string]interface{}),
		gclient: gclient,
	}
}

type globalRepo struct {
	db      *xorm.Engine
	gclient GoogleAPIClient
	repos   map[string]interface{}
	mutex   *sync.RWMutex
}

func (gr *globalRepo) factory(key string, fn func(db *xorm.Engine, gclient GoogleAPIClient) interface{}) interface{} {
	// First, check for an existing repository with a read lock for better performance.
	gr.mutex.RLock()
	val, exists := gr.repos[key]
	gr.mutex.RUnlock()
	if exists {
		return val
	}

	// If it doesn't exist, acquire a write lock to create it.
	gr.mutex.Lock()
	defer gr.mutex.Unlock()

	// Double-check in case another goroutine created it while we were waiting for the lock.
	val, exists = gr.repos[key]
	if exists {
		return val
	}

	// Create and store the new repository.
	newRepo := fn(gr.db, gr.gclient)
	gr.repos[key] = newRepo

	return newRepo
}

func (gr *globalRepo) Addresses() AddressesRepo {
	return gr.factory("Addresses", func(
		db *xorm.Engine, gclient GoogleAPIClient) interface{} {
		return NewAddressesRepo(db, gclient)
	}).(AddressesRepo)
}

func (gr *globalRepo) Users() UsersRepo {
	return gr.factory("Users", func(db *xorm.Engine, _ GoogleAPIClient) interface{} { return NewUsersRepo(db) }).(UsersRepo)
}

func (gr *globalRepo) Companies() CompaniesRepo {
	return gr.factory("Companies", func(db *xorm.Engine, _ GoogleAPIClient) interface{} { return NewCompaniesRepo(db) }).(CompaniesRepo)
}

func (gr *globalRepo) Locations() LocationsRepo {
	return gr.factory("Locations", func(db *xorm.Engine, _ GoogleAPIClient) interface{} { return NewLocationsRepo(db) }).(LocationsRepo)
}

func (gr *globalRepo) CommodityAttributes() CommodityAttributesRepo {
	return gr.factory("CommodityAttributes", func(db *xorm.Engine, _ GoogleAPIClient) interface{} { return NewCommodityAttributesRepo(db) }).(CommodityAttributesRepo)
}

func (gr *globalRepo) Commodities() CommoditiesRepo {
	return gr.factory("Commodities", func(db *xorm.Engine, _ GoogleAPIClient) interface{} { return NewCommoditiesRepo(db) }).(CommoditiesRepo)
}

func (gr *globalRepo) CompanyAttributes() CompanyAttributesRepo {
	return gr.factory("CompanyAttributes", func(db *xorm.Engine, _ GoogleAPIClient) interface{} { return NewCompanyAttributesRepo(db) }).(CompanyAttributesRepo)
}
