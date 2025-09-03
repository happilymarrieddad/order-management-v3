package repos_test

import (
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UsersRepo", func() {
	var (
		repo    repos.UsersRepo
		company *types.Company
		address *types.Address
	)

	BeforeEach(func() {
		// Assuming 'gr' (global repo), 'db', and 'ctx' are initialized in the _suite_test.go file
		repo = gr.Users()

		// Create a company and address for each test to ensure isolation
		// and avoid assuming existing IDs.
		companyRepo := gr.Companies()
		addressRepo := gr.Addresses()

		address = &types.Address{
			Line1:      "123 Main St",
			City:       "Anytown",
			State:      "CA",
			Country:    "USA",
			PostalCode: "12345",
		}
		address, err := addressRepo.Create(ctx, address)
		Expect(err).NotTo(HaveOccurred())

		company = &types.Company{Name: "Test Co Inc.", AddressID: address.ID}
		err = companyRepo.Create(ctx, company)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Create and Get", func() {
		It("should create a new user with a single role", func() {
			user := &types.User{
				CompanyID: company.ID,
				Email:     "single.role@example.com",
				Password:  "password123",
				AddressID: address.ID,
				FirstName: "Single",
				LastName:  "Role",
				Roles:     types.Roles{types.RoleUser},
			}

			err := repo.Create(ctx, user)
			Expect(err).NotTo(HaveOccurred())
			Expect(user.ID).NotTo(BeZero())

			retrievedUser, found, err := repo.Get(ctx, user.CompanyID, user.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(retrievedUser.Email).To(Equal(user.Email))
			Expect(retrievedUser.FirstName).To(Equal(user.FirstName))
			Expect(retrievedUser.LastName).To(Equal(user.LastName))
			Expect(retrievedUser.Roles).To(Equal(types.Roles{types.RoleUser}))
			Expect(retrievedUser.Address).NotTo(BeNil())
			Expect(retrievedUser.Address.ID).To(Equal(address.ID))
			Expect(retrievedUser.Address.Line1).To(Equal(address.Line1))
			Expect(retrievedUser.Address.City).To(Equal(address.City))
			Expect(retrievedUser.Address.State).To(Equal(address.State))
			Expect(retrievedUser.Address.Country).To(Equal(address.Country))
			Expect(retrievedUser.Address.PostalCode).To(Equal(address.PostalCode))
		})

		It("should create a new user with multiple roles", func() {
			user := &types.User{
				CompanyID: company.ID,
				Email:     "multi.role@example.com",
				Password:  "password123",
				AddressID: address.ID,
				FirstName: "Multi",
				LastName:  "Role",
				Roles:     types.Roles{types.RoleAdmin, types.RoleUser},
			}

			err := repo.Create(ctx, user)
			Expect(err).NotTo(HaveOccurred())
			Expect(user.ID).NotTo(BeZero())

			retrievedUser, found, err := repo.Get(ctx, user.CompanyID, user.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(retrievedUser.Email).To(Equal(user.Email))
			Expect(retrievedUser.FirstName).To(Equal(user.FirstName))
			Expect(retrievedUser.LastName).To(Equal(user.LastName))
			Expect(retrievedUser.Roles).To(ConsistOf(types.RoleAdmin, types.RoleUser))
		})

		It("should create a new user with no roles", func() {
			user := &types.User{
				CompanyID: company.ID,
				Email:     "no.role@example.com",
				Password:  "password123",
				AddressID: address.ID,
				FirstName: "No",
				LastName:  "Role",
				Roles:     types.Roles{}, // Explicitly empty
			}

			err := repo.Create(ctx, user)
			Expect(err).NotTo(HaveOccurred())
			Expect(user.ID).NotTo(BeZero())

			retrievedUser, found, err := repo.Get(ctx, user.CompanyID, user.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(retrievedUser.Email).To(Equal(user.Email))
			Expect(retrievedUser.FirstName).To(Equal(user.FirstName))
		})

		It("should not get a user from another company", func() {
			user := &types.User{
				CompanyID: company.ID,
				Email:     "other.company@example.com",
				Password:  "password123",
				AddressID: address.ID,
				FirstName: "Other",
				LastName:  "Company",
				Roles:     types.Roles{types.RoleUser},
			}
			err := repo.Create(ctx, user)
			Expect(err).NotTo(HaveOccurred())

			// Attempt to get the user with a different company ID
			_, found, err := repo.Get(ctx, user.CompanyID+999, user.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())
		})
	})

	Describe("Update", func() {
		var existingUser *types.User

		BeforeEach(func() {
			// Create a user to be updated in the test
			existingUser = &types.User{
				CompanyID: company.ID,
				Email:     "update.me@example.com",
				Password:  "password123",
				AddressID: address.ID,
				FirstName: "Update",
				LastName:  "Me",
				Roles:     types.Roles{types.RoleUser},
			}
			err := repo.Create(ctx, existingUser)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should update a user's fields, including roles", func() {
			existingUser.FirstName = "Updated"
			existingUser.Roles = types.Roles{types.RoleUser, types.RoleAdmin}
			err := repo.Update(ctx, existingUser)
			Expect(err).NotTo(HaveOccurred())

			updatedUser, found, err := repo.Get(ctx, existingUser.CompanyID, existingUser.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(updatedUser.FirstName).To(Equal("Updated"))
			Expect(updatedUser.Roles).To(ConsistOf(types.RoleUser, types.RoleAdmin))
		})
	})

	Describe("Find", func() {
		var otherCompany *types.Company

		BeforeEach(func() {
			// Create a second company for isolation testing
			companyRepo := gr.Companies()
			otherCompany = &types.Company{Name: "Other Test Co", AddressID: address.ID}
			Expect(companyRepo.Create(ctx, otherCompany)).To(Succeed())

			// Create a set of users with different roles and companies
			users := []*types.User{
				{CompanyID: company.ID, FirstName: "Find", LastName: "Admin", Email: "find.admin@test.com", Password: "password123", AddressID: address.ID, Roles: types.Roles{types.RoleAdmin}},
				{CompanyID: company.ID, FirstName: "Find", LastName: "User", Email: "find.user@test.com", Password: "password123", AddressID: address.ID, Roles: types.Roles{types.RoleUser}},
				{CompanyID: otherCompany.ID, FirstName: "Find", LastName: "Other", Email: "find.other@test.com", Password: "password123", AddressID: address.ID, Roles: types.Roles{}},
			}
			for _, u := range users {
				Expect(repo.Create(ctx, u)).To(Succeed())
			}
		})

		It("should find all users when no company is specified", func() {
			foundUsers, count, err := repo.Find(ctx, &repos.UserFindOpts{})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(3)))
			Expect(foundUsers).To(HaveLen(3))
		})

		It("should find only users belonging to the specified company", func() {
			foundUsers, count, err := repo.Find(ctx, &repos.UserFindOpts{CompanyID: company.ID})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(2)))
			Expect(foundUsers).To(HaveLen(2))

			// Verify roles for each user
			for _, u := range foundUsers {
				Expect(u.CompanyID).To(Equal(company.ID))
			}
		})

		It("should find users by email using LIKE (case-insensitive and partial match)", func() {
			// Create a user with a specific email for testing
			user1 := &types.User{CompanyID: company.ID, Email: "test.user.one@example.com", Password: "password123", AddressID: address.ID, FirstName: "Test", LastName: "UserOne", Roles: types.Roles{types.RoleUser}}
			Expect(repo.Create(ctx, user1)).To(Succeed())
			user2 := &types.User{CompanyID: company.ID, Email: "another.user@example.com", Password: "password123", AddressID: address.ID, FirstName: "Another", LastName: "User", Roles: types.Roles{types.RoleUser}}
			Expect(repo.Create(ctx, user2)).To(Succeed())

			// Test partial and case-insensitive email search
			foundUsers, count, err := repo.Find(ctx, &repos.UserFindOpts{Emails: []string{"user.one"}})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(1)))
			Expect(foundUsers).To(HaveLen(1))
			Expect(foundUsers[0].Email).To(Equal(user1.Email))

			foundUsers, count, err = repo.Find(ctx, &repos.UserFindOpts{Emails: []string{"EXAMPLE.COM"}})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(2))) // Should find both test.user.one and another.user
			Expect(foundUsers).To(HaveLen(2))
		})

		It("should find users by first name using LIKE (case-insensitive and partial match)", func() {
			user1 := &types.User{CompanyID: company.ID, Email: "fname1@example.com", Password: "password123", AddressID: address.ID, FirstName: "FirstNameTest", LastName: "User", Roles: types.Roles{types.RoleUser}}
			Expect(repo.Create(ctx, user1)).To(Succeed())
			user2 := &types.User{CompanyID: company.ID, Email: "fname2@example.com", Password: "password123", AddressID: address.ID, FirstName: "AnotherFirstName", LastName: "User", Roles: types.Roles{types.RoleUser}}
			Expect(repo.Create(ctx, user2)).To(Succeed())

			foundUsers, count, err := repo.Find(ctx, &repos.UserFindOpts{FirstNames: []string{"firstnametest"}})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(1)))
			Expect(foundUsers).To(HaveLen(1))
			Expect(foundUsers[0].FirstName).To(Equal(user1.FirstName))

			foundUsers, count, err = repo.Find(ctx, &repos.UserFindOpts{FirstNames: []string{"another"}})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(1)))
			Expect(foundUsers).To(HaveLen(1))
			Expect(foundUsers[0].FirstName).To(Equal(user2.FirstName))
		})

		It("should find users by last name using LIKE (case-insensitive and partial match)", func() {
			user1 := &types.User{CompanyID: company.ID, Email: "lname1@example.com", Password: "password123", AddressID: address.ID, FirstName: "User", LastName: "LastNameTest", Roles: types.Roles{types.RoleUser}}
			Expect(repo.Create(ctx, user1)).To(Succeed())
			user2 := &types.User{CompanyID: company.ID, Email: "lname2@example.com", Password: "password123", AddressID: address.ID, FirstName: "User", LastName: "AnotherLastName", Roles: types.Roles{types.RoleUser}}
			Expect(repo.Create(ctx, user2)).To(Succeed())

			foundUsers, count, err := repo.Find(ctx, &repos.UserFindOpts{LastNames: []string{"lastnametest"}})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(1)))
			Expect(foundUsers).To(HaveLen(1))
			Expect(foundUsers[0].LastName).To(Equal(user1.LastName))

			foundUsers, count, err = repo.Find(ctx, &repos.UserFindOpts{LastNames: []string{"another"}})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(1)))
			Expect(foundUsers).To(HaveLen(1))
			Expect(foundUsers[0].LastName).To(Equal(user2.LastName))
		})

		It("should find users by IDs", func() {
			user1 := &types.User{CompanyID: company.ID, Email: "idtest1@example.com", Password: "password123", AddressID: address.ID, FirstName: "ID", LastName: "Test1", Roles: types.Roles{types.RoleUser}}
			Expect(repo.Create(ctx, user1)).To(Succeed())
			user2 := &types.User{CompanyID: company.ID, Email: "idtest2@example.com", Password: "password123", AddressID: address.ID, FirstName: "ID", LastName: "Test2", Roles: types.Roles{types.RoleUser}}
			Expect(repo.Create(ctx, user2)).To(Succeed())

			foundUsers, count, err := repo.Find(ctx, &repos.UserFindOpts{IDs: []int64{user1.ID}})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(1)))
			Expect(foundUsers).To(HaveLen(1))
			Expect(foundUsers[0].ID).To(Equal(user1.ID))

			foundUsers, count, err = repo.Find(ctx, &repos.UserFindOpts{IDs: []int64{user1.ID, user2.ID}})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(2)))
			Expect(foundUsers).To(HaveLen(2))
			Expect(foundUsers).To(ConsistOf(
				WithTransform(func(u *types.User) int64 { return u.ID }, Equal(user1.ID)),
				WithTransform(func(u *types.User) int64 { return u.ID }, Equal(user2.ID)),
			))
		})
	})

	Describe("UpdateUserCompany", func() {
		var (
			userToMove *types.User
			newCompany *types.Company
			newAddress *types.Address
		)

		BeforeEach(func() {
			// This user starts in the global `company`
			userToMove = &types.User{
				CompanyID: company.ID,
				Email:     "move.me@example.com",
				Password:  "password123",
				AddressID: address.ID,
				FirstName: "Move",
				LastName:  "Me",
				Roles:     types.Roles{types.RoleUser},
			}
			err := repo.Create(ctx, userToMove)
			Expect(err).NotTo(HaveOccurred())

			// Create a second company to move the user to
			addressRepo := gr.Addresses()
			newAddress = &types.Address{
				Line1:      "456 New Ave",
				City:       "Newville",
				State:      "TX",
				Country:    "USA",
				PostalCode: "54321",
			}
			newAddress, err = addressRepo.Create(ctx, newAddress)
			Expect(err).NotTo(HaveOccurred())

			companyRepo := gr.Companies()
			newCompany = &types.Company{Name: "New Company LLC", AddressID: newAddress.ID}
			err = companyRepo.Create(ctx, newCompany)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should move a user to a new company", func() {
			// Move the user
			err := repo.UpdateUserCompany(ctx, userToMove.ID, newCompany.ID)
			Expect(err).NotTo(HaveOccurred())

			// 1. Check if user is findable in the NEW company
			foundUser, found, err := repo.Get(ctx, newCompany.ID, userToMove.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(foundUser.CompanyID).To(Equal(newCompany.ID))

			// 2. Check if user is NO LONGER findable in the OLD company
			_, found, err = repo.Get(ctx, company.ID, userToMove.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())
		})
	})

	Describe("Delete", func() {
		var (
			userToDelete *types.User
		)

		BeforeEach(func() {
			userToDelete = &types.User{
				CompanyID: company.ID,
				Email:     "delete.me@example.com",
				Password:  "password123",
				AddressID: address.ID,
				FirstName: "Delete",
				LastName:  "Me",
				Roles:     types.Roles{types.RoleUser},
			}
			err := repo.Create(ctx, userToDelete)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should soft delete a user", func() {
			err := repo.Delete(ctx, userToDelete.ID)
			Expect(err).NotTo(HaveOccurred())

			// The user should not be found with Get
			_, found, err := repo.Get(ctx, userToDelete.CompanyID, userToDelete.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())

			// The user should be found with GetIncludeInvisible
			invisibleUser, found, err := repo.GetIncludeInvisible(ctx, userToDelete.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(invisibleUser.Visible).To(BeFalse())
		})
	})
})
