package illustrations

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"time"

	_ "golang.org/x/image/webp"

	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	repo "github.com/fillipgms/portfolio-api/internal/adapters/postgresql/sqlc"
	"github.com/fillipgms/portfolio-api/internal/helpers"
	"github.com/fillipgms/portfolio-api/internal/json"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
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
	r.Body = http.MaxBytesReader(w, r.Body, 10 << 20)
	
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.MultipartForm.RemoveAll()

	uploadedFile, uploadedFileHeader, err := r.FormFile("image")

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer uploadedFile.Close()

	buffer := make([]byte, 512)
    if _, err := uploadedFile.Read(buffer); err != nil {
        log.Println(err)
        http.Error(w, "Failed to read file", http.StatusInternalServerError)
        return
    }
    mimeType := http.DetectContentType(buffer)
    uploadedFile.Seek(0, 0)

	img, _, err := image.DecodeConfig(uploadedFile)
    if err != nil {
        log.Println(err)
        http.Error(w, "Invalid image file", http.StatusBadRequest)
        return
    }
    uploadedFile.Seek(0, 0)
	
	hasher := sha256.New()
    fileSize := int64(0)
    if n, err := io.Copy(hasher, uploadedFile); err != nil {
        log.Println(err)
        http.Error(w, "Failed to calculate checksum", http.StatusInternalServerError)
        return
    } else {
        fileSize = n
    }
    uploadedFile.Seek(0, 0)

	checksum := hex.EncodeToString(hasher.Sum(nil))
    ext := filepath.Ext(uploadedFileHeader.Filename)
    filename := checksum + ext

    _, err = helpers.BunnyClient.Upload(r.Context(), "illustrations", filename, checksum, uploadedFile)

	if err != nil {
        log.Println(err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	imageURL := fmt.Sprintf("https://%s.b-cdn.net/illustrations/%s", 
        "fillipsportfolio", 
        filename,
    )

	formFinishedAt := r.FormValue("finishedAt")

	var finishedAt pgtype.Timestamptz

	if formFinishedAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, formFinishedAt)
		if err != nil {
			finishedAt = pgtype.Timestamptz{
				Time:  time.Now(),
				Valid: true,
			}
		} else {
			finishedAt = pgtype.Timestamptz{
				Time:  parsedTime,
				Valid: true,
			}
		}
	} else {
		finishedAt = pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		}
	}

	tempIllustration := repo.CreateIllustrationParams{
        Title:       r.FormValue("title"),
        Description: r.FormValue("description"),
        Imageurl:    imageURL,
        Imageheight: pgtype.Int4{
            Int32: int32(img.Height),
            Valid: true,
        },
        Imagewidth: pgtype.Int4{
            Int32: int32(img.Width),
            Valid: true,
        },
        Imagemimetype: pgtype.Text{
            String: mimeType,
            Valid:  true,
        },
        Imagefilesize: pgtype.Int4{
            Int32: int32(fileSize),
            Valid: true,
        },
		FinishedAt: finishedAt,
    }
	
	createdIllustration, err := h.service.CreateIllustration(r.Context(), tempIllustration)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	titleSlug := helpers.Slugify(createdIllustration.Title, createdIllustration.ID)
    pgText := pgtype.Text{
        String: titleSlug,
        Valid:  true,
    }

	updatedIllustration, err := h.service.UpdateSlug(r.Context(), pgText, createdIllustration.ID)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]any{
        "id":          updatedIllustration.ID,
        "title":       updatedIllustration.Title,
        "slug":        updatedIllustration.Slug,
        "description": updatedIllustration.Description,
        "image": map[string]any{
            "url":      updatedIllustration.Imageurl,
            "width":    updatedIllustration.Imagewidth.Int32,
            "height":   updatedIllustration.Imageheight.Int32,
            "mimeType": updatedIllustration.Imagemimetype.String,
            "fileSize": updatedIllustration.Imagefilesize.Int32,
        },
        "post":        updatedIllustration.Post,
        "finishedAt":  updatedIllustration.FinishedAt,
        "createdAt":   updatedIllustration.CreatedAt,
    }
    
	json.Write(w, http.StatusOK, response)
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

	responses := make([]map[string]any, 0, len(illustrations))

	for _, illustration := range illustrations {
		responses = append(responses, map[string]any{
			"id":          illustration.ID,
			"title":       illustration.Title,
			"slug":        illustration.Slug,
			"description": illustration.Description,
			"image": map[string]any{
				"url":      illustration.Imageurl,
				"width":    illustration.Imagewidth.Int32,
				"height":   illustration.Imageheight.Int32,
				"mimeType": illustration.Imagemimetype.String,
				"fileSize": illustration.Imagefilesize.Int32,
			},
			"post":       illustration.Post,
			"finishedAt": illustration.FinishedAt,
			"createdAt":  illustration.CreatedAt,
		})
	}

	res := helpers.PaginationFormat(count, responses, limit, offset, page)

	json.Write(w, http.StatusOK, res)
}

func (h *handler) FindIllustrationById(w http.ResponseWriter, r *http.Request) {
	idOrSlug := chi.URLParam(r, "illustrationId")

	if id, err := strconv.ParseInt(idOrSlug, 10, 64); err == nil {
		illustration, err := h.service.FindIllustrationById(r.Context(), id)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.Write(w, http.StatusOK, illustration)
		return
	}

	slugText := pgtype.Text{
		String: idOrSlug,
		Valid:  true,
	}

	illustration, err := h.service.FindIllustrationByName(r.Context(), slugText)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]any{
        "id":          illustration.ID,
        "title":       illustration.Title,
        "slug":        illustration.Slug,
        "description": illustration.Description,
        "image": map[string]any{
            "url":      illustration.Imageurl,
            "width":    illustration.Imagewidth.Int32,
            "height":   illustration.Imageheight.Int32,
            "mimeType": illustration.Imagemimetype.String,
            "fileSize": illustration.Imagefilesize.Int32,
        },
        "post":        illustration.Post,
        "finishedAt":  illustration.FinishedAt,
        "createdAt":   illustration.CreatedAt,
    }

	json.Write(w, http.StatusOK, response)
}
