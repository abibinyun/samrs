import { RouterProvider } from "@tanstack/react-router";
import { router } from "./routes";
import { Provider } from "react-redux";
import { store } from "./store";
import { Toaster } from "sonner";
import "@/index.css";

export const App = () => (
  <Provider store={store}>
    {/* RouterProvider sekarang punya akses ke Redux Store */}
    <RouterProvider router={router} defaultViewTransition={true} />
    
    {/* Toaster global untuk notifikasi (Ganti ToastContainer lama) */}
    <Toaster position="top-right" richColors closeButton />
  </Provider>
);