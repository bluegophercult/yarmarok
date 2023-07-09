<template>
    <div>
        <button type="button" @click="openModal"
                class="grid h-full w-10 place-content-center rounded-lg bg-white text-gray-600 shadow-md ring-1 ring-black ring-opacity-5 transition duration-200 hover:text-teal-400">
            <Icon name="heroicons:plus" class="h-6 w-6"/>
        </button>
    </div>

    <TheModal :is-open="isOpen" :close-modal="closeModal">
        <template #title>Додати новий розіграш</template>

        <form @submit.prevent="addRaffle">
            <div class="flex flex-col gap-2">
                <TheInput v-model="newRaffle.name" label="Назва" required/>
                <TheTextArea v-model="newRaffle.note" label="Опис"/>
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
import { Ref } from "@vue/reactivity"
import { ValidationError } from "yup"
import { NewRaffle, newRaffleSchema } from "~/types/raffle"
import { useRaffleStore } from "~/store/raffle"

const raffleStore = useRaffleStore()
const newRaffle: Ref<NewRaffle> = ref(<NewRaffle>{
    name: "",
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
        newRaffle.value = <NewRaffle>{
            name: "",
            note: "",
        }
    }, 200)
}

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
