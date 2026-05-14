package tistory

type FetchCategoriesRequest struct {
	BlogURL string `json:"blogUrl"`
}

type CategoryOption struct {
	CategoryID string `json:"categoryId"`
	Label      string `json:"label"`
}

type FetchCategoriesResponse struct {
	Items []CategoryOption `json:"items"`
}

type SessionStartResponse struct {
	Connected bool   `json:"connected"`
	Message   string `json:"message"`
}

type SessionStatusResponse struct {
	Connected bool   `json:"connected"`
	Message   string `json:"message"`
}
