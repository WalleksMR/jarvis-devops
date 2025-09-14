import "./index.css";
import { ThemeProvider } from "./components/theme/theme-provider";
import { RouterProvider } from "react-router-dom";
import { routes } from "./routes";
import { QueryClientProvider } from "@tanstack/react-query";
import { queryClient } from "./lib/react-query";

export function App() {
  return (
    <ThemeProvider defaultTheme="light" storageKey="@jarvis-devops-theme">
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={routes} />
      </QueryClientProvider>
    </ThemeProvider>
  );
}
