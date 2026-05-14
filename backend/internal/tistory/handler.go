package tistory

import (
	"net/http"
	"strings"
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
	if strings.TrimSpace(req.BlogURL) == "" || strings.TrimSpace(req.Login.ID) == "" || req.Login.Password == "" {
		h.writeErr(w, http.StatusBadRequest, "VALIDATION_ERROR", "blogUrl과 login 정보를 입력해주세요.")
		return
	}

	result, err := h.service.FetchCategories(r.Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), "kakao authentication is still required") {
			h.writeErr(w, http.StatusBadRequest, "TISTORY_AUTH_REQUIRED", "카카오 추가 인증이 필요합니다. Selenium 브라우저에서 인증을 완료한 뒤 다시 시도해주세요.")
			return
		}
		if strings.Contains(err.Error(), "tistory login failed") {
			h.writeErr(w, http.StatusBadRequest, "TISTORY_LOGIN_FAILED", "Tistory 로그인에 실패했습니다. 카카오 로그인 정보를 확인해주세요.")
			return
		}
		h.writeErr(w, http.StatusInternalServerError, "TISTORY_CATEGORY_FETCH_FAILED", "카테고리 목록을 불러오지 못했습니다.")
		return
	}

	h.writeJSON(w, http.StatusOK, result)
}
