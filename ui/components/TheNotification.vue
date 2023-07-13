<template>
    <client-only>
        <transition name="m-fade">
            <div v-if="shown"
                 class="fixed flex gap-2 items-center max-w-80 z-10 rounded-lg bg-white shadow-md ring-1 ring-black ring-opacity-5 right-4 p-2">
                <button type="button" @click="hideNotification">
                    <icon name="heroicons:x-mark"
                          class="w-7 h-7 text-gray-600 transition duration-200 hover:text-teal-400"/>
                </button>
                <p>{{ message }}</p>

                <icon v-if="type === 'info'" name="heroicons:information-circle" class="w-7 h-7 text-teal-400"/>
                <icon v-if="type === 'error'" name="heroicons:exclamation-triangle" class="w-7 h-7 text-red-400"/>
            </div>
        </transition>
    </client-only>
</template>

<script setup lang="ts">
import { useNotificationStore } from "~/store/notification"

const notificationStore = useNotificationStore()
const { shown, type, message } = storeToRefs(notificationStore)
const { hideNotification } = notificationStore
</script>
