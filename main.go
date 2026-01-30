package main

import (
	"category-api/database"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// core
var Categories = []Category{
	{
		ID:          1,
		Name:        "Makanan",
		Description: "Makanan adalah sesuatu yg bisa dimakan",
	},
}

func CreateCategory(data *Category) error {
	if data.Name == "" || data.Description == "" {
		return errors.New("Invalid Payload")
	}
	data.ID = len(Categories) + 1
	Categories = append(Categories, *data)
	return nil
}

func GeCategories() []Category {
	return Categories
}

func GetCategoryByID(ID int) (Category, error) {
	for i := 0; i < len(Categories); i++ {
		currentCategory := Categories[i]
		if currentCategory.ID == ID {
			return currentCategory, nil
		}
	}
	return Category{}, errors.New("NotFound")
}

func UpdateCategory(ID int, data *Category) error {
	if data.Name == "" || data.Description == "" {
		return errors.New("Invalid Payload")
	}

	for i := 0; i < len(Categories); i++ {
		currentCategory := Categories[i]
		if currentCategory.ID == ID {
			data.ID = ID
			Categories[i] = *data
			return nil
		}
	}
	return errors.New("NotFound")
}

func DeleteCategory(ID int) error {
	for i := range Categories {
		currentCategory := Categories[i]
		if currentCategory.ID == ID {
			Categories = append(Categories[:i], Categories[i+1:]...)
			return nil
		}
	}
	return errors.New("NotFound")
}

// end core

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

	http.HandleFunc("/categories/", func(response http.ResponseWriter, request *http.Request) {
		method := request.Method
		id, err := GetIDFromUrl(request.URL.Path, "/categories/")

		if err != nil {
			http.Error(response, "Invalid Category ID", http.StatusBadRequest)
			return
		}

		if method == http.MethodGet {
			result, err := GetCategoryByID(id)

			if err != nil {
				http.Error(response, "Not Found", http.StatusNotFound)
				return
			}

			response.Header().Set("Content-Type", "application/json")
			json.NewEncoder(response).Encode(result)
			return

		}

		if method == http.MethodPut {
			var updatedCategory Category

			json.NewDecoder(request.Body).Decode(&updatedCategory)
			err := UpdateCategory(id, &updatedCategory)

			if err != nil {
				if err.Error() == "NotFound" {
					http.Error(response, "Not Found", http.StatusNotFound)
					return
				}
				http.Error(response, "Invalid Payload", http.StatusBadRequest)
				return
			}

			response.Header().Set("Content-Type", "application/json")
			json.NewEncoder(response).Encode(map[string]string{
				"status":  "OK",
				"message": "Category updated",
			})
			return

		}

		if method == http.MethodDelete {
			err := DeleteCategory(id)

			if err != nil {
				http.Error(response, "Not Found", http.StatusNotFound)
				return
			}

			response.Header().Set("Content-Type", "application/json")
			json.NewEncoder(response).Encode(map[string]string{
				"status":  "OK",
				"message": "Category deleted",
			})
			return

		}

		http.Error(response, "Invalid Method", http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/categories", func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")

		if request.Method == http.MethodGet {
			result := GeCategories()
			json.NewEncoder(response).Encode(result)
			return
		}

		if request.Method == http.MethodPost {
			var newCategory Category
			err := json.NewDecoder(request.Body).Decode(&newCategory)
			if err != nil {
				http.Error(response, "Payload not valid", http.StatusBadRequest)
				return
			}

			createErr := CreateCategory(&newCategory)

			if createErr != nil {
				http.Error(response, "Failed to store category", http.StatusInternalServerError)
				return
			}

			response.WriteHeader(http.StatusCreated)
			json.NewEncoder(response).Encode(newCategory)
			return
		}

		http.Error(response, "Invalid Method", http.StatusMethodNotAllowed)

	})

	fmt.Printf("Server running on port %v", config.Port)
	serverError := http.ListenAndServe(config.Port, nil)
	if serverError != nil {
		fmt.Println("Error starting server:", serverError)
	}

}
