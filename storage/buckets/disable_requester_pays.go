// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package buckets

// [START storage_disable_requester_pays]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// disableRequesterPays sets requester pays flag to false.
func disableRequesterPays(w io.Writer, bucketName string) error {
	// bucketName := "bucket-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	bucket := client.Bucket(bucketName)

	// Set a metageneration-match precondition. The request to update is aborted
	// if the bucket's metageneration number does not match your precondition
	// criteria. This avoids race conditions and data corruption.
	attrs, err := bucket.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("bucket.Attrs: %v", err)
	}
	bucket = bucket.If(storage.BucketConditions{MetagenerationMatch: attrs.MetaGeneration})

	// Update the bucket to disable requester pays.
	bucketAttrsToUpdate := storage.BucketAttrsToUpdate{
		RequesterPays: false,
	}
	if _, err := bucket.Update(ctx, bucketAttrsToUpdate); err != nil {
		return fmt.Errorf("Bucket(%q).Update: %v", bucketName, err)
	}
	fmt.Fprintf(w, "Requester pays disabled for bucket %v\n", bucketName)
	return nil
}

// [END storage_disable_requester_pays]
