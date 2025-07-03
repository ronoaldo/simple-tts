
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

func main() {
	staticContent, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", http.FileServer(http.FS(staticContent)))

	http.HandleFunc("/say", func(w http.ResponseWriter, r *http.Request) {
		say := r.URL.Query().Get("say")
		if say == "" {
			http.Error(w, "The 'say' query parameter is required", http.StatusBadRequest)
			return
		}

		ctx := context.Background()
		client, err := texttospeech.NewClient(ctx)
		if err != nil {
			log.Printf("Failed to create client: %v", err)
			http.Error(w, "Failed to create texttospeech client", http.StatusInternalServerError)
			return
		}
		defer client.Close()

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

		resp, err := client.SynthesizeSpeech(ctx, req)
		if err != nil {
			log.Printf("Failed to synthesize speech: %v", err)
			http.Error(w, "Failed to synthesize speech", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "audio/mpeg")
		w.Write(resp.AudioContent)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
