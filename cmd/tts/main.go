
package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
)

//go:embed static
var staticFiles embed.FS

type ttsClient interface {
	SynthesizeSpeech(context.Context, *texttospeechpb.SynthesizeSpeechRequest) (*texttospeechpb.SynthesizeSpeechResponse, error)
}

type realTTSClient struct {
	client *texttospeech.Client
}

func (c *realTTSClient) SynthesizeSpeech(ctx context.Context, req *texttospeechpb.SynthesizeSpeechRequest) (*texttospeechpb.SynthesizeSpeechResponse, error) {
	return c.client.SynthesizeSpeech(ctx, req)
}

type server struct {
	tts ttsClient
}

func (s *server) sayHandler(w http.ResponseWriter, r *http.Request) {
	say := r.URL.Query().Get("say")
	if say == "" {
		http.Error(w, "The 'say' query parameter is required", http.StatusBadRequest)
		return
	}

	req := &texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: say},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "pt-BR",
			Name:         "pt-BR-Wavenet-B",
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
		},
	}

	resp, err := s.tts.SynthesizeSpeech(r.Context(), req)
	if err != nil {
		log.Printf("Failed to synthesize speech: %v", err)
		http.Error(w, "Failed to synthesize speech", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "audio/mpeg")
	w.Write(resp.AudioContent)
}

func main() {
	staticContent, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", http.FileServer(http.FS(staticContent)))

	ctx := context.Background()
	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create texttospeech client: %v", err)
	}
	defer client.Close()

	s := &server{
		tts: &realTTSClient{client: client},
	}

	http.HandleFunc("/say", s.sayHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
