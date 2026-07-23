# AI Interview Trainer

An interactive technical interview simulator built with **Vue 3**, **Go** and local LLM inference through **Ollama**.

- Live demo: https://ai-interview-trainer.netlify.app  
- Repository: https://github.com/DiSTdio/ai-interview-trainer

---

## Overview

AI Interview Trainer explores two different approaches to building an AI-style interview interface:

- **MiniAI** — a lightweight browser-only simulation with predefined technical topics and deterministic evaluation logic.  
- **W Edition** — a local-first LLM implementation using a Go API and Ollama for real streamed model responses.

The goal of the project was not only to create a chat UI, but to experiment with the full request flow between a reactive frontend, an API service, prompt construction, session memory and streamed LLM output.

---

## Features

### Two execution modes

#### MiniAI

- Runs entirely in the browser and requires no backend or model installation.  
- Uses predefined question sets and lightweight client-side evaluation logic to reproduce the application flow without external AI infrastructure.  
- Makes the public demo immediately accessible and allows the interface and interview workflow to be tested independently from the LLM environment.

#### W Edition

- Runs with a local LLM through Ollama.  
- The Vue frontend sends session data to a Go API; the backend constructs a mode-specific prompt, forwards it to Ollama and streams generated tokens back to the browser.  
- Intended to run locally because model inference remains on the user's machine.

---

## Session modes

### 💬 Practice Chat

- Casual technical discussion mode focused on coding, frontend, backend and architecture.  
- In MiniAI, responses are selected from lightweight discussion prompts.  
- In W Edition, the user's input is forwarded to the local LLM as an open technical conversation.

### 🎯 Interview

- Simulates a technical interview.  
- The user selects a qualification or technology area and answers interview questions related to that topic.  
- Available topics include: **Frontend, Backend, QA, AWS, DevOps, React, Vue, JavaScript, TypeScript, Golang**.  
- The LLM edition uses conversation history to preserve interview context between requests.

### ⚡ Judge

- Stricter interview mode with scoring and feedback.  
- The backend instructs the model to evaluate the candidate's answer, explain weaknesses and continue with a harder question.  
- Scores are extracted from the generated response and reflected in the chat UI.

---

## Architecture

```text
┌─────────────────────┐
│      Vue 3 UI       │
│       + Vite        │
└──────────┬──────────┘
           │
 ┌─────────┴─────────┐
 │                   │
MiniAI           W Edition
 │                   │
Browser evaluator   HTTP POST
Predefined Qs          │
No backend             ▼
               ┌─────────────────┐
               │     Go API      │
               │    /ai endpoint │
               └────────┬────────┘
                        │
      ┌─────────────────┼─────────────────┐
      │                 │                 │
Prompt Builder    Session Memory     Rate Limiter
      │                 │                 │
      └─────────────────┼─────────────────┘
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
```

### Request flow

1. The user selects an edition and session mode.  
2. Interview and Judge modes allow the user to select a technical topic.  
3. The frontend sends `user_id`, `mode`, `judge`, `topic` and the current input to the Go API.  
4. The backend retrieves recent session history and builds a prompt for the selected mode.  
5. The prompt is sent to the Ollama `/api/generate` endpoint with streaming enabled.  
6. The Go API forwards generated chunks to the browser as a Server-Sent Events style stream.  
7. Vue incrementally appends incoming text to the active AI message.  
8. In Judge mode, the backend extracts the score from the completed model response.  
9. The latest conversation messages are retained in bounded in-memory session history.

---

## Frontend

Built with:

- Vue 3  
- Composition API and `<script setup>`  
- Vite  
- JavaScript  
- Scoped component CSS

Stage-based UI flow:

1. Edition selection  
2. W Edition setup guide (Full Edition only)  
3. Session configuration  
4. Chat / Interview UI

Main frontend components:

```text
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
```

`ChatContainer.vue` handles both MiniAI and W Edition interaction flows:

- Reactive message state  
- Topic selection  
- Streamed response rendering  
- Typing state  
- Abortable requests  
- Auto-scroll  
- Dynamic textarea height  
- Score-based visual feedback

---

## Backend

Backend is written in Go and split into focused packages:

```text
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
```

### API layer (`/ai` handler)

- Decodes incoming JSON and validates the user ID.  
- Applies per-user rate limiting.  
- Retrieves session history and selects the appropriate prompt strategy.  
- Starts LLM generation and streams output to the client.  
- Extracts Judge scores and updates conversation history.

### LLM client

- Communicates with Ollama via `/api/generate`.  
- Uses streaming mode; each response chunk is forwarded to the browser immediately instead of waiting for the full response.

### Prompt builder

- Separates prompt behavior from HTTP handling.  
- Builds different instructions for:
  - Normal technical chat  
  - Helpful technical interview  
  - Strict Judge interview  
- Interview prompts include the selected topic and recent conversation context.

### Session memory

