interface ILoginRequest {
  username: string;
  password: string;
}

interface ILoginResponse {
  success: boolean;
  message: string;
  data: { token: string };
}

interface ILoginForm {
  username: string;
  password: string;
}

interface IMeResponse {
  success: boolean;
  message: string;
  data: {
    user: {
      id: string;
      username: string;
      tenant_id: string;
      role_id: number;
      is_system?: boolean; // Super Admin flag
      tenant: { id: string; slug: string; name: string };
      role: { 
        id: number; 
        name: string; 
        is_tenant_admin: boolean;
        is_system: boolean; // Super Admin flag from role
      };
    };
    permissions: string[];
  };
}

export type { ILoginRequest, ILoginResponse, IMeResponse, ILoginForm };