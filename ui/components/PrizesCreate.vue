<template>
    <TheButton @click="openModal" full-width>Додати приз</TheButton>

    <TheModal :is-open="isOpen" :close-modal="closeModal">
        <template #title>Додати новий приз</template>

        <form @submit.prevent="addPrize">
            <div class="flex flex-col gap-2">
                <TheInput v-model="newPrize.name" label="Назва" required/>
                <TheInput v-model="newPrize.ticketCost" number min="1" label="Ціна купону" required/>
                <TheTextArea v-model="newPrize.note" label="Опис"/>
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
import { usePrizeStore } from "~/store/prize"
import { Ref } from "@vue/reactivity"
import { ValidationError } from "yup"
import { NewPrize, newPrizeSchema } from "~/types/prize"

const prizeStore = usePrizeStore()
const newPrize: Ref<NewPrize> = ref(<NewPrize>{
    name: "",
    ticketCost: 10,
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
        newPrize.value = <NewPrize>{
            name: "",
            ticketCost: 10,
            note: "",
        }
    }, 200)
}

function addPrize() {
    newPrizeSchema.validate(newPrize.value)
        .then(() => {
            prizeStore.addPrize(newPrize.value)
            closeModal()
        })
        .catch((e: ValidationError) => {
            errorMsg.value = e.message
        })
}
</script>
