<template>
    <div class="grid gap-2 p-4">
        <div class="flex justify-between gap-2">
            <Raffles class="flex-grow"/>
            <!-- TODO: Add if there will be a way to log out -->
            <!-- <LogOut/> -->
        </div>
        <div class="grid grid-cols-1 gap-2 md:grid-cols-4">
            <div class="w-full">
                <Prizes/>
            </div>
            <div class="col-span-1 w-full md:col-span-2">
                <PrizeDetails/>
            </div>
            <div class="w-full">
                <Participants/>
            </div>
        </div>
    </div>

    <TheNotification/>
</template>

<script setup lang="ts">
import { useStateStore } from "~/store/state"

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

const stateStore = useStateStore()
stateStore.init()
</script>
