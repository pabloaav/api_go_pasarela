package prismadtos

import "fmt"

type StatusDetails struct {
	Ticket                string       `json:"ticket,omitempty"`
	CardAuthorizationCode string       `json:"card_authorization_code,omitempty"`
	AddressValidationCode string       `json:"address_validation_code,omitempty"`
	Error                 ErrorDetails `json:"error,omitempty"`
}

type ErrorDetails struct {
	Type   string `json:"type"`
	Reason Reason `json:"reason"`
}

type Reason struct {
	ID                    int64  `json:"id"`
	Description           string `json:"description"`
	AdditionalDescription string `json:"additional_description"`
}

func (e ErrorDetails) Error() string {
	if e.Reason.AdditionalDescription != "" {
		return fmt.Sprintf(" error %s, %s, %s", e.Type, e.Reason.Description, e.Reason.AdditionalDescription)
	}
	return fmt.Sprintf(" error %s, %s", e.Type, e.Reason.Description)
}
