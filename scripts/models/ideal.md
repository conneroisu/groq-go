
Current Response: 
```json
{
    "object": "list",
    "data": [
        {
            "id": "llama3-8b-8192",
            "object": "model",
            "created": 1693721698,
            "owned_by": "Meta",
            "active": true,
            "context_window": 8192,
            "public_apps": null
        },
        {
            "id": "llama3-groq-8b-8192-tool-use-preview",
            "object": "model",
            "created": 1693721698,
            "owned_by": "Groq",
            "active": true,
            "context_window": 8192,
            "public_apps": null
        },
        {
            "id": "whisper-large-v3",
            "object": "model",
            "created": 1693721698,
            "owned_by": "OpenAI",
            "active": true,
            "context_window": 448,
            "public_apps": null
        },
        {
            "id": "gemma-7b-it",
            "object": "model",
            "created": 1693721698,
            "owned_by": "Google",
            "active": true,
            "context_window": 8192,
            "public_apps": null
        },
        {
            "id": "llama3-70b-8192",
            "object": "model",
            "created": 1693721698,
            "owned_by": "Meta",
            "active": true,
            "context_window": 8192,
            "public_apps": null
        },
        {
            "id": "mixtral-8x7b-32768",
            "object": "model",
            "created": 1693721698,
            "owned_by": "Mistral AI",
            "active": true,
            "context_window": 32768,
            "public_apps": null
        },
        {
            "id": "llama3-groq-70b-8192-tool-use-preview",
            "object": "model",
            "created": 1693721698,
            "owned_by": "Groq",
            "active": true,
            "context_window": 8192,
            "public_apps": null
        },
        {
            "id": "llama-3.1-8b-instant",
            "object": "model",
            "created": 1693721698,
            "owned_by": "Meta",
            "active": true,
            "context_window": 131072,
            "public_apps": null
        },
        {
            "id": "gemma2-9b-it",
            "object": "model",
            "created": 1693721698,
            "owned_by": "Google",
            "active": true,
            "context_window": 8192,
            "public_apps": null
        },
        {
            "id": "llama-3.1-70b-versatile",
            "object": "model",
            "created": 1693721698,
            "owned_by": "Meta",
            "active": true,
            "context_window": 131072,
            "public_apps": null
        },
        {
            "id": "distil-whisper-large-v3-en",
            "object": "model",
            "created": 1693721698,
            "owned_by": "Hugging Face",
            "active": true,
            "context_window": 448,
            "public_apps": null
        },
        {
            "id": "llama-guard-3-8b",
            "object": "model",
            "created": 1693721698,
            "owned_by": "Meta",
            "active": true,
            "context_window": 8192,
            "public_apps": null
        }
    ]
}
```

But the converting of each of the models to a golang type is not done correctly.
Each variable should a be a valid Golang type and should be commented correctly to work with godoc and existing linters.