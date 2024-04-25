package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	gql "github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

type (
	Motorcycle struct {
		Year     time.Time `json:"year"`
		Make     string    `json:"make"`
		Model    string    `json:"model"`
		ImageURL string    `json:"imageUrl"`
		ID       int       `json:"name"`
	}
)

var (
	rootQuery = gql.NewObject(gql.ObjectConfig{
		Name: "RootQuery",
		Fields: gql.Fields{
			"motorcycle": &gql.Field{
				Type:        motorcycleType,
				Description: "Get single motorcycle",
				Args: gql.FieldConfigArgument{
					"make": &gql.ArgumentConfig{
						Type: gql.String,
					},
					"model": &gql.ArgumentConfig{
						Type: gql.String,
					},
					"year": &gql.ArgumentConfig{
						Type: gql.DateTime,
					},
				},
				Resolve: func(params gql.ResolveParams) (interface{}, error) {
					makeQuery, isOk := params.Args["make"].(string)
					if !isOk {
						return Motorcycle{}, fmt.Errorf("failed to get 'make' argument from query")
					}
					modelQuery, isOk := params.Args["model"].(string)
					if !isOk {
						return Motorcycle{}, fmt.Errorf("failed to get 'model' argument from query")
					}
					yearQuery, isOk := params.Args["year"].(time.Time)
					if !isOk {
						return Motorcycle{}, fmt.Errorf("failed to get 'year' argument from query")
					}
					for _, mc := range MotorcycleList {
						if mc.Make == makeQuery && mc.Model == modelQuery && mc.Year.Equal(yearQuery) {
							return mc, nil
						}
					}
					return Motorcycle{}, nil
				},
			},
			"motorcycleList": &gql.Field{
				Type:        gql.NewList(motorcycleType),
				Description: "List of motorcycles",
				Resolve: func(params gql.ResolveParams) (interface{}, error) {
					return MotorcycleList, nil
				},
			},
		},
	})

	motorcycleType = gql.NewObject(gql.ObjectConfig{
		Name: "Motorcycle",
		Fields: gql.Fields{
			"id": &gql.Field{
				Type: gql.Int,
			},
			"make": &gql.Field{
				Type: gql.String,
			},
			"model": &gql.Field{
				Type: gql.String,
			},
			"year": &gql.Field{
				Type: gql.DateTime,
			},
			"imageUrl": &gql.Field{
				Type: gql.String,
			},
		},
	})

	MotorcycleList []Motorcycle
	_ = importJSONDataFromFile("./motorcycleData.json", &MotorcycleList)

	MotorcycleSchema, _ = gql.NewSchema(gql.SchemaConfig{
		Query: rootQuery,
	})
)

func main() {
	h := handler.New(&handler.Config{
		Schema: &MotorcycleSchema,
		Pretty: true,
		GraphiQL: false,
	})

	http.Handle("/graphql", h)
	http.ListenAndServe(":8080", nil)
}

func importJSONDataFromFile(fileName string, result interface{}) (isOK bool) {
	isOK = true
	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Print("Error:", err)
		isOK = false
	}
	err = json.Unmarshal(content, result)
	if err != nil {
		isOK = false
		fmt.Print("Error:", err)
	}
	return
}
