package common

type IDNameExtensions struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

type Department struct {
	ID       int    `json:"id"`
	Name     string `json:"name,omitempty"`
	IdpID    int    `json:"idpId"`
	Comments string `json:"comments,omitempty"`
	Deleted  bool   `json:"deleted"`
}

type Groups struct {
	ID       int    `json:"id"`
	Name     string `json:"name,omitempty"`
	IdpID    int    `json:"idpId"`
	Comments string `json:"comments,omitempty"`
}
