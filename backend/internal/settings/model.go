package settings

type TistorySettings struct {
	BlogURL            string             `json:"blogUrl"`
	CategoryID         string             `json:"categoryId"`
	CategoryLabel      string             `json:"categoryLabel"`
	DefaultVisibility  int                `json:"defaultVisibility"`
	DefaultTags        []string           `json:"defaultTags"`
	Categories         []TistoryCategory  `json:"categories"`
	CategoriesSyncedAt string             `json:"categoriesSyncedAt"`
	Session            TistorySessionInfo `json:"session"`
}

type TistoryCategory struct {
	CategoryID string `json:"categoryId"`
	Label      string `json:"label"`
}

type TistorySessionInfo struct {
	Connected bool   `json:"connected"`
	Message   string `json:"message"`
}

type Response struct {
	Tistory TistorySettings `json:"tistory"`
}
