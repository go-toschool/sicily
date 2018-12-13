package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"google.golang.org/grpc"

	"github.com/graphql-go/graphql"
	"github.com/go-toschool/sicily/graph"
	"github.com/go-toschool/sicily/graph/queries"
	"github.com/go-toschool/syracuse/citizens"
)

func main() {
	conn, err := grpc.Dial("localhost:8001", grpc.WithInsecure())
	check(err)

	ctx := &graph.Context{
		UserService: citizens.NewCitizenshipClient(conn),
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: queries.Users(ctx),
	})
	check(err)

	http.HandleFunc("/users", disableCors(getUser(schema)))

	log.Println("Now server is running on port 3000")
	http.ListenAndServe(":3000", nil)
}

// ContentTypeGraphQL graphql content type
const ContentTypeGraphQL = "application/graphql"

type GraphRequest struct {
	Query string `json:"query"`
}

func getUser(schema graphql.Schema) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "This server does not support that HTTP method", http.StatusBadRequest)
			return
		}
		contentTypeStr := r.Header.Get("Content-Type")
		contentTypeTokens := strings.Split(contentTypeStr, ";")
		contentType := contentTypeTokens[0]

		var result *graphql.Result
		switch contentType {
		case ContentTypeGraphQL:
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "could not read body", http.StatusInternalServerError)
				return
			}

			gr := &GraphRequest{}

			if err := json.Unmarshal(body, gr); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			result = executeQuery(gr.Query, schema)
		default:
			http.Error(w, "bad content type", http.StatusBadRequest)
		}

		w.Header().Set("Accept-Encoding", "gzip")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

func disableCors(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, Accept-Encoding")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusOK)
			return
		}
		h.ServeHTTP(w, r)
	}
}

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v\n", result.Errors)
	}
	return result
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
