package repos_test

import (
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("LocationsRepo", func() {
	var (
		repo    repos.LocationsRepo
		company *types.Company
		address *types.Address
	)

	BeforeEach(func() {
		// From suite: gr, db, ctx
		repo = gr.Locations()

		address = &types.Address{
			Line1:      "123 Location St",
			City:       "Locoville",
			State:      "LC",
			PostalCode: "54321",
			Country:    "USA",
		}
		address, err := gr.Addresses().Create(ctx, address)
		Expect(err).NotTo(HaveOccurred())

		// Create dependencies for each test, following established patterns.
		company = &types.Company{Name: "Test Company for Locations", AddressID: address.ID}
		err = gr.Companies().Create(ctx, company)
		Expect(err).NotTo(HaveOccurred())

		// Update company with real address ID
		company.AddressID = address.ID
		err = gr.Companies().Update(ctx, company)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Create and Get", func() {
		It("should create a location and retrieve it successfully", func() {
			newLocation := &types.Location{
				CompanyID: company.ID,
				AddressID: address.ID,
				Name:      "Main Warehouse",
			}

			err := repo.Create(ctx, newLocation)
			Expect(err).NotTo(HaveOccurred())
			Expect(newLocation.ID).NotTo(BeZero())

			retrieved, found, err := repo.Get(ctx, company.ID, newLocation.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(retrieved.Name).To(Equal("Main Warehouse"))
			Expect(retrieved.CompanyID).To(Equal(company.ID))
			Expect(retrieved.Address).NotTo(BeNil())
			Expect(retrieved.Address.ID).To(Equal(address.ID))
		})

		It("should not get a location from another company", func() {
			newLocation := &types.Location{
				CompanyID: company.ID,
				AddressID: address.ID,
				Name:      "Another Warehouse",
			}
			err := repo.Create(ctx, newLocation)
			Expect(err).NotTo(HaveOccurred())

			_, found, err := repo.Get(ctx, company.ID+1, newLocation.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())
		})

		It("should fail to create a location with a duplicate name for the same company", func() {
			location1 := &types.Location{
				CompanyID: company.ID,
				AddressID: address.ID,
				Name:      "Duplicate Name Warehouse",
			}
			err := repo.Create(ctx, location1)
			Expect(err).NotTo(HaveOccurred())

			anotherAddress, err := gr.Addresses().Create(ctx, &types.Address{
				Line1: "456 Other St", City: "Otherville", State: "OT", PostalCode: "67890", Country: "USA",
			})
			Expect(err).NotTo(HaveOccurred())

			location2 := &types.Location{
				CompanyID: company.ID,
				AddressID: anotherAddress.ID,
				Name:      "Duplicate Name Warehouse", // Same name, same company
			}
			err = repo.Create(ctx, location2)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(repos.ErrLocationNameExists))
		})

		It("should successfully create locations with the same name for different companies", func() {
			company2 := &types.Company{Name: "Another Company", AddressID: address.ID}
			err := gr.Companies().Create(ctx, company2)
			Expect(err).NotTo(HaveOccurred())

			location1 := &types.Location{
				CompanyID: company.ID,
				AddressID: address.ID,
				Name:      "Shared Name Warehouse",
			}
			err = repo.Create(ctx, location1)
			Expect(err).NotTo(HaveOccurred())

			location2 := &types.Location{
				CompanyID: company2.ID,
				AddressID: address.ID,
				Name:      "Shared Name Warehouse", // Same name, different company
			}
			err = repo.Create(ctx, location2)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Update", func() {
		var existingLocation *types.Location

		BeforeEach(func() {
			existingLocation = &types.Location{
				CompanyID: company.ID,
				AddressID: address.ID,
				Name:      "Updatable Location",
			}
			err := repo.Create(ctx, existingLocation)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should update a location's name successfully", func() {
			existingLocation.Name = "Updated Location Name"
			err := repo.Update(ctx, existingLocation)
			Expect(err).NotTo(HaveOccurred())

			retrieved, found, err := repo.Get(ctx, company.ID, existingLocation.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(retrieved.Name).To(Equal("Updated Location Name"))
		})

		It("should fail to update a location name to a duplicate within the same company", func() {
			conflictingLocation := &types.Location{
				CompanyID: company.ID,
				AddressID: address.ID,
				Name:      "Existing Name",
			}
			err := repo.Create(ctx, conflictingLocation)
			Expect(err).NotTo(HaveOccurred())

			existingLocation.Name = "Existing Name"
			err = repo.Update(ctx, existingLocation)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(repos.ErrLocationNameExists))
		})
	})

		Describe("Delete", func() {
		It("should soft delete a location successfully", func() {
			locationToDelete := &types.Location{
				CompanyID: company.ID,
				AddressID: address.ID,
				Name:      "To Be Deleted",
			}
			err := repo.Create(ctx, locationToDelete)
			Expect(err).NotTo(HaveOccurred())

			err = repo.Delete(ctx, locationToDelete.ID)
			Expect(err).NotTo(HaveOccurred())

			_, found, err := repo.Get(ctx, company.ID, locationToDelete.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())

			// The location should be found with GetIncludeInvisible
			invisibleLocation, found, err := repo.GetIncludeInvisible(ctx, locationToDelete.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(invisibleLocation.Visible).To(BeFalse())
		})
	})

	Describe("CountByCompanyID", func() {
		It("should return the correct count of locations for a company", func() {
			Expect(repo.Create(ctx, &types.Location{CompanyID: company.ID, AddressID: address.ID, Name: "Count-1"})).To(Succeed())
			Expect(repo.Create(ctx, &types.Location{CompanyID: company.ID, AddressID: address.ID, Name: "Count-2"})).To(Succeed())

			company2 := &types.Company{Name: "Count Company 2", AddressID: address.ID}
			Expect(gr.Companies().Create(ctx, company2)).To(Succeed())
			Expect(repo.Create(ctx, &types.Location{CompanyID: company2.ID, AddressID: address.ID, Name: "Count-3"})).To(Succeed())

			count, err := repo.CountByCompanyID(ctx, company.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(2)))

			count, err = repo.CountByCompanyID(ctx, company2.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(1)))
		})

		It("should return 0 for a company with no locations", func() {
			company3 := &types.Company{Name: "No Locations Co", AddressID: address.ID}
			Expect(gr.Companies().Create(ctx, company3)).To(Succeed())

			count, err := repo.CountByCompanyID(ctx, company3.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(0)))
		})
	})

	Describe("Find", func() {
		var (
			company2 *types.Company
			address2 *types.Address
			location1 *types.Location
			location2 *types.Location
			location3 *types.Location
		)

		BeforeEach(func() {
			// Create additional company and address for diverse test data
			address2 = &types.Address{
				Line1:      "789 Test Ave",
				City:       "Testville",
				State:      "TS",
				PostalCode: "98765",
				Country:    "USA",
			}
			address2, err := gr.Addresses().Create(ctx, address2)
			Expect(err).NotTo(HaveOccurred())

			company2 = &types.Company{Name: "Another Test Company", AddressID: address2.ID}
			err = gr.Companies().Create(ctx, company2)
			Expect(err).NotTo(HaveOccurred())

			// Create test locations
			location1 = &types.Location{CompanyID: company.ID, AddressID: address.ID, Name: "Main Office"}
			Expect(repo.Create(ctx, location1)).To(Succeed())

			location2 = &types.Location{CompanyID: company.ID, AddressID: address2.ID, Name: "Branch Office A"}
			Expect(repo.Create(ctx, location2)).To(Succeed())

			location3 = &types.Location{CompanyID: company2.ID, AddressID: address.ID, Name: "Remote Site"}
			Expect(repo.Create(ctx, location3)).To(Succeed())
		})

		It("should find locations by name using LIKE (case-insensitive and partial match)", func() {
			foundLocations, count, err := repo.Find(ctx, &repos.LocationFindOpts{Names: []string{"office"}})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(2)))
			Expect(foundLocations).To(HaveLen(2))
			Expect(foundLocations).To(ConsistOf(
				WithTransform(func(l *types.Location) string { return l.Name }, Equal(location1.Name)),
				WithTransform(func(l *types.Location) string { return l.Name }, Equal(location2.Name)),
			))

			foundLocations, count, err = repo.Find(ctx, &repos.LocationFindOpts{Names: []string{"remote"}})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(1)))
			Expect(foundLocations).To(HaveLen(1))
			Expect(foundLocations[0].Name).To(Equal(location3.Name))
		})

		It("should find locations by IDs", func() {
			foundLocations, count, err := repo.Find(ctx, &repos.LocationFindOpts{IDs: []int64{location1.ID, location3.ID}})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(2)))
			Expect(foundLocations).To(HaveLen(2))
			Expect(foundLocations).To(ConsistOf(
				WithTransform(func(l *types.Location) int64 { return l.ID }, Equal(location1.ID)),
				WithTransform(func(l *types.Location) int64 { return l.ID }, Equal(location3.ID)),
			))
		})

		It("should find locations by CompanyIDs", func() {
			foundLocations, count, err := repo.Find(ctx, &repos.LocationFindOpts{CompanyIDs: []int64{company.ID}})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(2)))
			Expect(foundLocations).To(HaveLen(2))
			Expect(foundLocations).To(ConsistOf(
				WithTransform(func(l *types.Location) int64 { return l.ID }, Equal(location1.ID)),
				WithTransform(func(l *types.Location) int64 { return l.ID }, Equal(location2.ID)),
			))
		})

		It("should find locations by AddressIDs", func() {
			foundLocations, count, err := repo.Find(ctx, &repos.LocationFindOpts{AddressIDs: []int64{address2.ID}})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(1)))
			Expect(foundLocations).To(HaveLen(1))
			Expect(foundLocations[0].ID).To(Equal(location2.ID))
		})

		It("should find locations by a combination of filters", func() {
			foundLocations, count, err := repo.Find(ctx, &repos.LocationFindOpts{
				CompanyIDs: []int64{company.ID},
				Names:      []string{"branch"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(1)))
			Expect(foundLocations).To(HaveLen(1))
			Expect(foundLocations[0].ID).To(Equal(location2.ID))
		})

		It("should apply limit and offset", func() {
			// Create more locations to test pagination
			for i := 0; i < 5; i++ {
				Expect(repo.Create(ctx, &types.Location{CompanyID: company.ID, AddressID: address.ID, Name: "Paginated Location " + string(rune('A'+i))})).To(Succeed())
			}

			foundLocations, count, err := repo.Find(ctx, &repos.LocationFindOpts{Limit: 2, Offset: 0})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(BeNumerically(">=", int64(2)))
			Expect(foundLocations).To(HaveLen(2))

			foundLocations, count, err = repo.Find(ctx, &repos.LocationFindOpts{Limit: 2, Offset: 2})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(BeNumerically(">=", int64(2)))
			Expect(foundLocations).To(HaveLen(2))
		})
	})
})