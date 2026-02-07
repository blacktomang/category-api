package main

import (
	"category-api/database"
	"category-api/handlers"
	"category-api/repositories"
	"category-api/services"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DBCONN"`
}

func GetIDFromUrl(path string, prefix string) (int, error) {
	idStr := strings.TrimPrefix(path, prefix)

	return strconv.Atoi(idStr)
}

const notSwagger = `
<ul>
categories
    <li>
        METHOD: GET, Link <a href="/api/categories">/api/categories</a>
    </li>
    <li>
        METHOD:  POST, Link <a href="/api/categories">/api/categories</a>
    </li>
    <li>
        METHOD:  GET, Link <a href="/api/categories/1">/api/categories/{id}</a>
    </li>
    <li>
        METHOD:  PUT, Link <a href="/api/categories/1">/api/categories/{id}</a>
    </li>
    <li>
        METHOD:  DELETE, Link <a href="/api/categories/1">/api/categories/{id}</a>
    </li>
products
	<li>
		METHOD: GET, Link <a href="/api/products">/api/products</a>
	</li>
	<li>
		METHOD:  POST, Link <a href="/api/products">/api/products</a>
	</li>
	<li>
		METHOD:  GET, Link <a href="/api/products/1">/api/products/{id}</a>
	</li>
	<li>
		METHOD:  PUT, Link <a href="/api/products/1">/api/products/{id}</a>
	</li>
	<li>
		METHOD:  DELETE, Link <a href="/api/products/1">/api/products/{id}</a>
	</li>
</ul>
`

const API_PREFIX = "/api"

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DBCONN"),
	}

	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodGet {
			response.Header().Set("Content-Type", "text/html")
			io.WriteString(response, notSwagger)
		}
	})

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	http.HandleFunc(API_PREFIX+"/categories/", categoryHandler.HandleCategoryByID)
	http.HandleFunc(API_PREFIX+"/categories", categoryHandler.HandleCategories)

	http.HandleFunc(API_PREFIX+"/products", productHandler.HandleProducts)
	http.HandleFunc(API_PREFIX+"/products/", productHandler.HandleProductByID)

	http.HandleFunc(API_PREFIX+"/checkout", transactionHandler.HandleCheckout)
	http.HandleFunc(API_PREFIX+"/report", transactionHandler.ReportRange)
	http.HandleFunc(API_PREFIX+"/report/hari-ini", transactionHandler.ReportToday)

	addr := "0.0.0.0:" + config.Port
	fmt.Printf("Server running on port %v", addr)
	serverError := http.ListenAndServe(addr, nil)
	if serverError != nil {
		fmt.Println("Error starting server:", serverError)
	}

}
