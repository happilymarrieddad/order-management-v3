package repos

import (
	"context"
	"fmt"
	"strings"

	"github.com/happilymarrieddad/order-management-v3/api/types"
	"xorm.io/xorm"
)

// ProductsRepo defines the interface for product data operations.
//
//go:generate mockgen -source=./products.go -destination=./mocks/products.go -package=mock_repos ProductsRepo
type ProductsRepo interface {
	Get(ctx context.Context, id int64) (*types.Product, bool, error)
	Create(ctx context.Context, product *types.Product, attrs []*types.ProductAttributeValue) error
	CreateTx(ctx context.Context, tx *xorm.Session, product *types.Product, attrs []*types.ProductAttributeValue) error
	Update(ctx context.Context, product *types.Product, attrs []*types.ProductAttributeValue) error
	UpdateTx(ctx context.Context, tx *xorm.Session, product *types.Product, attrs []*types.ProductAttributeValue) error
	Delete(ctx context.Context, id int64) error
	DeleteTx(ctx context.Context, tx *xorm.Session, id int64) error
	Find(ctx context.Context, opts *ProductFindOpts) ([]*types.Product, int64, error) // Added
}

type productsRepo struct {
	db *xorm.Engine
}

// NewProductsRepo creates a new ProductsRepo.
func NewProductsRepo(db *xorm.Engine) ProductsRepo {
	return &productsRepo{db: db}
}

// ProductFindOpts provides options for finding products.
type ProductFindOpts struct {
	CompanyID int64
	IDs       []int64
	Name      string
	Limit     int
	Offset    int
}

func (r *productsRepo) Get(ctx context.Context, id int64) (*types.Product, bool, error) {
	product := new(types.Product)
	has, err := r.db.Context(ctx).ID(id).Get(product)
	return product, has, err
}

// Create inserts a new product into the database.
func (r *productsRepo) Create(ctx context.Context, product *types.Product, attrs []*types.ProductAttributeValue) error {
	_, err := wrapInSession(r.db, func(tx *xorm.Session) (*struct{}, error) {
		return nil, r.CreateTx(ctx, tx, product, attrs)
	})
	return err
}

func (r *productsRepo) CreateTx(ctx context.Context, tx *xorm.Session, product *types.Product, attrs []*types.ProductAttributeValue) error {
	if err := types.Validate(product); err != nil {
		return err
	}

	product.Visible = true

	// Insert the base product record first to get an ID
	if _, err := tx.Context(ctx).Insert(product); err != nil {
		return err
	}

	// Insert the attribute values
	for _, attr := range attrs {
		attr.ProductID = product.ID
		attr.CompanyID = product.CompanyID
		if _, err := tx.Context(ctx).Insert(attr); err != nil {
			return err
		}
	}

	// Derive and save the product name based on its attributes
	return r.deriveAndSaveNameTx(ctx, tx, product)
}

// Update updates an existing product and its attributes.
func (r *productsRepo) Update(ctx context.Context, product *types.Product, attrs []*types.ProductAttributeValue) error {
	_, err := wrapInSession(r.db, func(tx *xorm.Session) (*struct{}, error) {
		return nil, r.UpdateTx(ctx, tx, product, attrs)
	})
	return err
}

func (r *productsRepo) UpdateTx(ctx context.Context, tx *xorm.Session, product *types.Product, attrs []*types.ProductAttributeValue) error {
	if err := types.Validate(product); err != nil {
		return err
	}

	// 1. Update the base product record (excluding the name)
	if _, err := tx.Context(ctx).ID(product.ID).Cols("commodity_id", "company_id", "visible").Update(product); err != nil {
		return err
	}

	// If new attributes are provided, delete existing ones and insert the new ones.
	// Otherwise, keep the existing attributes.
	if len(attrs) > 0 {
		// 2. Delete existing attribute values
		if _, err := tx.Context(ctx).Where("product_id = ?", product.ID).Delete(&types.ProductAttributeValue{}); err != nil {
			return err
		}

		// 3. Insert the new attribute values
		for _, attr := range attrs {
			attr.ProductID = product.ID
			attr.CompanyID = product.CompanyID
			if _, err := tx.Context(ctx).Insert(attr); err != nil {
				return err
			}
		}
	}

	// 4. Derive and save the product name based on its new attributes
	return r.deriveAndSaveNameTx(ctx, tx, product)
}

