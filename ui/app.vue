<template>
    <div class="pt-4 pl-5">
        <HeadlessListbox class="w-72" v-model="selectedPerson">
            <div class="relative mt-1">
                <HeadlessListboxButton
                        class="relative w-full cursor-default rounded-lg bg-white py-2 pr-10 pl-3 text-left shadow-md focus:outline-none focus-visible:border-indigo-500 focus-visible:ring-2 focus-visible:ring-white focus-visible:ring-opacity-75 focus-visible:ring-offset-2 focus-visible:ring-offset-orange-300 sm:text-sm">
                    <span class="block truncate">{{ selectedPerson.name }}</span>
                    <span class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
                        <Icon name="heroicons:chevron-up-down" class="h-5 w-5 text-gray-400"/>
                    </span>
                </HeadlessListboxButton>

                <transition leave-active-class="transition duration-100 ease-in" leave-from-class="opacity-100"
                            leave-to-class="opacity-0">
                    <HeadlessListboxOptions
                            class="absolute mt-1 max-h-40 w-full overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none sm:text-sm">
                        <HeadlessListboxOption v-slot="{ active, selected }" v-for="person in people" :key="person.name"
                                               :value="person" as="template">
                            <li :class="[ active ? 'bg-amber-100 text-amber-900' : 'text-gray-900', 'relative cursor-default select-none py-2 pl-10 pr-4',]">
                                <span :class="[ selected ? 'font-medium' : 'font-normal', 'block truncate',]">{{
                                        person.name
                                    }}</span>
                                <span v-if="selected"
                                      class="absolute inset-y-0 left-0 flex items-center pl-3 text-amber-600">
                                    <Icon name="heroicons:check" class="h-5 w-5"/>
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
import { useUserStore } from "~/store/user"

const appConfig = useAppConfig()
useHead({
    title: appConfig.title,
    link: [
        { rel: "apple-touch-icon", sizes: "180x180", href: "/apple-touch-icon.png" },
        { rel: "icon", type: "image/png", sizes: "32x32", href: "/favicon-32x32.png" },
        { rel: "icon", type: "image/png", sizes: "16x16", href: "/favicon-16x16.png" },
        { rel: "manifest", href: "/site.webmanifest" },
    ],
})

const userStore = useUserStore()
const { user } = storeToRefs(userStore)

const people = [
    { name: user.value },
    { name: "Wade Cooper" },
    { name: "Arlene Mccoy" },
    { name: "Devon Webb" },
    { name: "Tom Cook" },
    { name: "Tanya Fox" },
    { name: "Hellen Schmidt" },
]
const selectedPerson = ref(people[0])
</script>
