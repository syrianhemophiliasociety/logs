import { test, expect } from "@playwright/test";
import {
  Account,
  createAccountAsSuperAdmin,
  loginAccount,
  resetCache,
  resetDB,
  seedAccounts,
} from "./factory";
import accounts from "./accounts.json";
import { randomUUID } from "crypto";

test.beforeAll(async ({ request }) => {
  await resetDB(request);
  await resetCache(request);
  await seedAccounts(request);
});

test.describe("Create Account", () => {
  test("creating an admin account with insufficient permissions fails", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.ziemowit);

    const resp = await request.post("/v1/accounts/admin", {
      headers: {
        Authorization: token,
      },
      data: JSON.stringify({
        new_account: {
          display_name: "Jack",
          username: "jackie420",
          password: "topsecret",
        },
      }),
    });

    expect(resp.status()).toBe(403);

    const respBody = await resp.json();
    expect(respBody.error_id).toBe("permission-denied");
  });

  test("creating a secritary account with insufficient permissions fails", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.ziemowit);

    const resp = await request.post("/v1/accounts/secritary", {
      headers: {
        Authorization: token,
      },
      data: JSON.stringify({
        new_account: {
          display_name: "Jack",
          username: "jackie420",
          password: "topsecret",
        },
      }),
    });

    expect(resp.status()).toBe(403);

    const respBody = await resp.json();
    expect(respBody.error_id).toBe("permission-denied");
  });

  test("creating an admin account with empty username fails", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.rex);

    const resp = await request.post("/v1/accounts/admin", {
      headers: {
        Authorization: token,
      },
      data: JSON.stringify({
        new_account: {
          display_name: "Jack",
          username: "",
          password: "topsecret",
        },
      }),
    });

    expect(resp.status()).toBe(400);

    const respBody = await resp.json();
    expect(respBody.error_id).toBe("invalid-account-username");
  });

  test("creating a secritary account with empty username fails", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.rex);

    const resp = await request.post("/v1/accounts/secritary", {
      headers: {
        Authorization: token,
      },
      data: JSON.stringify({
        new_account: {
          display_name: "Jack",
          username: "",
          password: "topsecret",
        },
      }),
    });

    expect(resp.status()).toBe(400);

    const respBody = await resp.json();
    expect(respBody.error_id).toBe("invalid-account-username");
  });

  test("creating an admin account with empty password fails", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.rex);

    const resp = await request.post("/v1/accounts/admin", {
      headers: {
        Authorization: token,
      },
      data: JSON.stringify({
        new_account: {
          display_name: "Jack",
          username: "jackie420",
          password: "",
        },
      }),
    });

    expect(resp.status()).toBe(400);

    const respBody = await resp.json();
    expect(respBody.error_id).toBe("invalid-account-password");
  });

  test("creating a secritary account with empty password fails", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.rex);

    const resp = await request.post("/v1/accounts/secritary", {
      headers: {
        Authorization: token,
      },
      data: JSON.stringify({
        new_account: {
          display_name: "Jack",
          username: "jackie420",
          password: "",
        },
      }),
    });

    expect(resp.status()).toBe(400);

    const respBody = await resp.json();
    expect(respBody.error_id).toBe("invalid-account-password");
  });

  test("creating an admin account with empty display name fails", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.rex);

    const resp = await request.post("/v1/accounts/admin", {
      headers: {
        Authorization: token,
      },
      data: JSON.stringify({
        new_account: {
          display_name: "",
          username: "jackie420",
          password: "topsecret",
        },
      }),
    });

    expect(resp.status()).toBe(400);

    const respBody = await resp.json();
    expect(respBody.error_id).toBe("invalid-account-display-name");
  });

  test("creating a secritary account with empty display name fails", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.rex);

    const resp = await request.post("/v1/accounts/secritary", {
      headers: {
        Authorization: token,
      },
      data: JSON.stringify({
        new_account: {
          display_name: "",
          username: "jackie420",
          password: "topsecret",
        },
      }),
    });

    expect(resp.status()).toBe(400);

    const respBody = await resp.json();
    expect(respBody.error_id).toBe("invalid-account-display-name");
  });

  test("creating an admin account with an existing username fails", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.rex);

    const resp = await request.post("/v1/accounts/admin", {
      headers: {
        Authorization: token,
      },
      data: JSON.stringify({
        new_account: accounts.harvey,
      }),
    });

    expect(resp.status()).toBe(409);

    const respBody = await resp.json();
    expect(respBody.error_id).toBe("account-exists");
  });

  test("creating a secritary account with an existing username fails", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.rex);

    const resp = await request.post("/v1/accounts/secritary", {
      headers: {
        Authorization: token,
      },
      data: JSON.stringify({
        new_account: accounts.cody,
      }),
    });

    expect(resp.status()).toBe(409);

    const respBody = await resp.json();
    expect(respBody.error_id).toBe("account-exists");
  });

  test("creating an admin account with valid parameters and permissions works", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.rex);

    const resp = await request.post("/v1/accounts/admin", {
      headers: {
        Authorization: token,
      },
      data: JSON.stringify({
        new_account: {
          display_name: "Jack",
          username: "jackie420",
          password: "topsecret",
        },
      }),
    });

    expect(resp).toBeOK();
  });

  test("creating a secritary account with valid parameters and permissions works", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.rex);

    const resp = await request.post("/v1/accounts/secritary", {
      headers: {
        Authorization: token,
      },
      data: JSON.stringify({
        new_account: {
          display_name: "Dank",
          username: "dankie444",
          password: "topsecret",
        },
      }),
    });

    expect(resp).toBeOK();
  });
});

