package types_test

import (
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("User Validation", func() {
	var user *types.User

	BeforeEach(func() {
		user = &types.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Password:  "password123",
			CompanyID: 1,
			AddressID: 1,
		}
	})

	Context("when all fields are valid", func() {
		It("should not return an error", func() {
			Expect(types.Validate(user)).To(Succeed())
		})
	})

	Context("when FirstName is invalid", func() {
		It("should return an error for missing FirstName", func() {
			user.FirstName = ""
			Expect(types.Validate(user)).NotTo(Succeed())
		})

		It("should return an error for short FirstName", func() {
			user.FirstName = "J"
			Expect(types.Validate(user)).NotTo(Succeed())
		})
	})

	Context("when LastName is invalid", func() {
		It("should return an error for missing LastName", func() {
			user.LastName = ""
			Expect(types.Validate(user)).NotTo(Succeed())
		})

		It("should return an error for short LastName", func() {
			user.LastName = "D"
			Expect(types.Validate(user)).NotTo(Succeed())
		})
	})

	Context("when Email is invalid", func() {
		It("should return an error for missing Email", func() {
			user.Email = ""
			Expect(types.Validate(user)).NotTo(Succeed())
		})

		It("should return an error for invalid Email format", func() {
			user.Email = "not-an-email"
			Expect(types.Validate(user)).NotTo(Succeed())
		})
	})

	Context("when Password is invalid", func() {
		It("should return an error for missing Password", func() {
			user.Password = ""
			Expect(types.Validate(user)).NotTo(Succeed())
		})

		It("should return an error for short Password", func() {
			user.Password = "1234567"
			Expect(types.Validate(user)).NotTo(Succeed())
		})
	})

	Context("when ID fields are invalid", func() {
		It("should return an error for missing CompanyID", func() {
			user.CompanyID = 0
			Expect(types.Validate(user)).NotTo(Succeed())
		})

		It("should return an error for missing AddressID", func() {
			user.AddressID = 0
			Expect(types.Validate(user)).NotTo(Succeed())
		})
	})
})
