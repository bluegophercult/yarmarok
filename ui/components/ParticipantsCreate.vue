<template>
    <TheButton @click="openModal" :disabled="!selectedRaffle" full-width>Додати учасника</TheButton>

    <TheModal :is-open="isOpen" :close-modal="closeModal">
        <template #title>Додати нового учасника</template>

        <form @submit.prevent="addParticipant">
            <div class="flex flex-col gap-2">
                <TheInput v-model="newParticipant.name" label="Ім'я" required/>
                <TheInput v-model="newParticipant.phone" label="Номер телефону" required/>
                <TheTextArea v-model="newParticipant.note" label="Нотатка"/>
            </div>

            <transition name="m-fade">
                <p v-show="errorMsg" class="mt-4 flex items-center gap-2 text-sm text-red-500 transition duration-200">
                    <Icon name="heroicons:exclamation-triangle" class="h-5 w-5"/>
                    {{ errorMsg }}
                </p>
            </transition>

            <div class="mt-4 flex gap-4">
                <TheButton submit full-width>Додати</TheButton>
                <TheButton :click="closeModal" secondary full-width>Закрити</TheButton>
            </div>
        </form>
    </TheModal>
</template>

<script setup lang="ts">
import { useParticipantStore } from "~/store/participant"
import { Ref } from "@vue/reactivity"
import { ValidationError } from "yup"
import { NewParticipant, newParticipantSchema } from "~/types/participant"
import { useRaffleStore } from "~/store/raffle"
import { useNotificationStore } from "~/store/notification"

const raffleStore = useRaffleStore()
const { selectedRaffle } = storeToRefs(raffleStore)

const { showError } = useNotificationStore()

const participantStore = useParticipantStore()
const newParticipant: Ref<NewParticipant> = ref(<NewParticipant>{
    name: "",
    phone: "",
    note: "",
})

const isOpen = ref(false)
const errorMsg = ref("")

function openModal() {
    errorMsg.value = ""
    isOpen.value = true
}

function closeModal() {
    isOpen.value = false
    setTimeout(() => {
        newParticipant.value = <NewParticipant>{
            name: "",
            phone: "",
            note: "",
        }
    }, 200)
}

function addParticipant() {
    newParticipantSchema.validate(newParticipant.value)
        .then(() => {
            participantStore.addParticipant(selectedRaffle.value!.id, newParticipant.value).catch(e => {
                console.error(e)
                showError("Не вдалося створити учасника!")
            })
            closeModal()
        })
        .catch((e: ValidationError) => {
            errorMsg.value = e.message
        })
}
</script>
