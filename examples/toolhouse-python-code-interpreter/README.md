# python-code-interpreter

This is an example of using groq-go to create a chat completion using the llama-3.1-70B-8192-tool-use-preview model.

It interacts with the [Toolhouse](https://app.toolhouse.ai/) API to run the code interpreter.

## Usage

Make sure you have a groq key set in the environment variable `GROQ_KEY`.
Also, make sure you have a toolhouse api key set in the environment variable `TOOLHOUSE_API_KEY`.
```bash
export GROQ_KEY=your-groq-key
export TOOLHOUSE_API_KEY=your-toolhouse-api-key
go run .
```
