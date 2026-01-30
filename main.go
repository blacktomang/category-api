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
    <li>
        METHOD: GET, Link <a href="/categories">/categories</a>
    </li>
    <li>
        METHOD:  POST, Link <a href="/categories">/categories</a>
    </li>
    <li>
        METHOD:  GET, Link <a href="/categories/1">/categories/{id}</a>
    </li>
    <li>
        METHOD:  PUT, Link <a href="/categories/1">/categories/{id}</a>
    </li>
    <li>
        METHOD:  DELETE, Link <a href="/categories/1">/categories/{id}</a>
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

	http.HandleFunc(API_PREFIX+"/categories/", categoryHandler.HandleCategoryByID)
	http.HandleFunc(API_PREFIX+"/categories", categoryHandler.HandleCategories)

	http.HandleFunc(API_PREFIX+"/products", productHandler.HandleProducts)
	http.HandleFunc(API_PREFIX+"/products/", productHandler.HandleProductByID)
	fmt.Printf("Server running on port %v", config.Port)
	serverError := http.ListenAndServe(config.Port, nil)
	if serverError != nil {
		fmt.Println("Error starting server:", serverError)
	}

}