test.describe("Get Account", () => {
  test("getting account with invalid token fails", async ({ request }) => {
    const token = "abc123";
    const accountId = await createAccountAsSuperAdmin(request, {
      username: "foo" + randomUUID(),
      password: "bar",
      display_name: "Baz",
      type: "secritary",
    } as Account);

    const resp = await request.get(`/v1/accounts/${accountId}`, {
      headers: {
        Authorization: token,
      },
    });

    expect(resp.status()).toBe(401);
  });

  test("getting account with empty token fails", async ({ request }) => {
    const token = "";
    const accountId = await createAccountAsSuperAdmin(request, {
      username: "foo" + randomUUID(),
      password: "bar",
      display_name: "Baz",
      type: "secritary",
    } as Account);

    const resp = await request.get(`/v1/accounts/${accountId}`, {
      headers: {
        Authorization: token,
      },
    });

    expect(resp.status()).toBe(401);
  });

  test("getting account with insufficient permissions fails", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.ziemowit);
    const accountId = await createAccountAsSuperAdmin(request, {
      username: "foo" + randomUUID(),
      password: "bar",
      display_name: "Baz",
      type: "secritary",
    } as Account);

    const resp = await request.get(`/v1/accounts/${accountId}`, {
      headers: {
        Authorization: token,
      },
    });

    expect(resp.status()).toBe(403);
  });

  test("getting account with sufficient permissions works", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.rex);
    const accountId = await createAccountAsSuperAdmin(request, {
      username: "foo" + randomUUID(),
      password: "bar",
      display_name: "Baz",
      type: "secritary",
    } as Account);

    const resp = await request.get(`/v1/accounts/${accountId}`, {
      headers: {
        Authorization: token,
      },
    });

    expect(resp.status()).toBe(200);
  });
});

test.describe("List Accounts", () => {
  test("listing accounts with invalid token fails", async ({ request }) => {
    const token = "abc123";
    const resp = await request.get("/v1/accounts", {
      headers: {
        Authorization: token,
      },
    });

    expect(resp.status()).toBe(401);
  });

  test("listing accounts with empty token fails", async ({ request }) => {
    const token = "";
    const resp = await request.get("/v1/accounts", {
      headers: {
        Authorization: token,
      },
    });

    expect(resp.status()).toBe(401);
  });

  test("listing account with insufficient permissions fails", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.ziemowit);

    const resp = await request.get("/v1/accounts", {
      headers: {
        Authorization: token,
      },
    });

    expect(resp.status()).toBe(403);
  });

  test("listing accounts with sufficient permissions returns a list of admins and secritaries", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.rex);
    const allowedTypes = ["admin", "secritary"];
    // using the existing list from ./accounts.json
    const resp = await request.get("/v1/accounts", {
      headers: {
        Authorization: token,
      },
    });

    expect(resp.status()).toBe(200);
    const respBody = await resp.json();
    for (const account of respBody.data) {
      expect(allowedTypes.includes(account.type)).toBe(true);
    }
  });
});

