"use client";

import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import { Brightness4, Brightness7 } from "@mui/icons-material";
import { useThemeMode } from "../providers/ThemeProvider";

export default function ThemeToggle() {
   const { mode, toggleTheme } = useThemeMode();

   return (
      <Tooltip title={`Switch to ${mode === "light" ? "dark" : "light"} mode`}>
         <IconButton
            onClick={toggleTheme}
            color="inherit"
            aria-label="toggle theme"
         >
            {mode === "dark" ? <Brightness7 /> : <Brightness4 />}
         </IconButton>
      </Tooltip>
   );
}
