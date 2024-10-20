# e2b-go-project

This is an example of using groq-go to create a simple golang project using the e2b and groq api powered by the groq-go library.

## Usage

Make sure you have a groq key set in the environment variable `GROQ_KEY`.
Also, make sure that you have a e2b api key set in the environment variable `E2B_API_KEY`.

```bash
export GROQ_KEY=your-groq-key
export E2B_API_KEY=your-e2b-api-key
go run .
```

### System Prompt

```txt
Given the tools given to you, create a golang project with the following files:

<files>
main.go
utils.go
<files>

The main function should call the `utils.run() error` function.

The project should, when run, print the following:

<output>
Hello, World!
<output>
```
