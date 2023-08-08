<template>
    <TheModal :is-open="isOpen" :close-modal="closeModal">
        <template #title>Змінити учасника</template>

        <form @submit.prevent="updateRaffle">
            <div class="flex flex-col gap-2">
                <TheInput v-model="updatedParticipant.name" label="Ім'я" required/>
                <TheInput v-model="updatedParticipant.phone" label="Номер телефону" required/>
                <TheTextArea v-model="updatedParticipant.note" label="Нотатка"/>
            </div>

            <transition name="m-fade">
                <p v-show="errorMsg" class="mt-4 flex items-center gap-2 text-sm text-red-500 transition duration-200">
                    <Icon name="heroicons:exclamation-triangle" class="h-5 w-5"/>
                    {{ errorMsg }}
                </p>
            </transition>

            <div class="mt-4 flex gap-4">
                <TheButton submit full-width>Зберегти</TheButton>
                <TheButton :click="closeModal" secondary full-width>Закрити</TheButton>
            </div>
        </form>
    </TheModal>
</template>

<script setup lang="ts">
import { useParticipantStore } from "~/store/participant"
import { newParticipantSchema, Participant } from "~/types/participant"
import { ValidationError } from "yup"

const participantStore = useParticipantStore()

const props = defineProps<{
    participant: Participant,
    isOpen: boolean,
    closeModal: () => void,
}>()

const errorMsg = ref("")
const updatedParticipant = ref(<Participant>{ ...props.participant })
onBeforeUpdate(() => {
    setTimeout(() => {
        updatedParticipant.value = { ...props.participant }
        errorMsg.value = ""
    }, 200)
})

function updateRaffle() {
    newParticipantSchema.validate(updatedParticipant.value)
        .then(() => {
            participantStore.updateParticipant(updatedParticipant.value)
            props.closeModal()
        })
        .catch((e: ValidationError) => {
            errorMsg.value = e.message
        })
}
</script>
