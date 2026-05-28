import { type ReactNode } from "react";
import { Navigate } from "@tanstack/react-router";

interface Props {
  children: ReactNode;
}

export function ProtectedRoute({ children }: Props) {
  const token = localStorage.getItem("samrs_token");

  // Jika tidak ada token, langsung arahkan ke login
  if (!token) {
    return <Navigate to="/login" replace />;
  }

  return <>{children}</>;
}