// Delete performs a soft delete on a product.
func (r *productsRepo) Delete(ctx context.Context, id int64) error {
	_, err := wrapInSession(r.db, func(tx *xorm.Session) (*struct{}, error) {
		return nil, r.DeleteTx(ctx, tx, id)
	})
	return err
}

// DeleteTx performs a soft delete on a product by setting their visible flag to false.
func (r *productsRepo) DeleteTx(ctx context.Context, tx *xorm.Session, id int64) error {
	_, err := tx.Context(ctx).ID(id).Cols("visible").Update(&types.Product{Visible: false})
	return err
}

// Find retrieves a list of visible products with pagination and filtering, and a total count.
func (r *productsRepo) Find(ctx context.Context, opts *ProductFindOpts) ([]*types.Product, int64, error) {
	s := r.db.NewSession().Context(ctx)
	defer s.Close()
	s.Where("visible = ?", true)
	applyProductFindOpts(s, opts)
	var products []*types.Product
	count, err := s.FindAndCount(&products)
	return products, count, err
}

// applyProductFindOpts is a helper function to build the query based on find options.
func applyProductFindOpts(s *xorm.Session, opts *ProductFindOpts) {
	if opts == nil {
		return
	}

	if opts.CompanyID > 0 {
		s.And("company_id = ?", opts.CompanyID)
	}
	if len(opts.IDs) > 0 {
		s.In("id", opts.IDs)
	}
	if opts.Name != "" {
		s.And("LOWER(name) LIKE LOWER(?)", "%"+opts.Name+"%")
	}

	if opts.Limit > 0 {
		s.Limit(opts.Limit, opts.Offset)
	}
}

// deriveAndSaveNameTx constructs the product's name from its attributes and commodity, then updates the record.
func (r *productsRepo) deriveAndSaveNameTx(ctx context.Context, tx *xorm.Session, product *types.Product) error {
	// 1. Get the base commodity
	commodity := new(types.Commodity)
	has, err := tx.Context(ctx).ID(product.CommodityID).Get(commodity)
	if err != nil {
		return fmt.Errorf("failed to get commodity %d: %w", product.CommodityID, err)
	}
	if !has {
		return fmt.Errorf("commodity %d not found", product.CommodityID)
	}

	// 2. Get the company-specific attribute order settings
	var settings []types.CompanyAttributeSetting
	if err = tx.Context(ctx).Where("company_id = ?", product.CompanyID).OrderBy("display_order ASC").Find(&settings); err != nil {
		return fmt.Errorf("failed to get company attribute settings for company %d: %w", product.CompanyID, err)
	}

	// 3. Get all attribute values for the product
	var productAttrs []types.ProductAttributeValue
	if err = tx.Context(ctx).Where("product_id = ?", product.ID).Find(&productAttrs); err != nil {
		return fmt.Errorf("failed to get product attribute values for product %d: %w", product.ID, err)
	}

	valuesMap := make(map[int64]string)
	for _, pav := range productAttrs {
		valuesMap[pav.CommodityAttributeID] = pav.Value
	}

	var nameParts []string
	// If a custom order is defined, use it.
	if len(settings) > 0 {
		for _, setting := range settings {
			if value, ok := valuesMap[setting.CommodityAttributeID]; ok {
				nameParts = append(nameParts, value)
			}
		}
	} else {
		// Fallback: if no order is defined, sort by commodity attribute ID for consistent naming
		var fallbackAttrs []types.ProductAttributeValue
		if err = tx.Context(ctx).Table("product_attribute_values").
			Join("INNER", "commodity_attributes", "commodity_attributes.id = product_attribute_values.commodity_attribute_id").
			Where("product_attribute_values.product_id = ?", product.ID).
			OrderBy("commodity_attributes.id").
			Find(&fallbackAttrs); err != nil {
			return fmt.Errorf("failed to get fallback product attribute values: %w", err)
		}
		for _, attr := range fallbackAttrs {
			nameParts = append(nameParts, attr.Value)
		}
	}

	nameParts = append(nameParts, commodity.Name)

	product.Name = strings.Join(nameParts, " ")

	// 4. Update the product with the new name
	if _, err = tx.Context(ctx).ID(product.ID).Cols("name").Update(product); err != nil {
		return fmt.Errorf("failed to update product name for product %d: %w", product.ID, err)
	}

	return nil
}