<!-- <template>
    <div class="mini">

        <h2>{{ question }}</h2>

        <textarea v-model="answer" />

        <button @click="submit">
            Send
        </button>

        <div v-if="feedback">
            {{ feedback }}
        </div>

    </div>
</template>

<script setup>
import { ref } from "vue"
import { FIRSTTOPIC } from "./questions"
import { evaluate } from "./evaluator"

const step = ref(0)
const topic = "FRONTEND"

const question = ref(FIRSTTOPIC[topic][0])

const answer = ref("")
const feedback = ref("")

const submit = () => {
    feedback.value =
        evaluate(answer.value, step.value)

    step.value++

    if (FIRSTTOPIC[topic][step.value]) {
        question.value =
            FIRSTTOPIC[topic][step.value]
    }

    answer.value = ""
}
</script> -->