// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    modules: [
        "@nuxtjs/tailwindcss",
        "@nuxtjs/google-fonts",
        "@pinia/nuxt",
        "nuxt-headlessui",
        "nuxt-icon",
    ],
    devtools: { enabled: true },
    tailwindcss: {
        cssPath: "~/assets/css/tailwind.scss",
    },
    googleFonts: {
        prefetch: true,
        preconnect: true,
        preload: true,
        families: {
            Roboto: true,
        },
    },
    pinia: {
        autoImports: [
            "defineStore",
            "storeToRefs",
        ],
    },
})
