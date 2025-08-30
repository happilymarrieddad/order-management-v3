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

			retrievedUser, found, err := repo.Get(ctx, user.ID)
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

			retrievedUser, found, err := repo.Get(ctx, user.ID)
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

			retrievedUser, found, err := repo.Get(ctx, user.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(retrievedUser.Email).To(Equal(user.Email))
			Expect(retrievedUser.FirstName).To(Equal(user.FirstName))
			Expect(retrievedUser.LastName).To(Equal(user.LastName))
			Expect(retrievedUser.Roles).To(BeEmpty())
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

			updatedUser, found, err := repo.Get(ctx, existingUser.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(updatedUser.FirstName).To(Equal("Updated"))
			Expect(updatedUser.Roles).To(ConsistOf(types.RoleUser, types.RoleAdmin))
		})
	})

	Describe("Find", func() {
		BeforeEach(func() {
			// Create a set of users with different roles
			users := []*types.User{
				{CompanyID: company.ID, FirstName: "Find", LastName: "Admin", Email: "find.admin@example.com", Password: "password123", AddressID: address.ID, Roles: types.Roles{types.RoleAdmin}},
				{CompanyID: company.ID, FirstName: "Find", LastName: "User", Email: "find.user@example.com", Password: "password123", AddressID: address.ID, Roles: types.Roles{types.RoleUser}},
				{CompanyID: company.ID, FirstName: "Find", LastName: "None", Email: "find.none@example.com", Password: "password123", AddressID: address.ID, Roles: types.Roles{}},
			}
			for _, u := range users {
				Expect(repo.Create(ctx, u)).To(Succeed())
			}
		})

		It("should find users and correctly populate their roles", func() {
			foundUsers, count, err := repo.Find(ctx, &repos.UserFindOpts{})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(3)))
			Expect(foundUsers).To(HaveLen(3))

			// Verify roles for each user
			for _, u := range foundUsers {
				switch u.Email {
				case "find.admin@example.com":
					Expect(u.Roles).To(Equal(types.Roles{types.RoleAdmin}))
				case "find.user@example.com":
					Expect(u.Roles).To(Equal(types.Roles{types.RoleUser}))
				case "find.none@example.com":
					Expect(u.Roles).To(BeEmpty())
				}
			}
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
			_, found, err := repo.Get(ctx, userToDelete.ID)
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