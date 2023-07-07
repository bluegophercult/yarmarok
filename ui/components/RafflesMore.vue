<template>
    <HeadlessMenu as="div" class="relative">
        <HeadlessMenuButton as="div"
                            class="grid h-full w-10 place-content-center rounded-lg bg-white text-gray-600 shadow-md transition duration-200 hover:text-teal-400 ring-1 ring-black ring-opacity-5">
            <Icon name="heroicons:ellipsis-vertical" class="h-6 w-6"/>
        </HeadlessMenuButton>

        <transition name="m-fade">
            <HeadlessMenuItems
                    class="absolute right-0 mt-2 w-48 origin-top-right rounded-md bg-white shadow-lg ring-1 ring-black ring-opacity-5 divide-y divide-gray-200">
                <HeadlessMenuItem as="div" class="px-1 py-1" v-slot="{ active }">
                    <button type="button" :class="[
                            active ? 'bg-teal-100 text-teal-950' : 'text-gray-900',
                            'group flex w-full items-center rounded-md px-2 py-2 text-sm',
                        ]">
                        <Icon :active="active" name="heroicons:pencil" class="mr-2 h-5 w-5 text-teal-400"/>
                        Змінити
                    </button>
                </HeadlessMenuItem>
                <HeadlessMenuItem as="div" class="px-1 py-1" v-slot="{ active }">
                    <button type="button" @click="openDeleteModal" :class="[
                          active ? 'bg-teal-100 text-teal-950' : 'text-gray-900',
                          'group flex w-full items-center rounded-md px-2 py-2 text-sm',
                        ]">
                        <Icon :active="active" name="heroicons:trash" class="mr-2 h-5 w-5 text-teal-400"/>
                        Видалити
                    </button>
                </HeadlessMenuItem>
            </HeadlessMenuItems>
        </transition>
    </HeadlessMenu>

    <RafflesDelete :raffle="selectedRaffle" :is-open="isOpenDelete" :close-modal="closeDeleteModal"/>
</template>

<script setup lang="ts">
import { useRaffleStore } from "~/store/raffle"

const raffleStore = useRaffleStore()
const { selectedRaffle } = storeToRefs(raffleStore)

const isOpenDelete = ref(false)

function closeDeleteModal() {
    isOpenDelete.value = false
}

function openDeleteModal() {
    isOpenDelete.value = true
}
</script>
