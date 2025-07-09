
package main

import (
	"context"
	"encoding/json"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"

	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
)

// mockTTSClient is a mock implementation of the ttsClient interface for testing.
type mockTTSClient struct {
	SynthesizeSpeechFunc func(context.Context, *texttospeechpb.SynthesizeSpeechRequest) (*texttospeechpb.SynthesizeSpeechResponse, error)
	ListVoicesFunc       func(context.Context, *texttospeechpb.ListVoicesRequest) (*texttospeechpb.ListVoicesResponse, error)
}

func (m *mockTTSClient) SynthesizeSpeech(ctx context.Context, req *texttospeechpb.SynthesizeSpeechRequest) (*texttospeechpb.SynthesizeSpeechResponse, error) {
	return m.SynthesizeSpeechFunc(ctx, req)
}

func (m *mockTTSClient) ListVoices(ctx context.Context, req *texttospeechpb.ListVoicesRequest) (*texttospeechpb.ListVoicesResponse, error) {
	return m.ListVoicesFunc(ctx, req)
}

func TestStaticFileServer(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// We can't easily test the embedded file system,
	// so we'll just check that the handler is registered
	// and returns a 200 status code.
	staticContent, err := fs.Sub(staticFiles, "static")
	if err != nil {
		t.Fatal(err)
	}
	handler := http.FileServer(http.FS(staticContent))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestSayHandler_MissingQueryParam(t *testing.T) {
	s := &server{} // Create a server with a nil tts client, as it won't be used in this test.
	req, err := http.NewRequest("GET", "/say", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.sayHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestSayHandler_Success(t *testing.T) {
	mockClient := &mockTTSClient{
		SynthesizeSpeechFunc: func(ctx context.Context, req *texttospeechpb.SynthesizeSpeechRequest) (*texttospeechpb.SynthesizeSpeechResponse, error) {
			return &texttospeechpb.SynthesizeSpeechResponse{
				AudioContent: []byte("fake audio content"),
			}, nil
		},
	}

	s := &server{tts: mockClient}
	req, err := http.NewRequest("GET", "/say?say=hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.sayHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Header().Get("Content-Type") != "audio/mpeg" {
		t.Errorf("handler returned wrong content type: got %v want %v",
			rr.Header().Get("Content-Type"), "audio/mpeg")
	}

	if string(rr.Body.Bytes()) != "fake audio content" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			string(rr.Body.Bytes()), "fake audio content")
	}
}

func TestVoicesHandler_Success(t *testing.T) {
	mockClient := &mockTTSClient{
		ListVoicesFunc: func(ctx context.Context, req *texttospeechpb.ListVoicesRequest) (*texttospeechpb.ListVoicesResponse, error) {
			return &texttospeechpb.ListVoicesResponse{
				Voices: []*texttospeechpb.Voice{
					{Name: "pt-BR-Chirp-A"},
					{Name: "pt-BR-Chirp-B"},
				},
			}, nil
		},
	}

	s := &server{tts: mockClient}
	req, err := http.NewRequest("GET", "/voices", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.voicesHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var voices []*texttospeechpb.Voice
	if err := json.NewDecoder(rr.Body).Decode(&voices); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if len(voices) != 2 {
		t.Errorf("Expected 2 voices, got %d", len(voices))
	}
}



