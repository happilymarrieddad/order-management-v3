package types_test

import (
	"testing"

	"github.com/happilymarrieddad/order-management-v3/api/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTypes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Types Suite")
}

var _ = Describe("Roles", func() {
	Describe("ToDB", func() {
		It("should convert a slice of roles to a DB string", func() {
			roles := types.Roles{types.RoleAdmin, types.RoleUser}
			dbBytes, err := roles.ToDB()
			Expect(err).NotTo(HaveOccurred())
			Expect(string(dbBytes)).To(Equal("{admin,user}"))
		})

		It("should handle an empty slice", func() {
			roles := types.Roles{}
			dbBytes, err := roles.ToDB()
			Expect(err).NotTo(HaveOccurred())
			Expect(string(dbBytes)).To(Equal("{}"))
		})

		It("should handle a nil slice", func() {
			var roles types.Roles
			dbBytes, err := roles.ToDB()
			Expect(err).NotTo(HaveOccurred())
			Expect(string(dbBytes)).To(Equal("{}"))
		})
	})

	Describe("FromDB", func() {
		It("should convert a DB string to a slice of roles", func() {
			var roles types.Roles
			dbString := []byte("{admin,user}")
			err := roles.FromDB(dbString)
			Expect(err).NotTo(HaveOccurred())
			Expect(roles).To(ConsistOf(types.RoleAdmin, types.RoleUser))
		})

		It("should handle an empty DB string", func() {
			var roles types.Roles
			dbString := []byte("{}")
			err := roles.FromDB(dbString)
			Expect(err).NotTo(HaveOccurred())
			Expect(roles).To(BeEmpty())
		})

		It("should return an error for an invalid role", func() {
			var roles types.Roles
			dbString := []byte("{admin,invalidrole}")
			err := roles.FromDB(dbString)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("HasRole", func() {
		It("should return true if the role exists", func() {
			roles := types.Roles{types.RoleAdmin, types.RoleUser}
			Expect(roles.HasRole(types.RoleAdmin)).To(BeTrue())
		})

		It("should return false if the role does not exist", func() {
			roles := types.Roles{types.RoleUser}
			Expect(roles.HasRole(types.RoleAdmin)).To(BeFalse())
		})
	})
})
