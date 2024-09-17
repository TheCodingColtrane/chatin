package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"
)

type Asset struct {
	ChatToken    string    `json:"chatToken"`
	FileName     string    `json:"fileName"`
	MimeType     string    `json:"mimeType"`
	Content      string    `json:"content"`
	LastModified time.Time `json:"lastModified"`
	Size         int64
}

const (
	ASSETSPATH     = "./assets"
	AUDIOSPATH     = "./assets/audios/"
	PICTURESPATH   = "./assets/pictures/"
	VIDEOSPATH     = "./assets/videos/"
	AUDIOMIMETYPE  = "audio/webm"
	JPEGMIMETYPE   = "image/jpeg"
	PNGMIMETYPE    = "image/png"
	MP4MIMETYPE    = "video/mp4"
	AUDIOEXTENSION = ".webm"
	JPEGEXTENSION  = ".jpeg"
	PNGEXTENSION   = ".png"
	MP4EXTENSION   = ".mp4"
)

func main() {
	http.HandleFunc("POST /asset", func(res http.ResponseWriter, req *http.Request) {
		var decoder = json.NewDecoder(req.Body)
		var createdAsset Asset
		err := decoder.Decode(&createdAsset)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
		}

		createdAsset, err = createAsset(&createdAsset)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
		}
		var fileCreation struct {
			Success      bool  `json:"success"`
			CreatedAsset Asset `json:"asset"`
		}
		fileCreation.Success = true
		fileCreation.CreatedAsset = createdAsset
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
		}
		// var header http.Header
		res.WriteHeader(200)
		var result = json.NewEncoder(res)
		result.Encode(&fileCreation)
	})

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.ListenAndServe("localhost:8001", nil)
}

func checkFolders(asset Asset) error {
	var basePath = ASSETSPATH + "/" + asset.ChatToken
	var _, err = os.Stat(basePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = os.Mkdir(basePath, os.ModePerm)
			if err != nil {
				return err
			}
			var folders struct {
				audioFolderChanErr error
				picsFolderChanErr  error
				videoFolderChanErr error
			}

			var (
				baseAudioPath = basePath + "/audios"
				baseVideoPath = basePath + "/videos"
				basePicsPath  = basePath + "/pictures"
			)

			folders.audioFolderChanErr = os.Mkdir(baseAudioPath, os.ModePerm)
			folders.picsFolderChanErr = os.Mkdir(basePicsPath, os.ModePerm)
			folders.videoFolderChanErr = os.Mkdir(baseVideoPath, os.ModePerm)

			if folders.audioFolderChanErr != nil || folders.picsFolderChanErr != nil || folders.videoFolderChanErr != nil {
				return folders.audioFolderChanErr
			}

			return nil

		}
	}

	return nil

}

func createAsset(asset *Asset) (Asset, error) {
	var err = checkFolders(*asset)
	if err != nil {
		return Asset{}, err
	}
	var basePath = ASSETSPATH + "/" + asset.ChatToken

	if asset.MimeType == AUDIOMIMETYPE {
		asset.FileName = basePath + "/audios/" + asset.FileName
	} else if asset.MimeType == JPEGMIMETYPE || asset.MimeType == PNGMIMETYPE {
		asset.FileName = basePath + "/pictures/" + asset.FileName
	} else {
		asset.FileName = basePath + "/videos/" + asset.FileName
	}

	file, err := os.Create(asset.FileName)
	if err != nil {
		return Asset{}, err
	}
	var base64Position = len("data:" + asset.MimeType + ";base64,")
	base64Data := asset.Content[base64Position:]
	fileData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return Asset{}, err
	}

	_, err = file.Write(fileData)
	if err != nil {
		return Asset{}, err
	}

	return *asset, err
}
