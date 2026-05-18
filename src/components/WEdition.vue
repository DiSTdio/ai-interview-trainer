<template>
    <div class="matrix">

        <h1>W Edition</h1>

        <pre>
Full Edition Setup

Full Edition runs with a local LLM through Ollama.

1. Clone the repository
git clone *your-repository-url*
cd *project-folder*
2. Install Ollama

Download and install Ollama:

https://ollama.com/download

Windows users may need WSL2 enabled.

After installation, verify:

ollama --version
3. Download a model

Recommended:

ollama pull qwen2.5-coder:3b

To list installed models:

ollama list

If you install another model, open:

ai-backend/cmd/main.go

Replace:

"qwen2.5-coder:3b"

with your installed model name.

Example:

"deepseek-coder:6.7b"
4. Start Ollama server

Usually starts automatically.

If needed:

ollama serve

Default endpoint:

http://localhost:11434
5. Start backend

Open terminal in:

ai-backend

Run:

docker build -t ai-go .
docker run -p 8080:8080 ai-go

Expected output:

🚀 :8080
6. Start frontend

Open another terminal in:

ai-ui

Install dependencies:

npm install

Run app:

npm run dev

Open browser:

http://localhost:5173
7. Launch Full Edition

Choose:

Full Edition

You now have:

Practice Chat
Interview Simulation
Judge Mode
Real streamed LLM responses
Persistent session memory
</pre>

        <button @click="$emit('proceed')">
            Proceed
        </button>

    </div>
</template>

<style scoped>
.matrix {
    max-width: 900px;
    margin: 40px auto;
    color: #00ff88;
    background: black;
    padding: 32px;
    border-radius: 18px;
    font-family: monospace;
    box-shadow: 0 0 40px rgba(0, 255, 136, .12);

    max-height: 85vh;
    overflow-y: auto;
}

pre {
    white-space: pre-wrap;
}

button {
    margin-top: 24px;
    padding: 14px 24px;
    width: 100%;

    background: rgba(0, 255, 136, .08);
    color: #00ff88;

    border: 1px solid rgba(0, 255, 136, .35);
    border-radius: 12px;

    font-family: monospace;
    font-size: 16px;
    letter-spacing: 2px;
    text-transform: uppercase;

    cursor: pointer;

    box-shadow:
        0 0 10px rgba(0, 255, 136, .08),
        inset 0 0 10px rgba(0, 255, 136, .04);

    transition: all .22s ease;
}

button:hover {
    background: rgba(0, 255, 136, .14);

    box-shadow:
        0 0 20px rgba(0, 255, 136, .18),
        0 0 40px rgba(0, 255, 136, .06),
        inset 0 0 16px rgba(0, 255, 136, .08);

    transform: translateY(-2px);
}

button:active {
    transform: scale(.985);
}

.matrix::-webkit-scrollbar {
    width: 8px;
}

.matrix::-webkit-scrollbar-thumb {
    background: rgba(0, 255, 136, .35);
    border-radius: 999px;
}

.matrix::-webkit-scrollbar-track {
    background: rgba(255, 255, 255, .03);
}
</style>