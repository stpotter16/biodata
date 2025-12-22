import { test, expect } from "@playwright/test";

test.describe("Login and entry management", () => {
  test("Login", async ({ page }) => {
    await page.goto("/");

    await expect(page).toHaveTitle(/Login/i);

    // Login
    const passphrase = process.env.PASSPHRASE;
    await page.getByLabel("Passphrase").fill(passphrase);
    await page.getByRole("button", { name: "Login" }).click();

    // Home page - no data
    await expect(page).toHaveTitle(/Biotrak/i);
    await expect(
      page.getByRole("heading", { name: "No entries yet. Add one!" }),
    ).toBeVisible();
    await page.getByRole("link", { name: "Add entry" }).click();

    // Add entry
    await expect(page).toHaveTitle(/Add Entry/i);
    await page.getByLabel("Date").fill("2025-12-22");
    await page.getByLabel("Weight").fill("185.5");
    await page.getByLabel("Waist").fill("33.2");
    await page.getByLabel("BP").fill("122/81");
    await page.getByRole("button", { name: "Add" }).click();

    // Home page - data
    await expect(page).toHaveTitle(/Biotrak/i);

    // Verify the table is visible (no longer showing "No entries yet")
    await expect(page.getByRole("table")).toBeVisible();

    // Verify all the entered data appears in the table
    await expect(page.getByRole("cell", { name: "2025-12-22" })).toBeVisible();
    await expect(page.getByRole("cell", { name: "185.5" })).toBeVisible();
    await expect(page.getByRole("cell", { name: "33.2" })).toBeVisible();
    await expect(page.getByRole("cell", { name: "122/81" })).toBeVisible();
    await expect(page.getByRole("cell", { name: "Edit" })).toBeVisible();

    // Edit entry
    await page.getByRole("link", { name: "Edit" }).click();

    await expect(page).toHaveTitle(/Edit Entry/i);
    await expect(
      page.getByRole("heading", { name: "2025-12-22" }),
    ).toBeVisible();
    await expect(page.getByLabel("Weight")).toHaveValue("185.5");
    await expect(page.getByLabel("Waist")).toHaveValue("33.2");
    await expect(page.getByLabel("BP")).toHaveValue("122/81");

    await page.getByLabel("Weight").fill("184.5");
    await page.getByRole("button", { name: "Update" }).click();

    // Home page - edited data
    await expect(page).toHaveTitle(/Biotrak/i);

    // Verify the table is visible (no longer showing "No entries yet")
    await expect(page.getByRole("table")).toBeVisible();

    // Verify all the entered data appears in the table
    await expect(page.getByRole("cell", { name: "2025-12-22" })).toBeVisible();
    await expect(page.getByRole("cell", { name: "184.5" })).toBeVisible();
    await expect(page.getByRole("cell", { name: "33.2" })).toBeVisible();
    await expect(page.getByRole("cell", { name: "122/81" })).toBeVisible();
    await expect(page.getByRole("cell", { name: "Edit" })).toBeVisible();
  });
});
