import { defineConfig, devices } from "@playwright/test";

export default defineConfig({
  testDir: "./tests",
  fullyParallel: false,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : 1,
  reporter: "dot",
  use: {
    baseURL: "http://127.0.0.1:3000",
    trace: "on-first-retry",
  },

  projects: [
    {
      name: "chromium",
      use: { ...devices["Desktop Chrome"] },
    },
  ],

  webServer: {
    env: {
      PORT: "3000",
      GO_ENV: "test",
      JWT_SECRET: "pashinahui",
      BLOBS_DIR: ".",
      DB_NAME: "shsdb_test",
      DB_HOST: "localhost:3306",
      DB_USERNAME: "root",
      DB_PASSWORD: "previetcomrade",
      CACHE_HOST: "localhost:6379",
      CACHE_PASSWORD: "previetcomrade",
      SUPERADMIN_USERNAME: "b",
      SUPERADMIN_PASSWORD: "kurwamatch",
    },
    command: "make shs-server",
    reuseExistingServer: true, //!process.env.CI,
    url: "http://127.0.0.1:3000/v1/status",
  },
});
