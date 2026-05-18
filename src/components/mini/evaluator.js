const practiceReplies = [
    "Good question. Think about how state propagates through reactive dependencies.",
    "Interesting direction. Try breaking the problem into smaller composable units.",
    "That’s close. Consider performance implications and rerender cost.",
    "A stronger answer would mention architecture tradeoffs explicitly.",
    "Nice thinking. Now connect that to production-scale maintainability."
]

const chatReplies = [
    "Interesting. How would that affect system scalability?",
    "Good direction. How would frontend and backend coordinate there?",
    "That touches architecture. What tradeoffs would you consider?",
    "Good start. What happens if that assumption breaks at scale?",
    "What would be the biggest bottleneck there?"
]

export const evaluate = (
    answer,
    score,
    nextQuestion,
    mode
) => {

    if (mode === "chat") {
        return {
            score: 8,
            text:
                chatReplies[
                    Math.floor(
                        Math.random() *
                        chatReplies.length
                    )
                ]
        }
    }

    if (mode === "practice") {
        return {
            score: 8,
            text:
                practiceReplies[
                    Math.floor(
                        Math.random() *
                        practiceReplies.length
                    )
                ]
        }
    }

    const len = answer.trim().length

    if (len < 40) {
        return {
            score: 4,
            text:
                `Not quite. Needs more depth.\n\n${nextQuestion}`
        }
    }

    if (len < 120) {
        return {
            score: 7,
            text:
                `Solid answer. Could go deeper.\n\n${nextQuestion}`
        }
    }

    return {
        score: 9,
        text:
            `Excellent depth and reasoning.\n\n${nextQuestion}`
    }
}