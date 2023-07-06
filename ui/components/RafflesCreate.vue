<template>
    <div>
        <button type="button" @click="openModal"
                class="grid h-full w-10 place-content-center rounded-lg bg-white text-gray-600 shadow-md transition duration-200 hover:text-teal-400">
            <Icon name="heroicons:plus" class="h-6 w-6"/>
        </button>
    </div>

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
                                Додати новий розіграш
                            </HeadlessDialogTitle>

                            <hr class="mt-2">

                            <div class="mt-2 flex flex-col gap-2">
                                <label class="flex flex-col gap-1">
                                    <span>Назва <span class="text-red-400 text-sm">*</span></span>
                                    <input type="text"
                                           class="rounded-md border border-gray-400 px-3 py-2 transition duration-100 focus:border-teal-400"/>
                                </label>
                                <label class="flex flex-col gap-1">
                                    <span>Опис</span>
                                    <textarea rows="3"
                                              class="rounded-md border border-gray-400 px-3 py-2 transition duration-100 focus:border-teal-400"/>
                                </label>
                            </div>

                            <div class="mt-4 flex gap-4">
                                <button type="button" @click="closeModal"
                                        class="rounded-md shadow w-full bg-teal-100 px-4 py-2 text-sm font-medium text-teal-950 hover:bg-teal-200 focus:bg-teal-200">
                                    Додати
                                </button>
                                <button type="button" @click="closeModal"
                                        class="rounded-md shadow w-full bg-gray-100 px-4 py-2 text-sm font-medium text-black hover:bg-gray-200 focus:bg-gray-200">
                                    Закрити
                                </button>
                            </div>
                        </HeadlessDialogPanel>
                    </HeadlessTransitionChild>
                </div>
            </div>
        </HeadlessDialog>
    </HeadlessTransitionRoot>
</template>

<script setup lang="ts">
const isOpen = ref(false)

function closeModal() {
    isOpen.value = false
}

function openModal() {
    isOpen.value = true
}
</script>
