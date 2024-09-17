package response

import (
	"encoding/json"
	"net/http"
)

// type responseModel struct{}

func JSON(res http.ResponseWriter, data any) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(200)
	var result = json.NewEncoder(res)
	result.Encode(data)
}
