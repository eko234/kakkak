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
	apikey := os.Getenv("GEMINI_API_KEY")

	println("KEY::::::::::::::::::::::::::::::::::")
	println(apikey)
	println(":::::::::::::::::::::::::::::::::::::")

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

	// BUFFERING:
	// 	for {
	// 		resp, err := stream.Recv()
	// 		if err != nil {
	// 			if err == io.EOF {
	// 				outputFile.WriteString("\n")
	// 				break BUFFERING
	// 			}
	// 			msg := fmt.Sprintf("ChatCompletion error: %v", err)
	// 			if debug {
	// 				outputFile.WriteString(msg)
	// 				outputFile.WriteString("\n")
	// 			}
	// 			fmt.Printf("%s\n", msg)
	// 		}
	// 		outputFile.WriteString(fmt.Sprintf("%s", resp.Choices[0].Delta.Content))
	// 		msg.Content = msg.Content + resp.Choices[0].Delta.Content
	// 	}

	// 	req.Messages = append(req.Messages, msg)

	// 	ran := len(req.Messages) - messageCap
	// 	if ran < 0 {
	// 		ran = 0
	// 	}

	// 	req.Messages = req.Messages[ran:]
	// }

	// s := bufio.NewScanner(inputFile)
	// for s.Scan() {
	// 	req.Messages = append(req.Messages, openai.ChatCompletionMessage{
	// 		Role:    openai.ChatMessageRoleUser,
	// 		Content: s.Text(),
	// 	})
	// 	stream, err := client.CreateChatCompletionStream(context.Background(), req)
	// 	if err != nil {
	// 		msg := fmt.Sprintf("ChatCompletion error: %v", err)
	// 		if debug {
	// 			outputFile.WriteString(msg)
	// 			outputFile.WriteString("\n")
	// 		}
	// 		fmt.Printf("%s\n", msg)
	// 		continue
	// 	}

	// 	msg := openai.ChatCompletionMessage{Role: openai.ChatMessageRoleAssistant}

	// BUFFERING:
	// 	for {
	// 		resp, err := stream.Recv()
	// 		if err != nil {
	// 			if err == io.EOF {
	// 				outputFile.WriteString("\n")
	// 				break BUFFERING
	// 			}
	// 			msg := fmt.Sprintf("ChatCompletion error: %v", err)
	// 			if debug {
	// 				outputFile.WriteString(msg)
	// 				outputFile.WriteString("\n")
	// 			}
	// 			fmt.Printf("%s\n", msg)
	// 		}
	// 		outputFile.WriteString(fmt.Sprintf("%s", resp.Choices[0].Delta.Content))
	// 		msg.Content = msg.Content + resp.Choices[0].Delta.Content
	// 	}

	// 	req.Messages = append(req.Messages, msg)

	// 	ran := len(req.Messages) - messageCap
	// 	if ran < 0 {
	// 		ran = 0
	// 	}

	// 	req.Messages = req.Messages[ran:]
	// }
}

func paniferr(err error) {
	if err != nil {
		panic(err)
	}
}
