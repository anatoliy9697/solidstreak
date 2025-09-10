<script setup>
import { ref } from 'vue'

const props = defineProps({
  mode: String,
  inputHabit: {
    id: String,
    name: String,
    description: String
    // color: String
  }
})

const draftHabit = ref(!!props.inputHabit ? { ...props.inputHabit } : { name: null, description: null })

</script>

<template>
  <div class="habit-editor">

    <div class="header-row">
      <h2>{{ mode === "create" ? "New habit" : "Edit habit" }}</h2>
      <button class="popup-close-btn" aria-label="Close">&times;</button>
    </div>

    <label>
        Name
        <input v-model="draftHabit.name" type="text" placeholder="Enter a habit name" />
    </label>

    <label>
        Description
        <textarea v-model=" draftHabit.description" placeholder="Add a short description (optional)"></textarea>
    </label>

    <div>
      <button class="popup-action-btn">{{ mode === "create" ? "Add habit" : "Save changes" }}</button>
    </div>

  </div>
</template>

<style scoped>
.habit-editor {
  text-align: center;
}

.header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

label {
  display: block;
  margin-bottom: 1em;
  text-align: left;
}

input,
textarea {
  display: block;
  width: 100%;
  box-sizing: border-box;
}
</style>