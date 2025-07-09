package main

import (
	"context"
	"flag"
	"testing"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
)

var integration = flag.Bool("integration", false, "run integration tests")

func TestIntegrationListVoices(t *testing.T) {
	if !*integration {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		t.Fatalf("Failed to create texttospeech client: %v", err)
	}
	defer client.Close()

	realClient := &realTTSClient{client: client}

	req := &texttospeechpb.ListVoicesRequest{
		LanguageCode: "pt-BR",
	}
	resp, err := realClient.ListVoices(ctx, req)
	if err != nil {
		t.Fatalf("Failed to list voices: %v", err)
	}

	if len(resp.Voices) == 0 {
		t.Errorf("Expected at least one voice, got 0")
	}
}
