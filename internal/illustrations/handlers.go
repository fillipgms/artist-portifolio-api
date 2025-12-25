package illustrations

import (
	"log"
	"net/http"
	"strconv"

	repo "github.com/fillipgms/portfolio-api/internal/adapters/postgresql/sqlc"
	"github.com/fillipgms/portfolio-api/internal/helpers"
	"github.com/fillipgms/portfolio-api/internal/json"
	"github.com/go-chi/chi/v5"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}


func (h *handler) CreateIllustration(w http.ResponseWriter, r *http.Request) {
	var tempIllustration repo.CreateIllustrationParams
	if err := json.Read(r, &tempIllustration); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	createdIllustration, err := h.service.CreateIllustration(r.Context(), tempIllustration)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, createdIllustration)

}

func (h *handler) ListIllustrations(w http.ResponseWriter, r *http.Request) {

	page := int64(1)
	limit := int64(10)

	if p := r.URL.Query().Get("page"); p != "" {
		pageInt, err := strconv.ParseInt(p, 10, 64)

		if err != nil || pageInt < 1 {
			http.Error(w, "invalid page", http.StatusBadRequest)
			return
		}

		page = pageInt
	}

	offset := (page - 1) * limit

	count, err := h.service.FindIllustrationsCount(r.Context())

	if (err != nil) {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	illustrations, err := h.service.ListIllustrations(r.Context(), int32(limit), int32(offset))

	if (err != nil) {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := helpers.PaginationFormat(count, illustrations, limit, offset, page)

	json.Write(w, http.StatusOK, res)
}

func (h *handler) FindIllustrationById(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "illustrationId")

	id, err := strconv.ParseInt(idString, 10, 64)

	if (err != nil) {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	illustration, err := h.service.FindIllustrationById(r.Context(), id)

	if (err != nil) {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, illustration)
}
