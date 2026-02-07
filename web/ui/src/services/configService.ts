import { RuntimeConfig, validateConfig } from "../types/config";

const CONFIG_ENDPOINT = "/config";

export async function fetchConfig(): Promise<RuntimeConfig> {
    const baseURL = import.meta.env.VITE_API_BASE_URL ?? ""
    try {
        const response = await fetch(baseURL + CONFIG_ENDPOINT);

        if (!response.ok) {
            throw new Error(
                `Config endpoint returned ${response.status}: ${response.statusText}`
            );
        }

        const data = await response.json();
        const config = validateConfig(data);

        return config;
    } catch (error) {
        const message =
            error instanceof Error
                ? error.message
                : "Failed to fetch config from server";
        throw new Error(`Failed to load runtime configuration: ${message}`);
    }
}
