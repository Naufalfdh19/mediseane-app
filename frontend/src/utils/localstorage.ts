import { TOKEN_KEY } from "../const/const";

export function SetAccessToken(token: string) {
  localStorage.setItem(TOKEN_KEY, token)
}

export function RemoveAccessToken() {
  localStorage.removeItem(TOKEN_KEY);
}

export function GetAccessToken() {
  return localStorage.getItem(TOKEN_KEY);
}