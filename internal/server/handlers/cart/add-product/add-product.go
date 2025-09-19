package cart_add_product

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/AlexMickh/coledzh-shop-backend/pkg/api"
	"github.com/AlexMickh/coledzh-shop-backend/pkg/logger"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	ProductId string `json:"product_id" validate:"required,uuid4"`
}

type Response struct {
	ID string `json:"id"`
}

type ProductAdder interface {
	AddProduct(ctx context.Context, userId, productId string) (string, error)
}

// New godoc
//
//	@Summary		add product to cart
//	@Description	add product to cart
//	@Tags			cart
//	@Accept			json
//	@Produce		json
//	@Param			product_id	body		string	true	"product id"
//	@Success		201			{object}	Response
//	@Failure		400			{object}	api.ErrorResponse
//	@Failure		401			{object}	api.ErrorResponse
//	@Failure		500			{object}	api.ErrorResponse
//	@Security		SessionAuth
//	@Router			/cart/add [post]
func New(validator *validator.Validate, productAdder ProductAdder) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		const op = "handlers.cart.add-product.New"
		ctx := r.Context()
		log := logger.FromCtx(ctx).With(slog.String("op", op))

		var req Request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", logger.Err(err))
			return api.Error("failed to decode request body", http.StatusBadRequest)
		}
		defer r.Body.Close()

		if err := validator.Struct(&req); err != nil {
			log.Error("failed to validate request body", logger.Err(err))
			return api.Error("failed to validate request body", http.StatusBadRequest)
		}

		userId, ok := ctx.Value("user_id").(string)
		if !ok {
			log.Error("failed to get user id")
			return api.Error("failed to get user id", http.StatusUnauthorized)
		}

		cartId, err := productAdder.AddProduct(ctx, userId, req.ProductId)
		if err != nil {
			log.Error("failed to add product", logger.Err(err))
			return api.Error("failed to add product", http.StatusInternalServerError)
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, Response{
			ID: cartId,
		})

		return nil
	}
}
