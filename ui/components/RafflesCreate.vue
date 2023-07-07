<template>
    <div>
        <button type="button" @click="openModal"
                class="grid h-full w-10 place-content-center rounded-lg bg-white text-gray-600 shadow-md ring-1 ring-black ring-opacity-5 transition duration-200 hover:text-teal-400">
            <Icon name="heroicons:plus" class="h-6 w-6"/>
        </button>
    </div>

    <TheModal :is-open="isOpen" :close-modal="closeModal">
        <template #title>Додати новий розіграш</template>

        <div class="flex flex-col gap-2">
            <TheInput v-model="newRaffle.name" label="Назва" required/>
            <TheTextArea v-model="newRaffle.note" label="Опис"/>
            <transition name="m-fade">
                <p v-show="errorMsg" class="transition duration-200 mt-3 text-sm text-red-400">{{ errorMsg }}</p>
            </transition>
        </div>

        <div class="mt-4 flex gap-4">
            <TheButton :click="addRaffle" full-width>Додати</TheButton>
            <TheButton :click="closeModal" full-width secondary>Закрити</TheButton>
        </div>
    </TheModal>
</template>

<script setup lang="ts">
import { Ref } from "@vue/reactivity"
import { ValidationError } from "yup"
import { NewRaffle, newRaffleSchema } from "~/types/raffle"
import { useRaffleStore } from "~/store/raffle"

const isOpen = ref(false)
const errorMsg = ref("")

function closeModal() {
    isOpen.value = false
    setTimeout(() => {
        newRaffle.value = <NewRaffle>{
            name: "",
            note: "",
        }
    }, 200)
}

function openModal() {
    errorMsg.value = ""
    isOpen.value = true
}

const raffleStore = useRaffleStore()
const newRaffle: Ref<NewRaffle> = ref(<NewRaffle>{
    name: "",
    note: "",
})

function addRaffle() {
    newRaffleSchema.validate(newRaffle.value)
        .then(() => {
            raffleStore.addRaffle(newRaffle.value)
            closeModal()
        })
        .catch((e: ValidationError) => {
            errorMsg.value = e.message
        })
}
</script>
