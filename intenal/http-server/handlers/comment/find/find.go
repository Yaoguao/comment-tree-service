package find

import (
	"comment-tree-service/intenal/models"
	"comment-tree-service/intenal/service"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type CommentGetter interface {
	GetThread(parentIDHex string, limit, offset int, sort string) ([]models.Comment, error)
}

type CommentSearcher interface {
	Search(query string, limit, offset int) ([]models.Comment, error)
}

type response struct {
	Comments []models.Comment `json:"comments,omitempty"`
	Error    string           `json:"error,omitempty"`
}

func NewGetHandler(log *logrus.Logger, getter CommentGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parentIDHex := r.URL.Query().Get("parent")
		if parentIDHex == "" {
			log.Warn("parent ID is missing")
			http.Error(w, "parent ID is required", http.StatusBadRequest)
			return
		}

		limit := 20
		offset := 0
		sort := "asc"

		if l := r.URL.Query().Get("limit"); l != "" {
			if parsed, err := strconv.Atoi(l); err == nil {
				limit = parsed
			}
		}
		if o := r.URL.Query().Get("offset"); o != "" {
			if parsed, err := strconv.Atoi(o); err == nil {
				offset = parsed
			}
		}
		if s := r.URL.Query().Get("sort"); s != "" {
			if s == "asc" || s == "desc" {
				sort = s
			}
		}

		comments, err := getter.GetThread(parentIDHex, limit, offset, sort)
		if err != nil {
			if errors.Is(err, service.InvalidID) {
				log.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response{Error: "UPS!!"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response{Comments: comments})
	}
}

func NewSearchHandler(log *logrus.Logger, searcher CommentSearcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")

		limit := 20
		offset := 0

		if l := r.URL.Query().Get("limit"); l != "" {
			if parsed, err := strconv.Atoi(l); err == nil {
				limit = parsed
			}
		}
		if o := r.URL.Query().Get("offset"); o != "" {
			if parsed, err := strconv.Atoi(o); err == nil {
				offset = parsed
			}
		}

		comments, err := searcher.Search(query, limit, offset)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response{Error: "UPS!!"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response{Comments: comments})
	}
}
