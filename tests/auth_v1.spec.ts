import { test, expect } from "@playwright/test";
import { loginAccount, resetCache, resetDB, seedAccounts } from "./factory";
import accounts from "./accounts.json";

test.beforeAll(async ({ request }) => {
  await resetDB(request);
  await resetCache(request);
  await seedAccounts(request);
});

test.describe("Login", () => {
  test("when user provides incorrect login credentials login fails", async ({
    request,
  }) => {
    const resp = await request.post("/v1/login/username", {
      data: JSON.stringify({
        username: "foo",
        password: "bar",
      }),
    });

    expect(resp.status()).toBe(401);
  });

  test("when user provides correct login credentials login fails", async ({
    request,
  }) => {
    const resp = await request.post("/v1/login/username", {
      data: JSON.stringify({
        username: "ziemowit",
        password: "1234",
      }),
    });

    expect(resp).toBeOK();

    const respBody = await resp.json();
    expect(respBody.session_token).not.toBeUndefined();
  });
});

test.describe("Me", () => {
  test("calling /me/auth with invalid token responds with unauthorized", async ({
    request,
  }) => {
    const token = "abc123";

    const resp = await request.get("/v1/me/auth", {
      headers: {
        Authorization: token,
      },
    });

    expect(resp.status()).toBe(401);
  });

  test("calling /me/auth with empty token responds with unauthorized", async ({
    request,
  }) => {
    const token = "";

    const resp = await request.get("/v1/me/auth", {
      headers: {
        Authorization: token,
      },
    });

    expect(resp.status()).toBe(401);
  });

  test("calling /me/auth with a valid token responds with the account details", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.rex);

    const resp = await request.get("/v1/me/auth", {
      headers: {
        Authorization: token,
      },
    });

    expect(resp).toBeOK();
    const respBody = await resp.json();
    expect(respBody.display_name).toBe("Rex");
    expect(respBody.username).toBe("rex");
    expect(respBody.password).toBeUndefined();
  });

  test("calling /me/logout with invalid token responds with unauthorized", async ({
    request,
  }) => {
    const token = "abc123";

    const resp = await request.get("/v1/me/logout", {
      headers: {
        Authorization: token,
      },
    });

    expect(resp.status()).toBe(401);
  });

  test("calling /me/logout with empty token responds with unauthorized", async ({
    request,
  }) => {
    const token = "";

    const resp = await request.get("/v1/me/logout", {
      headers: {
        Authorization: token,
      },
    });

    expect(resp.status()).toBe(401);
  });
});
