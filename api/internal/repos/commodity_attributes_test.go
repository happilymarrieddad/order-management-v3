package repos_test

import (
	"context"

	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CommodityAttributesRepo", func() {
	var (
		repo repos.CommodityAttributesRepo
		ctx  context.Context
	)

	BeforeEach(func() {
		repo = gr.CommodityAttributes()
		ctx = context.Background()
	})

	Context("Create and Get", func() {
		It("should create a new commodity attribute and retrieve it", func() {
			newAttribute := &types.CommodityAttribute{
				Name:          "Color",
				CommodityType: types.CommodityTypeProduce,
			}

			err := repo.Create(ctx, newAttribute)
			Expect(err).NotTo(HaveOccurred())
			Expect(newAttribute.ID).NotTo(BeZero())

			fetchedAttribute, found, err := repo.Get(ctx, newAttribute.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(fetchedAttribute.ID).To(Equal(newAttribute.ID))
			Expect(fetchedAttribute.Name).To(Equal(newAttribute.Name))
			Expect(fetchedAttribute.CommodityType).To(Equal(newAttribute.CommodityType))
		})

		It("should return not found for a non-existent ID", func() {
			_, found, err := repo.Get(ctx, 99999)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())
		})

		It("should fail to create an attribute with a duplicate name", func() {
			newAttribute := &types.CommodityAttribute{
				Name:          "Size",
				CommodityType: types.CommodityTypeProduce,
			}
			err := repo.Create(ctx, newAttribute)
			Expect(err).NotTo(HaveOccurred())

			duplicateAttribute := &types.CommodityAttribute{
				Name:          "Size", // Duplicate name
				CommodityType: types.CommodityTypeProduce,
			}
			err = repo.Create(ctx, duplicateAttribute)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("duplicate key value violates unique constraint"))
		})
	})

	Context("Update", func() {
		var existingAttribute *types.CommodityAttribute

		BeforeEach(func() {
			existingAttribute = &types.CommodityAttribute{
				Name:          "Weight",
				CommodityType: types.CommodityTypeProduce,
			}
			err := repo.Create(ctx, existingAttribute)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should update an existing commodity attribute", func() {
			existingAttribute.Name = "Updated Weight"

			err := repo.Update(ctx, existingAttribute)
			Expect(err).NotTo(HaveOccurred())

			updatedAttribute, found, err := repo.Get(ctx, existingAttribute.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(updatedAttribute.Name).To(Equal("Updated Weight"))
		})
	})

	Context("Find", func() {
		BeforeEach(func() {
			// Seed data for find tests
			attributesToCreate := []*types.CommodityAttribute{
				{Name: "Find-Color-Produce", CommodityType: types.CommodityTypeProduce},
				{Name: "Find-Size-Produce", CommodityType: types.CommodityTypeProduce},
				{Name: "Find-Color-Unknown", CommodityType: types.CommodityTypeUnknown},
			}
			for _, attr := range attributesToCreate {
				err := repo.Create(ctx, attr)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("should find attributes by commodity type", func() {
			foundAttributes, count, err := repo.Find(ctx, &repos.CommodityAttributeFindOpts{
				CommodityTypes: []types.CommodityType{types.CommodityTypeProduce},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(BeNumerically(">=", 2)) // At least 2 produce attributes
			Expect(foundAttributes).To(ContainElements(
				WithTransform(func(attr *types.CommodityAttribute) string { return attr.Name }, Equal("Find-Color-Produce")),
				WithTransform(func(attr *types.CommodityAttribute) string { return attr.Name }, Equal("Find-Size-Produce")),
			))
			Expect(foundAttributes).NotTo(ContainElement(
				WithTransform(func(attr *types.CommodityAttribute) string { return attr.Name }, Equal("Find-Color-Unknown")),
			))
		})

		It("should find attributes by multiple commodity types", func() {
			foundAttributes, count, err := repo.Find(ctx, &repos.CommodityAttributeFindOpts{
				CommodityTypes: []types.CommodityType{types.CommodityTypeProduce, types.CommodityTypeUnknown},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(BeNumerically(">=", 3)) // At least 3 total attributes
			Expect(foundAttributes).To(ContainElements(
				WithTransform(func(attr *types.CommodityAttribute) string { return attr.Name }, Equal("Find-Color-Produce")),
				WithTransform(func(attr *types.CommodityAttribute) string { return attr.Name }, Equal("Find-Size-Produce")),
				WithTransform(func(attr *types.CommodityAttribute) string { return attr.Name }, Equal("Find-Color-Unknown")),
			))
		})

		It("should respect limit and offset", func() {
			// Assuming at least 3 attributes are seeded
			foundAttributes, count, err := repo.Find(ctx, &repos.CommodityAttributeFindOpts{
				Limit:  1,
				Offset: 0,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(BeNumerically(">=", 3))
			Expect(foundAttributes).To(HaveLen(1))

			foundAttributes, count, err = repo.Find(ctx, &repos.CommodityAttributeFindOpts{
				Limit:  1,
				Offset: 1,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(BeNumerically(">=", 3))
			Expect(foundAttributes).To(HaveLen(1))
		})
	})
})
