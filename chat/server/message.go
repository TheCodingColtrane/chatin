package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"messenger/models"
	"net/http"
	"time"
)

type AssetUpload struct {
	ChatToken    string    `json:"chatToken"`
	FileName     string    `json:"fileName"`
	MimeType     string    `json:"mimeType"`
	Content      string    `json:"content"`
	LastModified time.Time `json:"lastModified"`
	Size         int64
}

func CreateAsset(message models.IncomingMessage) (AssetUpload, error) {
	assetServer := "http://localhost:8001/asset" // Include the protocol
	if message.Asset.Name == "" {
		message.Asset.Name = createFileName(message)
	}
	asset := AssetUpload{
		ChatToken:    message.ChatCode,
		FileName:     message.Asset.Name,
		MimeType:     message.MimeType,
		Content:      message.Asset.Content,
		LastModified: message.Asset.LastModified,
		Size:         message.Asset.Size,
	}

	jsonData, err := json.Marshal(&asset)
	if err != nil {
		return AssetUpload{}, err
	}

	data := bytes.NewReader(jsonData)
	req, err := http.NewRequest(http.MethodPost, assetServer, data)
	if err != nil {
		return AssetUpload{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return AssetUpload{}, err
	}
	var fileCreated struct {
		Success      bool        `json:"success"`
		CreatedAsset AssetUpload `json:"asset"`
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return AssetUpload{}, err
	}

	err = json.Unmarshal(body, &fileCreated)
	if err != nil {
		return AssetUpload{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AssetUpload{}, errors.New("failed to upload asset")
	}
	if fileCreated.Success {
		return fileCreated.CreatedAsset, nil

	}

	return AssetUpload{}, nil
}

func createFileName(message models.IncomingMessage) string {
	year := fmt.Sprint(time.Now().Year())
	month := fmt.Sprint(int(time.Now().Month()))
	day := fmt.Sprint(time.Now().Day())
	minute := fmt.Sprint(time.Now().Minute())
	hour := fmt.Sprint(time.Now().Hour())
	second := fmt.Sprint(time.Now().Second())
	var extension = ""
	if message.MimeType == "audio/webm" {
		extension = ".webm"
	} else if message.MimeType == "image/jpeg" {
		extension = ".jpeg"
	} else if message.MimeType == "image/png" {
		extension = ".png"
	} else {
		extension = ".mp4"
	}
	message.Asset.Name = year + month + day + "_" + hour + minute + second + extension
	return message.Asset.Name
}
