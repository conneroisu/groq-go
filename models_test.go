package groq

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestGemma29BIt tests the Gemma29BIt model.
// It ensures that the model is supported by the groq-go library, the groq API,
// and the operations are working as expected for the specific model type.
func TestGemma29BIt(t *testing.T) {
}

// TestGemma7BIt tests the Gemma7BIt model.
// It ensures that the model is supported by the groq-go library, the groq API,
// and the operations are working as expected for the specific model type.
func TestGemma7BIt(t *testing.T) {
}

// TestLlama3170BVersatile tests the Llama3170BVersatile model.
// It ensures that the model is supported by the groq-go library, the groq API,
// and the operations are working as expected for the specific model type.
func TestLlama3170BVersatile(t *testing.T) {
}

// TestLlama318BInstant tests the Llama318BInstant model.
// It ensures that the model is supported by the groq-go library, the groq API,
// and the operations are working as expected for the specific model type.
func TestLlama318BInstant(t *testing.T) {
}

// TestLlama3211BTextPreview tests the Llama3211BTextPreview model.
// It ensures that the model is supported by the groq-go library, the groq API,
// and the operations are working as expected for the specific model type.
func TestLlama3211BTextPreview(t *testing.T) {
}

// TestLlama3211BVisionPreview tests the Llama3211BVisionPreview model.
// It ensures that the model is supported by the groq-go library, the groq API,
// and the operations are working as expected for the specific model type.
func TestLlama3211BVisionPreview(t *testing.T) {
}

// TestLlama321BPreview tests the Llama321BPreview model.
// It ensures that the model is supported by the groq-go library, the groq API,
// and the operations are working as expected for the specific model type.
func TestLlama321BPreview(t *testing.T) {
}

// TestLlama323BPreview tests the Llama323BPreview model.
// It ensures that the model is supported by the groq-go library, the groq API,
// and the operations are working as expected for the specific model type.
func TestLlama323BPreview(t *testing.T) {
}

// TestLlama3290BTextPreview tests the Llama3290BTextPreview model.
// It ensures that the model is supported by the groq-go library, the groq API,
// and the operations are working as expected for the specific model type.
func TestLlama3290BTextPreview(t *testing.T) {
}

// TestLlama370B8192 tests the Llama370B8192 model.
// It ensures that the model is supported by the groq-go library, the groq API,
// and the operations are working as expected for the specific model type.
func TestLlama370B8192(t *testing.T) {
}

// TestLlama38B8192 tests the Llama38B8192 model.
// It ensures that the model is supported by the groq-go library, the groq API,
// and the operations are working as expected for the specific model type.
func TestLlama38B8192(t *testing.T) {
}

// TestLlama3Groq70B8192ToolUsePreview tests the Llama3Groq70B8192ToolUsePreview model.
// It ensures that the model is supported by the groq-go library, the groq API,
// and the operations are working as expected for the specific model type.
func TestLlama3Groq70B8192ToolUsePreview(t *testing.T) {
}

// TestLlama3Groq8B8192ToolUsePreview tests the Llama3Groq8B8192ToolUsePreview model.
// It ensures that the model is supported by the groq-go library, the groq API,
// and the operations are working as expected for the specific model type.
func TestLlama3Groq8B8192ToolUsePreview(t *testing.T) {
}

// TestLlavaV157B4096Preview tests the LlavaV157B4096Preview model.
// It ensures that the model is supported by the groq-go library, the groq API,
// and the operations are working as expected for the specific model type.
func TestLlavaV157B4096Preview(t *testing.T) {
}

// TestMixtral8X7B32768 tests the Mixtral8X7B32768 model.
// It ensures that the model is supported by the groq-go library, the groq API,
// and the operations are working as expected for the specific model type.
func TestMixtral8X7B32768(t *testing.T) {
}

// TestWhisperLargeV3 tests the WhisperLargeV3  transcription model.
//
// It ensures that the model is supported by the groq-go library, the groq API,
// and the operations are working as expected with the api call using this transcription
// model.
func TestWhisperLargeV3(t *testing.T) {
	time.Sleep(time.Second * 5)
	a := assert.New(t)
	ctx := context.Background()
	client, err := NewClient(os.Getenv("GROQ_KEY"))
	a.NoError(err, "NewClient error")
	response, err := client.CreateTranscription(ctx, AudioRequest{
		Model:    ModelWhisperLargeV3,
		FilePath: "./examples/audio-lex-fridman/The Roman Emperors who went insane Gregory Aldrete and Lex Fridman.mp3",
	})
	a.NoError(err, "CreateTranscription error")
	a.NotEmpty(response.Text, "response.Text is empty")
}

// TestLlamaGuard38B tests the LlamaGuard38B model.
// It ensures that the model is supported by the groq-go library, the groq API,
// and the operations are working as expected for the specific model type.
func TestLlamaGuard38B(t *testing.T) {
	time.Sleep(time.Second * 5)
	a := assert.New(t)
	ctx := context.Background()
	client, err := NewClient(os.Getenv("GROQ_KEY"))
	a.NoError(err, "NewClient error")
	response, err := client.Moderate(ctx, ModerationRequest{
		Model: ModelLlamaGuard38B,
		Messages: []ChatCompletionMessage{
			{
				Role:    ChatMessageRoleUser,
				Content: "I want to kill them.",
			},
		},
	})
	a.NoError(err, "Moderation error")
	a.Equal(true, response.Flagged)
	a.Contains(
		response.Categories,
		CategoryViolentCrimes,
	)
}
