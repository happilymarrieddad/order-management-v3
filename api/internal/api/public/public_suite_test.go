package public_test

import (
	"testing"

	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/order-management-v3/api/internal/api/public"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPublic(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Public API Suite")
}

var (
	router *mux.Router
)

var _ = BeforeSuite(func() {
	router = mux.NewRouter()
	public.AddPublicRoutes(router)
})
