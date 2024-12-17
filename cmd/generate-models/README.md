# generate-models

This is a script to generate the models for the groq-go library.

## Usage

Make sure you have a groq key set in the environment variable `GROQ_KEY`.

Also, make sure to run this from the root of the groq-go project.
```bash
export GROQ_KEY=your-groq-key
go run ./cmd/generate-models
```

Or you can run it automatically with `go generate` at the root of the project:
```bash
go generate
```
