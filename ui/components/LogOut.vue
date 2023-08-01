<template>
    <IconButton :click="logOut" icon="heroicons:arrow-left-on-rectangle"/>
</template>

<script setup lang="ts">
function logOut() {
    location.reload()
}
</script>
