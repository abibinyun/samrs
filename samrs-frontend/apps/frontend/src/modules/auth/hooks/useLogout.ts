import { useDispatch } from "react-redux";
import { clearCredentials } from "../actions/slice";
import { api } from "@/store/api";
import { useNavigate } from "@tanstack/react-router";

export const useLogout = () => {
  const dispatch = useDispatch();
  const navigate = useNavigate();

  const handleLogout = () => {
    dispatch(clearCredentials());
    dispatch(api.util.resetApiState());
    navigate({ to: "/login" });
  };

  return { handleLogout };
};
