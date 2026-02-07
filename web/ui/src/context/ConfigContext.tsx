import React, { createContext, useContext } from "react";
import { RuntimeConfig } from "../types/config";

interface ConfigContextType {
    config: RuntimeConfig;
}

const ConfigContext = createContext<ConfigContextType | undefined>(undefined);

interface ConfigProviderProps {
    config: RuntimeConfig;
    children: React.ReactNode;
}

export function ConfigProvider({ config, children }: ConfigProviderProps) {
    return (
        <ConfigContext.Provider value={{ config }}>
            {children}
        </ConfigContext.Provider>
    );
}

export function useConfig(): RuntimeConfig {
    const context = useContext(ConfigContext);
    if (!context) {
        throw new Error("useConfig must be used within ConfigProvider");
    }
    return context.config;
}
