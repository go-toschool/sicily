package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-toschool/sicily"
	"github.com/graphql-go/graphql"
)

// API ...
func API(ctx *Context, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "This server does not support that HTTP method", http.StatusBadRequest)
		return
	}

	id, ok := r.Context().Value(sicily.AuthUserIDContextKey).(string)
	if !ok {
		http.Error(w, "missing user id", http.StatusBadRequest)
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

		result = ctx.ExecuteQuery(gr.Query, id)
	default:
		http.Error(w, "bad content type", http.StatusBadRequest)
		return
	}

	w.Header().Set("Accept-Encoding", "gzip")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
