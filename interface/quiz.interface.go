package schemaInterface

import "time"

type QuestionType string

const (
	SingleSelect QuestionType = "singleSelect"
	MultiSelect  QuestionType = "multiSelect"
)

type QuestionStatus string

const (
	Active   QuestionStatus = "active"
	Archived QuestionStatus = "archived"
)

type BaseSchema struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type OptionSchema struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	IsCorrect bool   `json:"isCorrect"`
}

type AttachmentSchema struct {
	URL string `json:"url"`
}

type ThemeSchema struct {
	BaseSchema
	Name string `json:"name"`
}

type QuestionSchema struct {
	BaseSchema
	QuestionText string            `json:"questionText"`
	QuestionType QuestionType      `json:"questionType"`
	Options      []OptionSchema    `json:"options"`
	Attachment   *AttachmentSchema `json:"attachment,omitempty"` // Optional attachment
	Status       QuestionStatus    `json:"status"`
	ThemeID      string            `json:"themeId"`
}
