package app

import (
	"tech-challenge-fase-1/internal/infra/config"
	"tech-challenge-fase-1/internal/infra/database"
	"tech-challenge-fase-1/internal/infra/events"
	httpserver "tech-challenge-fase-1/internal/infra/http"
	"tech-challenge-fase-1/internal/infra/queries"
	"tech-challenge-fase-1/internal/infra/repositories"
	"tech-challenge-fase-1/internal/infra/services"

	"github.com/gin-contrib/cors"
)

type APIApp struct {
	httpServer         *httpserver.GinHTTPServerAdapter
	connection         *database.PGXConnectionAdapter
	customerRepository *repositories.CustomerRepositoryDB
	customerService *services.AwsCustomerService
	productRepository  *repositories.ProductRepositoryDB
	orderRepository    *repositories.OrderRepositoryDB
	orderDisplayListQuery *queries.OrderDisplayListQueryDB
	mercadoPagoGateway *services.MercadoPagoGateway
	eventManager *events.EventManager
}

func NewAPIApp() *APIApp {
	app := &APIApp{}
	app.initGin()
	app.configCors()
	app.initConnectionDB()
	app.configRoutes()
	return app
}

func (app *APIApp) initGin() {
	app.httpServer = httpserver.NewGinHTTPServerAdapter()
	app.httpServer.SetTrustedProxies(nil)
}

func (app *APIApp) configCors() {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"*"}
	config.AllowHeaders = []string{"Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"}
	app.httpServer.Engine.Use(cors.New(config))
}

func (app *APIApp) initConnectionDB() {
	app.connection = database.NewPGXConnectionAdapter()

	app.customerRepository = repositories.NewCustomerRepositoryDB(app.connection)
	app.productRepository = repositories.NewProductRepositoryDB(app.connection)
	app.orderRepository = repositories.NewOrderRepositoryDB(app.connection)
	app.orderDisplayListQuery = queries.NewOrderDisplayListQueryDB(app.connection)

	app.eventManager = events.NewEventManager()

	app.mercadoPagoGateway = services.NewMercadoPagoGateway(app.eventManager)
	var err error
	app.customerService, err = services.NewAwsCustomerService(
		config.AWS_REGION,
		config.AWS_USER_POOL_ID,
	)
	if err != nil {
		panic(err)
	}
}

func (app *APIApp) configRoutes() {
	registerRouters(app)
}

func (app *APIApp) Run() {
	app.httpServer.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func (app *APIApp) Shutdown() {
	app.connection.Close()
}
