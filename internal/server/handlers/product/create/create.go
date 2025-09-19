package create_product

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/AlexMickh/coledzh-shop-backend/pkg/api"
	"github.com/AlexMickh/coledzh-shop-backend/pkg/logger"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Name        string   `validate:"required,min=3"`
	Description string   `validate:"required,min=3"`
	price       float32  `validate:"required,gt=0,lte=1000000"`
	CategoryIds []string `validate:"required"`
}

type Response struct {
	ID string `json:"id"`
}

type ProductCreator interface {
	CreateProduct(
		ctx context.Context,
		categoryIds []string,
		name string,
		description string,
		price float32,
		image []byte,
	) (string, error)
}

// New godoc
//
//	@Summary		create new product
//	@Description	create new product
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@Param			name		formData	string	true	"product name"
//	@Param			description	formData	string	true	"product description"
//	@Param			price		formData	number	true	"product price"
//	@Param			category_id	formData	string	true	"product category id"
//	@Param			image		formData	file	true	"product image"
//	@Success		201			{object}	Response
//	@Failure		400			{object}	api.ErrorResponse
//	@Failure		500			{object}	api.ErrorResponse
//	@Security		SessionAuth
//	@Router			/admin/create-product [post]
func New(validator *validator.Validate, productCreator ProductCreator) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		const op = "handlers.product.create.New"
		ctx := r.Context()
		log := logger.FromCtx(ctx).With(slog.String("op", op))

		name := r.FormValue("name")
		description := r.FormValue("description")
		priceStr := r.FormValue("price")
		price, err := strconv.ParseFloat(priceStr, 32)
		if err != nil {
			log.Error("failed convert price", logger.Err(err))
			return api.Error("failed to get price", http.StatusBadRequest)
		}
		categoryId := r.FormValue("category_id")
		image, _, err := r.FormFile("image")
		if err != nil {
			log.Error("failed to get image", logger.Err(err))
			return api.Error("failed to get image", http.StatusBadRequest)
		}

		categoryIds := strings.Split(categoryId, " ")

		req := Request{
			Name:        name,
			Description: description,
			price:       float32(price),
			CategoryIds: categoryIds,
		}
		if err = validator.Struct(&req); err != nil {
			log.Error("failed to validate request", logger.Err(err))
			return api.Error("failed to validate request", http.StatusBadRequest)
		}

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, image)
		if err != nil {
			log.Error("failed to process image", logger.Err(err))
			return api.Error("failed to process image", http.StatusBadRequest)
		}

		id, err := productCreator.CreateProduct(
			ctx,
			req.CategoryIds,
			req.Name,
			req.Description,
			req.price,
			buf.Bytes(),
		)
		if err != nil {
			log.Error("failed to create product", logger.Err(err))
			return api.Error("failed to create product", http.StatusInternalServerError)
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, Response{
			ID: id,
		})

		return nil
	}
}
