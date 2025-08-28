package repos_test

import (
	"fmt"

	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CompanyRepo Integration", func() {
	var (
		companyRepo repos.CompaniesRepo
		addressRepo repos.AddressesRepo
	)

	BeforeEach(func() {
		// The global repo 'gr' is initialized in repos_suite_test.go
		// and connected to a real test database.
		companyRepo = gr.Companies()
		addressRepo = gr.Addresses()
	})

	Context("Create and Get", func() {
		It("should create a new company and then retrieve it", func() {
			// A company requires an address, so create one first.
			addr := &types.Address{Line1: "1 Corporate Way", City: "Biz Town", State: "BT", Country: "USA", PostalCode: "12345"}
			createdAddr, err := addressRepo.Create(ctx, addr)
			Expect(err).NotTo(HaveOccurred())

			newCompany := &types.Company{
				Name:      "TestCo",
				AddressID: createdAddr.ID,
			}

			err = companyRepo.Create(ctx, newCompany)
			Expect(err).NotTo(HaveOccurred())
			Expect(newCompany.ID).To(BeNumerically(">", 0))

			// Now, get the company by its new ID
			fetchedCompany, found, err := companyRepo.Get(ctx, newCompany.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(fetchedCompany).NotTo(BeNil())
			Expect(fetchedCompany.ID).To(Equal(newCompany.ID))
			Expect(fetchedCompany.Name).To(Equal("TestCo"))
		})

		It("should return not found for a non-existent company ID", func() {
			_, found, err := companyRepo.Get(ctx, 99999)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())
		})

		It("should fail to create a company with an empty name", func() {
			// A company requires an address, so create one first to isolate the validation error.
			addr := &types.Address{Line1: "4 Invalid Way", City: "Fail Town", State: "FT", Country: "USA", PostalCode: "00000"}
			createdAddr, err := addressRepo.Create(ctx, addr)
			Expect(err).NotTo(HaveOccurred())

			// Assuming 'Name' is a required field in types.Company
			invalidCompany := &types.Company{Name: "", AddressID: createdAddr.ID}
			err = companyRepo.Create(ctx, invalidCompany)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Field validation for 'Name' failed on the 'required' tag"))
		})
	})

	Context("Update", func() {
		It("should update an existing company", func() {
			// First, create a company to update
			addr := &types.Address{Line1: "2 Original Pl", City: "Old City", State: "OC", Country: "USA", PostalCode: "54321"}
			createdAddr, err := addressRepo.Create(ctx, addr)
			Expect(err).NotTo(HaveOccurred())

			company := &types.Company{Name: "Original Name", AddressID: createdAddr.ID}
			err = companyRepo.Create(ctx, company)
			Expect(err).NotTo(HaveOccurred())

			// Now, update it
			company.Name = "Updated Name"
			err = companyRepo.Update(ctx, company)
			Expect(err).NotTo(HaveOccurred())

			// Retrieve it again to verify the update
			updatedCompany, found, err := companyRepo.Get(ctx, company.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(updatedCompany.Name).To(Equal("Updated Name"))
		})
	})

	Context("Delete", func() {
		It("should delete an existing company", func() {
			// Create a company to delete
			addr := &types.Address{Line1: "3 Deletion Dr", City: "Goneville", State: "GV", Country: "USA", PostalCode: "98765"}
			createdAddr, err := addressRepo.Create(ctx, addr)
			Expect(err).NotTo(HaveOccurred())

			company := &types.Company{Name: "To Be Deleted Inc.", AddressID: createdAddr.ID}
			err = companyRepo.Create(ctx, company)
			Expect(err).NotTo(HaveOccurred())

			// Delete it
			err = companyRepo.Delete(ctx, company.ID)
			Expect(err).NotTo(HaveOccurred())

			// Try to get it again, it should not be found
			_, found, err := companyRepo.Get(ctx, company.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())
		})
	})

	Context("Find", func() {
		BeforeEach(func() {
			// Seed the database with some companies for find tests
			companiesToCreate := []*types.Company{
				{Name: "Alpha Corp"}, {Name: "Beta LLC"}, {Name: "Gamma Inc"},
			}
			for i, company := range companiesToCreate {
				// Each company needs a unique address
				addr := &types.Address{Line1: fmt.Sprintf("%d Find St", i+1), City: "Findsburge", State: "FS", Country: "USA", PostalCode: "55555"}
				createdAddr, err := addressRepo.Create(ctx, addr)
				Expect(err).NotTo(HaveOccurred())
				company.AddressID = createdAddr.ID

				err = companyRepo.Create(ctx, company)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("should respect limit and offset", func() {
			// Get the 2nd and 3rd companies
			companies, count, err := companyRepo.Find(ctx, &repos.CompanyFindOpts{
				Limit: 2, Offset: 1,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(3))) // Total count should be 3
			Expect(len(companies)).To(Equal(2))
			Expect(companies[0].Name).To(Equal("Beta LLC"))
			Expect(companies[1].Name).To(Equal("Gamma Inc"))
		})
	})
})
