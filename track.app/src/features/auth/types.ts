
export interface User {
  id: string;
  name: string;
  email: string;
  role: 'user' | 'admin';
  profilePicture: string;
}

export type AuthUser = User & {
  token: string;
}
