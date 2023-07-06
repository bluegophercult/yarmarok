<template>
    <HeadlessTransitionRoot :show="isOpen" as="template">
        <HeadlessDialog as="div" @close="closeModal" class="relative z-10">
            <HeadlessTransitionChild
                    as="template"
                    enter="duration-200 ease-out"
                    enter-from="opacity-0"
                    enter-to="opacity-100"
                    leave="duration-200 ease-in"
                    leave-from="opacity-100"
                    leave-to="opacity-0">
                <div class="fixed inset-0 bg-black bg-opacity-50"/>
            </HeadlessTransitionChild>

            <div class="fixed inset-0 overflow-y-auto">
                <div class="flex min-h-full items-center justify-center p-4 text-center">
                    <HeadlessTransitionChild
                            as="template"
                            enter="duration-300 ease-out"
                            enter-from="opacity-0 scale-95"
                            enter-to="opacity-100 scale-100"
                            leave="duration-200 ease-in"
                            leave-from="opacity-100 scale-100"
                            leave-to="opacity-0 scale-95">
                        <HeadlessDialogPanel
                                class="w-full max-w-md transform overflow-hidden rounded-xl bg-white p-6 text-left align-middle shadow-xl transition-all">
                            <HeadlessDialogTitle as="h3" class="text-lg font-medium leading-6 text-gray-900">
                                <slot name="title"/>
                            </HeadlessDialogTitle>
                            <hr class="my-2">
                            <slot/>
                        </HeadlessDialogPanel>
                    </HeadlessTransitionChild>
                </div>
            </div>
        </HeadlessDialog>
    </HeadlessTransitionRoot>
</template>

<script setup lang="ts">
defineProps<{
    isOpen: boolean,
    closeModal: (_close: boolean) => void,
}>()
</script>
