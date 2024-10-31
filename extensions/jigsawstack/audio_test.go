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
	if !test.IsUnitTest() {
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
	t.Log("ereatedc client")
	a.NoError(err)
	response, err := client.AudioTTS(ctx, jigsawstack.TTSRequest{
		Text: "Hello, world!",
	})
	a.NoError(err)
	t.Logf("response: %s", response)
	// write the io.reader to a file
	f, err := os.Create("tts.mp3")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	_, err = io.Copy(f, strings.NewReader(response))
	if err != nil {
		t.Fatal(err)
	}
}
