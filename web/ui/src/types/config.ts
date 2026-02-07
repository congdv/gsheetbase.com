export interface RuntimeConfig {
    apiBaseUrl: string;
    workerBaseUrl: string;
    landingPageUrl: string;
    forceProd: boolean;
}

export function validateConfig(obj: unknown): RuntimeConfig {
    if (typeof obj !== "object" || obj === null) {
        throw new Error("Config must be an object");
    }

    const config = obj as Record<string, unknown>;

    if (typeof config.apiBaseUrl !== "string" || !config.apiBaseUrl) {
        throw new Error("apiBaseUrl must be a non-empty string");
    }
    if (typeof config.workerBaseUrl !== "string" || !config.workerBaseUrl) {
        throw new Error("workerBaseUrl must be a non-empty string");
    }
    if (typeof config.landingPageUrl !== "string" || !config.landingPageUrl) {
        throw new Error("landingPageUrl must be a non-empty string");
    }
    if (typeof config.forceProd !== "boolean") {
        throw new Error("forceProd must be a boolean");
    }

    return {
        apiBaseUrl: config.apiBaseUrl,
        workerBaseUrl: config.workerBaseUrl,
        landingPageUrl: config.landingPageUrl,
        forceProd: config.forceProd,
    };
}
