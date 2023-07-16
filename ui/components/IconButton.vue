<template>
    <button type="button" @click="click"
            class="grid h-10 w-10 place-content-center rounded-lg bg-white text-gray-600 shadow-md ring-1 ring-black ring-opacity-5 transition duration-200 hover:text-teal-400">
        <Icon :name="icon" class="h-6 w-6"/>
    </button>
</template>

<script setup lang="ts">
defineProps<{
    click: () => void,
    icon: string
}>()
</script>
