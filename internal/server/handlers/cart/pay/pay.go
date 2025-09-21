package pay_cart

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"net/http"

	"github.com/AlexMickh/coledzh-shop-backend/pkg/api"
	"github.com/AlexMickh/coledzh-shop-backend/pkg/logger"
	"github.com/go-chi/render"
	"github.com/rvinnie/yookassa-sdk-go/yookassa"
	yoocommon "github.com/rvinnie/yookassa-sdk-go/yookassa/common"
	yoopayment "github.com/rvinnie/yookassa-sdk-go/yookassa/payment"
)

type Response struct {
	PaymentId   string `josn:"payment_id"`
	RedirectURL string `json:"redirect_url"`
}

type CartProvider interface {
	CartPriceByUserId(ctx context.Context, userId string) (float32, error)
}

// Pay godoc
//
//	@Summary		returns users cart
//	@Description	returns users cart
//	@Tags			cart
//	@Accept			json
//	@Produce		json
//	@Success		201	{object}	Response
//	@Failure		401	{object}	api.ErrorResponse
//	@Failure		500	{object}	api.ErrorResponse
//	@Security		SessionAuth
//	@Router			/cart/pay [post]
func Pay(paymentHandler *yookassa.PaymentHandler, cartProvider CartProvider) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		const op = "handlers.cart.pay.Pay"
		ctx := r.Context()
		log := logger.FromCtx(ctx).With(slog.String("op", op))

		userId, ok := ctx.Value("user_id").(string)
		if !ok {
			log.Error("failed to get user id")
			return api.Error("failed to get user id", http.StatusUnauthorized)
		}

		price, err := cartProvider.CartPriceByUserId(ctx, userId)
		if err != nil {
			log.Error("failed to get cart price", logger.Err(err))
			return api.Error("failed to get cart price", http.StatusInternalServerError)
		}
		price = roundFloat(price, 2)

		payment, err := paymentHandler.CreatePayment(&yoopayment.Payment{
			Amount: &yoocommon.Amount{
				Value:    fmt.Sprint(price),
				Currency: "RUB",
			},
			PaymentMethod: yoopayment.PaymentTypeBankCard,
			Confirmation: yoopayment.Redirect{
				Type:      "redirect",
				ReturnURL: "https://www.google.com",
			},
			Description: "Test payment",
			Metadata: map[string]string{
				"user_id": userId,
			},
		})
		if err != nil {
			log.Error("failed to create payment", logger.Err(err))
			return api.Error("failed to create payment", http.StatusInternalServerError)
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, Response{
			PaymentId:   payment.ID,
			RedirectURL: payment.Confirmation.(map[string]interface{})["confirmation_url"].(string),
		})

		return nil
	}
}

type YookassaRequest struct {
	Event  string `json:"event"`
	Object object `json:"object"`
}

type object struct {
	Paid     bool           `json:"paid"`
	Metadata map[string]any `json:"metadata"`
}

type CartDeleter interface {
	DeleteCartByUserId(ctx context.Context, userId string) error
}

func Webhook(cartDeleter CartDeleter) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		const op = "handlers.cart.pay.Pay"
		ctx := r.Context()
		log := logger.FromCtx(ctx).With(slog.String("op", op))

		var req YookassaRequest
		if err := render.Decode(r, &req); err != nil {
			log.Error("failed to decode request", logger.Err(err))
			return api.Error("failed to decode request", http.StatusBadRequest)
		}

		if req.Event == "payment.waiting_for_capture" && req.Object.Paid {
			userId, ok := req.Object.Metadata["user_id"].(string)
			if !ok {
				log.Error("failed to get user id from metadata")
				return api.Error("failed to get user id from metadata", http.StatusBadRequest)
			}

			err := cartDeleter.DeleteCartByUserId(ctx, userId)
			if err != nil {
				log.Error("failed to delete cart", logger.Err(err))
				return api.Error("failed to delete cart", http.StatusInternalServerError)
			}

			render.Status(r, http.StatusOK)
		}

		render.Status(r, http.StatusInternalServerError)

		return nil
	}
}

func roundFloat(val float32, precision uint) float32 {
	ratio := math.Pow(10, float64(precision))
	return float32(math.Round(float64(val)*ratio) / ratio)
}