test.describe("Update Account", () => {
  test("updating an account with invalid token fails", async ({ request }) => {
    const token = "abc123";
    const accountId = await createAccountAsSuperAdmin(request, {
      username: "foo" + randomUUID(),
      password: "bar",
      display_name: "Baz",
      type: "secritary",
    } as Account);

    const resp = await request.put(`/v1/accounts/${accountId}`, {
      headers: {
        Authorization: token,
      },
      data: JSON.stringify({
        new_account: {
          display_name: "Foo",
        },
      }),
    });

    expect(resp.status()).toBe(401);
  });

  test("updating an account with empty token fails", async ({ request }) => {
    const token = "";
    const accountId = await createAccountAsSuperAdmin(request, {
      username: "foo" + randomUUID(),
      password: "bar",
      display_name: "Baz",
      type: "secritary",
    } as Account);

    const resp = await request.put(`/v1/accounts/${accountId}`, {
      headers: {
        Authorization: token,
      },
      data: JSON.stringify({
        new_account: {
          display_name: "Foo",
        },
      }),
    });

    expect(resp.status()).toBe(401);
  });

  test("updating an account with insufficient permissions fails", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.ziemowit);
    const accountId = await createAccountAsSuperAdmin(request, {
      username: "foo" + randomUUID(),
      password: "bar",
      display_name: "Baz",
      type: "secritary",
    } as Account);

    const resp = await request.put(`/v1/accounts/${accountId}`, {
      headers: {
        Authorization: token,
      },
      data: JSON.stringify({
        new_account: {
          display_name: "Foo",
        },
      }),
    });

    expect(resp.status()).toBe(403);
  });

  test("updating an account's username with sufficient permissions works", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.rex);
    const accountId = await createAccountAsSuperAdmin(request, {
      username: "foo" + randomUUID(),
      password: "bar",
      display_name: "Baz",
      type: "secritary",
    } as Account);

    const newAccountUsername = "foofi" + randomUUID();
    const resp = await request.put(`/v1/accounts/${accountId}`, {
      headers: {
        Authorization: token,
      },
      data: JSON.stringify({
        new_account: {
          username: newAccountUsername,
        },
      }),
    });
    expect(resp.status()).toBe(200);

    const resp2 = await request.get(`/v1/accounts/${accountId}`, {
      headers: {
        Authorization: token,
      },
    });
    expect(resp2).toBeOK();
    const respBody = await resp2.json();

    expect(respBody.data.username).toBe(newAccountUsername);
  });

  test("updating an account's password with sufficient permissions works", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.rex);
    const accountUsername = "foo" + randomUUID();
    const accountId = await createAccountAsSuperAdmin(request, {
      username: accountUsername,
      password: "bar",
      display_name: "Baz",
      type: "secritary",
    } as Account);

    const newAccountPassword = "bibi";
    const resp = await request.put(`/v1/accounts/${accountId}`, {
      headers: {
        Authorization: token,
      },
      data: JSON.stringify({
        new_account: {
          password: newAccountPassword,
        },
      }),
    });
    expect(resp.status()).toBe(200);

    await loginAccount(request, {
      username: accountUsername,
      password: newAccountPassword,
    });
  });

  test("updating an account's display name with sufficient permissions works", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.rex);
    const accountId = await createAccountAsSuperAdmin(request, {
      username: "foo" + randomUUID(),
      password: "bar",
      display_name: "Baz",
      type: "secritary",
    } as Account);

    const newAccountDisplayName = "Foo";
    const resp = await request.put(`/v1/accounts/${accountId}`, {
      headers: {
        Authorization: token,
      },
      data: JSON.stringify({
        new_account: {
          display_name: newAccountDisplayName,
        },
      }),
    });
    expect(resp.status()).toBe(200);

    const resp2 = await request.get(`/v1/accounts/${accountId}`, {
      headers: {
        Authorization: token,
      },
    });
    expect(resp2).toBeOK();
    const respBody = await resp2.json();

    expect(respBody.data.display_name).toBe(newAccountDisplayName);
  });
});

test.describe("Delete Account", () => {
  test("deleting an account with invalid token fails", async ({ request }) => {
    const token = "abc123";
    const accountId = await createAccountAsSuperAdmin(request, {
      username: "foo" + randomUUID(),
      password: "bar",
      display_name: "Baz",
      type: "secritary",
    } as Account);

    const resp = await request.delete(`/v1/accounts/${accountId}`, {
      headers: {
        Authorization: token,
      },
    });

    expect(resp.status()).toBe(401);
  });

  test("deleting an account with empty token fails", async ({ request }) => {
    const token = "";
    const accountId = await createAccountAsSuperAdmin(request, {
      username: "foo" + randomUUID(),
      password: "bar",
      display_name: "Baz",
      type: "secritary",
    } as Account);

    const resp = await request.delete(`/v1/accounts/${accountId}`, {
      headers: {
        Authorization: token,
      },
    });

    expect(resp.status()).toBe(401);
  });

  test("deleting an account with insufficient permissions fails", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.ziemowit);
    const accountId = await createAccountAsSuperAdmin(request, {
      username: "foo" + randomUUID(),
      password: "bar",
      display_name: "Baz",
      type: "secritary",
    } as Account);

    const resp = await request.delete(`/v1/accounts/${accountId}`, {
      headers: {
        Authorization: token,
      },
    });

    expect(resp.status()).toBe(403);
  });

  test("deleting an account with sufficient permissions works", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.rex);
    const accountId = await createAccountAsSuperAdmin(request, {
      username: "foo" + randomUUID(),
      password: "bar",
      display_name: "Baz",
      type: "secritary",
    } as Account);

    const resp = await request.delete(`/v1/accounts/${accountId}`, {
      headers: {
        Authorization: token,
      },
    });

    expect(resp.status()).toBe(200);
  });
});
