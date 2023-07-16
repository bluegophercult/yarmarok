import { useFetch } from "#app"

export const useApiFetch: typeof useFetch = (request, options) => {
    const { apiBaseURL } = useAppConfig()
    return useFetch(request, {
        baseURL: apiBaseURL,
        headers: {
            "Content-Type": "application/json",
        },
        ...options,
    })
}
