import { jwtDecode } from "jwt-decode";

export function getUsernameFromToken(token: string | null): string {
  if (!token) return "";
  try {
    const decoded: any = jwtDecode(token);
    return decoded.data.username;
  } catch {
    return "";
  }
}