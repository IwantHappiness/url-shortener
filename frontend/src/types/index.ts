export interface User {
  id: number;
  version: number;
  nickname: string;
  email: string;
  created_at: string;
}

export interface CreateUserRequest {
  nickname: string;
  email: string;
}

export interface PatchUserRequest {
  nickname?: string | null;
  email?: string | null;
}

export interface Url {
  id: number;
  version: number;
  user_id: number;
  original_url: string;
  short_url: string;
  created_at: string;
}

export interface CreateUrlRequest {
  url: string;
  user_id: number;
}

export interface Stats {
  short_url: string;
  original_url: string;
  created_at: string;
  total_clicks: number;
  unique_ips: number;
  last_clicked_at: string | null;
}

export interface ApiError {
  error: string;
  message: string;
}
