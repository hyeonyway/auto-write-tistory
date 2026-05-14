package tistory

type Login struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

type FetchCategoriesRequest struct {
	BlogURL string `json:"blogUrl"`
	Login   Login  `json:"login"`
}

type CategoryOption struct {
	CategoryID string `json:"categoryId"`
	Label      string `json:"label"`
}

type FetchCategoriesResponse struct {
	Items []CategoryOption `json:"items"`
}
