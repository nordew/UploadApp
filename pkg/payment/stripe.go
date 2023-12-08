package payment

import (
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
)

type stripePayment struct {
}

func NewPayement() Payment {
	return &stripePayment{}
}

func (s *stripePayment) NewProduct(params *stripe.PaymentIntentParams) (*stripe.PaymentIntent, error) {
	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, err
	}

	return pi, nil
}
