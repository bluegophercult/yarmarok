<template>
    <label class="flex flex-col gap-1">
        <span>{{ label }} <span v-if="required" class="text-red-400 text-sm">*</span></span>
        <textarea rows="3" :value="modelValue" @input="$emit('update:modelValue', $event.target.value)"
                  :required="required"
                  class="rounded-md border border-gray-400 px-3 py-2 transition duration-100 focus:border-teal-400"/>
    </label>
</template>

<script setup lang="ts">
defineProps<{
    modelValue: string,
    label: string,
    required: boolean,
}>()
defineEmits([ "update:modelValue" ])
</script>

<style scoped lang="scss">
textarea {
  @apply focus:outline-none focus:ring-0;
}
</style>
