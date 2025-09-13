import { createBrowserRouter } from "react-router-dom";
import { AppLayout } from "./pages/_layouts/app";

export const routes = createBrowserRouter([
  {
    path: "/",
    element: <AppLayout />,
  },
]);
