package jigsawstack

import (
	"context"
	"net/http"

	"github.com/conneroisu/groq-go/pkg/builders"
)

const (
	sentimentSuffix Endpoint = "/v1/ai/sentiment"
	// EmotionAnger is the anger emotion.
	EmotionAnger Emotion = "anger"
	// EmotionFear is the fear emotion.
	EmotionFear Emotion = "fear"
	// EmotionSadness is the sadness emotion.
	EmotionSadness Emotion = "sadness"
	// EmotionHappiness is the happiness emotion.
	EmotionHappiness Emotion = "happiness"
	// EmotionAnxiety is the anxiety emotion.
	EmotionAnxiety Emotion = "anxiety"
	// EmotionDisgust is the disgust emotion.
	EmotionDisgust Emotion = "disgust"
	// EmotionEmbarrassment is the embarrassment emotion.
	EmotionEmbarrassment Emotion = "embarrassment"
	// EmotionLove is the love emotion.
	EmotionLove Emotion = "love"
	// EmotionSurprise is the surprise emotion.
	EmotionSurprise Emotion = "surprise"
	// EmotionShame is the shame emotion.
	EmotionShame Emotion = "shame"
	// EmotionEnvy is the envy emotion.
	EmotionEnvy Emotion = "envy"
	// EmotionSatisfaction is the satisfaction emotion.
	EmotionSatisfaction Emotion = "satisfaction"
	// EmotionSelfConfidence is the self-confidence emotion.
	EmotionSelfConfidence Emotion = "self-confidence"
	// EmotionAnnoyance is the annoyance emotion.
	EmotionAnnoyance Emotion = "annoyance"
	// EmotionBoredom is the boredom emotion.
	EmotionBoredom Emotion = "boredom"
	// EmotionHatred is the hatred emotion.
	EmotionHatred Emotion = "hatred"
	// EmotionCompassion is the compassion emotion.
	EmotionCompassion Emotion = "compassion"
	// EmotionGuilt is the guilt emotion.
	EmotionGuilt Emotion = "guilt"
	// EmotionLoneliness is the loneliness emotion.
	EmotionLoneliness Emotion = "loneliness"
	// EmotionDepression is the depression emotion.
	EmotionDepression Emotion = "depression"
	// EmotionPride is the pride emotion.
	EmotionPride Emotion = "pride"
	// EmotionNeutral is the neutral emotion.
	EmotionNeutral Emotion = "neutral"
)

type (
	// Emotion is an emotion.
	Emotion string
	// SentimentRequest represents a request structure for sentiment API.
	SentimentRequest struct {
		Text string `json:"text"`
	}
	// SentimentResponse represents a response structure for sentiment API.
	SentimentResponse struct {
		Success   bool `json:"success"`
		Sentiment struct {
			Emotion   Emotion `json:"emotion"`
			Sentiment string  `json:"sentiment"`
			Score     float64 `json:"score"`
			Sentences []struct {
				Text      string  `json:"text"`
				Emotion   Emotion `json:"emotion"`
				Sentiment string  `json:"sentiment"`
				Score     float64 `json:"score"`
			} `json:"sentences"`
		} `json:"sentiment"`
	}
)

// Sentiment performs a sentiment api call over a string.
func (c *JigsawStack) Sentiment(
	ctx context.Context,
	request SentimentRequest,
) (SentimentResponse, error) {
	req, err := builders.NewRequest(
		ctx,
		c.header,
		http.MethodPost,
		c.baseURL+string(sentimentSuffix),
		builders.WithBody(request),
	)
	if err != nil {
		return SentimentResponse{}, err
	}
	var respH SentimentResponse
	err = c.sendRequest(req, &respH)
	if err != nil {
		return SentimentResponse{}, err
	}
	return respH, nil
}
