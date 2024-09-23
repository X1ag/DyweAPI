package models

type Collection struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Creator     string `json:"creator"`
	ImageURL    string `json:"image_url"`
}
