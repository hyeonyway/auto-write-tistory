package settings

type TistorySettings struct {
	BlogURL           string   `json:"blogUrl"`
	CategoryID        string   `json:"categoryId"`
	DefaultVisibility int      `json:"defaultVisibility"`
	DefaultTags       []string `json:"defaultTags"`
}

type Response struct {
	Tistory TistorySettings `json:"tistory"`
}
