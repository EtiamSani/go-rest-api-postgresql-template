package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
	"github.com/stripe/stripe-go/v78/customer"
)

var PriceId = os.Getenv("STRIPE_PRICE_ID")

func checkout(email string) (*stripe.CheckoutSession, error) {
	// discounts := []*stripe.CheckoutSessionDiscountParams{{Coupon: stripe.String("FMARC"),},}

	customerParams := &stripe.CustomerParams{
		Email: stripe.String(email),
	}
	customerParams.AddMetadata("FinalEmail", email)
	newCustomer, err := customer.New(customerParams)

	if err != nil {
		return nil, err
	}

	meta := map[string]string{
		"FinalEmail" : email,
	}

	log.Println("Creating meta for user: ", meta)

	params := &stripe.CheckoutSessionParams{
		Customer: &newCustomer.ID,
		SuccessURL: stripe.String(os.Getenv("STRIPE_SUCCESS_URL")),
		CancelURL: stripe.String(os.Getenv("STRIPE_CANCEL_URL")),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		// Discounts: discounts,
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(PriceId),
				Quantity: stripe.Int64(1),
			},
		},
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			TrialPeriodDays: stripe.Int64(7),
			Metadata: meta,
		},
	}
	return session.New(params)
}

type EmailInput struct {
	Email string `json:"email"`
}

type SessionOutput struct {
	RedirectUrl string `json:"redirectUrl"`
}

func CheckoutCreator(c *gin.Context){
	input := &EmailInput{}
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stripeSession, err := checkout(input.Email)
	if err != nil{
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &SessionOutput{RedirectUrl: stripeSession.URL})
}

func HandleEvent(w http.ResponseWriter, req * http.Request)  {
	event, err := getEvent(w, req)

	if err != nil{
		log.Fatal(err)
	}

	log.Println(event.Type)

	if event.Type == "customer.subscription.created" {
		c, err := customer.Get(event.Data.Object["customer"].(string), nil)
		if err != nil {
			log.Fatal(err)
		}
		email := c.Metadata["FinalEmail"]
		log.Println("Subscription created by", email)
	}


}

func getEvent(w http.ResponseWriter, req * http.Request) (eventRes * stripe.Event, err error){
	const MaxBodyBytes = int64(65536)
	req.Body = http.MaxBytesReader(w, req.Body, MaxBodyBytes)
	payload, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	event := stripe.Event{}
	err = json.Unmarshal(payload, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}