package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/nordew/UploadApp/internal/controller/http/dto"
	"github.com/stripe/stripe-go/v76"
	"net/http"
)

func (h *Handler) createPaymentIntent(c *gin.Context) {
	var payment dto.PaymentDTO

	if err := c.ShouldBindJSON(&payment); err != nil {
		invalidJSONResponse(c)
	}

	if payment.Amount < 1000 {
		writeErrorResponse(c, http.StatusBadRequest, "Error in amount", "Amount must be greater than 1$")
		return
	}

	params := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(payment.Amount),
		Currency:           stripe.String(string(stripe.CurrencyUSD)),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
	}

	pi, err := h.payment.NewProduct(params)
	if err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, "failed to create payment", err.Error())
	}

	response := gin.H{
		"secret": pi.ClientSecret,
	}

	writeResponse(c, http.StatusOK, response)
}
