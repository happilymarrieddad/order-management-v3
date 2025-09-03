package repos_test

import (
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
			Expect(fetchedCompany.Address).NotTo(BeNil())
			Expect(fetchedCompany.Address.ID).To(Equal(createdAddr.ID))
			Expect(fetchedCompany.Address.Line1).To(Equal(addr.Line1))
			Expect(fetchedCompany.Address.City).To(Equal(addr.City))
			Expect(fetchedCompany.Address.State).To(Equal(addr.State))
			Expect(fetchedCompany.Address.Country).To(Equal(addr.Country))
			Expect(fetchedCompany.Address.PostalCode).To(Equal(addr.PostalCode))
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
			// This checks for the validation error from the repository
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
		It("should soft delete an existing company", func() {
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

			// The company should be found with GetIncludeInvisible
			invisibleCompany, found, err := companyRepo.GetIncludeInvisible(ctx, company.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(invisibleCompany.Visible).To(BeFalse())
		})
	})

	Context("Find", func() {
		var (
			company1 *types.Company
			company2 *types.Company
			company3 *types.Company
		)

		BeforeEach(func() {
			// Seed the database with some companies for find tests
			addr1 := &types.Address{Line1: "1 Find St", City: "Findsburge", State: "FS", Country: "USA", PostalCode: "55555"}
			createdAddr1, err := addressRepo.Create(ctx, addr1)
			Expect(err).NotTo(HaveOccurred())
			company1 = &types.Company{Name: "Alpha Corp", AddressID: createdAddr1.ID}
			Expect(companyRepo.Create(ctx, company1)).To(Succeed())

			addr2 := &types.Address{Line1: "2 Find St", City: "Findsburge", State: "FS", Country: "USA", PostalCode: "55555"}
			createdAddr2, err := addressRepo.Create(ctx, addr2)
			Expect(err).NotTo(HaveOccurred())
			company2 = &types.Company{Name: "Beta LLC", AddressID: createdAddr2.ID}
			Expect(companyRepo.Create(ctx, company2)).To(Succeed())

			addr3 := &types.Address{Line1: "3 Find St", City: "Findsburge", State: "FS", Country: "USA", PostalCode: "55555"}
			createdAddr3, err := addressRepo.Create(ctx, addr3)
			Expect(err).NotTo(HaveOccurred())
			company3 = &types.Company{Name: "Gamma Inc", AddressID: createdAddr3.ID}
			Expect(companyRepo.Create(ctx, company3)).To(Succeed())
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

		It("should find companies by name using LIKE (case-insensitive and partial match)", func() {
			foundCompanies, count, err := companyRepo.Find(ctx, &repos.CompanyFindOpts{Names: []string{"alpha"}})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(1)))
			Expect(foundCompanies).To(HaveLen(1))
			Expect(foundCompanies[0].Name).To(Equal(company1.Name))

			foundCompanies, count, err = companyRepo.Find(ctx, &repos.CompanyFindOpts{Names: []string{"llc"}})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(1)))
			Expect(foundCompanies).To(HaveLen(1))
			Expect(foundCompanies[0].Name).To(Equal(company2.Name))

			foundCompanies, count, err = companyRepo.Find(ctx, &repos.CompanyFindOpts{Names: []string{"corp", "inc"}})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(2)))
			Expect(foundCompanies).To(HaveLen(2))
			Expect(foundCompanies).To(ConsistOf(
				WithTransform(func(c *types.Company) string { return c.Name }, Equal(company1.Name)),
				WithTransform(func(c *types.Company) string { return c.Name }, Equal(company3.Name)),
			))
		})

		It("should find companies by IDs", func() {
			foundCompanies, count, err := companyRepo.Find(ctx, &repos.CompanyFindOpts{IDs: []int64{company1.ID, company3.ID}})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(2)))
			Expect(foundCompanies).To(HaveLen(2))
			Expect(foundCompanies).To(ConsistOf(
				WithTransform(func(c *types.Company) int64 { return c.ID }, Equal(company1.ID)),
				WithTransform(func(c *types.Company) int64 { return c.ID }, Equal(company3.ID)),
			))
		})

		It("should find companies by a combination of filters", func() {
			foundCompanies, count, err := companyRepo.Find(ctx, &repos.CompanyFindOpts{
				Names: []string{"beta"},
				IDs:   []int64{company2.ID},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(1)))
			Expect(foundCompanies).To(HaveLen(1))
			Expect(foundCompanies[0].ID).To(Equal(company2.ID))
		})
	})
})
