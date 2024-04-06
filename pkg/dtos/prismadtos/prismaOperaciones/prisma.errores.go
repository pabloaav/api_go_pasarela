package prismadtos

import "fmt"

type PrismaErrorResponse struct {
	Type   string `json:"type,omitempty"`
	Reason Reason `json:"reason,omitempty"`
}

type Reason struct {
	ID                    int64  `json:"id,omitempty"`
	Description           string `json:"description,omitempty"`
	AdditionalDescription string `json:"additional_description,omitempty"`
}

func (e PrismaErrorResponse) Error() string {
	if e.Reason.AdditionalDescription != "" {
		return fmt.Sprintf(" error %s, %s, %s", e.Type, e.Reason.Description, e.Reason.AdditionalDescription)
	}
	return fmt.Sprintf(" error %s, %s", e.Type, e.Reason.Description)
}
