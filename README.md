AI Interview Trainer

An interactive technical interview simulator built with Vue 3, Go and local LLM inference through Ollama.

Live demo: ai-interview-trainer.netlify.app
Repository: github.com/DiSTdio/ai-interview-trainer

Overview

AI Interview Trainer explores two different approaches to building an AI-style interview interface:

MiniAI — a lightweight browser-only simulation with predefined technical topics and deterministic evaluation logic.
W Edition — a local-first LLM implementation using a Go API and Ollama for real streamed model responses.

The goal of the project was not only to create a chat UI, but to experiment with the full request flow between a reactive frontend, an API service, prompt construction, session memory and streamed LLM output.

Features
Two execution modes
MiniAI

Runs entirely in the browser and requires no backend or model installation.

It uses predefined question sets and lightweight client-side evaluation logic to reproduce the application flow without external AI infrastructure.

This makes the public demo immediately accessible and allows the interface and interview workflow to be tested independently from the LLM environment.

W Edition

Runs with a local LLM through Ollama.

The Vue frontend sends session data to a Go API. The backend constructs a mode-specific prompt, forwards it to Ollama and streams generated tokens back to the browser.

W Edition is intended to run locally because model inference remains on the user's machine.

Session modes
💬 Practice Chat

A casual technical discussion mode focused on coding, frontend, backend and architecture.

In MiniAI, responses are selected from lightweight discussion prompts.

In W Edition, the user's input is forwarded to the local LLM as an open technical conversation.

🎯 Interview

Simulates a technical interview.

The user selects a qualification or technology area and answers interview questions related to that topic.

Available topics currently include:

Frontend
Backend
QA
AWS
DevOps
React
Vue
JavaScript
TypeScript
Golang

The LLM edition uses conversation history to preserve interview context between requests.

⚡ Judge

A stricter interview mode with scoring and feedback.

The backend instructs the model to evaluate the candidate's answer, explain weaknesses and continue with a harder question.

Scores are extracted from the generated response and reflected in the chat UI.

Architecture
                    ┌─────────────────────┐
                    │      Vue 3 UI       │
                    │       + Vite        │
                    └──────────┬──────────┘
                               │
                 ┌─────────────┴─────────────┐
                 │                           │
            MiniAI                      W Edition
                 │                           │
      Browser-side evaluator             HTTP POST
      Predefined questions                  │
      No backend required                   ▼
                                    ┌─────────────────┐
                                    │     Go API      │
                                    │    /ai endpoint │
                                    └────────┬────────┘
                                             │
                         ┌───────────────────┼───────────────────┐
                         │                   │                   │
                    Prompt Builder      Session Memory      Rate Limiter
                         │                   │                   │
                         └───────────────────┼───────────────────┘
                                             │
                                             ▼
                                    ┌─────────────────┐
                                    │     Ollama      │
                                    │ Local LLM model │
                                    └────────┬────────┘
                                             │
                                      streamed chunks
                                             │
                                             ▼
                                         Vue UI
Request flow
The user selects an edition and session mode.
Interview and Judge modes allow the user to select a technical topic.
The frontend sends user_id, mode, judge, topic and the current input to the Go API.
The backend retrieves recent session history and builds a prompt for the selected mode.
The prompt is sent to the Ollama /api/generate endpoint with streaming enabled.
The Go API forwards generated chunks to the browser as a Server-Sent Events style stream.
Vue incrementally appends incoming text to the active AI message.
In Judge mode, the backend extracts the score from the completed model response.
The latest conversation messages are retained in bounded in-memory session history.
Frontend

The frontend is built with:

Vue 3
Composition API
<script setup>
Vite
JavaScript
Scoped component CSS

The application uses a simple stage-based UI flow:

Edition selection
        ↓
W Edition setup guide (Full Edition only)
        ↓
Session configuration
        ↓
Chat / Interview UI

Main frontend components:

src/
├── App.vue
└── components/
    ├── EditionSelect.vue
    ├── WEdition.vue
    ├── SessionSetupCard.vue
    ├── ChatContainer.vue
    └── mini/
        ├── evaluator.js
        └── questions.js

ChatContainer.vue handles both MiniAI and W Edition interaction flows.

The component includes:

reactive message state
topic selection
streamed response rendering
typing state
abortable requests
auto-scroll
dynamic textarea height
score-based visual feedback
Backend

The backend is written in Go and deliberately split into small packages with focused responsibilities.

ai-backend/
├── cmd/
│   └── main.go
├── internal/
│   ├── api/
│   │   └── handler.go
│   ├── llm/
│   │   └── ollama.go
│   ├── memory/
│   │   └── history.go
│   ├── model/
│   │   └── types.go
│   ├── prompt/
│   │   └── builder.go
│   └── rate/
│       └── limiter.go
├── Dockerfile
└── go.mod
API layer

The /ai handler:

decodes incoming JSON
validates the user ID
applies per-user rate limiting
retrieves session history
selects the appropriate prompt strategy
starts LLM generation
streams output to the client
extracts Judge scores
updates conversation history
LLM client

The Ollama client communicates with:

/api/generate

Generation uses streaming mode.

