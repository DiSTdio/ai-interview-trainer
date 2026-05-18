<template>
    <div class="wrapper">

        <h1>AI Interview Trainer</h1>
        <p class="subtitle">
            Choose your edition
        </p>

        <div class="cards">

            <div class="card" @click="$emit('select', 'mini')">
                <h2>MiniAI</h2>
                <p>Instant mock interview</p>
                <p>No installation required</p>
                <button>Launch</button>
            </div>

            <div class="card matrix" @click="$emit('select', 'w')">
                <h2>W Edition</h2>
                <p>Run local LLM with Ollama</p>
                <p>Advanced real inference</p>
                <button>Install</button>
            </div>

        </div>

    </div>
</template>

<script setup>
defineEmits(["select"])
</script>

<style scoped>
.wrapper {
    max-width: 900px;
    margin: 60px auto;
    text-align: center;
    color: white;
}

.subtitle {
    color: #94a3b8;
    margin-bottom: 32px;
}

.cards {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 24px;
}

.card {
    background: #111827;
    border: 1px solid #1f2937;
    border-radius: 18px;
    padding: 28px;
    cursor: pointer;
    transition: .25s;
}

.card:hover {
    transform: translateY(-4px);
    border-color: #4fc3f7;
}

.matrix:hover {
    border-color: #00ff88;
    box-shadow: 0 0 20px rgba(0, 255, 136, .15);
}

button {
    margin-top: 18px;
    padding: 12px 18px;
    border: none;
    border-radius: 12px;
    background: #334155;
    color: white;
    cursor: pointer;
}
</style>