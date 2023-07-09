<template>
    <label class="flex flex-col gap-1">
        <span>{{ label }} <span v-if="required" class="text-sm text-red-400">*</span></span>
        <input :type="number ? 'number' : 'text'" :value="modelValue ? modelValue : ''"
               :placeholder="placeholder ? placeholder : ''" :required="required" :min="min" :max="max"
               @input="$emit('update:modelValue', $event.target.value)"
               class="rounded-md border-0 px-3 py-2 ring-1 ring-black ring-opacity-30 transition duration-100 focus:ring-1 focus:ring-teal-400"/>
    </label>
</template>

<script setup lang="ts">
defineProps<{
    modelValue?: string,
    number?: boolean,
    min?: number,
    max?: number,
    label: string,
    required?: boolean,
    placeholder?: string
}>()
defineEmits([ "update:modelValue" ])
</script>

<style scoped lang="scss">
input {
  @apply focus:outline-none;
}
</style>
