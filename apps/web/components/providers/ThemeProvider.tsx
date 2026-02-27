"use client";

import { createContext, useContext, useEffect, useMemo, useState } from "react";
import { ThemeProvider as MUIThemeProvider } from "@mui/material/styles";
import CssBaseline from "@mui/material/CssBaseline";
import { PaletteMode } from "@mui/material";
import { getTheme } from "../theme";

interface ThemeContextType {
   mode: PaletteMode;
   toggleTheme: () => void;
   setThemeMode: (mode: PaletteMode) => void;
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined);

export const useThemeMode = () => {
   const context = useContext(ThemeContext);
   if (!context) {
      throw new Error("useThemeMode must be used within ThemeProvider");
   }
   return context;
};

interface ThemeProviderProps {
   children: React.ReactNode;
}

export default function ThemeProvider({ children }: ThemeProviderProps) {
   const [mode, setMode] = useState<PaletteMode>("light");
   const [mounted, setMounted] = useState(false);

   // Load theme from localStorage on mount
   useEffect(() => {
      setMounted(true);
      const savedMode = localStorage.getItem("themeMode") as PaletteMode | null;
      if (savedMode) {
         setMode(savedMode);
      } else {
         // Check system preference
         const prefersDark = window.matchMedia(
            "(prefers-color-scheme: dark)",
         ).matches;
         setMode(prefersDark ? "dark" : "light");
      }
   }, []);

   // Save theme to localStorage when it changes
   useEffect(() => {
      if (mounted) {
         localStorage.setItem("themeMode", mode);
      }
   }, [mode, mounted]);

   const toggleTheme = () => {
      setMode((prevMode) => (prevMode === "light" ? "dark" : "light"));
   };

   const setThemeMode = (newMode: PaletteMode) => {
      setMode(newMode);
   };

   const theme = useMemo(() => getTheme(mode), [mode]);

   const contextValue = useMemo(
      () => ({
         mode,
         toggleTheme,
         setThemeMode,
      }),
      [mode],
   );

   // Prevent flash of wrong theme on initial load
   if (!mounted) {
      return null;
   }

   return (
      <ThemeContext.Provider value={contextValue}>
         <MUIThemeProvider theme={theme}>
            <CssBaseline />
            {children}
         </MUIThemeProvider>
      </ThemeContext.Provider>
   );
}
