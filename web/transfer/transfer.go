package transfer

import "net/http"

type (
	Handler struct {
	}

	payload struct {
		From   string `json:"from"`
		To     string `json:"to"`
		Amount int    `json:"amount"`
	}

)

func NewHandler() Handler {
	return Handler{}
}

func (Handler) Handle(w http.ResponseWriter, r *http.Request) {

}
