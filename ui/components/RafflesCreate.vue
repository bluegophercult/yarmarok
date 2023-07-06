<template>
    <div>
        <button type="button" @click="openModal"
                class="grid h-full w-10 place-content-center rounded-lg bg-white text-gray-600 shadow-md transition duration-200 hover:text-teal-400">
            <Icon name="heroicons:plus" class="h-6 w-6"/>
        </button>
    </div>

    <Modal :is-open="isOpen" :close-modal="closeModal">
        <template #title>Додати новий розіграш</template>

        <div class="flex flex-col gap-2">
            <Input v-model="newRaffle.name" label="Назва" required/>
            <TextArea v-model="newRaffle.note" label="Опис"/>
        </div>

        <div class="mt-4 flex gap-4">
            <Button :click="addRaffle" full-width>Додати</Button>
            <Button :click="closeModal" full-width secondary>Закрити</Button>
        </div>
    </Modal>
</template>

<script setup lang="ts">
import { Ref } from "@vue/reactivity"
import { NewRaffle } from "~/types/raffle"
import { useRaffleStore } from "~/store/raffle"

const isOpen = ref(false)

function closeModal() {
    isOpen.value = false
}

function openModal() {
    isOpen.value = true
}

const raffleStore = useRaffleStore()
const newRaffle: Ref<NewRaffle> = ref(<NewRaffle>{})

function addRaffle() {
    raffleStore.addRaffle(newRaffle.value)
    closeModal()
    newRaffle.value = <NewRaffle>{}
}
</script>
