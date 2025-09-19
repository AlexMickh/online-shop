package get_cart

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/AlexMickh/coledzh-shop-backend/internal/models"
	"github.com/AlexMickh/coledzh-shop-backend/pkg/api"
	"github.com/AlexMickh/coledzh-shop-backend/pkg/logger"
	"github.com/go-chi/render"
)

type Response struct {
	ID       string        `json:"id"`
	Price    float32       `json:"price"`
	Products []productInfo `json:"products"`
}

type productInfo struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	ImageUrl string  `json:"image_url"`
}

type CartProvider interface {
	CartByUserId(ctx context.Context, userId string) (models.Cart, error)
}

// New godoc
//
//	@Summary		returns users cart
//	@Description	returns users cart
//	@Tags			cart
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Response
//	@Failure		401	{object}	api.ErrorResponse
//	@Failure		500	{object}	api.ErrorResponse
//	@Security		SessionAuth
//	@Router			/cart [get]
func New(cartProvider CartProvider) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		const op = "handlers.cart.get.New"
		ctx := r.Context()
		log := logger.FromCtx(ctx).With(slog.String("op", op))

		userId, ok := ctx.Value("user_id").(string)
		if !ok {
			log.Error("failed to get user id")
			return api.Error("failed to get user id", http.StatusUnauthorized)
		}

		cart, err := cartProvider.CartByUserId(ctx, userId)
		if err != nil {
			log.Error("failed to get users cart", logger.Err(err))
			return api.Error("failed to get users cart", http.StatusInternalServerError)
		}

		productsInfo := make([]productInfo, 0, len(cart.Products))
		for _, product := range cart.Products {
			productInfo := productInfo{
				ID:       product.ID,
				Name:     product.Name,
				Price:    product.Price,
				ImageUrl: product.ImageUrl,
			}
			productsInfo = append(productsInfo, productInfo)
		}

		render.JSON(w, r, Response{
			ID:       cart.ID,
			Price:    cart.Price,
			Products: productsInfo,
		})

		return nil
	}
}
