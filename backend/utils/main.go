package utils

import (
	"chatin/enums"
	"database/sql"
)

func GetMIMETypeInt(mimeType string) uint8 {
	var mime = 0
	switch mimeType {
	case "audio/webm":
		mime = enums.AUDIO
	case "image/jpeg":
		mime = enums.IMAGE
	case "image/png":
		mime = enums.IMAGE
	case "video/mp4":
		mime = enums.VIDEO
	}

	return uint8(mime)
}

func GetMIMETypeString(mimeType uint8) string {
	var mime = ""
	switch mimeType {
	case 1:
		mime = "audio/webm"
	case 2:
		mime = "image/jpeg"
	case 3:
		mime = "image/png"
	case 4:
		mime = "video/mp4"
	}

	return mime
}

func GetNullString(text sql.NullString) string {
	if text.Valid {
		return text.String
	}

	return ""
}

func GetNullInt(number sql.NullInt32) int {
	if number.Valid {
		return int(number.Int32)
	}
	return 0
}
