package gae

import (
	"code.google.com/p/go.net/context"
	"encoding/json"
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/blobstore"
	"google.golang.org/appengine/file"
	"google.golang.org/appengine/log"
	"google.golang.org/cloud/storage"
	"io/ioutil"
	"mime/multipart"
	"strings"
)

type GStorage struct {
	ctx    context.Context
	client *storage.Client
	// bucket is the Google Cloud Storage bucket name used for the GStorage.
	bucket *storage.BucketHandle

	bname string
}

func NewGStorage(c context.Context) (*GStorage, error) {
	bucketName, err := file.DefaultBucketName(c)
	if err != nil {
		log.Errorf(c, "failed to get default GCS bucket name: %v", err)
		return nil, err
	}

	client, err := storage.NewClient(c)
	if err != nil {
		return nil, err
	}

	gs := &GStorage{
		ctx:    c,
		client: client,
		bucket: client.Bucket(bucketName),
		bname:  bucketName,
	}
	return gs, nil
}

func (gs GStorage) BucketName() string {
	return gs.bname
}

func (gs *GStorage) CreateImageFile(fileName string, img multipart.File, ct string) error {
	wc := gs.bucket.Object(fileName).NewWriter(gs.ctx)
	wc.ContentType = ct

	b, err := ioutil.ReadAll(img)
	if err != nil {
		return err
	}

	if _, err := wc.Write(b); err != nil {
		return fmt.Errorf("createFile: unable to write data to bucket %q, file %q: %v", gs.BucketName(), fileName, err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("createFile: unable to close bucket %q, file %q: %v", gs.BucketName(), fileName, err)
	}

	return nil
}

// CreateFile creates a file in Google Cloud Storage.
func (gs *GStorage) CreateJsonFile(fileName string, v interface{}) error {
	wc := gs.bucket.Object(fileName).NewWriter(gs.ctx)
	wc.ContentType = "application/json"
	wc.CacheControl = "private, max-age=0, no-transform"

	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("putVendorStaticFile failed, due to json.Marshal", err)
	}

	if _, err := wc.Write(b); err != nil {
		return fmt.Errorf("createFile: unable to write data to bucket %q, file %q: %v", gs.BucketName(), fileName, err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("createFile: unable to close bucket %q, file %q: %v", gs.BucketName(), fileName, err)
	}

	return nil
}

func (gs *GStorage) CopyBlob(src appengine.BlobKey, v ImageBlober) (appengine.BlobKey, error) {
	sObj, err := gs.ReadBlobKey(src)
	if err != nil {
		return "", err
	}

	var dest = gs.bucket.Object(v.BucketPath())
	if _, err = gs.bucket.Object(sObj.Name).CopyTo(gs.ctx, dest, nil); err != nil {
		return "", err
	}

	gcsFilename := "/" + strings.Join([]string{"gs", sObj.Bucket, v.BucketPath()}, "/")
	return blobstore.BlobKeyForFile(gs.ctx, gcsFilename)
}

func (gs *GStorage) ReadBlobKey(src appengine.BlobKey) (*storage.ObjectAttrs, error) {
	info, err := blobstore.Stat(gs.ctx, src)
	if err != nil {
		return nil, err
	}

	bucket := appengine.DefaultVersionHostname(gs.ctx)
	name := strings.TrimPrefix(info.ObjectName, "/"+bucket+"/")

	return gs.bucket.Object(name).Attrs(gs.ctx)
}
