package repos_test

import (
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CompanyAttributeSettingsRepo", func() {
	var (
		repo                 repos.CompanyAttributeSettingsRepo
		company              *types.Company
		commodityAttribute1  *types.CommodityAttribute
		commodityAttribute2  *types.CommodityAttribute
	)

	BeforeEach(func() {
		repo = gr.CompanyAttributeSettings()

		// Create dependencies
		address, err := gr.Addresses().Create(ctx, &types.Address{
			Line1: "123 Main St", City: "Anytown", State: "CA", Country: "USA", PostalCode: "12345",
		})
		Expect(err).NotTo(HaveOccurred())

		company = &types.Company{Name: "Test Co", AddressID: address.ID}
		Expect(gr.Companies().Create(ctx, company)).To(Succeed())

		commodityAttribute1 = &types.CommodityAttribute{Name: "Color", CommodityType: types.CommodityTypeProduce}
		Expect(gr.CommodityAttributes().Create(ctx, commodityAttribute1)).To(Succeed())

		commodityAttribute2 = &types.CommodityAttribute{Name: "Size", CommodityType: types.CommodityTypeProduce}
		Expect(gr.CommodityAttributes().Create(ctx, commodityAttribute2)).To(Succeed())
	})

	Describe("Create and Get", func() {
		It("should create and retrieve a setting", func() {
			setting := &types.CompanyAttributeSetting{
				CompanyID:            company.ID,
				CommodityAttributeID: commodityAttribute1.ID,
				DisplayOrder:         1,
			}

			Expect(repo.Create(ctx, setting)).To(Succeed())
			Expect(setting.ID).NotTo(BeZero())

			retrieved, found, err := repo.Get(ctx, setting.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(retrieved.DisplayOrder).To(Equal(1))
			Expect(retrieved.CommodityAttributeID).To(Equal(commodityAttribute1.ID))
		})
	})

	Describe("Update", func() {
		It("should update a setting", func() {
			setting := &types.CompanyAttributeSetting{
				CompanyID:            company.ID,
				CommodityAttributeID: commodityAttribute1.ID,
				DisplayOrder:         1,
			}
			Expect(repo.Create(ctx, setting)).To(Succeed())

			setting.DisplayOrder = 2
			setting.CommodityAttributeID = commodityAttribute2.ID
			Expect(repo.Update(ctx, setting)).To(Succeed())

			retrieved, found, err := repo.Get(ctx, setting.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(retrieved.DisplayOrder).To(Equal(2))
			Expect(retrieved.CommodityAttributeID).To(Equal(commodityAttribute2.ID))
		})
	})

	Describe("Delete", func() {
		It("should delete settings", func() {
			setting1 := &types.CompanyAttributeSetting{CompanyID: company.ID, CommodityAttributeID: commodityAttribute1.ID, DisplayOrder: 1}
			setting2 := &types.CompanyAttributeSetting{CompanyID: company.ID, CommodityAttributeID: commodityAttribute2.ID, DisplayOrder: 2}
			Expect(repo.Create(ctx, setting1)).To(Succeed())
			Expect(repo.Create(ctx, setting2)).To(Succeed())

			Expect(repo.Delete(ctx, setting1.ID, setting2.ID)).To(Succeed())

			_, found, err := repo.Get(ctx, setting1.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())

			_, found, err = repo.Get(ctx, setting2.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())
		})
	})

	Describe("Find", func() {
		It("should find settings by company id", func() {
			setting1 := &types.CompanyAttributeSetting{CompanyID: company.ID, CommodityAttributeID: commodityAttribute1.ID, DisplayOrder: 1}
			setting2 := &types.CompanyAttributeSetting{CompanyID: company.ID, CommodityAttributeID: commodityAttribute2.ID, DisplayOrder: 2}
			Expect(repo.Create(ctx, setting1)).To(Succeed())
			Expect(repo.Create(ctx, setting2)).To(Succeed())

			found, count, err := repo.Find(ctx, &repos.CompanyAttributeSettingFindOpts{CompanyIDs: []int64{company.ID}})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(2)))
			Expect(found).To(HaveLen(2))
		})

		It("should return an empty slice for a company with no settings", func() {
			found, count, err := repo.Find(ctx, &repos.CompanyAttributeSettingFindOpts{CompanyIDs: []int64{999}})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(0)))
			Expect(found).To(BeEmpty())
		})
	})
})
