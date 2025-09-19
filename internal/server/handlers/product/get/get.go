package get_product

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/AlexMickh/coledzh-shop-backend/internal/models"
	"github.com/AlexMickh/coledzh-shop-backend/pkg/api"
	"github.com/AlexMickh/coledzh-shop-backend/pkg/logger"
	"github.com/go-chi/render"
)

type Response struct {
	Products []productInfo `json:"products"`
}

type productInfo struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	ImageUrl string  `json:"image_url"`
}

type ProductProvider interface {
	ProductsCard(ctx context.Context, categoryId string, page int) ([]models.ProductCard, error)
}

// New godoc
//
//	@Summary		get products
//	@Description	get products
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			category_id	query		string	false	"product category id"
//	@Param			page		query		int		true	"page for pagination"
//	@Success		200			{object}	Response
//	@Failure		400			{object}	api.ErrorResponse
//	@Failure		500			{object}	api.ErrorResponse
//	@Router			/products [get]
func New(productProvider ProductProvider) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		const op = "handlers.product.get.New"
		ctx := r.Context()
		log := logger.FromCtx(ctx).With(slog.String("op", op))

		pageStr := r.URL.Query().Get("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			log.Error("failed to convert page", logger.Err(err))
			return api.Error("page must be int", http.StatusBadRequest)
		}
		if page < 0 {
			log.Error("page is negative number")
			return api.Error("page must be non negative", http.StatusBadRequest)
		}

		categoryId := r.URL.Query().Get("category_id")
		products, err := productProvider.ProductsCard(ctx, categoryId, page)
		if err != nil {
			log.Error("failed to get products", logger.Err(err))
			return api.Error("failed to get products", http.StatusInternalServerError)
		}

		productsInfo := make([]productInfo, 0, len(products))
		for _, product := range products {
			productInfo := productInfo{
				ID:       product.ID,
				Name:     product.Name,
				Price:    product.Price,
				ImageUrl: product.ImageUrl,
			}
			productsInfo = append(productsInfo, productInfo)
		}

		render.JSON(w, r, Response{
			Products: productsInfo,
		})

		return nil
	}
}
