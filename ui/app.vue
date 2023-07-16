<template>
    <div class="grid gap-2 p-4">
        <div class="flex gap-2 justify-between">
            <Raffles class="flex-grow"/>
            <LogOut/>
        </div>
        <Prizes/>
    </div>

    <TheNotification/>
</template>

<script setup lang="ts">
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
</script>
