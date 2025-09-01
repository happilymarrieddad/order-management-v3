package repos_test

import (
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CommodityRepo Integration", func() {
	var commodityRepo repos.CommoditiesRepo

	BeforeEach(func() {
		commodityRepo = gr.Commodities()
	})

	Context("Create and Get", func() {
		It("should create a new commodity and then retrieve it", func() {
			newCommodity := &types.Commodity{
				Name:          "Test Commodity",
				CommodityType: types.CommodityTypeProduce,
				Visible:       true,
			}

			err := commodityRepo.Create(ctx, newCommodity)
			Expect(err).NotTo(HaveOccurred())
			Expect(newCommodity.ID).To(BeNumerically(">", 0))

			fetchedCommodity, found, err := commodityRepo.Get(ctx, newCommodity.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(fetchedCommodity).NotTo(BeNil())
			Expect(fetchedCommodity.ID).To(Equal(newCommodity.ID))
			Expect(fetchedCommodity.Name).To(Equal(newCommodity.Name))
		})

		It("should return not found for a non-existent commodity ID", func() {
			_, found, err := commodityRepo.Get(ctx, 99999)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())
		})

		It("should return an error for a duplicate name", func() {
			commodity1 := &types.Commodity{
				Name:          "Duplicate Name",
				CommodityType: types.CommodityTypeProduce,
				Visible:       true,
			}
			err := commodityRepo.Create(ctx, commodity1)
			Expect(err).NotTo(HaveOccurred())

			commodity2 := &types.Commodity{
				Name:          "Duplicate Name",
				CommodityType: types.CommodityTypeProduce,
				Visible:       true,
			}
			err = commodityRepo.Create(ctx, commodity2)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Update", func() {
		It("should update an existing commodity", func() {
			commodity := &types.Commodity{
				Name:          "Original Name",
				CommodityType: types.CommodityTypeProduce,
				Visible:       true,
			}
			err := commodityRepo.Create(ctx, commodity)
			Expect(err).NotTo(HaveOccurred())

			commodity.Name = "Updated Name"
			err = commodityRepo.Update(ctx, commodity)
			Expect(err).NotTo(HaveOccurred())

			updatedCommodity, found, err := commodityRepo.Get(ctx, commodity.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(updatedCommodity.Name).To(Equal("Updated Name"))
		})
	})

	Context("Find", func() {
		BeforeEach(func() {
			commoditiesToCreate := []*types.Commodity{
				{Name: "Find Apple", CommodityType: types.CommodityTypeProduce, Visible: true},
				{Name: "Find Banana", CommodityType: types.CommodityTypeProduce, Visible: true},
				{Name: "Find Not Visible", CommodityType: types.CommodityType(99), Visible: false},
				{Name: "Find Carrot", CommodityType: types.CommodityType(99), Visible: true}, // Different type
			}
			for _, c := range commoditiesToCreate {
				err := commodityRepo.Create(ctx, c)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("should find all commodities when no options are provided", func() {
			commodities, count, err := commodityRepo.Find(ctx, &repos.FindCommoditiesOpts{})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(BeNumerically(">=", 3))
			Expect(len(commodities)).To(BeNumerically(">=", 3))
		})

		It("should find commodities by name", func() {
			opts := &repos.FindCommoditiesOpts{Name: "Apple"}
			commodities, count, err := commodityRepo.Find(ctx, opts)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(1)))
			Expect(len(commodities)).To(Equal(1))
			Expect(commodities[0].Name).To(Equal("Find Apple"))
		})

		It("should find commodities by commodity type", func() {
			opts := &repos.FindCommoditiesOpts{CommodityType: types.CommodityTypeProduce}
			commodities, count, err := commodityRepo.Find(ctx, opts)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(2)))
			Expect(len(commodities)).To(Equal(2))
		})

		It("should respect limit and offset for visible items", func() {
			opts := &repos.FindCommoditiesOpts{Limit: 1, Offset: 1}
			commodities, count, err := commodityRepo.Find(ctx, opts)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(3))) // Total visible count
			Expect(len(commodities)).To(Equal(1))
		})
	})
})
