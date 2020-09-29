package bdd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/corvinusz/echo-xorm/app/server/version"
)

var _ = Describe("Test /version", func() {
	Context("GET /version", func() {
		It("should respond properly", func() {
			result := new(version.Result)
			resp, err := suite.rc.R().SetResult(result).Get("/version")
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(200))
			Expect(result.Result).To(Equal("OK"))
			Expect(result.Version).To(Equal(suite.app.Ctx.Config.Version))
			Expect(result.ServerTime).NotTo(BeZero())
		})
	})
})
