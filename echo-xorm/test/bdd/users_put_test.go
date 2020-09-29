package bdd_test

import (
	"encoding/json"
	"strconv"

	"github.com/corvinusz/echo-xorm/app/server/users"

	. "github.com/Benjamintf1/unmarshalledmatchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//------------------------------------------------------------------------------
func testPutUsers() {
	testData := initPutUsersTestData()
	Context("PUT /users with correct data", func() {
		for i := range testData {
			i := i // dont remove, golang closure variable workaround
			data := testData[i]
			It("should respond OK for "+data.Comment, func() {
				result := new(users.User)
				// take payload for TestData and make request
				resp, err := suite.rc.R().
					SetBody(data.JsonIn).
					SetResult(result).
					Put("/users/" + strconv.Itoa(data.ID))
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

func initPutUsersTestData() []TestData {
	return []TestData{
		// ok data
		{
			Comment: "full update",
			ID:      6,
			JsonIn: `{
					"email": "put_users_test_email_001@example.com",
					"password": "put_users_test_name_001",
					"displayName": "put_users_test_name_001",
					"passwordEtime": 100500600
				}`,
			JsonOut: `{
					"id": 6,
					"email": "put_users_test_email_001@example.com",
					"displayName": "put_users_test_name_001",
					"passwordEtime": 100500600
				 }`,
			HttpCode: 200,
		},
		{
			Comment:  "one field update",
			ID:       5,
			JsonIn:   `{"email": "put_users_test_email_002@example.com"}`,
			JsonOut:  `{ "id": 5, "email": "put_users_test_email_002@example.com"}`,
			HttpCode: 200,
		},
		{
			Comment:  "one field update",
			ID:       5,
			JsonIn:   `{"displayName": "put_users_test_name_002@example.com"}`,
			JsonOut:  `{ "id": 5, "displayName": "put_users_test_name_002@example.com" }`,
			HttpCode: 200,
		},
		{
			Comment:  "set multiple values to null",
			ID:       4,
			JsonIn:   `{"passwordEtime": 0, "displayName":""}`,
			JsonOut:  `{"id": 4, "passwordEtime": 0, "displayName":""}`,
			HttpCode: 200,
		},
		// errors check
		{
			Comment:  "bad json",
			ID:       6,
			JsonIn:   `this is not a json; ' select 1;`,
			HttpCode: 400,
		},
		{
			Comment:  "bad id",
			ID:       -6,
			JsonIn:   `{"displayName": "a_updated_test_user_003"}`,
			HttpCode: 400,
		},
		{
			Comment:  "non-existing user id",
			ID:       100506,
			JsonIn:   `{"displayName": "a_updated_test_user_003"}`,
			HttpCode: 404,
		},
		{
			Comment:  "always existing email",
			ID:       5,
			JsonIn:   `{"email":"put_users_test_email_001@example.com"}`,
			HttpCode: 409,
		},
	}
}
