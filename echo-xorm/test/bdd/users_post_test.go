package bdd_test

import (
	"encoding/json"

	"github.com/corvinusz/echo-xorm/app/server/users"

	. "github.com/Benjamintf1/unmarshalledmatchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//------------------------------------------------------------------------------
func testPostUsers() {
	testData := initPostUsersTestData()
	Context("POST /users", func() {
		for i := range testData {
			i := i // dont remove, golang closure variable workaround
			data := testData[i]
			It("should respond OK for "+data.Comment, func() {
				result := new(users.User)
				// take payload for TestData and make request
				resp, err := suite.rc.R().
					SetBody(data.JsonIn).
					SetResult(result).
					Post("/users")
				Expect(err).NotTo(HaveOccurred())
				// verify response
				if data.HttpCode != 0 {
					Expect(resp.StatusCode()).To(Equal(data.HttpCode))
				}
				if data.JsonOut != "" {
					Expect(resp.Body()).Should(ContainUnorderedJSON(data.JsonOut))
				}
				// verify database
				if data.HaveToCheckDb {
					// extract object from storage
					fromDb := users.User{}
					found, err := suite.app.Ctx.Orm.ID(result.ID).Get(&fromDb)
					Expect(err).NotTo(HaveOccurred())
					Expect(found).To(BeTrue())
					// convert object to json
					jsonUser, err := json.Marshal(fromDb)
					Expect(err).NotTo(HaveOccurred())
					// verify that database contains expected data
					Expect(jsonUser).Should(ContainUnorderedJSON(data.JsonOut))
				}
			})
		}
	})
}

func initPostUsersTestData() []TestData {
	return []TestData{
		// correct responses
		{
			Comment: "full post",
			JsonIn: `{
			"email": "post_users_test_email_001@example.com",
			"password": "post_users_test_password_001",
			"displayName": "post_users_test_name_001",
			"passwordEtime":1
			}`,
			JsonOut: `{
				"email": "post_users_test_email_001@example.com",
				"displayName": "post_users_test_name_001",
				"passwordEtime":1
			}`,
			HttpCode: 201,
		},
		{
			Comment: "minimal post",
			JsonIn: `{
				"email": "post_users_test_email_002@example.com",
				"password": "post_users_test_password_002"
			}`,
			JsonOut: `{
				"email": "post_users_test_email_002@example.com"
			}`,
			HttpCode: 201,
		},
		// error checks
		{
			Comment:  "bad payload",
			JsonIn:   `this is not a json; ' select 1;`,
			HttpCode: 400,
		},
		{
			Comment: "always existing email",
			JsonIn: `{
				"email": "admin",
				"password":"adminx"
			}`,
			HttpCode: 409,
		},
		{
			Comment: "deficient payload",
			JsonIn: `{
				"displayName": "post_users_test_name_003",
				"password": "post_users_test_password"
			}`,
			HttpCode: 400,
		},
		{
			Comment: "invalid email",
			JsonIn: `{
				"email": "post_users_test_email_003",
				"password": "post_users_test_password"
			}`,
			HttpCode: 400,
		},
		{
			Comment: "short password",
			JsonIn: `{
				"email": "post_users_test_email_003@example.com",
				"password": "123"
			}`,
			HttpCode: 400,
		},
	}

}
