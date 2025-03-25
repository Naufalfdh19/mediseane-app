export interface TokenPayload {
  type: string;
  role: string;
  iss: string;
  sub: string;
  exp: number;
  iat: number;
}
