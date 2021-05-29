package quote

import "encoding/json"

func UnmarshalQuote(data []byte) (Quote, error) {
	var r Quote
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Quote) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Quote struct {
	Title   *string `json:"title,omitempty"`
	Content *string `json:"content,omitempty"`
}
