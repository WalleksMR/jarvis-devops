import { api } from "@/lib/axios";

interface GetStatusResponse {
  is_installed: boolean;
  is_running: boolean;
  version: string;
  config_valid: boolean;
  config_error?: string;
  last_reload: Date;
}

export async function getStatus() {
  const response = await api.get<GetStatusResponse>("/status");
  return response.data;
}
