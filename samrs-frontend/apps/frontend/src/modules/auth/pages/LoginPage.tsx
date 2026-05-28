import { useLogin } from "../hooks/useLogin";
import { LoginForm } from "../layouts/LoginForm";

export default function LoginPage() {
  const { handleLogin, isLoggingIn } = useLogin();

  return <LoginForm onSubmit={handleLogin} isLoading={isLoggingIn} />;
};