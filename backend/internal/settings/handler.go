package settings

import (
	"net/http"
)

type writeJSONFunc func(http.ResponseWriter, int, any)
type writeErrorFunc func(http.ResponseWriter, int, string, string)
type decodeJSONFunc func(*http.Request, any) error

type Handler struct {
	service   *Service
	writeJSON writeJSONFunc
	writeErr  writeErrorFunc
	decode    decodeJSONFunc
}

type updateRequest struct {
	Tistory TistorySettings `json:"tistory"`
}

func NewHandler(service *Service, writeJSON writeJSONFunc, writeErr writeErrorFunc, decode decodeJSONFunc) *Handler {
	return &Handler{
		service:   service,
		writeJSON: writeJSON,
		writeErr:  writeErr,
		decode:    decode,
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.Get(r.Context())
	if err != nil {
		h.writeErr(w, http.StatusInternalServerError, "DATABASE_ERROR", "설정을 불러오지 못했습니다.")
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var req updateRequest
	if err := h.decode(r, &req); err != nil {
		h.writeErr(w, http.StatusBadRequest, "VALIDATION_ERROR", "요청 형식이 올바르지 않습니다.")
		return
	}
	result, err := h.service.Update(r.Context(), req.Tistory)
	if err != nil {
		h.writeErr(w, http.StatusInternalServerError, "DATABASE_ERROR", "설정을 저장하지 못했습니다.")
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}
