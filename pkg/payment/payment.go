package payment

import "github.com/stripe/stripe-go/v76"

type Payment interface {
	NewProduct(params *stripe.PaymentIntentParams) (*stripe.PaymentIntent, error)
}
