package post

import (
	"errors"
	"net/http"
	"strconv"
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

func NewHandler(service *Service, writeJSON writeJSONFunc, writeErr writeErrorFunc, decode decodeJSONFunc) *Handler {
	return &Handler{
		service:   service,
		writeJSON: writeJSON,
		writeErr:  writeErr,
		decode:    decode,
	}
}

func (h *Handler) Preview(w http.ResponseWriter, r *http.Request) {
	var req PreviewRequest
	if err := h.decode(r, &req); err != nil {
		h.writeErr(w, http.StatusBadRequest, "VALIDATION_ERROR", "요청 형식이 올바르지 않습니다.")
		return
	}
	result, err := h.service.Preview(r.Context(), req)
	if err != nil {
		h.handleError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := h.decode(r, &req); err != nil {
		h.writeErr(w, http.StatusBadRequest, "VALIDATION_ERROR", "요청 형식이 올바르지 않습니다.")
		return
	}
	result, err := h.service.Create(r.Context(), req)
	if err != nil {
		h.handleError(w, err)
		return
	}
	h.writeJSON(w, http.StatusCreated, result)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, h.writeErr)
	if !ok {
		return
	}

	var req UpdateRequest
	if err := h.decode(r, &req); err != nil {
		h.writeErr(w, http.StatusBadRequest, "VALIDATION_ERROR", "요청 형식이 올바르지 않습니다.")
		return
	}
	result, err := h.service.Update(r.Context(), id, req)
	if err != nil {
		h.handleError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page := parseInt(r.URL.Query().Get("page"), 1)
	size := parseInt(r.URL.Query().Get("size"), 20)
	result, err := h.service.List(r.Context(), page, size)
	if err != nil {
		h.writeErr(w, http.StatusInternalServerError, "DATABASE_ERROR", "글 목록을 불러오지 못했습니다.")
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, h.writeErr)
	if !ok {
		return
	}
	result, err := h.service.Get(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, h.writeErr)
	if !ok {
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, map[string]bool{"deleted": true})
}

func (h *Handler) PublishTistory(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, h.writeErr)
	if !ok {
		return
	}

	var req PublishRequest
	if err := h.decode(r, &req); err != nil {
		h.writeErr(w, http.StatusBadRequest, "VALIDATION_ERROR", "요청 형식이 올바르지 않습니다.")
		return
	}
	result, err := h.service.PublishTistory(r.Context(), id, req)
	if err != nil {
		h.handleError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	if validationErr, ok := IsValidationError(err); ok {
		h.writeErr(w, http.StatusBadRequest, validationErr.Code, validationErr.Message)
		return
	}
	if errors.Is(err, ErrNotFound) {
		h.writeErr(w, http.StatusNotFound, "POST_NOT_FOUND", "글을 찾을 수 없습니다.")
		return
	}
	h.writeErr(w, http.StatusInternalServerError, "DATABASE_ERROR", "데이터 처리 중 오류가 발생했습니다.")
}

func pathID(w http.ResponseWriter, r *http.Request, writeErr writeErrorFunc) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		writeErr(w, http.StatusBadRequest, "VALIDATION_ERROR", "id가 올바르지 않습니다.")
		return 0, false
	}
	return id, true
}

func parseInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
