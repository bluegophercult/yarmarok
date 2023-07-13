import { useFetch } from "#app"

export const useApiFetch: typeof useFetch = (request, options) => {
    const { apiBaseURL } = useAppConfig()
    return useFetch(request, {
        baseURL: apiBaseURL,
        headers: {
            "Content-Type": "application/json",
            "X-Goog-Authenticated-User-Id": "test", // TODO: Get from Cookies
        },
        ...options,
    })
}
