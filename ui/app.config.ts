export default defineAppConfig({
    title: "Ярмарок",
    apiBaseURL: process.env.NODE_ENV === "production" ? "https://yarmarock.com.ua" : "http://localhost:8081",
})