- Conversation history stored in-memory and grouped by `user_id`.  
- Uses `sync.RWMutex` for concurrency safety.  
- Retains only the latest 10 messages per user to keep context bounded.  
- Intentionally avoids a database for this prototype.

### Rate limiting

- Lightweight per-user sliding-window limiter.  
- Current config: 5 requests within a 10-second window per user ID.

### Judge score extraction

- Judge mode asks the LLM to return `Score: X/10`.  
- After generation, the API extracts the numeric score via regex and sends it as a separate stream event.  
- The UI maps the score to visual feedback states.

---

## Tech stack

| Layer              | Technology                          |
|--------------------|-------------------------------------|
| Frontend           | Vue 3                               |
| Build tool         | Vite                                |
| Frontend language  | JavaScript                          |
| Backend            | Go                                  |
| LLM runtime        | Ollama                              |
| Default model      | Qwen2.5-Coder 3B                    |
| Streaming          | HTTP streaming / SSE-style events   |
| Session state      | In-memory Go store                  |
| Concurrency safety | `sync.RWMutex` / `sync.Mutex`       |
| Backend packaging  | Docker                              |
| Public demo        | Netlify                             |

---

## Running MiniAI

MiniAI requires only the frontend:

```bash
git clone https://github.com/DiSTdio/ai-interview-trainer.git
cd ai-interview-trainer

npm install
npm run dev
```

Open: http://localhost:5173 and select **MiniAI**.  
No Ollama installation or backend is required.

---

## Running W Edition

W Edition requires Ollama, the Go backend and the Vue frontend.

1. Install Ollama  
   - Download: https://ollama.com/download  
   - Verify: `ollama --version`

2. Download the model  
   - Default: `ollama pull qwen2.5-coder:3b`  
   - List: `ollama list`  
   - To change the model, update the name in `ai-backend/cmd/main.go`.

3. Start Ollama  
   - Usually starts automatically.  
   - If needed: `ollama serve`

4. Start the backend  

   ```bash
   cd ai-backend
   docker build -t ai-go .
   docker run -p 8080:8080 ai-go
   ```

   API listens on `http://localhost:8080`.  
   Dockerized backend reaches Ollama via `host.docker.internal:11434`.

5. Start the frontend  

   ```bash
   cd ..
   npm install
   npm run dev
   ```

   Open: http://localhost:5173 and select **W Edition**.

---

## Deployment note

The Netlify deployment primarily demonstrates the browser-only MiniAI edition.  
W Edition uses a local backend and local Ollama instance and is designed for local execution rather than public hosted inference.  
A production Full Edition would need a remotely accessible backend and LLM runtime.

This separation is intentional: MiniAI provides a zero-install public demo, while W Edition demonstrates the complete **frontend → Go API → local LLM** architecture.

---

## Design decisions

### Why two editions?

- LLM execution is fundamentally different from static frontend hosting.  
- MiniAI keeps the project instantly testable in a browser.  
- W Edition demonstrates real model integration without paid external AI APIs.

### Why Go?

- Small, efficient backend surface for HTTP, streaming and concurrent session state.  
- Packages separate API handling, model communication, memory, prompts and rate limiting instead of a single monolithic server file.

### Why local Ollama?

- Designed to experiment with local model inference.  
- Avoids coupling to proprietary hosted LLM APIs.  
- Model can be replaced by changing Ollama configuration.

### Why in-memory history?

- Only short-lived interview context is needed.  
- Avoids database persistence and operational overhead for a prototype.  
- Bounded history prevents unbounded prompt growth.

### Why streaming?

- Waiting for full LLM responses makes chat UIs feel unresponsive.  
- Backend forwards Ollama chunks as they arrive and the frontend renders incrementally.

---

## Current limitations

- W Edition requires a locally running Ollama instance.  
- Public deployment does not provide hosted LLM inference.  
- Session history is lost when the backend restarts.  
- User IDs are session identifiers, not authenticated accounts.  
- MiniAI evaluation is lightweight and heuristic-based.  
- Judge score extraction depends on model following the requested format.  
- CORS is open for development/demo.  
- Backend and model configuration are defined in code, not environment variables.

---

## Possible improvements

Future iterations could add:

- Environment-based backend/model configuration.  
- Persistent session storage.  
- Authenticated user sessions.  
- Structured JSON output for Judge results.  
- Stronger score parsing and validation.  
- Configurable interview length and difficulty.  
- Custom job-role prompts.  
- Resume or job-description context.  
- Automated frontend and backend tests.  
- Hosted inference support.  
- Docker Compose for frontend, backend and Ollama.

---

## What I explored

This project is a practical exploration of:

- Vue reactive UI state and component-driven frontend flow.  
- Go HTTP services and frontend/backend integration.  
- Local LLM inference, prompt design and streamed model responses.  
- Short-term conversation memory and concurrency-safe shared state.  
- Rate limiting and lightweight response parsing.  
- Dockerized backend execution.

The main engineering challenge was connecting two different execution paths — a zero-install browser simulation and a real local LLM workflow — behind the same interview interface.
