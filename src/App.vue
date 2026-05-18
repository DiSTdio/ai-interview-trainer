<template>
  <EditionSelect v-if="stage === 'edition'" @select="selectEdition" />

  <WEdition v-else-if="stage === 'winstall'" @proceed="stage = 'login'" />

  <SessionSetupCard v-else-if="stage === 'login'" @login="login" />

  <ChatContainer v-else-if="stage === 'chat'" :user-id="userId" :mode="mode" :judge="judge" :edition="edition"
    @logout="logout" />
</template>

<script setup>
import { ref } from "vue"

import EditionSelect from "./components/EditionSelect.vue"
import WEdition from "./components/WEdition.vue"
import SessionSetupCard from "./components/SessionSetupCard.vue"
import ChatContainer from "./components/ChatContainer.vue"

const stage = ref("edition")
const edition = ref("w")
const userId = ref("")
const mode = ref("chat")
const judge = ref(false)

const selectEdition = (selected) => {
  edition.value = selected

  if (selected === "mini") {
    stage.value = "login"
  } else {
    stage.value = "winstall"
  }
}

const login = ({ name, mode: selectedMode }) => {
  userId.value = name
  stage.value = "chat"

  if (selectedMode === "judge") {
    mode.value = "interview"
    judge.value = true
  } else {
    mode.value = selectedMode
    judge.value = false
  }
}

const logout = () => {
  stage.value = "edition"
  userId.value = ""
  mode.value = "chat"
  judge.value = false
}
</script>