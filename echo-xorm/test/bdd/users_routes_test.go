package bdd_test

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/corvinusz/echo-xorm/app/server/users"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test /users", func() {
	Context("Test GET /users", testGetUsers)
	Context("Test POST /users", testPostUsers)
	Context("Test PUT /users", testPutUsers)
	// Context("Test DELETE /users", testDeleteUsers)
})

func testGetUsers() {
	Context("Get all users", func() {
		It("should respond properly", func() {
			var fromDb, result []users.User
			// get fromDb
			err := suite.app.Ctx.Orm.Omit("password").Find(&fromDb)
			Expect(err).NotTo(HaveOccurred())
			// get resp
			resp, err := suite.rc.R().SetResult(&result).Get("/users")
			Expect(err).NotTo(HaveOccurred())
			Expect(http.StatusOK).To(Equal(resp.StatusCode()))
			Expect(len(fromDb)).To(BeNumerically(">=", 5))
			Expect(len(result)).To(Equal(len(fromDb)))
			Expect(result).To(BeEquivalentTo(fromDb))
		})
	})
	Context("GET /users/{id} with 3 random id", func() {
		It("should respond properly", func() {
			for i := 0; i < 3; i++ {
				id := rand.Int()%7 + 1
				fromDb := new(users.User)
				result := new(users.User)
				// get fromDb
				found, err := suite.app.Ctx.Orm.ID(id).Omit("password").Get(fromDb)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				// get resp
				resp, err := suite.rc.R().SetResult(result).Get("/users/" + strconv.Itoa(id))
				Expect(err).NotTo(HaveOccurred())
				Expect(http.StatusOK).To(Equal(resp.StatusCode()))
				Expect(result).To(BeEquivalentTo(fromDb))
			}
		})
	})
}
