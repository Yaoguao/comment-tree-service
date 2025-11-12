package save

import (
	"comment-tree-service/intenal/models"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type response struct {
	Comment *models.Comment `json:"comment,omitempty"`
	Error   string          `json:"error,omitempty"`
}

type CommentSaver interface {
	SaveComment(c *models.Comment) error
}

func New(log *logrus.Logger, saver CommentSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var com models.Comment
		if err := json.NewDecoder(r.Body).Decode(&com); err != nil {
			log.WithError(err).Warn("failed to decode request body")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response{Error: "invalid request body"})
			return
		}

		if com.Content == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response{Error: "content is empty"})
			return
		}

		if com.Author == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response{Error: "author name is empty"})
			return
		}

		err := saver.SaveComment(&com)
		if err != nil {
			log.WithError(err).Error("failed to save comment")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response{Error: "failed to save comment"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response{Comment: &com})
	}
}
