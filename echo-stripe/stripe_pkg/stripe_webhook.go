package stripetest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"echo-stripe/common"

	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/webhook"
)

type OrgChargeInvoice struct {
	OrgSubID             int64     `xorm:"org_sub_id"`
	InvoiceID            int64     `xorm:"invoice_id"`
	OrganizationId       int64     `xorm:"organization_id"`
	StripeCustomerId     string    `xorm:"stripe_customer_id"`
	StripeSubscriptionId string    `xorm:"stripe_subscription_id"`
	SubscriptionStatus   string    `xorm:"subscription_status"`
	StripePriceID        string    `xorm:"stripe_price_id"`
	InvoicePaid          bool      `xorm:"invoice_paid"`
	CreatedAt            time.Time `xorm:"created_at"`
}

func StripeWebhook(c echo.Context) error {
	//cc := c.(*common.AppContext)

	reqBytes, err := ioutil.ReadAll(c.Request().Body)
	common.CheckError(err)

	event, err := webhook.ConstructEvent(reqBytes, c.Request().Header.Get("Stripe-Signature"), TestSecretKey)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "signature not authenticated")
	}

	utnNow := time.Now().UTC()
	fmt.Println(utnNow)

	if strings.EqualFold(event.Type, "price.created") {
		var price stripe.Price
		err := json.Unmarshal(event.Data.Raw, &price)
		common.CheckError(err)

		if strings.EqualFold(price.Product.ID, TestSecretKey) {

		}
	} else if strings.EqualFold(event.Type, "price.deleted") {
		var price stripe.Price
		err := json.Unmarshal(event.Data.Raw, &price)
		common.CheckError(err)

	} else if strings.EqualFold(event.Type, "price.updated") {
		var price stripe.Price

		err := json.Unmarshal(event.Data.Raw, &price)
		common.CheckError(err)

		if strings.EqualFold(price.Product.ID, TestSecretKey) {

		}
	} else if strings.EqualFold(event.Type, "invoice.payment_succeeded") {
		// invoice is paid successfully now update period end time
		var invoice stripe.Invoice
		err := json.Unmarshal(event.Data.Raw, &invoice)
		common.CheckError(err)

		// invoice belongs to subscription
		if invoice.Subscription == nil {

		}

	} else if strings.EqualFold(event.Type, "invoice.payment_failed") {
		// invoice is failed so keep in application to show the user
		var invoice stripe.Invoice
		err := json.Unmarshal(event.Data.Raw, &invoice)
		common.CheckError(err)

	} else if strings.EqualFold(event.Type, "customer.subscription.deleted") {
		// either user cancelled subscription or stripe cancelled after all payment attempts are failed
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		common.CheckError(err)

	} else if strings.EqualFold(event.Type, "customer.updated") {
		var customer stripe.Customer
		err := json.Unmarshal(event.Data.Raw, &customer)
		common.CheckError(err)

	}

	return c.JSON(http.StatusOK, "ok")
}
