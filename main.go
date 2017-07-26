// Copyright 2016 Google Inc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/pubsub"
)

// Publish request body to a Google Cloud Pub/Sub topic.
type pubSubHandler struct {
	topic        string
	pubsubClient *pubsub.Client
}

// PubSubHandler returns a request handler that publishes
// each request body it receives to the given pub/sub topic.
func PubSubHandler(topic string, pubsubClient *pubsub.Client) http.Handler {
	return &pubSubHandler{topic, pubsubClient}
}

func (ph *pubSubHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	topic := ph.pubsubClient.Topic(ph.topic)

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to extract message from request: %v", err)
		http.Error(w, "Failed to extract message from request", http.StatusInternalServerError)
		return
	}

	msgIDs, err := topic.Publish(context.Background(), &pubsub.Message{
		Data: data,
	}).Get(context.Background())

	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		http.Error(w, "Failed to publish message", http.StatusInternalServerError)
		return
	}

	log.Printf("Published a message with a message ID: %s\n", msgIDs[0])
}

func main() {
	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		log.Fatalf("PROJECT_ID must be set and non-empty.")
	}

	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		log.Fatalf("GOOGLE_APPLICATION_CREDENTIALS must be set and non-empty.")
	}

	topic := os.Getenv("TOPIC")
	if topic == "" {
		log.Fatalf("TOPIC must be set and non-empty")
	}

	pubsubClient, err := pubsub.NewClient(context.Background(), projectID)
	if err != nil {
		log.Fatalf("Failed to create pubsub client: %v", err)
	}

	http.Handle("/", PubSubHandler(topic, pubsubClient))

	server := &http.Server{Addr: ":8080"}
	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("Shutdown signal received, exiting...")
	server.Shutdown(context.Background())
}
