import "./index.css";
import { ThemeProvider } from "./components/theme/theme-provider";
import { RouterProvider } from "react-router-dom";
import { routes } from "./routes";

export function App() {
  return (
    <ThemeProvider defaultTheme="dark" storageKey="@jarvis-devops-theme">
      <RouterProvider router={routes} />
    </ThemeProvider>
  );
}
