import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { loginSchema, type LoginFields } from "../schemas";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label"

interface Props {
  onSubmit: (data: LoginFields) => void;
  isLoading: boolean;
}

export const LoginForm = ({ onSubmit, isLoading }: Props) => {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFields>({
    resolver: zodResolver(loginSchema),
  });

  return (
    <div className="flex items-center justify-center h-screen bg-background p-4">
      <div className="w-full max-w-md p-8 bg-card rounded-xl shadow-lg border border-border">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-primary">SAMRS</h1>
          <p className="text-muted-foreground mt-2">Professional Hospital System</p>
        </div>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label>Username</Label>
            <Input {...register("username")} placeholder="Masukkan username" />
            {errors.username && (
              <p className="text-xs text-destructive">{errors.username.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label>Password</Label>
            <Input {...register("password")} type="password" placeholder="••••••" />
            {errors.password && (
              <p className="text-xs text-destructive">{errors.password.message}</p>
            )}
          </div>

          <Button type="submit" className="w-full" disabled={isLoading}>
            {isLoading ? "Memverifikasi..." : "Masuk ke Sistem"}
          </Button>
        </form>
      </div>
    </div>
  );
};
