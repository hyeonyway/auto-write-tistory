package tistory

import (
	"errors"
	"net/http"
	"strings"

	publishertistory "devlog-studio/backend/internal/publisher/tistory"
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

func (h *Handler) FetchCategories(w http.ResponseWriter, r *http.Request) {
	var req FetchCategoriesRequest
	if err := h.decode(r, &req); err != nil {
		h.writeErr(w, http.StatusBadRequest, "VALIDATION_ERROR", "요청 형식이 올바르지 않습니다.")
		return
	}
	if strings.TrimSpace(req.BlogURL) == "" {
		h.writeErr(w, http.StatusBadRequest, "VALIDATION_ERROR", "blogUrl을 입력해주세요.")
		return
	}

	result, err := h.service.FetchCategories(r.Context(), req)
	if err != nil {
		if errors.Is(err, publishertistory.ErrSessionRequired) {
			h.writeErr(w, http.StatusBadRequest, "TISTORY_SESSION_REQUIRED", "Tistory 세션을 먼저 연결해주세요.")
			return
		}
		h.writeErr(w, http.StatusInternalServerError, "TISTORY_CATEGORY_FETCH_FAILED", "카테고리 목록을 불러오지 못했습니다.")
		return
	}

	h.writeJSON(w, http.StatusOK, result)
}

func (h *Handler) StartSession(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.StartSession(r.Context())
	if err != nil {
		h.writeErr(w, http.StatusInternalServerError, "TISTORY_SESSION_START_FAILED", "Tistory 로그인 브라우저를 열지 못했습니다.")
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

func (h *Handler) ConfirmSession(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.ConfirmSession(r.Context())
	if err != nil {
		h.writeErr(w, http.StatusBadRequest, "TISTORY_SESSION_INVALID", "Tistory 세션을 확인하지 못했습니다.")
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

func (h *Handler) SessionStatus(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.SessionStatus(r.Context())
	if err != nil {
		h.writeErr(w, http.StatusInternalServerError, "TISTORY_SESSION_STATUS_FAILED", "Tistory 세션 상태를 확인하지 못했습니다.")
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

func (h *Handler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	if err := h.service.DeleteSession(); err != nil {
		h.writeErr(w, http.StatusInternalServerError, "TISTORY_SESSION_DELETE_FAILED", "Tistory 세션을 삭제하지 못했습니다.")
		return
	}
	h.writeJSON(w, http.StatusOK, map[string]bool{"deleted": true})
}