Each Ollama response chunk is immediately forwarded to the browser instead of waiting for the complete model response.

Prompt builder

Prompt behavior is separated from the HTTP handler.

The backend currently builds different instructions for:

normal technical chat
helpful technical interview
strict Judge interview

Interview prompts also include the selected topic and recent conversation context.

Session memory

Conversation history is stored in memory and grouped by user_id.

The history implementation uses sync.RWMutex to protect concurrent access.

Only the latest 10 messages are retained for each user, keeping prompt context bounded.

The current implementation intentionally does not use a database.

Rate limiting

The backend contains a lightweight per-user sliding-window rate limiter.

The current server configuration allows five requests within a ten-second window for each user ID.

Judge score extraction

Judge mode asks the LLM to return a score using a Score: X/10 style format.

After generation completes, the API uses a regular expression to extract the numeric score and sends it to the frontend as a separate stream event.

The UI maps the score to visual feedback states.

Tech stack
Layer	Technology
Frontend	Vue 3
Build tool	Vite
Frontend language	JavaScript
Backend	Go
LLM runtime	Ollama
Default model	Qwen2.5-Coder 3B
Streaming	HTTP streaming / SSE-style events
Session state	In-memory Go store
Concurrency safety	sync.RWMutex / sync.Mutex
Backend packaging	Docker
Public frontend demo	Netlify
Running MiniAI

MiniAI requires only the frontend.

git clone https://github.com/DiSTdio/ai-interview-trainer.git
cd ai-interview-trainer

npm install
npm run dev

Open:

http://localhost:5173

Select MiniAI.

No Ollama installation or backend is required.

Running W Edition

W Edition requires Ollama, the Go backend and the Vue frontend.

1. Install Ollama

Download Ollama from:

https://ollama.com/download

Verify the installation:

ollama --version
2. Download the model

The default backend configuration uses:

ollama pull qwen2.5-coder:3b

List installed models with:

ollama list

To use another model, update the model name in:

ai-backend/cmd/main.go
3. Start Ollama

Ollama normally starts automatically.

If necessary:

ollama serve
4. Start the backend

Open a terminal in:

ai-backend

Build the Docker image:

docker build -t ai-go .

Run the backend:

docker run -p 8080:8080 ai-go

The API listens on:

http://localhost:8080

The Dockerized backend is configured to reach Ollama through host.docker.internal:11434.

5. Start the frontend

From the repository root:

npm install
npm run dev

Open:

http://localhost:5173

Select W Edition and follow the session setup flow.

Deployment note

The public Netlify deployment is primarily a demonstration of the browser-only MiniAI edition.

W Edition currently uses a local backend endpoint and a locally running Ollama instance. It is therefore designed for local execution rather than public hosted inference.

A production-hosted Full Edition would require a remotely accessible backend and an LLM runtime available to that backend.

This separation is intentional: MiniAI provides a zero-install public demonstration, while W Edition demonstrates the complete frontend → Go API → local LLM architecture.

Design decisions
Why two editions?

Running an LLM is fundamentally different from deploying a static frontend.

MiniAI keeps the project immediately testable in a browser, while W Edition demonstrates real model integration without requiring a paid external AI API.

Why Go?

Go provides a small backend surface for HTTP handling, streaming and concurrent session state.

The backend packages separate API handling, model communication, memory, prompts and rate limiting rather than placing the complete flow in a single server file.

Why local Ollama?

The project was designed to experiment with local model inference and avoid coupling the application to a proprietary hosted LLM API.

The model can be replaced by changing the Ollama model configuration.

Why in-memory history?

The application only needs short-lived interview context.

A database would add persistence and operational complexity that the current prototype does not require.

The bounded history also prevents conversation context from growing indefinitely.

Why streaming?

Waiting for a complete LLM response makes chat interfaces feel unresponsive.

The backend forwards Ollama chunks as they arrive, allowing the frontend to render the answer incrementally.

Current limitations
W Edition requires a locally running Ollama instance.
The public deployment does not provide hosted LLM inference.
Session history is lost when the backend restarts.
User IDs are session identifiers, not authenticated accounts.
MiniAI evaluation is intentionally lightweight and heuristic-based.
Judge score extraction depends on the model following the requested score format.
CORS is currently open for development and demonstration purposes.
Backend and model configuration are currently defined in code rather than environment variables.
Possible improvements

Future iterations could include:

environment-based backend and model configuration
persistent session storage
authenticated user sessions
structured JSON output for Judge results
stronger score parsing and validation
configurable interview length and difficulty
custom job-role prompts
resume or job-description context
automated frontend and backend tests
hosted inference support
Docker Compose for frontend, backend and Ollama orchestration
What I explored with this project

This project was built as a practical exploration of:

Vue reactive UI state
component-driven frontend flow
Go HTTP services
frontend/backend integration
local LLM inference
prompt design
streamed model responses
short-term conversation memory
concurrency-safe shared state
rate limiting
lightweight response parsing
Dockerized backend execution

The main engineering challenge was connecting two very different execution paths — a zero-install browser simulation and a real local LLM workflow — behind the same interview interface.
