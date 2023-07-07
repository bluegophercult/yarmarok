<template>
    <label class="flex flex-col gap-1">
        <span>{{ label }} <span v-if="required" class="text-sm text-red-400">*</span></span>
        <input type="text" :value="modelValue" @input="$emit('update:modelValue', $event.target.value)"
               :required="required"
               class="rounded-md border-0 ring-1 ring-opacity-30 ring-black px-3 py-2 transition duration-100 focus:!ring-1 focus:ring-teal-400"/>
    </label>
</template>

<script setup lang="ts">
defineProps<{
    modelValue: string,
    label: string,
    required?: boolean,
}>()
defineEmits([ "update:modelValue" ])
</script>

<style scoped lang="scss">
input {
  @apply focus:outline-none focus:ring-0;
}
</style>
