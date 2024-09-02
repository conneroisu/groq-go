# Changelog

## 0.1.0-alpha.1 (2024-09-02)

Full Changelog: [v0.0.1-alpha.0...v0.1.0-alpha.1](https://github.com/conneroisu/groq-go/compare/v0.0.1-alpha.0...v0.1.0-alpha.1)

### Features

* **api:** automatically fetch and manage available models ([0f6bfe5](https://github.com/conneroisu/groq-go/commit/0f6bfe5c4a94bab7a8c8de3f1b7a20d3be9de334))
* **chat:** add support for response formats ([0f6bfe5](https://github.com/conneroisu/groq-go/commit/0f6bfe5c4a94bab7a8c8de3f1b7a20d3be9de334))
* **chat:** enrich model error details with available models ([7bf6ad0](https://github.com/conneroisu/groq-go/commit/7bf6ad0ad4f1c5b180fd643301ce3f7d42cdf87a))
* **client:** add error response for unavailable model in Chat request ([3968c8e](https://github.com/conneroisu/groq-go/commit/3968c8e1fdc0c285a99123e958f2d580515c317d))
* **client:** add new fields to ModelResponse struct ([0c9e81a](https://github.com/conneroisu/groq-go/commit/0c9e81a07bbe9f521f2f7deca13b27e56ee6b680))
* **client:** include created and active fields in model response ([0c9e81a](https://github.com/conneroisu/groq-go/commit/0c9e81a07bbe9f521f2f7deca13b27e56ee6b680))


### Bug Fixes

* **chat:** correct JSON tag syntax in ChatRequest struct ([0f6bfe5](https://github.com/conneroisu/groq-go/commit/0f6bfe5c4a94bab7a8c8de3f1b7a20d3be9de334))
* **chat:** maintain single source of truth for chat URL ([c6369c3](https://github.com/conneroisu/groq-go/commit/c6369c34791b374f2887beb4c3eb0ffc8b235d82))
* **chat:** remove duplicate URL declaration in Chat function ([c6369c3](https://github.com/conneroisu/groq-go/commit/c6369c34791b374f2887beb4c3eb0ffc8b235d82))
* **client:** adjust GetModels method to parse ModelResponse ([0c9e81a](https://github.com/conneroisu/groq-go/commit/0c9e81a07bbe9f521f2f7deca13b27e56ee6b680))
* **client:** correct JSON annotation in ModelResponse struct ([0c9e81a](https://github.com/conneroisu/groq-go/commit/0c9e81a07bbe9f521f2f7deca13b27e56ee6b680))
* **client:** enhance error message to include available models list ([7bf6ad0](https://github.com/conneroisu/groq-go/commit/7bf6ad0ad4f1c5b180fd643301ce3f7d42cdf87a))
* **client:** rename Models struct to ModelResponse ([0c9e81a](https://github.com/conneroisu/groq-go/commit/0c9e81a07bbe9f521f2f7deca13b27e56ee6b680))
* **error:** handle unavailable models in Chat method ([3968c8e](https://github.com/conneroisu/groq-go/commit/3968c8e1fdc0c285a99123e958f2d580515c317d))
* **error:** remove redundant error handling in client ([0f6bfe5](https://github.com/conneroisu/groq-go/commit/0f6bfe5c4a94bab7a8c8de3f1b7a20d3be9de334))
* **mod:** update module path to reflect new package name ([0f6bfe5](https://github.com/conneroisu/groq-go/commit/0f6bfe5c4a94bab7a8c8de3f1b7a20d3be9de334))


### Chores

* **client:** clean up whitespace and formatting ([0c9e81a](https://github.com/conneroisu/groq-go/commit/0c9e81a07bbe9f521f2f7deca13b27e56ee6b680))
* go live ([#11](https://github.com/conneroisu/groq-go/issues/11)) ([bf69524](https://github.com/conneroisu/groq-go/commit/bf69524cc2cd7f35ddf56bb7fbdd967bb1cb0544))
* remove unused internal package and refactor code ([0f6bfe5](https://github.com/conneroisu/groq-go/commit/0f6bfe5c4a94bab7a8c8de3f1b7a20d3be9de334))


### Documentation

* **client:** add comments and improve readability ([0f6bfe5](https://github.com/conneroisu/groq-go/commit/0f6bfe5c4a94bab7a8c8de3f1b7a20d3be9de334))


### Styles

* rename package from gogroq to groq ([0f6bfe5](https://github.com/conneroisu/groq-go/commit/0f6bfe5c4a94bab7a8c8de3f1b7a20d3be9de334))


### Refactors

* **client:** remove unused Models struct ([0c9e81a](https://github.com/conneroisu/groq-go/commit/0c9e81a07bbe9f521f2f7deca13b27e56ee6b680))
* **client:** simplify Chat function implementation ([0f6bfe5](https://github.com/conneroisu/groq-go/commit/0f6bfe5c4a94bab7a8c8de3f1b7a20d3be9de334))
* **client:** streamline Client struct for better readability ([0c9e81a](https://github.com/conneroisu/groq-go/commit/0c9e81a07bbe9f521f2f7deca13b27e56ee6b680))
* **client:** update contains method to match ModelResponse format ([0c9e81a](https://github.com/conneroisu/groq-go/commit/0c9e81a07bbe9f521f2f7deca13b27e56ee6b680))
* **constants:** move chat URL to a constant ([c6369c3](https://github.com/conneroisu/groq-go/commit/c6369c34791b374f2887beb4c3eb0ffc8b235d82))
* delete unused files and simplify the codebase ([0f6bfe5](https://github.com/conneroisu/groq-go/commit/0f6bfe5c4a94bab7a8c8de3f1b7a20d3be9de334))
