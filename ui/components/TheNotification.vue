<template>
    <client-only>
        <transition name="m-fade">
            <div v-if="shown"
                 class="absolute top-16 right-4 z-10 flex items-center gap-2 rounded-lg bg-white p-2 shadow-md ring-1 ring-black ring-opacity-5 max-w-80">
                <button type="button" @click="hideNotification">
                    <Icon name="heroicons:x-mark"
                          class="h-7 w-7 text-gray-600 transition duration-200 hover:text-teal-400"/>
                </button>
                <p>{{ message }}</p>

                <Icon v-if="type === 'info'" name="heroicons:information-circle" class="h-7 w-7 text-teal-400"/>
                <Icon v-if="type === 'error'" name="heroicons:exclamation-triangle" class="h-7 w-7 text-red-400"/>
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
