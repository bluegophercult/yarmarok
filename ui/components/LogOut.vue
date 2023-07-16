<template>
    <IconButton :click="logOut" icon="heroicons:arrow-left-on-rectangle"/>
</template>

<script setup lang="ts">
function logOut() {
    const auth = useCookie("GCP_IAP_UID")
    auth.value = null
    location.reload()
}
</script>
