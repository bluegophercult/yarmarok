<template>
    <div class="w-full sm:w-72">
        <HeadlessListbox v-model="selectedRaffle">
            <div class="relative">
                <HeadlessListboxButton
                        class="relative w-full cursor-default rounded-lg bg-white py-2 pr-10 pl-3 text-left shadow-md ring-1 ring-black ring-opacity-5 group hover:cursor-pointer">
                    <span class="block truncate">{{ selectedRaffle ? selectedRaffle.name : "Немає розіграшів" }}</span>
                    <span class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
                        <Icon name="heroicons:chevron-up-down"
                              class="h-5 w-5 text-gray-600 transition duration-200 group-hover:text-teal-400"/>
                    </span>
                </HeadlessListboxButton>

                <transition name="m-fade">
                    <HeadlessListboxOptions
                            class="absolute mt-2 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
                        <HeadlessListboxOption v-if="!selectedRaffle" disabled>
                            <li class="select-none px-4 py-2 text-gray-400">
                                Пусто
                            </li>
                        </HeadlessListboxOption>
                        <HeadlessListboxOption v-slot="{ active, selected, disabled }"
                                               v-for="raffle in raffles.slice().reverse()"
                                               :key="raffle.id" :value="raffle" :disabled="raffle.id === ''"
                                               as="template">
                            <li :class="[
                                active && !disabled ? 'bg-teal-100 text-teal-950' : 'text-gray-900',
                                disabled ? 'text-gray-600' : '',
                                 'relative cursor-default hover:cursor-pointer select-none py-2 pl-10 pr-4',
                                ]">
                                <span :class="[selected ? 'font-medium' : 'font-normal', 'block truncate',]">
                                    {{ raffle.name }}
                                </span>
                                <span v-if="selected"
                                      class="absolute inset-y-0 left-0 flex items-center pl-3 text-teal-400">
                                    <Icon name="heroicons:chevron-right-20-solid" class="h-5 w-5"/>
                                </span>
                            </li>
                        </HeadlessListboxOption>
                    </HeadlessListboxOptions>
                </transition>
            </div>
        </HeadlessListbox>
    </div>
</template>

<script setup lang="ts">
import { useRaffleStore } from "~/store/raffle"

const raffleStore = useRaffleStore()
const { raffles, selectedRaffle } = storeToRefs(raffleStore)
</script>
