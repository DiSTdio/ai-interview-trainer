<template>
    <div class="login-card">
        <h1>AI Interview Trainer</h1>
        <p class="subtitle">Practice technical communication with AI</p>

        <input v-model="name" maxlength="20" laceholder="Your name" :class="{ 'input-error': error }"
            @input="error = ''" />
        <div v-if="error" class="error-message">
            {{ error }}
        </div>

        <div class="modes">
            <div v-for="m in modes" :key="m.value" class="mode-card" :class="{ active: selected === m.value }"
                @click="selected = m.value">
                <div class="icon">{{ m.icon }}</div>
                <h3>{{ m.label }}</h3>
                <p>{{ m.desc }}</p>
            </div>
        </div>

        <transition name="fade">
            <div v-if="warning" class="warning-box">
                ⚠️ {{ warning }}
            </div>
        </transition>

        <button :disabled="!name.trim()" @click="start">
            Start Session
        </button>
    </div>
</template>

<script setup>
import { ref, watch } from "vue"

const emit = defineEmits(["login"])

const name = ref("")
const selected = ref("chat")
const error = ref("")
const warning = ref("")

const modes = [
    {
        value: "chat",
        icon: "💬",
        label: "Practice Chat",
        desc: "casual coding discussion"
    },
    {
        value: "interview",
        icon: "🎯",
        label: "Interview",
        desc: "real interview simulation"
    },
    {
        value: "judge",
        icon: "⚡",
        label: "Judge",
        desc: "strict scoring + feedback"
    }
]

const start = () => {
    error.value = ""

    if (!name.value.trim()) {
        error.value = "Please enter your name"
        return
    }

    emit("login", {
        name: name.value.trim(),
        mode: selected.value
    })
}

watch(selected, () => {
    warning.value =
        selected.value === "judge"
            ? "Judge mode is experimental. Scores may be inconsistent."
            : ""
})
</script>

<style scoped>
.login-card {
    max-width: 600px;
    margin: 17px auto;
    padding: 19px;
    background: #111827;
    border-radius: 20px;
    color: white;
    text-align: center;
    border: 1px solid #1f2937;
}

.subtitle {
    color: #94a3b8;
    margin-bottom: 24px;
}

input {
    width: 100%;
    padding: 14px;
    background: #1f2937;
    border: 1px solid #374151;
    border-radius: 12px;
    color: white;
    margin-bottom: 8px;
    box-sizing: border-box;
}

.input-error {
    border-color: #ff5a5a;
}

.error-message {
    color: #ff8c8c;
    font-size: 0.85rem;
    text-align: left;
    margin-bottom: 12px;
}

.modes {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 12px;
    margin: 20px 0;
}

.mode-card {
    padding: 16px;
    border-radius: 12px;
    background: #1e293b;
    cursor: pointer;
    border: 2px solid transparent;
    transition: all 0.2s;
}

.mode-card.active {
    border-color: #4fc3f7;
    background: rgba(79, 195, 247, 0.1);
}

.warning-box {
    background: rgba(255, 210, 100, 0.1);
    border: 1px solid #ffd264;
    color: #ffd264;
    padding: 12px;
    border-radius: 10px;
    margin-bottom: 20px;
    font-size: 0.85rem;
}

button {
    width: 100%;
    padding: 14px;
    background: #4fc3f7;
    color: #000;
    font-weight: bold;
    border: none;
    border-radius: 12px;
    cursor: pointer;
}

button:disabled {
    background: #334155;
    color: #64748b;
    cursor: not-allowed;
}

.fade-enter-active,
.fade-leave-active {
    transition: opacity 0.3s;
}

.fade-enter-from,
.fade-leave-to {
    opacity: 0;
}
</style>
