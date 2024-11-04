package jigsawstack_test

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/conneroisu/groq-go/extensions/jigsawstack"
	"github.com/conneroisu/groq-go/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestAudioTTS(t *testing.T) {
	if !test.IsIntegrationTest() {
		t.Skip()
	}
	a := assert.New(t)
	apiKey, err := test.GetAPIKey("JIGSAWSTACK_API_KEY")
	a.NoError(err)
	ctx := context.Background()
	client, err := jigsawstack.NewJigsawStack(
		apiKey,
		jigsawstack.WithLogger(test.DefaultLogger),
	)
	a.NoError(err)
	response, err := client.AudioTTS(ctx,
		"Hello, world! Welcome to Groq!",
		jigsawstack.WithAccent("zh-TW-female-19"),
	)
	a.NoError(err)
	// write the io.reader to a file
	f, err := os.Create("tts.mp3")
	a.NoError(err)
	defer f.Close()
	_, err = io.Copy(f, strings.NewReader(response))
	a.NoError(err)
}
