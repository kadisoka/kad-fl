package rest

type ErrorResponse struct {
	Code string `json:"code,omitempty"`

	// We use the term description because it describes the error
	// to the developer rather than a message for the end user.
	Description string `json:"description,omitempty"`

	Fields []ErrorResponseField `json:"fields,omitempty"`
	DocURL string               `json:"doc_url,omitempty"`
}

func (err ErrorResponse) Error() string {
	if code := err.Code; code != "" {
		if desc := err.Description; desc != "" {
			return "[" + code + "] " + desc
		}
		return "code " + code
	}
	if desc := err.Description; desc != "" {
		return desc
	}
	return "<empty error>"
}

type ErrorResponseField struct {
	Field       string `json:"field"`
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	DocURL      string `json:"doc_url,omitempty"`
}

func (err ErrorResponseField) Error() string {
	if field := err.Field; field != "" {
		if code := err.Code; code != "" {
			if desc := err.Description; desc != "" {
				return field + ": [" + code + "] " + desc
			}
			return field + ": " + code
		}
		if desc := err.Description; desc != "" {
			return field + ": " + desc
		}
		return "field " + field
	}
	if code := err.Code; code != "" {
		if desc := err.Description; desc != "" {
			return "[" + code + "] " + desc
		}
		return "code " + code
	}
	if desc := err.Description; desc != "" {
		return desc
	}
	return "<empty field error>"
}
