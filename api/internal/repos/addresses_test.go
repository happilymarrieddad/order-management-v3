package repos_test

import (
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AddressRepo Integration", func() {
	var addressRepo repos.AddressesRepo

	BeforeEach(func() {
		// The global repo 'gr' is initialized in repos_suite_test.go
		// and connected to a real test database.
		addressRepo = gr.Addresses()
	})

	Context("Create and Get", func() {
		It("should create a new address and then retrieve it", func() {
			newAddr := &types.Address{
				Line1:      "123 Test St",
				City:       "Testville",
				State:      "TS",
				PostalCode: "12345",
				Country:    "USA",
				GlobalCode: "849VCWC8+R9",
			}

			createdAddr, err := addressRepo.Create(ctx, newAddr)
			Expect(err).NotTo(HaveOccurred())
			Expect(createdAddr).NotTo(BeNil())
			Expect(createdAddr.ID).To(BeNumerically(">", 0))
			Expect(createdAddr.Line1).To(Equal(newAddr.Line1))
			Expect(createdAddr.GlobalCode).To(Equal(newAddr.GlobalCode))

			// Now, get the address by its new ID
			fetchedAddr, found, err := addressRepo.Get(ctx, createdAddr.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(fetchedAddr).NotTo(BeNil())
			Expect(fetchedAddr.ID).To(Equal(createdAddr.ID))
			Expect(fetchedAddr.City).To(Equal("Testville"))
		})

		It("should return not found for a non-existent address ID", func() {
			_, found, err := addressRepo.Get(ctx, 99999)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())
		})
	})

	Context("Update", func() {
		It("should update an existing address", func() {
			// First, create an address to update
			addr := &types.Address{Line1: "Original Ave", City: "Old Town",
				State:      "TS",
				PostalCode: "12345",
				Country:    "USA",
				GlobalCode: "849VCWC8+R9"}
			createdAddr, err := addressRepo.Create(ctx, addr)
			Expect(err).NotTo(HaveOccurred())

			// Now, update it
			createdAddr.Line1 = "Updated Blvd"
			createdAddr.City = "New City"
			err = addressRepo.Update(ctx, createdAddr)
			Expect(err).NotTo(HaveOccurred())

			// Retrieve it again to verify the update
			updatedAddr, found, err := addressRepo.Get(ctx, createdAddr.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(updatedAddr.Line1).To(Equal("Updated Blvd"))
			Expect(updatedAddr.City).To(Equal("New City"))
		})
	})

	Context("Delete", func() {
		It("should delete an existing address", func() {
			// Create an address to delete
			addr := &types.Address{Line1: "To Be Deleted Dr", City: "Old Town",
				State:      "TS",
				PostalCode: "12345",
				Country:    "USA",
				GlobalCode: "849VCWC8+R9"}
			createdAddr, err := addressRepo.Create(ctx, addr)
			Expect(err).NotTo(HaveOccurred())

			// Delete it
			err = addressRepo.Delete(ctx, createdAddr.ID)
			Expect(err).NotTo(HaveOccurred())

			// Try to get it again, it should not be found
			_, found, err := addressRepo.Get(ctx, createdAddr.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())
		})
	})

	Context("Find", func() {
		BeforeEach(func() {
			// Seed the database with some addresses for find tests
			addressesToCreate := []*types.Address{
				{Line1: "1 Main St", City: "FindMeVille", State: "TS",
					PostalCode: "12345",
					Country:    "USA"},
				{Line1: "2 Main St", City: "FindMeVille", State: "TS",
					PostalCode: "12345",
					Country:    "USA"},
				{Line1: "3 Side St", City: "OtherPlace", State: "TS",
					PostalCode: "12345",
					Country:    "USA"},
				{Line1: "4 Main St", City: "FindMeVille", State: "TS",
					PostalCode: "12345",
					Country:    "USA"},
				{Line1: "5 Another Ave", City: "OtherPlace", State: "TS",
					PostalCode: "12345",
					Country:    "USA"},
			}
			for _, addr := range addressesToCreate {
				_, err := addressRepo.Create(ctx, addr)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("should find all addresses when no options are provided", func() {
			addrs, count, err := addressRepo.Find(ctx, &repos.AddressFindOpts{})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(5)))
			Expect(len(addrs)).To(Equal(5))
		})

		It("should find addresses by a list of IDs", func() {
			opts := &repos.AddressFindOpts{IDs: []int64{1, 3, 5}}
			addrs, count, err := addressRepo.Find(ctx, opts)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(3)))
			Expect(len(addrs)).To(Equal(3))
			Expect(addrs[0].ID).To(Equal(int64(1)))
			Expect(addrs[1].ID).To(Equal(int64(3)))
			Expect(addrs[2].ID).To(Equal(int64(5)))
		})

		It("should respect limit and offset", func() {
			opts := &repos.AddressFindOpts{Limit: 2, Offset: 1} // Get the 2nd and 3rd addresses
			addrs, count, err := addressRepo.Find(ctx, opts)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(5))) // Total count should still be 5
			Expect(len(addrs)).To(Equal(2))
			Expect(addrs[0].ID).To(Equal(int64(2)))
			Expect(addrs[1].ID).To(Equal(int64(3)))
		})
	})
})
