// Copyright 2022 Google LLC
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

package notifications

import (
	"bytes"
	"context"
	"log"
	"strings"
	"testing"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func createTestTopic(t *testing.T, projectID string) string {
	t.Helper()
	ctx := context.Background()

	pubsubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}
	defer pubsubClient.Close()

	//topicName := MustGetEnv(t, "PUBSUB_TOPIC")
	topicName := testutil.MustGetEnv(t, "PUBSUB_TOPIC")
	topic := pubsubClient.Topic(topicName)

	// Create the topic if it doesn't exist.
	exists, err := topic.Exists(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if !exists {
		//log.Printf("Topic %v doesn't exist - creating it", topicName)
		_, err = pubsubClient.CreateTopic(ctx, topicName)
		if err != nil {
			log.Fatal(err)
		}
	}
	return topicName
}

func TestNotifications(t *testing.T) {
	tc := testutil.SystemTest(t)
	projectID := tc.ProjectID
	ctx := context.Background()

	topic := createTestTopic(t, projectID)

	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	defer client.Close()

	bucket, err := testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, "storage-buckets-test")
	if err != nil {
		t.Fatalf("creating bucket: %v", err)
	}
	defer testutil.DeleteBucketIfExists(ctx, client, bucket)

	var buf bytes.Buffer
	if err := createBucketNotification(&buf, tc.ProjectID, bucket, topic); err != nil {
		t.Errorf("createBucketNotification: %v", err)
	}

	if got, want := buf.String(), "created notification with ID"; !strings.Contains(got, want) {
		t.Errorf("createBucketNotification: got %q; want to contain %q", got, want)
	}
}
