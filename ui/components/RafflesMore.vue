<template>
    <HeadlessMenu as="div" class="relative">
        <HeadlessMenuButton as="div"
                            class="grid h-full w-10 place-content-center rounded-lg bg-white text-gray-600 shadow-md ring-1 ring-black ring-opacity-5 transition duration-200 hover:cursor-pointer hover:text-teal-400">
            <Icon name="heroicons:ellipsis-vertical" class="h-6 w-6"/>
        </HeadlessMenuButton>

        <transition name="m-fade">
            <HeadlessMenuItems
                    class="absolute right-0 mt-2 w-48 origin-top-right rounded-md bg-white shadow-lg ring-1 ring-black ring-opacity-5 divide-y divide-gray-200 focus:outline-none">
                <HeadlessMenuItem as="div" v-for="option in options" class="px-1 py-1" v-slot="{ active, disabled }"
                                  :disabled="!selectedRaffle">
                    <button type="button" @click="option.click" :class="[
                            active && !disabled ? 'bg-teal-100 text-teal-950' : 'text-gray-900',
                            disabled ? 'text-gray-300' : '',
                            'group flex w-full items-center rounded-md px-2 py-2 text-sm',
                        ]">
                        <Icon :active="active" :name="option.icon"
                              :class="disabled ? 'text-gray-400' : 'text-teal-400'" class="mr-2 h-5 w-5"/>
                        {{ option.text }}
                    </button>
                </HeadlessMenuItem>
            </HeadlessMenuItems>
        </transition>
    </HeadlessMenu>

    <RafflesUpdate v-if="selectedRaffle" :raffle="selectedRaffle" :is-open="isOpenUpdate"
                   :close-modal="() => isOpenUpdate = false"/>
    <RafflesDelete v-if="selectedRaffle" :raffle="selectedRaffle" :is-open="isOpenDelete"
                   :close-modal="() => isOpenDelete = false"/>
</template>

<script setup lang="ts">
import { useRaffleStore } from "~/store/raffle"

const raffleStore = useRaffleStore()
const { selectedRaffle } = storeToRefs(raffleStore)

const options = [
    {
        text: "Змінити",
        icon: "heroicons:pencil",
        click: () => isOpenUpdate.value = true,
    },
    {
        text: "Видалити",
        icon: "heroicons:trash",
        click: () => isOpenDelete.value = true,
    },
]

const isOpenUpdate = ref(false)
const isOpenDelete = ref(false)
</script>
