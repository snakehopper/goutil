package gae

import (
	"appengine"
	"appengine/blobstore"
	"appengine/datastore"
	"appengine/file"
	"appengine/urlfetch"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
	"strings"
)

type ImageBlober interface {
	BucketPath() string
	FileName() string
	DatastoreKey() *datastore.Key
	UploadHandler() string
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func createFromPart(name string, w *multipart.Writer, h *multipart.Part) (io.Writer, error) {
	ext := filepath.Ext(h.FileName())
	filename := name + ext
	fieldname := "file"
	contentType := h.Header.Get("Content-Type")

	newHeader := make(textproto.MIMEHeader)
	newHeader.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))
	newHeader.Set("Content-Type", contentType)
	return w.CreatePart(newHeader)
}

func PutBlob(c appengine.Context, v ImageBlober, r *http.Request) error {
	bucket, err := file.DefaultBucketName(c)
	if err != nil {
		return err
	}

	opt := &blobstore.UploadURLOptions{StorageBucket: bucket + "/" + v.BucketPath()}
	u, err := blobstore.UploadURL(c, v.UploadHandler(), opt)
	if err != nil {
		return err
	}

	mr, err := r.MultipartReader()
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	for p, err := mr.NextPart(); err != io.EOF; p, err = mr.NextPart() {
		fw, _ := createFromPart(v.FileName(), w, p)
		if _, err := io.Copy(fw, p); err != nil {
			return err
		}
	}

	k := v.DatastoreKey()
	w.WriteField("datastoreKey", k.Encode())
	w.Close()

	req, err := http.NewRequest("POST", u.String(), &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := urlfetch.Client(c)
	if _, err := client.Do(req); err != nil {
		return err
	}

	return nil
}
