import { createBrowserRouter } from "react-router-dom";
import { AppLayout } from "./pages/_layouts/app";
import { Nginx } from "./pages/app/nginx/nginx";

export const routes = createBrowserRouter([
  {
    path: "/",
    element: <AppLayout />,
    children: [
      {
        path: "/nginx",
        element: <Nginx />,
      },
    ],
  },
]);
