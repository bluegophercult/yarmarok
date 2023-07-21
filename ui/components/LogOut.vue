<template>
    <IconButton :click="logOut" icon="heroicons:arrow-left-on-rectangle"/>
</template>

<script setup lang="ts">
const appConfig = useAppConfig()

function logOut() {
    const auth = useCookie(appConfig.authCookieName)
    auth.value = null
    location.reload()
}
</script>
