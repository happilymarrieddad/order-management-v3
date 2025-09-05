package repos_test

import (
	"fmt"

	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ProductsRepo", func() {
	var (
		repo                repos.ProductsRepo
		casRepo             repos.CompanyAttributeSettingsRepo
		company1            *types.Company
		company2            *types.Company
		company3            *types.Company // The company with no specific order
		commodity           *types.Commodity
		commodityAttribute1 *types.CommodityAttribute // e.g., Variety
		commodityAttribute2 *types.CommodityAttribute // e.g., Color
		commodityAttribute3 *types.CommodityAttribute // e.g., Size
	)

	BeforeEach(func() {
		repo = gr.Products()
		casRepo = gr.CompanyAttributeSettings()

		// Dependency setup
		address, err := gr.Addresses().Create(ctx, &types.Address{
			Line1: "123 Main St", City: "Anytown", State: "CA", Country: "USA", PostalCode: "12345",
		})
		Expect(err).NotTo(HaveOccurred())

		company1 = &types.Company{Name: "Test Co 1", AddressID: address.ID}
		Expect(gr.Companies().Create(ctx, company1)).To(Succeed())

		company2 = &types.Company{Name: "Test Co 2", AddressID: address.ID}
		Expect(gr.Companies().Create(ctx, company2)).To(Succeed())

		company3 = &types.Company{Name: "Test Co 3", AddressID: address.ID}
		Expect(gr.Companies().Create(ctx, company3)).To(Succeed())

		commodity = &types.Commodity{Name: "Apple", CommodityType: types.CommodityTypeProduce}
		Expect(gr.Commodities().Create(ctx, commodity)).To(Succeed())

		// These must be created in a specific order for the fallback test to be predictable
		commodityAttribute1 = &types.CommodityAttribute{Name: "Variety", CommodityType: types.CommodityTypeProduce}
		Expect(gr.CommodityAttributes().Create(ctx, commodityAttribute1)).To(Succeed())

		commodityAttribute2 = &types.CommodityAttribute{Name: "Color", CommodityType: types.CommodityTypeProduce}
		Expect(gr.CommodityAttributes().Create(ctx, commodityAttribute2)).To(Succeed())

		commodityAttribute3 = &types.CommodityAttribute{Name: "Size", CommodityType: types.CommodityTypeProduce}
		Expect(gr.CommodityAttributes().Create(ctx, commodityAttribute3)).To(Succeed())

		// Setup company-specific ordering
		// Company 1: Size, then Color, then Variety
		orderC1 := []*types.CompanyAttributeSetting{
			{CompanyID: company1.ID, CommodityAttributeID: commodityAttribute3.ID, DisplayOrder: 1},
			{CompanyID: company1.ID, CommodityAttributeID: commodityAttribute2.ID, DisplayOrder: 2},
			{CompanyID: company1.ID, CommodityAttributeID: commodityAttribute1.ID, DisplayOrder: 3},
		}
		for _, attr := range orderC1 {
			Expect(casRepo.Create(ctx, attr)).To(Succeed())
		}

		// Company 2: Color, then Variety, then Size
		orderC2 := []*types.CompanyAttributeSetting{
			{CompanyID: company2.ID, CommodityAttributeID: commodityAttribute2.ID, DisplayOrder: 1},
			{CompanyID: company2.ID, CommodityAttributeID: commodityAttribute1.ID, DisplayOrder: 2},
			{CompanyID: company2.ID, CommodityAttributeID: commodityAttribute3.ID, DisplayOrder: 3},
		}
		for _, attr := range orderC2 {
			Expect(casRepo.Create(ctx, attr)).To(Succeed())
		}
		// Company 3 has no specific order, so it should use the fallback.
	})

	It("should create products with derived names based on company-specific and fallback ordering", func() {
		// Test Data
		testCases := []struct {
			company      *types.Company
			attributes   []*types.ProductAttributeValue
			expectedName string
		}{
			// Company 1 Products (Order: Size, Color, Variety)
			{
				company: company1,
				attributes: []*types.ProductAttributeValue{
					{CommodityAttributeID: commodityAttribute1.ID, Value: "Honeycrisp"},
					{CommodityAttributeID: commodityAttribute2.ID, Value: "Red"},
					{CommodityAttributeID: commodityAttribute3.ID, Value: "Large"},
				},
				expectedName: "Large Red Honeycrisp Apple",
			},
			{
				company: company1,
				attributes: []*types.ProductAttributeValue{
					{CommodityAttributeID: commodityAttribute1.ID, Value: "Granny Smith"},
					{CommodityAttributeID: commodityAttribute2.ID, Value: "Green"},
					{CommodityAttributeID: commodityAttribute3.ID, Value: "Small"},
				},
				expectedName: "Small Green Granny Smith Apple",
			},
			// Company 2 Products (Order: Color, Variety, Size)
			{
				company: company2,
				attributes: []*types.ProductAttributeValue{
					{CommodityAttributeID: commodityAttribute1.ID, Value: "Fuji"},
					{CommodityAttributeID: commodityAttribute2.ID, Value: "Pink"},
					{CommodityAttributeID: commodityAttribute3.ID, Value: "Medium"},
				},
				expectedName: "Pink Fuji Medium Apple",
			},
			// Company 3 Products (Fallback Order: Variety, Color, Size - based on ID)
			{
				company: company3,
				attributes: []*types.ProductAttributeValue{
					{CommodityAttributeID: commodityAttribute1.ID, Value: "Gala"},
					{CommodityAttributeID: commodityAttribute2.ID, Value: "Red-Yellow"},
					{CommodityAttributeID: commodityAttribute3.ID, Value: "Medium"},
				},
				expectedName: "Gala Red-Yellow Medium Apple",
			},
		}

		// Adding more products to reach the ~20 goal
		for i := 0; i < 15; i++ {
			company := company1
			attrs := []*types.ProductAttributeValue{
				{CommodityAttributeID: commodityAttribute1.ID, Value: "Honeycrisp"},
				{CommodityAttributeID: commodityAttribute2.ID, Value: "Red"},
				{CommodityAttributeID: commodityAttribute3.ID, Value: "Large"},
			}
			expectedName := "Large Red Honeycrisp Apple" // Company 1 order: Size, Color, Variety

			if i%2 == 0 {
				company = company2
				attrs = []*types.ProductAttributeValue{
					{CommodityAttributeID: commodityAttribute1.ID, Value: "Fuji"},
					{CommodityAttributeID: commodityAttribute2.ID, Value: "Pink"},
					{CommodityAttributeID: commodityAttribute3.ID, Value: "Medium"},
				}
				expectedName = "Pink Fuji Medium Apple" // Company 2 order: Color, Variety, Size
			} else if i%3 == 0 {
				company = company3
				attrs = []*types.ProductAttributeValue{
					{CommodityAttributeID: commodityAttribute1.ID, Value: "Gala"},
					{CommodityAttributeID: commodityAttribute2.ID, Value: "Red-Yellow"},
					{CommodityAttributeID: commodityAttribute3.ID, Value: "Medium"},
				}
				expectedName = "Gala Red-Yellow Medium Apple" // Fallback order: Variety, Color, Size
			}

			testCases = append(testCases, struct {
				company         *types.Company
				attributes      []*types.ProductAttributeValue
				expectedName    string
			}{
				company:      company,
				attributes:   attrs,
				expectedName: expectedName,
			})
		}

		for i, tc := range testCases {
			By(fmt.Sprintf("Running test case %d for company %s", i+1, tc.company.Name), func() {
				product := &types.Product{
					CompanyID:   tc.company.ID,
					CommodityID: commodity.ID,
				}

				err := repo.Create(ctx, product, tc.attributes)
				Expect(err).NotTo(HaveOccurred())

				retrieved, found, err := repo.Get(ctx, product.ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(retrieved.Name).To(Equal(tc.expectedName))
			})
		}
	})

	Context("Update method", func() {
		var (
			product *types.Product
			initialAttrs []*types.ProductAttributeValue
		)

		BeforeEach(func() {
			// Create a product with initial attributes
			product = &types.Product{
				CompanyID:   company1.ID,
				CommodityID: commodity.ID,
			}
			initialAttrs = []*types.ProductAttributeValue{
				{CommodityAttributeID: commodityAttribute1.ID, Value: "InitialVariety"},
				{CommodityAttributeID: commodityAttribute2.ID, Value: "InitialColor"},
			}
			Expect(repo.Create(ctx, product, initialAttrs)).To(Succeed())
		})

		It("should retain existing attributes if new attributes are nil or empty", func() {
			// Update the product without providing new attributes
			updatedProduct := &types.Product{
				ID:          product.ID,
				CompanyID:   product.CompanyID,
				CommodityID: product.CommodityID,
			}
			Expect(repo.Update(ctx, updatedProduct, nil)).To(Succeed())

			// Retrieve the product and verify its attributes are still the initial ones
			retrieved, found, err := repo.Get(ctx, product.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())

			// Verify the name is still derived from initial attributes
			Expect(retrieved.Name).To(ContainSubstring("InitialVariety"))
			Expect(retrieved.Name).To(ContainSubstring("InitialColor"))

			// To be more robust, we should ideally fetch the ProductAttributeValues directly
			// and assert their count and values. This requires a Find method for ProductAttributeValue.
			// For now, relying on the derived name is a reasonable proxy.
		})

		It("should update attributes if new attributes are provided", func() {
			newAttrs := []*types.ProductAttributeValue{
				{CommodityAttributeID: commodityAttribute1.ID, Value: "NewVariety"},
				{CommodityAttributeID: commodityAttribute2.ID, Value: "NewColor"},
			}
			updatedProduct := &types.Product{
				ID:          product.ID,
				CompanyID:   product.CompanyID,
				CommodityID: product.CommodityID,
			}
			Expect(repo.Update(ctx, updatedProduct, newAttrs)).To(Succeed())

			retrieved, found, err := repo.Get(ctx, product.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())

			Expect(retrieved.Name).To(ContainSubstring("NewVariety"))
			Expect(retrieved.Name).To(ContainSubstring("NewColor"))
			Expect(retrieved.Name).NotTo(ContainSubstring("InitialVariety"))
		})
	})
})