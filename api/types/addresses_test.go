package types_test

import (
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Address Validation", func() {
	var address *types.Address

	BeforeEach(func() {
		address = &types.Address{
			Line1:      "123 Main St",
			City:       "Anytown",
			State:      "CA",
			PostalCode: "12345",
			Country:    "USA",
		}
	})

	It("should not return an error for a valid address", func() {
		Expect(types.Validate(address)).To(Succeed())
	})

	It("should return an error if Line1 is missing", func() {
		address.Line1 = ""
		Expect(types.Validate(address)).NotTo(Succeed())
	})

	It("should return an error if City is missing", func() {
		address.City = ""
		Expect(types.Validate(address)).NotTo(Succeed())
	})

	It("should return an error if State is missing", func() {
		address.State = ""
		Expect(types.Validate(address)).NotTo(Succeed())
	})

	It("should return an error if PostalCode is missing", func() {
		address.PostalCode = ""
		Expect(types.Validate(address)).NotTo(Succeed())
	})

	It("should return an error if Country is missing", func() {
		address.Country = ""
		Expect(types.Validate(address)).NotTo(Succeed())
	})
})
