const isProduction = process.env.NODE_ENV === "production"

export default defineAppConfig({
    title: "Ярмарок",
    apiBaseURL: isProduction ? "https://yarmarock.com.ua" : "http://localhost:8081",
})
