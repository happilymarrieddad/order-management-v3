package types_test

import (
	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Company Validation", func() {
	var company *types.Company

	BeforeEach(func() {
		company = &types.Company{
			Name:      "Test Company",
			AddressID: 1,
		}
	})

	It("should not return an error for a valid company", func() {
		Expect(types.Validate(company)).To(Succeed())
	})

	It("should return an error if Name is missing", func() {
		company.Name = ""
		Expect(types.Validate(company)).NotTo(Succeed())
	})

	It("should return an error if AddressID is missing", func() {
		company.AddressID = 0
		Expect(types.Validate(company)).NotTo(Succeed())
	})
})
