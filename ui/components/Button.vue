<template>
    <button type="button" @click="click"
            :class="[
                secondary ? 'bg-gray-100 text-black hover:bg-gray-200 focus:bg-gray-200' : 'bg-teal-100 text-teal-950 hover:bg-teal-200 focus:bg-teal-200',
                fullWidth ? 'w-full' : '',
                'rounded-md shadow px-4 py-2 text-sm font-medium'
            ]">
        <slot/>
    </button>
</template>

<script setup lang="ts">
defineProps<{
    click: () => void,
    secondary: boolean,
    fullWidth: boolean,
}>()
</script>
