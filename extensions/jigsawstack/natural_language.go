package jigsawstack

import (
	"context"
	"net/http"

	"github.com/conneroisu/groq-go/pkg/builders"
)

const (
	summaryEndpoint   Endpoint = "/v1/ai/summarize"
	sentimentSuffix   Endpoint = "/v1/ai/sentiment"
	translateEndpoint Endpoint = "/v1/ai/translate"

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

	// Language is a language.
	Language string
	// TranslateRequest represents a request structure for translate API.
	TranslateRequest struct {
		CurrentLanguage Language `json:"current_language"`
		TargetLanguage  Language `json:"target_language"`
		Text            string   `json:"text"`
	}
	// TranslateResponse represents a response structure for translate API.
	TranslateResponse struct {
		Success        bool   `json:"success"`
		TranslatedText string `json:"translated_text"`
	}
	// Emotion is an emotion.
	Emotion string
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
func (j *JigsawStack) Sentiment(
	ctx context.Context,
	text string,
) (SentimentResponse, error) {
	var request = struct {
		Text string `json:"text"`
	}{Text: text}
	req, err := builders.NewRequest(
		ctx,
		j.header,
		http.MethodPost,
		j.baseURL+string(sentimentSuffix),
		builders.WithBody(request),
	)
	if err != nil {
		return SentimentResponse{}, err
	}
	var respH SentimentResponse
	err = j.sendRequest(req, &respH)
	if err != nil {
		return SentimentResponse{}, err
	}
	return respH, nil
}

type (
	// SummaryRequest represents a request structure for summary API.
	SummaryRequest struct {
		Text string `json:"text"`
	}
	// SummaryResponse represents a response structure for summary API.
	SummaryResponse struct {
		Success bool   `json:"success"`
		Summary string `json:"summary"`
	}
)

// Summarize summarizes the give text.
//
// Max text character is 5000.
func (j *JigsawStack) Summarize(
	ctx context.Context,
	request SummaryRequest,
) (response SummaryResponse, err error) {
	req, err := builders.NewRequest(
		ctx,
		j.header,
		http.MethodPost,
		j.baseURL+string(summaryEndpoint),
		builders.WithBody(request),
	)
	if err != nil {
		return
	}
	var resp SummaryResponse
	err = j.sendRequest(req, &resp)
	if err != nil {
		return
	}
	return resp, nil
}

// Translate translates the text from the current language to the target language.
//
// Max text character is 5000.
func (j *JigsawStack) Translate(
	ctx context.Context,
	request TranslateRequest,
) (response TranslateResponse, err error) {
	req, err := builders.NewRequest(
		ctx,
		j.header,
		http.MethodPost,
		j.baseURL+string(translateEndpoint),
		builders.WithBody(request),
	)
	if err != nil {
		return
	}
	var resp TranslateResponse
	err = j.sendRequest(req, &resp)
	if err != nil {
		return
	}
	return resp, nil
}
