package rest

import (
	"context"
	"fmt"
	"go-clean/docs/swagger"
	"go-clean/src/business/usecase"
	"go-clean/src/lib/auth"
	"go-clean/src/lib/configreader"
	"go-clean/src/utils/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var once = &sync.Once{}

type REST interface {
	Run()
}

type rest struct {
	http         *gin.Engine
	conf         config.ApplicationMeta
	configreader configreader.Interface
	uc           *usecase.Usecase
	auth         auth.Interface
}

func Init(conf config.ApplicationMeta, confReader configreader.Interface, uc *usecase.Usecase, auth auth.Interface) REST {
	r := &rest{}
	once.Do(func() {
		httpServ := gin.Default()

		r = &rest{
			conf:         conf,
			configreader: confReader,
			http:         httpServ,
			uc:           uc,
			auth:         auth,
		}

		r.http.Use(cors.New(cors.Config{
			AllowAllOrigins: true,
			AllowHeaders:    []string{"*"},
			AllowMethods: []string{
				http.MethodHead,
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
			},
		}))

		// Set Recovery
		r.http.Use(gin.Recovery())

		r.http.Use(gin.Logger())

		r.Register()
	})

	return r
}

func (r *rest) Run() {
	port := ":8080"

	server := &http.Server{
		Addr:    port,
		Handler: r.http,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println(fmt.Sprintf("Serving HTTP error: %s", err.Error()))
		}
	}()
	fmt.Println(fmt.Sprintf("Listening and Serving HTTP on %s", server.Addr))

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(fmt.Sprintf("Server forced to shutdown: %v", err))
	}

	log.Println("Server exiting")
}

func (r *rest) Register() {
	r.registerSwaggerRoutes()
	publicApi := r.http.Group("/public")
	publicApi.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "hello world",
		})
	})

	r.http.Static("/public/assets/umkm", "./public/assets/umkm")
	r.http.Static("/public/assets/menu", "./public/assets/menu")
	api := r.http.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/", r.VerifyUser, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "hello mail",
		})
	})

	auth := v1.Group("/auth")
	auth.POST("/register", r.VerifyUser, r.VerifyAdmin, r.RegisterUser)
	auth.POST("/login", r.LoginUser)
	auth.POST("/guest", r.LoginGuestUser)

	umkm := v1.Group("/umkm")
	umkm.POST("/create", r.VerifyUser, r.VerifyAdmin, r.CreateUmkm)
	umkm.GET("/:umkm_id", r.VerifyUser, r.GetUmkmByID)
	umkm.GET("", r.VerifyUser, r.GetUmkmList)
	umkm.PUT("/:umkm_id", r.VerifyUser, r.VerifyUmkm, r.UpdateUmkm)
	umkm.DELETE("/:umkm_id", r.VerifyUser, r.VerifyUmkm, r.DeleteUmkm)
	umkm.POST("/:umkm_id/upload-image", r.VerifyUser, r.VerifyUmkm, r.UploadImageUmkm)

	// menu
	menu := v1.Group("/menu")
	umkm.POST("/:umkm_id/menu/create", r.VerifyUser, r.VerifyUmkm, r.CreateMenu)
	menu.GET("/:menu_id", r.VerifyUser, r.GetMenuByID)
	menu.GET("", r.VerifyUser, r.GetMenuList)
	menu.PUT("/:menu_id", r.VerifyUser, r.VerifyMenu, r.UpdateMenu)
	menu.DELETE("/:menu_id", r.VerifyUser, r.VerifyMenu, r.DeleteMenu)
	umkm.POST("/:umkm_id/menu/:menu_id/upload-image", r.VerifyUser, r.VerifyUmkm, r.UploadImageMenu)

	cart := v1.Group("/cart")
	cart.POST("/create", r.VerifyUser, r.AddMenuToCart)
	cart.GET("", r.VerifyUser, r.GetListCartByUser)
	cart.PUT("/:cart_id/decrease", r.VerifyUser, r.DecreaseItem)
	cart.DELETE("/:cart_id", r.VerifyUser, r.VerifyCart, r.DeleteItemCart)
	cart.DELETE("/clear-cart", r.VerifyUser, r.ClearCart)

	admin := v1.Group("/admin")

	// transaction
	umkm.GET("/:umkm_id/transactions", r.VerifyUser, r.VerifyUmkm, r.GetTransactionListUmkm)
	admin.GET("/transactions", r.VerifyUser, r.VerifyAdmin, r.GetTransactionList)
	admin.GET("/transactions/recap", r.VerifyUser, r.VerifyAdmin, r.GetRecapSalesList)
	transaction := v1.Group("/transaction")
	transaction.POST("/create", r.VerifyUser, r.CreateOrder)
	transaction.GET("/:transaction_id/payment-detail", r.VerifyUser, r.GetPaymentDetail)
	transaction.GET("/:transaction_id", r.GetOrderDetail)
	transaction.GET("/me", r.VerifyUser, r.GetMyTransaction)
	umkm.PUT("/:umkm_id/transaction/:transaction_id/mark-as-done", r.VerifyUser, r.VerifyUmkm, r.CompleteOrder)
	umkm.PUT("/:umkm_id/transaction/:transaction_id/cancel-order", r.VerifyUser, r.VerifyUmkm, r.CancelOrder)
	admin.PUT("/transaction/:order_id/mark-as-paid", r.VerifyUser, r.VerifyAdmin, r.MarkAsPaid)
	admin.GET("/transactions/recap/download", r.VerifyUser, r.VerifyAdmin, r.DownloadMonthlyRecap)

	midtransTransaction := v1.Group("/midtrans-transaction")
	midtransTransaction.POST("/handle", r.HandleNotification)

	user := v1.Group("/user")
	user.GET("/cart-count", r.VerifyUser, r.GetCartCount)
	user.GET("/me", r.VerifyUser, r.GetMe)

	// analytic
	umkm.GET("/:umkm_id/analytic/dashboard-widget", r.VerifyUser, r.VerifyUmkm, r.GetDashboardWidget)
	admin.GET("/analytic/dashboard-widget", r.VerifyUser, r.VerifyAdmin, r.GetAllDashboardWidget)

	// withdraw
	admin.GET("/withdraw", r.VerifyUser, r.VerifyAdmin, r.GetWithdrawList)
	admin.POST("/withdraw", r.VerifyUser, r.VerifyAdmin, r.CreateWithdraw)
	admin.PUT("/withdraw/:withdraw_id", r.VerifyUser, r.VerifyAdmin, r.UpdateWithdraw)
}

func (r *rest) registerSwaggerRoutes() {
	swagger.SwaggerInfo.Title = r.conf.Title
	swagger.SwaggerInfo.Description = r.conf.Description
	swagger.SwaggerInfo.Version = r.conf.Version
	swagger.SwaggerInfo.Host = r.conf.Host
	swagger.SwaggerInfo.BasePath = r.conf.BasePath

	r.http.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
