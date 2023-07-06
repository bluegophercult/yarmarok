<template>
    <div class="w-72">
        <HeadlessListbox v-model="selectedPerson">
            <div class="relative">
                <HeadlessListboxButton
                        class="relative w-full cursor-default hover:cursor-pointer group rounded-lg bg-white py-2 pl-3 pr-10 text-left shadow-md">
                    <span class="block truncate">{{ selectedPerson.name }}</span>
                    <span class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
                        <Icon name="heroicons:chevron-up-down"
                              class="h-5 w-5 text-gray-600 group-hover:text-teal-400 transition duration-200"/>
                    </span>
                </HeadlessListboxButton>

                <transition name="m-fade">
                    <HeadlessListboxOptions
                            class="absolute mt-2 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5">
                        <HeadlessListboxOption v-slot="{ active, selected, disabled }" v-for="raffle in raffles"
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
let raffles: Array<{
    id: string,
    name: string,
}> = [
    { id: "1", name: "Фестиваль їжі" },
    { id: "2", name: "Atlas weekend" },
]

if (raffles.length === 0) {
    raffles.push({ id: "", name: "Немає ярмарок" })
}
const selectedPerson = ref(raffles[0])
</script>

<style scoped>

</style>