<template>
    <label class="flex flex-col gap-1">
        <span>{{ label }} <span v-if="required" class="text-sm text-red-400">*</span></span>
        <input :type="number ? 'number' : 'text'" :value="modelValue ? modelValue.toString() : ''"
               :placeholder="placeholder ? placeholder.toString() : ''" :min="min" :max="max"
               @input="$emit('update:modelValue', number ? Number($event.target.value.trim()) : $event.target.value.trim())"
               class="rounded-md border-0 px-3 py-2 ring-1 ring-black ring-opacity-30 transition duration-100 focus:outline-none focus:ring-1 focus:ring-teal-400"/>
    </label>
</template>

<script setup lang="ts">
defineProps<{
    modelValue?: string | number,
    number?: boolean,
    min?: number,
    max?: number,
    label: string,
    required?: boolean,
    placeholder?: string | number
}>()
defineEmits([ "update:modelValue" ])
</script>
