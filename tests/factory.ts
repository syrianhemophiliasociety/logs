import { APIRequestContext, expect } from "@playwright/test";
import accounts from "./accounts.json";

export type AccountKey = keyof typeof accounts;

export interface Account {
  display_name: string;
  username: string;
  password: string;
  type: "superadmin" | "admin" | "secritary" | "patient";
  permissions: number;
}

export async function loginAccount(
  request: APIRequestContext,
  params: {
    username: string;
    password: string;
  },
): Promise<string> {
  const resp = await request.post("/v1/login/username", {
    data: JSON.stringify({
      username: params.username,
      password: params.password,
    }),
  });
  expect(resp).toBeOK();

  const respBody = await resp.json();

  return respBody.session_token;
}

export async function createAccountAsSuperAdmin(
  request: APIRequestContext,
  account: Account,
): Promise<number> {
  if (!["secritary", "admin"].includes(account.type)) {
    return Promise.reject(
      "invalid account type, can only create admin and secritary",
    );
  }

  const adminToken = await loginAccount(request, accounts.b);

  const resp = await request.post(`/v1/accounts/${account.type}`, {
    data: JSON.stringify({
      new_account: account,
    }),
    headers: {
      Authorization: adminToken,
    },
  });
  expect(resp).toBeOK();

  return (await resp.json()).id;
}

export async function resetDB(request: APIRequestContext) {
  const resp = await request.post("/v1/tests/reset/db");
  expect(resp).toBeOK();
}

export async function resetCache(request: APIRequestContext) {
  const resp = await request.post("/v1/tests/reset/cache");
  expect(resp).toBeOK();
}

export async function seedAccounts(request: APIRequestContext): Promise<void> {
  for (const account in accounts) {
    if (!["secritary", "admin"].includes(accounts[account].type)) {
      // skip patient and superadmin creation
      continue;
    }

    await createAccountAsSuperAdmin(request, accounts[account]);
  }
}

export function rfc3339NanoToUnixMicros(rfc3339NanoTime: string) {
  const date = new Date(rfc3339NanoTime.slice(0, -4) + "Z");
  const unixMillis = date.getTime();
  const nanoseconds = parseInt(rfc3339NanoTime.slice(-13, -1));
  return unixMillis * 1000 + Math.floor(nanoseconds / 1000);
}
