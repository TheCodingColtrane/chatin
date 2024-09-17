package request

import (
	"encoding/json"
	"net/http"
)

// type requestModel struct{}

func GetBody(req *http.Request, data any) error {
	jsonData := json.NewDecoder(req.Body)
	err := jsonData.Decode(data)
	if err != nil {
		return err
	}

	return nil

}
