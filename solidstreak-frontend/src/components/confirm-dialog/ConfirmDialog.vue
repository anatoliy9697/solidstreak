<script setup lang="ts">

import ConfirmDialog from 'primevue/confirmdialog';

</script>

<template>
    <ConfirmDialog group="headless" class="px-4 pb-4 pt-2">
        <template #container="{ message, acceptCallback, rejectCallback }">
            <h2 v-if="message.header" class="mb-2">{{ message.header }}</h2>
            <p class="mb-2">{{ message.message }}</p>
            <div class="flex gap-2">
                <button
                    @click="acceptCallback"
                    class="w-1/2 px-4 py-2 rounded-md border border-gray-300 bg-gray-100 text-gray-800 font-medium hover:bg-gray-200"
                >Yes</button>
                <button
                    @click="rejectCallback"
                    class="w-1/2 px-4 py-2 rounded-md border border-gray-300 bg-gray-100 text-red-800 font-medium hover:bg-gray-200"
                >No</button>
            </div>
        </template>
    </ConfirmDialog>
</template>