<template>
    <button :type="submit ? 'submit' : 'button'" @click="click" :disabled="disabled"
            :class="[
                disabled ?
                     'bg-gray-100 text-gray-400 ring-gray-600' :
                secondary ?
                    'bg-gray-100 text-black hover:bg-gray-200 focus:bg-gray-200 ring-gray-600' :
                    'bg-teal-100 text-teal-950 hover:bg-teal-200 focus:bg-teal-200 ring-teal-600',
                fullWidth ? 'w-full' : '',
            ]"
            class="rounded-md shadow px-4 py-2 text-sm font-medium ring-1 ring-opacity-5 focus:outline-none">
        <slot/>
    </button>
</template>

<script setup lang="ts">
defineProps<{
    click?: () => void,
    submit?: boolean,
    secondary?: boolean,
    disabled?: boolean,
    fullWidth?: boolean,
}>()
</script>
