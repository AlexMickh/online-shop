package server

import (
	"context"
	"fmt"
	"net/http"

	_ "github.com/AlexMickh/coledzh-shop-backend/docs"
	"github.com/AlexMickh/coledzh-shop-backend/internal/config"
	"github.com/AlexMickh/coledzh-shop-backend/internal/lib/email"
	"github.com/AlexMickh/coledzh-shop-backend/internal/models"
	"github.com/AlexMickh/coledzh-shop-backend/internal/server/handlers/auth/login"
	"github.com/AlexMickh/coledzh-shop-backend/internal/server/handlers/auth/register"
	"github.com/AlexMickh/coledzh-shop-backend/internal/server/handlers/auth/verify"
	cart_add_product "github.com/AlexMickh/coledzh-shop-backend/internal/server/handlers/cart/add-product"
	get_cart "github.com/AlexMickh/coledzh-shop-backend/internal/server/handlers/cart/get"
	create_category "github.com/AlexMickh/coledzh-shop-backend/internal/server/handlers/category/create"
	get_category "github.com/AlexMickh/coledzh-shop-backend/internal/server/handlers/category/get"
	create_product "github.com/AlexMickh/coledzh-shop-backend/internal/server/handlers/product/create"
	get_product "github.com/AlexMickh/coledzh-shop-backend/internal/server/handlers/product/get"
	get_product_by_id "github.com/AlexMickh/coledzh-shop-backend/internal/server/handlers/product/get-by-id"
	"github.com/AlexMickh/coledzh-shop-backend/internal/server/middlewares"
	"github.com/AlexMickh/coledzh-shop-backend/pkg/api"
	"github.com/AlexMickh/coledzh-shop-backend/pkg/logger"
	"github.com/AlexMickh/coledzh-shop-backend/pkg/session"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Server struct {
	srv *http.Server
}

type AuthService interface {
	Register(ctx context.Context, login, email, password string) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type TokenService interface {
	CreateToken(ctx context.Context, userId, tokenType string) (string, error)
	VerifyEmail(ctx context.Context, token string) error
}

type CategoryService interface {
	CreateCategory(ctx context.Context, name string) (string, error)
	AllCategories(ctx context.Context) ([]models.Category, error)
}

type UserService interface {
	ValidateAdminSession(ctx context.Context, sessionId string) error
	ValidateUserSession(ctx context.Context, sessionId string) (string, error)
}

type ProductService interface {
	CreateProduct(
		ctx context.Context,
		categoryIds []string,
		name string,
		description string,
		price float32,
		image []byte,
	) (string, error)
	ProductsCard(ctx context.Context, categoryId string, page int) ([]models.ProductCard, error)
	ProductById(ctx context.Context, productId string) (models.Product, error)
}

type CartService interface {
	AddProduct(ctx context.Context, userId, productId string) (string, error)
	CartByUserId(ctx context.Context, userId string) (models.Cart, error)
}

// @title						Your API
// @version					1.0
// @description				Your API description
// @securityDefinitions.apikey	SessionAuth
// @in							cookie
// @name						session_id
func New(
	ctx context.Context,
	cfg config.ServerConfig,
	authService AuthService,
	mailCfg config.MailConfig,
	tokenService TokenService,
	categoryService CategoryService,
	userService UserService,
	productService ProductService,
	cartService CartService,
) *Server {
	const op = "server.New"

	r := chi.NewRouter()

	validator := validator.New()
	email := email.New(mailCfg)
	session := session.New("session_id", true, false, 60*60*24*5)

	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	}).Handler)

	r.Use(middleware.RequestID)
	r.Use(logger.ChiMiddleware(ctx))
	r.Use(middleware.Recoverer)
	// r.Use(middleware.URLFormat)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s/swagger/doc.json", cfg.Addr)), //The url pointing to API definition
	))

	r.Get("/health-check", api.ErrorWrapper(func(w http.ResponseWriter, r *http.Request) error {
		logger.FromCtx(r.Context()).Info("hello")
		w.WriteHeader(200)
		return nil
	}))

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", api.ErrorWrapper(register.New(validator, authService, tokenService, email)))
		r.Post("/login", api.ErrorWrapper(login.New(authService, *validator, session)))
		r.Get("/verify/{token}", api.ErrorWrapper(verify.New(tokenService)))
	})

	r.Route("/category", func(r chi.Router) {
		r.Get("/", api.ErrorWrapper(get_category.New(categoryService)))
	})

	r.Route("/products", func(r chi.Router) {
		r.Get("/", api.ErrorWrapper(get_product.New(productService)))
		r.Get("/{id}", api.ErrorWrapper(get_product_by_id.New(productService)))
	})

	r.Route("/admin", func(r chi.Router) {
		r.Use(middlewares.Admin(userService))
		r.Post("/create-category", api.ErrorWrapper(create_category.New(categoryService, validator)))
		r.Post("/create-product", api.ErrorWrapper(create_product.New(validator, productService)))
	})

	r.Route("/cart", func(r chi.Router) {
		r.Use(middlewares.User(userService))
		r.Post("/add", api.ErrorWrapper(cart_add_product.New(validator, cartService)))
		r.Get("/", api.ErrorWrapper(get_cart.New(cartService)))
	})

	return &Server{
		srv: &http.Server{
			Addr:         cfg.Addr,
			Handler:      r,
			ReadTimeout:  cfg.Timeout,
			WriteTimeout: cfg.Timeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	const op = "server.Run"

	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Server) GracefulStop(ctx context.Context) error {
	const op = "server.GracefulStop"

	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Server) Addr() string {
	return s.srv.Addr
}
