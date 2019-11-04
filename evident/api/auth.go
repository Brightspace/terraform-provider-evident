package api

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

func NewHMAC(message []byte, key []byte) string {
	hash := hmac.New(sha1.New, key)
	hash.Write(message)
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func NewHTTPSignature(url string, method string, contents []byte, date time.Time, public string, secret string) (map[string]interface{}, error) {
	headers := make(map[string]interface{})

	contentType := "application/vnd.api+json"
	utc := date.Format(http.TimeFormat)
	md5sum := md5.Sum(contents)
	md5 := base64.StdEncoding.EncodeToString(md5sum[:])
	request := fmt.Sprintf("%s,%s,%s,%s,%s", method, contentType, md5, url, utc)
	encodedAuth := NewHMAC([]byte(request), []byte(secret))

	headers["Accept"] = contentType
	headers["Content-Type"] = contentType
	headers["Content-MD5"] = md5
	headers["Date"] = utc
	headers["Authorization"] = fmt.Sprintf("APIAuth %s:%s", public, encodedAuth)

	return headers, nil
}
