package delete

import (
	"comment-tree-service/intenal/service"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
)

type CommentDeleter interface {
	DeleteThread(IDHex string) error
}

func New(log *logrus.Logger, deleter CommentDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idHex := chi.URLParam(r, "id")
		if idHex == "" {
			log.Warn("id is missing")
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		err := deleter.DeleteThread(idHex)
		if err != nil {
			if errors.Is(err, service.InvalidID) {
				log.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode("UPS!!")
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("delete success")
	}
}
