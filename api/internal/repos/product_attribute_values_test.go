package repos_test

import (
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ProductAttributeValuesRepo", func() {
	var (
		repo                repos.ProductAttributeValuesRepo
		company1            *types.Company
		company2            *types.Company
		commodity           *types.Commodity
		commodityAttribute1 *types.CommodityAttribute
		commodityAttribute2 *types.CommodityAttribute
		product1            *types.Product
		product2            *types.Product
		product3            *types.Product
	)

	BeforeEach(func() {
		repo = gr.ProductAttributeValues()

		// Create dependencies
		address, err := gr.Addresses().Create(ctx, &types.Address{
			Line1: "123 Main St", City: "Anytown", State: "CA", Country: "USA", PostalCode: "12345",
		})
		Expect(err).NotTo(HaveOccurred())

		company1 = &types.Company{Name: "Test Co 1", AddressID: address.ID}
		Expect(gr.Companies().Create(ctx, company1)).To(Succeed())

		company2 = &types.Company{Name: "Test Co 2", AddressID: address.ID}
		Expect(gr.Companies().Create(ctx, company2)).To(Succeed())

		commodity = &types.Commodity{Name: "Test Commodity", CommodityType: types.CommodityTypeProduce}
		Expect(gr.Commodities().Create(ctx, commodity)).To(Succeed())

		commodityAttribute1 = &types.CommodityAttribute{Name: "Color", CommodityType: types.CommodityTypeProduce}
		Expect(gr.CommodityAttributes().Create(ctx, commodityAttribute1)).To(Succeed())

		commodityAttribute2 = &types.CommodityAttribute{Name: "Size", CommodityType: types.CommodityTypeProduce}
		Expect(gr.CommodityAttributes().Create(ctx, commodityAttribute2)).To(Succeed())

		product1 = &types.Product{CompanyID: company1.ID, CommodityID: commodity.ID, Name: "p1"}
		_, err = db.Insert(product1) // Simplified creation for testing
		Expect(err).NotTo(HaveOccurred())

		product2 = &types.Product{CompanyID: company1.ID, CommodityID: commodity.ID, Name: "p2"}
		_, err = db.Insert(product2) // Simplified creation for testing
		Expect(err).NotTo(HaveOccurred())

		product3 = &types.Product{CompanyID: company2.ID, CommodityID: commodity.ID, Name: "p3"}
		_, err = db.Insert(product3) // Simplified creation for testing
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Create and Get", func() {
		It("should create and retrieve a product attribute value", func() {
			attr := &types.ProductAttributeValue{
				ProductID:            product1.ID,
				CommodityAttributeID: commodityAttribute1.ID,
				CompanyID:            company1.ID,
				Value:                "Red",
			}

			Expect(repo.Create(ctx, attr)).To(Succeed())
			Expect(attr.ID).NotTo(BeZero())

			retrieved, found, err := repo.Get(ctx, attr.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(retrieved.Value).To(Equal("Red"))
			Expect(retrieved.ProductID).To(Equal(product1.ID))
		})
	})

	Describe("Update", func() {
		It("should update a product attribute value without affecting others", func() {
			attr1 := &types.ProductAttributeValue{
				ProductID:            product1.ID,
				CommodityAttributeID: commodityAttribute1.ID,
				CompanyID:            company1.ID,
				Value:                "Blue",
			}
			Expect(repo.Create(ctx, attr1)).To(Succeed())

			attr2 := &types.ProductAttributeValue{
				ProductID:            product1.ID,
				CommodityAttributeID: commodityAttribute2.ID,
				CompanyID:            company1.ID,
				Value:                "Small",
			}
			Expect(repo.Create(ctx, attr2)).To(Succeed())

			attr1.Value = "Green"
			Expect(repo.Update(ctx, attr1)).To(Succeed())

			retrieved1, found, err := repo.Get(ctx, attr1.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(retrieved1.Value).To(Equal("Green"))

			retrieved2, found, err := repo.Get(ctx, attr2.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(retrieved2.Value).To(Equal("Small")) // Unchanged
		})
	})

	Describe("Delete", func() {
		It("should delete multiple product attribute values", func() {
			attr1 := &types.ProductAttributeValue{ProductID: product1.ID, CommodityAttributeID: commodityAttribute1.ID, CompanyID: company1.ID, Value: "Val1"}
			attr2 := &types.ProductAttributeValue{ProductID: product2.ID, CommodityAttributeID: commodityAttribute1.ID, CompanyID: company1.ID, Value: "Val2"}
			Expect(repo.Create(ctx, attr1)).To(Succeed())
			Expect(repo.Create(ctx, attr2)).To(Succeed())

			Expect(repo.Delete(ctx, attr1.ID, attr2.ID)).To(Succeed())

			_, found, err := repo.Get(ctx, attr1.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())

			_, found, err = repo.Get(ctx, attr2.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())
		})
	})

	Describe("Find", func() {
		BeforeEach(func() {
			// Create a set of values for testing Find
			attrs := []*types.ProductAttributeValue{
				{ProductID: product1.ID, CommodityAttributeID: commodityAttribute1.ID, CompanyID: company1.ID, Value: "Red"},
				{ProductID: product1.ID, CommodityAttributeID: commodityAttribute2.ID, CompanyID: company1.ID, Value: "Large"},
				{ProductID: product2.ID, CommodityAttributeID: commodityAttribute1.ID, CompanyID: company1.ID, Value: "Blue"},
				{ProductID: product3.ID, CommodityAttributeID: commodityAttribute1.ID, CompanyID: company2.ID, Value: "Green"},
			}
			for _, attr := range attrs {
				Expect(repo.Create(ctx, attr)).To(Succeed())
			}
		})

		It("should find by ProductIDs", func() {
			found, count, err := repo.Find(ctx, &repos.ProductAttributeValueFindOpts{ProductIDs: []int64{product1.ID, product3.ID}})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(3)))
			Expect(found).To(HaveLen(3))
		})

		It("should find by CompanyIDs", func() {
			found, count, err := repo.Find(ctx, &repos.ProductAttributeValueFindOpts{CompanyIDs: []int64{company1.ID}})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(3)))
			Expect(found).To(HaveLen(3))
		})

		It("should find by a combination of ProductIDs and CompanyIDs", func() {
			found, count, err := repo.Find(ctx, &repos.ProductAttributeValueFindOpts{
				ProductIDs: []int64{product1.ID},
				CompanyIDs: []int64{company1.ID},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(2)))
			Expect(found).To(HaveLen(2))
		})

		It("should respect pagination", func() {
			found, count, err := repo.Find(ctx, &repos.ProductAttributeValueFindOpts{Limit: 2, Offset: 1})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(4)))
			Expect(found).To(HaveLen(2))
		})
	})
})
