<template>
    <div class="container">

        <div class="controls">
            <div class="user">{{ userId }}</div>

            <div class="mode-badge">
                {{ judge ? "judge" : mode }}
            </div>

            <label class="checkbox">
                <input type="checkbox" v-model="autoScroll" />
                Auto-scroll
            </label>

            <button @click="$emit('logout')">
                Exit
            </button>

        </div>
        <div v-if="
            props.edition === 'mini' &&
            props.mode !== 'chat' &&
            !topicSelected
        " class="topic-picker">

            <div class="ai-msg">
                Choose role / qualification
            </div>

            <div class="topic-cloud">
                <button v-for="t in topics" :key="t" @click="chooseTopic(t)" class="topic-pill">
                    {{ t }}
                </button>
            </div>

        </div>

        <div class="chat" :class="scoreClass" ref="chatRef">
            <div v-for="(m, i) in messages" :key="i" :class="['msg', m.role]">
                <b>{{ m.role }}:</b>

                <span v-if="m.thinking" class="typing">
                    <span></span>
                    <span></span>
                    <span></span>
                </span>

                <span v-else>
                    {{ m.text }}
                </span>
            </div>
        </div>

        <div class="input">
            <textarea v-model="input" ref="inputRef" :disabled="loading" @input="adjustHeight"
                @keydown.enter.exact.prevent="send" placeholder="Type message..." rows="1"></textarea>

            <button v-if="!loading" @click="send" class="send-btn">
                Send
            </button>
            <button v-else @click="stopGeneration" class="stop-btn">
                Stop
            </button>
        </div>

    </div>
</template>

<script setup>
import {
    ref, computed, nextTick,
    watch
} from "vue"
import { FIRSTTOPIC } from "./mini/questions"
import { evaluate } from "./mini/evaluator"


const topic = ref("")
const topicSelected = ref(false)
const miniStep = ref(0)

const topics = [
    "Frontend",
    "Backend",
    "QA",
    "AWS",
    "DevOps",
    "React",
    "Vue",
    "JavaScript",
    "TypeScript",
    "Golang"
]

const props = defineProps([
    "userId",
    "mode",
    "judge",
    "edition"
])

const chooseTopic = async (t) => {
    topic.value = t
    topicSelected.value = true
    miniStep.value = 0

    if (props.edition === "mini") {
        messages.value = []

        messages.value.push({
            role: "ai",
            text: FIRSTTOPIC[t][0]
        })

        miniStep.value = 1
        return
    }

    input.value = "Start"
    await send()
}




const input = ref("")
const messages = ref([])
const loading = ref(false)
if (
    props.edition === "mini" &&
    props.mode === "chat"
) {
    topicSelected.value = true

    messages.value.push({
        role: "ai",
        text:
            "Hi. Tell me how coding, architecture, frontend and backend connect in real systems."
    })
}


const autoScroll = ref(true)
const chatRef = ref(null)
const score = ref(null)
const inputRef = ref(null)

let controller = null

watch(autoScroll, async (enabled) => {
    if (enabled) {
        await scrollToBottom()
    }
})

const adjustHeight = () => {
    const el = inputRef.value
    if (!el) return
    el.style.height = 'auto'
    el.style.height = el.scrollHeight + 'px'
}

const scoreClass = computed(() => {
    if (score.value === null) return ""
    if (score.value <= 5) return "bad"
    if (score.value <= 8) return "mid"
    return "good"
})

const scrollToBottom = async () => {
    if (!autoScroll.value) return

    await nextTick()

    const el = chatRef.value
    if (el) el.scrollTop = el.scrollHeight
}

const stopGeneration = () => {
    controller?.abort()
    loading.value = false
}

