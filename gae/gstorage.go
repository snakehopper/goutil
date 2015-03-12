package gae

import (
	"appengine"
	"appengine/blobstore"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"

	"golang.org/x/net/context"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/appengine/file"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/cloud"
	"google.golang.org/cloud/storage"
)

type GStorage struct {
	c   appengine.Context
	ctx context.Context
	// bucket is the Google Cloud Storage bucket name used for the GStorage.
	bucket string
	// failed indicates that one or more of the GStorage steps failed.
	failed bool
}

func NewGStorage(c appengine.Context) (*GStorage, error) {
	ctxFromGAE := c.(context.Context)
	bucketName, err := file.DefaultBucketName(ctxFromGAE)
	if err != nil {
		c.Errorf("failed to get default GCS bucket name: %v", err)
		return nil, err
	}

	hc := &http.Client{
		Transport: &oauth2.Transport{
			Source: google.AppEngineTokenSource(ctxFromGAE, storage.ScopeFullControl),
			Base:   &urlfetch.Transport{Context: ctxFromGAE},
		},
	}
	ctx := cloud.NewContext(appengine.AppID(c), hc)

	gs := &GStorage{
		c:      c,
		ctx:    ctx,
		bucket: bucketName,
	}
	return gs, nil
}

func (gs *GStorage) CreateImageFile(fileName string, img multipart.File, ct string) error {
	wc := storage.NewWriter(gs.ctx, gs.bucket, fileName)
	wc.ContentType = ct

	b, err := ioutil.ReadAll(img)
	if err != nil {
		return err
	}

	if _, err := wc.Write(b); err != nil {
		return fmt.Errorf("createFile: unable to write data to bucket %q, file %q: %v", gs.bucket, fileName, err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("createFile: unable to close bucket %q, file %q: %v", gs.bucket, fileName, err)
	}
	// Wait for the file to be fully written.
	if obj := wc.Object(); obj == nil {
		return fmt.Errorf("createFile: unable to finalize file from bucket %q, file %q", gs.bucket, fileName)
	}

	return nil
}

// CreateFile creates a file in Google Cloud Storage.
func (gs *GStorage) CreateJsonFile(fileName string, v interface{}) error {
	wc := storage.NewWriter(gs.ctx, gs.bucket, fileName)
	wc.ContentType = "application/json"
	wc.CacheControl = "private, max-age=0, no-transform"

	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("putVendorStaticFile failed, due to json.Marshal", err)
	}

	if _, err := wc.Write(b); err != nil {
		return fmt.Errorf("createFile: unable to write data to bucket %q, file %q: %v", gs.bucket, fileName, err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("createFile: unable to close bucket %q, file %q: %v", gs.bucket, fileName, err)
	}
	// Wait for the file to be fully written.
	if obj := wc.Object(); obj == nil {
		return fmt.Errorf("createFile: unable to finalize file from bucket %q, file %q", gs.bucket, fileName)
	}

	return nil
}

func (gs *GStorage) CopyBlob(src appengine.BlobKey, v ImageBlober) (appengine.BlobKey, error) {
	sObj, err := gs.ReadBlobKey(src)
	if err != nil {
		return "", err
	}

	var srcName = sObj.Name
	attrs := storage.ObjectAttrs{
		Name:            v.BucketPath(),
		ContentType:     sObj.ContentType,
		ContentLanguage: sObj.ContentLanguage,
		ContentEncoding: sObj.ContentEncoding,
		CacheControl:    sObj.CacheControl,
		ACL:             sObj.ACL,
		Metadata:        sObj.Metadata,
	}
	_, err = storage.CopyObject(gs.ctx, gs.bucket, srcName, gs.bucket, attrs)

	gcsFilename := "/" + strings.Join([]string{"gs", sObj.Bucket, sObj.Name}, "/")
	return blobstore.BlobKeyForFile(gs.c, gcsFilename)
}

func (gs *GStorage) ReadBlobKey(src appengine.BlobKey) (*storage.Object, error) {
	info, err := blobstore.Stat(gs.c, src)
	if err != nil {
		return nil, err
	}

	bucket := appengine.DefaultVersionHostname(gs.c)
	name := strings.TrimPrefix(info.ObjectName, "/"+bucket+"/")

	return storage.StatObject(gs.ctx, bucket, name)
}
