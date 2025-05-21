package main

import (
	"bufio"
	"log"
	"context"
	"fmt"
	"os"

	"google.golang.org/genai"
)

const USER = "user"
const MODEL = "model"

func main() {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GEMINI_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})

	paniferr(err)

	rawMessageCap := os.Getenv("MESSAGE_CAP")

	if rawMessageCap == "" {
		rawMessageCap = "50"
	}

	model := "gemini-2.0-flash"
	content := []*genai.Content{}

	infifo := os.Args[1]
	outfifo := os.Args[2]

	inputFile, err := os.OpenFile(infifo, os.O_RDWR, os.ModeNamedPipe)
	if err != nil {
		log.Fatal("Failed to open input FIFO:", err)
	}
	defer inputFile.Close()

	outputFile, err := os.OpenFile(outfifo, os.O_RDWR, os.ModeNamedPipe)
	if err != nil {
		log.Fatal("Failed to open output FIFO:", err)
	}
	defer outputFile.Close()

	s := bufio.NewScanner(inputFile)
	for s.Scan() {
		content = append(content, &genai.Content{Parts: []*genai.Part{
			{Text: s.Text()},
		}, Role: USER})
		r, err := client.Models.GenerateContent(ctx, model, content, nil)
		paniferr(err)
		result := r.Candidates[0].Content.Parts[0].Text

		content = append(content, &genai.Content{Parts: []*genai.Part{
			{Text: result},
		}, Role: MODEL})

		outputFile.WriteString(fmt.Sprintf("%s", result))
	}

}

func paniferr(err error) {
	if err != nil {
		panic(err)
	}
}