const send = async () => {

    if (props.edition === "mini") {

        if (!input.value.trim()) return

        const text = input.value
        input.value = ""

        messages.value.push({
            role: "user",
            text
        })

        loading.value = true

        const aiIndex = messages.value.length

        messages.value.push({
            role: "ai",
            text: "",
            thinking: true
        })

        await scrollToBottom()

        const nextQuestion =
            props.mode === "chat"
                ? ""
                : (
                    FIRSTTOPIC[topic.value]?.[miniStep.value]
                    || "Interview complete."
                )

        const result =
            evaluate(
                text,
                score.value,
                nextQuestion,
                props.mode)

        score.value = result.score

        messages.value[aiIndex].thinking = false

        const words = result.text.split(" ")

        for (const word of words) {
            messages.value[aiIndex].text += word + " "
            await scrollToBottom()
            await new Promise(r => setTimeout(r, 55))
        }

        if (props.mode === "chat") {

            miniStep.value++

            if (miniStep.value >= 3) {
                messages.value.push({
                    role: "ai",
                    text:
                        "Good discussion. That covers the core architectural reasoning."
                })

                loading.value = false
                return
            }

            loading.value = false
            return
        }
        miniStep.value++
        if (
            FIRSTTOPIC[topic.value]?.[miniStep.value]
        ) {
            loading.value = false
            return
        }

        messages.value.push({
            role: "ai",
            text: "Interview complete."
        })
        loading.value = false
        return
    }

    if (!input.value.trim() || loading.value) return

    controller = new AbortController()

    const text = input.value
    input.value = ""

    if (inputRef.value) {
        inputRef.value.style.height = 'auto'
    }

    loading.value = true
    score.value = null

    messages.value.push({
        role: "user",
        text
    })

    const aiIndex = messages.value.length

    messages.value.push({
        role: "ai",
        text: "",
        thinking: true
    })

    await scrollToBottom()

    try {
        const res = await fetch(
            "http://localhost:8080/ai",
            {
                method: "POST",
                signal: controller.signal,
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    user_id: props.userId,
                    mode: props.mode,
                    judge: props.judge,
                    topic: topic.value,
                    input: text
                })
            }
        )

        let buffer = ""

        const reader = res.body.getReader()
        const decoder = new TextDecoder()

        while (true) {
            const { value, done } = await reader.read()

            if (done) break

            buffer += decoder.decode(value, { stream: true })

            const parts = buffer.split("\n\n")

            for (let i = 0; i < parts.length - 1; i++) {
                const block = parts[i]

                const lines = block.split("\n")

                const eventLine =
                    lines.find(l => l.startsWith("event:"))

                const dataLine =
                    lines.find(l => l.startsWith("data:"))

                if (!dataLine) continue

                const data =
                    dataLine.replace("data: ", "")

                if (eventLine === "event: score") {
                    score.value = Number(data)
                    continue
                }

                if (eventLine === "event: done") {
                    continue
                }

                messages.value[aiIndex].thinking = false
                messages.value[aiIndex].text += data

                await scrollToBottom()
            }

            buffer = parts[parts.length - 1]
        }

    } catch (err) {
        console.error(err)
    }

    loading.value = false

}
</script>

<style scoped>
.container {
    height: calc(100vh - 48px);
    margin: 24px auto;
    display: flex;
    flex-direction: column;
    gap: 12px;
    /* overflow: hidden; */
    padding: 0 4px;
}

.controls {
    display: flex;
    align-items: center;
    gap: 10px;
}

.user,
.mode-badge {
    background: #1e293b;
    color: #cbd5e1;
    padding: 8px 14px;
    border-radius: 10px;
    font-size: 14px;
}

.chat {
    flex: 1;

    background: #111827;
    padding: 14px;
    overflow-y: auto;
    border-radius: 14px;
    border: 1px solid #1f2937;
    margin: 0;



}

.chat.bad {
    box-shadow:
        0 0 0 3px rgba(255, 90, 90, .22),
        0 0 25px rgba(255, 90, 90, .08);
}

.chat.mid {
    box-shadow:
        0 0 0 3px rgba(255, 210, 100, .18),
        0 0 25px rgba(255, 210, 100, .07);
}

.chat.good {
    box-shadow:
        0 0 0 3px rgba(100, 255, 180, .18),
        0 0 25px rgba(100, 255, 180, .07);
}

.msg {
    margin-bottom: 12px;
    white-space: pre-wrap;
    line-height: 1.5;
}

.msg.user {
    color: #4fc3f7;
}

.msg.ai {
    color: #81c784;
}

.input {
    flex-shrink: 0;
    display: flex;
    align-items: flex-end;
    gap: 10px;
    background: #1e293b;
    padding: 10px 14px;
    border-radius: 12px;
    border: 1px solid #334155;
    margin-bottom: 2px;
}

textarea {
    flex: 1;
    background: transparent;
    border: none;
    color: white;
    outline: none;
    resize: none;
    padding: 4px 0;
    font-family: inherit;
    font-size: 16px;
    line-height: 1.5;
    max-height: 30vh;
    overflow-y: auto;
    display: block;
}

textarea:focus {
    border-color: #4fc3f7;
}

textarea::-webkit-scrollbar {
    width: 4px;
}

textarea::-webkit-scrollbar-thumb {
    background: #334155;
    border-radius: 10px;
}

button {
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;

    padding: 12px 16px;
    border: none;
    border-radius: 10px;
    background: #334155;
    color: white;
    cursor: pointer;
    transition: .2s;
}

button:hover {
    background: #475569;
}

button:active {
    transform: scale(.98);
}

.checkbox {
    display: flex;
    align-items: center;
    gap: 8px;
    color: #cbd5e1;
    font-size: 14px;
}

.checkbox input {
    accent-color: #4fc3f7;
    width: 16px;
    height: 16px;
}

.typing {
    display: inline-flex;
    gap: 6px;
    align-items: center;
}

.typing span {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: #81c784;
    animation: bounce 1.4s infinite ease-in-out;
}

.typing span:nth-child(2) {
    animation-delay: .2s;
}

.typing span:nth-child(3) {
    animation-delay: .4s;
}

@keyframes bounce {

    0%,
    80%,
    100% {
        transform: translateY(0);
        opacity: .35;
    }

    40% {
        transform: translateY(-6px);
        opacity: 1;
    }

}

.topic-cloud {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
}

.topic-pill {
    border: none;
    background: #1e293b;
    color: #4fc3f7;
    padding: 10px 16px;
    border-radius: 999px;
    cursor: pointer;
    transition: .2s;
}

.topic-pill:hover {
    transform: translateY(-2px);
    background: #334155;
}
</style>