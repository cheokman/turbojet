package content

type PropertyContent struct {
	PropertyID string    `json:"property_id"`
	Contents   []Content `json:"contents"`
}

func NewPropertyContent() *PropertyContent {
	return &PropertyContent{}
}
