import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import App from "./App.tsx";
import { RouterProvider, createBrowserRouter } from "react-router-dom";
import { ConfirmationPage } from "./ConfirmationPage.tsx";

export const API_URL =
  import.meta.env.VITE_API_URL || "http://localhost:3030/v1";

const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,
  },
  {
    path: "/confirm/:token",
    element: <ConfirmationPage />,
  },
]);

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <RouterProvider router={router} />
  </StrictMode>
);
