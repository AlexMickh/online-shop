package get_product_by_id

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
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       float32    `json:"price"`
	ImageUrl    string     `json:"image"`
	Categories  []category `json:"categories"`
}

type category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProductProvider interface {
	ProductById(ctx context.Context, productId string) (models.Product, error)
}

// New godoc
//
//	@Summary		get product by id
//	@Description	get product by id
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			product_id	path		string	true	"product category id"
//	@Success		200			{object}	Response
//	@Failure		400			{object}	api.ErrorResponse
//	@Failure		500			{object}	api.ErrorResponse
//	@Router			/products/{id} [get]
func New(productProvider ProductProvider) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		const op = "handlers.product.get_by_id.New"
		ctx := r.Context()
		log := logger.FromCtx(ctx).With(slog.String("op", op))

		productId := r.PathValue("id")

		product, err := productProvider.ProductById(ctx, productId)
		if err != nil {
			log.Error("failed to get product", logger.Err(err))
			return api.Error("failed to get product", http.StatusInternalServerError)
		}

		categories := make([]category, 0, len(product.Categories))
		for _, categoryItem := range product.Categories {
			c := category{
				ID:   categoryItem.ID,
				Name: categoryItem.Name,
			}
			categories = append(categories, c)
		}

		render.JSON(w, r, Response{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			ImageUrl:    product.ImageUrl,
			Categories:  categories,
		})

		return nil
	}
}
