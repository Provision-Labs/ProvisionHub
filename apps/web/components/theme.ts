import { createTheme, PaletteMode } from "@mui/material/styles";

// Common theme options shared between light and dark modes
const getCommonTheme = (mode: PaletteMode) => ({
   typography: {
      fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
      h1: {
         fontSize: "2.5rem",
         fontWeight: 600,
      },
      h2: {
         fontSize: "2rem",
         fontWeight: 600,
      },
      h3: {
         fontSize: "1.75rem",
         fontWeight: 600,
      },
      h4: {
         fontSize: "1.5rem",
         fontWeight: 600,
      },
      h5: {
         fontSize: "1.25rem",
         fontWeight: 600,
      },
      h6: {
         fontSize: "1rem",
         fontWeight: 600,
      },
      body1: {
         fontSize: "1rem",
      },
   },
   shape: {
      borderRadius: 8,
   },
   components: {
      MuiButton: {
         styleOverrides: {
            root: {
               textTransform: "none",
               fontWeight: 500,
            },
         },
      },
      MuiCard: {
         styleOverrides: {
            root: {
               boxShadow:
                  mode === "light"
                     ? "0 2px 8px rgba(0,0,0,0.1)"
                     : "0 2px 8px rgba(0,0,0,0.3)",
            },
         },
      },
   },
});

// Light theme configuration
export const lightTheme = createTheme({
   palette: {
      mode: "light",
      primary: {
         main: "#1976d2",
         light: "#42a5f5",
         dark: "#1565c0",
      },
      secondary: {
         main: "#dc004e",
         light: "#f73378",
         dark: "#9a0036",
      },
      background: {
         default: "#fafafa",
         paper: "#ffffff",
      },
      text: {
         primary: "#000000",
         secondary: "#666666",
      },
   },
   ...getCommonTheme("light"),
});

// Dark theme configuration
export const darkTheme = createTheme({
   palette: {
      mode: "dark",
      primary: {
         main: "#90caf9",
         light: "#b3d9fc",
         dark: "#648dae",
      },
      secondary: {
         main: "#f48fb1",
         light: "#f6a5c1",
         dark: "#aa647b",
      },
      background: {
         default: "#121212",
         paper: "#1e1e1e",
      },
      text: {
         primary: "#ffffff",
         secondary: "#b0b0b0",
      },
   },
   ...getCommonTheme("dark"),
});

// Function to get theme based on mode
export const getTheme = (mode: PaletteMode) => {
   return mode === "light" ? lightTheme : darkTheme;
};

// Default export for backward compatibility
export default lightTheme;
