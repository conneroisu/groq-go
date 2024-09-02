# Audio

## Transcriptions

Response Types:

- <a href="https://pkg.go.dev/github.com/conneroisu/groq-go">groq</a>.<a href="https://pkg.go.dev/github.com/conneroisu/groq-go#AudioTranscriptionNewResponse">AudioTranscriptionNewResponse</a>

Methods:

- <code title="post /openai/v1/audio/transcriptions">client.Audio.Transcriptions.<a href="https://pkg.go.dev/github.com/conneroisu/groq-go#AudioTranscriptionService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, body <a href="https://pkg.go.dev/github.com/conneroisu/groq-go">groq</a>.<a href="https://pkg.go.dev/github.com/conneroisu/groq-go#AudioTranscriptionNewParams">AudioTranscriptionNewParams</a>) (<a href="https://pkg.go.dev/github.com/conneroisu/groq-go">groq</a>.<a href="https://pkg.go.dev/github.com/conneroisu/groq-go#AudioTranscriptionNewResponse">AudioTranscriptionNewResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

## Translations

Response Types:

- <a href="https://pkg.go.dev/github.com/conneroisu/groq-go">groq</a>.<a href="https://pkg.go.dev/github.com/conneroisu/groq-go#AudioTranslationNewResponse">AudioTranslationNewResponse</a>

Methods:

- <code title="post /openai/v1/audio/translations">client.Audio.Translations.<a href="https://pkg.go.dev/github.com/conneroisu/groq-go#AudioTranslationService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, body <a href="https://pkg.go.dev/github.com/conneroisu/groq-go">groq</a>.<a href="https://pkg.go.dev/github.com/conneroisu/groq-go#AudioTranslationNewParams">AudioTranslationNewParams</a>) (<a href="https://pkg.go.dev/github.com/conneroisu/groq-go">groq</a>.<a href="https://pkg.go.dev/github.com/conneroisu/groq-go#AudioTranslationNewResponse">AudioTranslationNewResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

# Chat

## Completions

Response Types:

- <a href="https://pkg.go.dev/github.com/conneroisu/groq-go">groq</a>.<a href="https://pkg.go.dev/github.com/conneroisu/groq-go#ChatCompletionNewResponse">ChatCompletionNewResponse</a>

Methods:

- <code title="post /openai/v1/chat/completions">client.Chat.Completions.<a href="https://pkg.go.dev/github.com/conneroisu/groq-go#ChatCompletionService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, body <a href="https://pkg.go.dev/github.com/conneroisu/groq-go">groq</a>.<a href="https://pkg.go.dev/github.com/conneroisu/groq-go#ChatCompletionNewParams">ChatCompletionNewParams</a>) (<a href="https://pkg.go.dev/github.com/conneroisu/groq-go">groq</a>.<a href="https://pkg.go.dev/github.com/conneroisu/groq-go#ChatCompletionNewResponse">ChatCompletionNewResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

# Embeddings

Response Types:

- <a href="https://pkg.go.dev/github.com/conneroisu/groq-go">groq</a>.<a href="https://pkg.go.dev/github.com/conneroisu/groq-go#EmbeddingNewResponse">EmbeddingNewResponse</a>

Methods:

- <code title="post /openai/v1/embeddings">client.Embeddings.<a href="https://pkg.go.dev/github.com/conneroisu/groq-go#EmbeddingService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, body <a href="https://pkg.go.dev/github.com/conneroisu/groq-go">groq</a>.<a href="https://pkg.go.dev/github.com/conneroisu/groq-go#EmbeddingNewParams">EmbeddingNewParams</a>) (<a href="https://pkg.go.dev/github.com/conneroisu/groq-go">groq</a>.<a href="https://pkg.go.dev/github.com/conneroisu/groq-go#EmbeddingNewResponse">EmbeddingNewResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

# Models

Response Types:

- <a href="https://pkg.go.dev/github.com/conneroisu/groq-go">groq</a>.<a href="https://pkg.go.dev/github.com/conneroisu/groq-go#Model">Model</a>
- <a href="https://pkg.go.dev/github.com/conneroisu/groq-go">groq</a>.<a href="https://pkg.go.dev/github.com/conneroisu/groq-go#ModelListResponse">ModelListResponse</a>
- <a href="https://pkg.go.dev/github.com/conneroisu/groq-go">groq</a>.<a href="https://pkg.go.dev/github.com/conneroisu/groq-go#ModelDeleteResponse">ModelDeleteResponse</a>

Methods:

- <code title="get /openai/v1/models/{model}">client.Models.<a href="https://pkg.go.dev/github.com/conneroisu/groq-go#ModelService.Get">Get</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, model <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/github.com/conneroisu/groq-go">groq</a>.<a href="https://pkg.go.dev/github.com/conneroisu/groq-go#Model">Model</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /openai/v1/models">client.Models.<a href="https://pkg.go.dev/github.com/conneroisu/groq-go#ModelService.List">List</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>) (<a href="https://pkg.go.dev/github.com/conneroisu/groq-go">groq</a>.<a href="https://pkg.go.dev/github.com/conneroisu/groq-go#ModelListResponse">ModelListResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="delete /openai/v1/models/{model}">client.Models.<a href="https://pkg.go.dev/github.com/conneroisu/groq-go#ModelService.Delete">Delete</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, model <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/github.com/conneroisu/groq-go">groq</a>.<a href="https://pkg.go.dev/github.com/conneroisu/groq-go#ModelDeleteResponse">ModelDeleteResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
