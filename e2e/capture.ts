import { chromium, type Page } from "@playwright/test";
import { execSync, spawn, type ChildProcess } from "child_process";
import { mkdirSync, rmSync, existsSync } from "fs";
import { join, dirname } from "path";
import { fileURLToPath } from "url";

const __dirname = dirname(fileURLToPath(import.meta.url));
const PORT = 9222;
const BASE = `http://localhost:${PORT}`;
const DB_PATH = "/tmp/poolvibes-capture.db";
const BIN_PATH = join(__dirname, "bin", "poolvibes");
const SCREENSHOTS_DIR = join(__dirname, "screenshots");
const PROJECT_ROOT = join(__dirname, "..");

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

async function waitForServer(url: string, timeoutMs = 15000): Promise<void> {
  const start = Date.now();
  while (Date.now() - start < timeoutMs) {
    try {
      const res = await fetch(url);
      if (res.ok) return;
    } catch {}
    await new Promise((r) => setTimeout(r, 200));
  }
  throw new Error(`Server at ${url} not ready after ${timeoutMs}ms`);
}

async function screenshot(page: Page, name: string): Promise<void> {
  await page.screenshot({ path: join(SCREENSHOTS_DIR, name), fullPage: true });
  console.log(`  📸 ${name}`);
}

async function seedViaFetch(
  page: Page,
  path: string,
  body: object,
): Promise<void> {
  const resp = await page.evaluate(
    async ({ url, data }) => {
      const r = await fetch(url, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      });
      return { status: r.status, ok: r.ok };
    },
    { url: `${BASE}${path}`, data: body },
  );
  if (!resp.ok) {
    console.warn(`  ⚠️  POST ${path} returned ${resp.status}`);
  }
}

async function clickTab(page: Page, label: string): Promise<void> {
  await page.locator(".tabs a", { hasText: label }).click();
  await page.waitForTimeout(800);
}

async function dismissModal(page: Page): Promise<void> {
  // Force-dismiss any active modal by resetting the #modal container
  await page.evaluate(() => {
    const modal = document.querySelector("#modal");
    if (modal) {
      modal.className = "";
      modal.innerHTML = "";
    }
  });
  await page.waitForTimeout(200);
}

async function captureModal(
  page: Page,
  tabLabel: string,
  addButtonText: string,
  screenshotName: string,
): Promise<void> {
  await clickTab(page, tabLabel);
  try {
    await page.locator(`button:has-text("${addButtonText}")`).click();
    await page.waitForSelector(".modal.is-active", { timeout: 3000 });
    await page.waitForTimeout(300);
    await screenshot(page, screenshotName);
  } catch (e) {
    console.warn(`  ⚠️  Could not capture modal: ${screenshotName}`);
  }
  await dismissModal(page);
}

// ---------------------------------------------------------------------------
// Main
// ---------------------------------------------------------------------------

async function main() {
  mkdirSync(SCREENSHOTS_DIR, { recursive: true });
  mkdirSync(join(__dirname, "bin"), { recursive: true });

  if (existsSync(DB_PATH)) rmSync(DB_PATH);

  // 1. Build
  console.log("🔨 Building Go binary...");
  execSync(`go build -o ${BIN_PATH} .`, {
    cwd: PROJECT_ROOT,
    stdio: "inherit",
  });

  // 2. Start server
  console.log("🚀 Starting server...");
  const server: ChildProcess = spawn(
    BIN_PATH,
    ["serve", "--addr", `:${PORT}`, "--db", DB_PATH, "--db-driver", "sqlite"],
    { stdio: "pipe" },
  );
  server.stderr?.on("data", (d: Buffer) => process.stderr.write(d));
  server.stdout?.on("data", (d: Buffer) => process.stdout.write(d));

  try {
    await waitForServer(`${BASE}/login`);
    console.log("✅ Server ready\n");

    // 3. Launch browser
    const browser = await chromium.launch();
    const context = await browser.newContext({
      viewport: { width: 1280, height: 720 },
    });
    const page = await context.newPage();

    // -----------------------------------------------------------------------
    // Auth pages
    // -----------------------------------------------------------------------
    console.log("📋 Capturing auth pages...");
    await page.goto(`${BASE}/login`);
    await page.waitForTimeout(300);
    await screenshot(page, "login.png");

    await page.goto(`${BASE}/signup`);
    await page.waitForTimeout(300);
    await screenshot(page, "signup.png");

    // -----------------------------------------------------------------------
    // Sign up
    // -----------------------------------------------------------------------
    console.log("\n🔐 Signing up test user...");
    await page.fill('input[name="email"]', "test@poolvibes.app");
    await page.fill('input[name="password"]', "password123");
    await page.fill('input[name="confirm"]', "password123");
    await page.click('button[type="submit"]');
    await page.waitForURL(`${BASE}/`);
    // Wait for Datastar to init and load the default tab
    await page.waitForTimeout(1500);
    console.log("✅ Signed up and redirected to home\n");

    // -----------------------------------------------------------------------
    // Empty states
    // -----------------------------------------------------------------------
    console.log("📋 Capturing empty states...");
    // Chemistry is the default tab
    await screenshot(page, "chemistry-empty.png");

    await clickTab(page, "Tasks");
    await screenshot(page, "tasks-empty.png");

    await clickTab(page, "Equipment");
    await screenshot(page, "equipment-empty.png");

    await clickTab(page, "Chemicals");
    await screenshot(page, "chemicals-empty.png");

    await clickTab(page, "Settings");
    await screenshot(page, "settings.png");

    // -----------------------------------------------------------------------
    // Seed data
    // -----------------------------------------------------------------------
    console.log("\n🌱 Seeding data...");

    // Chemistry logs
    for (const log of [
      {
        ph: 7.4,
        freeChlorine: 2.5,
        combinedChlorine: 0.1,
        totalAlkalinity: 100,
        cya: 40,
        calciumHardness: 300,
        temperature: 82,
        notes: "All readings in range",
        testedAt: "2026-02-14T10:00",
      },
      {
        ph: 8.2,
        freeChlorine: 0.5,
        combinedChlorine: 0.8,
        totalAlkalinity: 60,
        cya: 90,
        calciumHardness: 150,
        temperature: 78,
        notes: "pH high, chlorine low — needs attention",
        testedAt: "2026-02-13T09:30",
      },
      {
        ph: 7.2,
        freeChlorine: 3.0,
        combinedChlorine: 0.2,
        totalAlkalinity: 110,
        cya: 50,
        calciumHardness: 280,
        temperature: 80,
        notes: "Looking good after treatment",
        testedAt: "2026-02-12T11:15",
      },
    ]) {
      await seedViaFetch(page, "/chemistry", log);
    }
    console.log("  ✅ 3 chemistry logs");

    // Tasks
    for (const task of [
      {
        taskName: "Backwash filter",
        taskDescription: "Backwash the sand filter for 3 minutes",
        recurrenceFrequency: "weekly",
        recurrenceInterval: 1,
        dueDate: "2026-02-10",
      },
      {
        taskName: "Check pump pressure",
        taskDescription: "Verify pump PSI is within 10-15 range",
        recurrenceFrequency: "daily",
        recurrenceInterval: 1,
        dueDate: "2026-02-14",
      },
      {
        taskName: "Clean skimmer baskets",
        taskDescription: "Remove debris from all skimmer baskets",
        recurrenceFrequency: "weekly",
        recurrenceInterval: 1,
        dueDate: "2026-02-21",
      },
      {
        taskName: "Inspect pool cover",
        taskDescription: "Check for tears or damage on the winter cover",
        recurrenceFrequency: "none",
        recurrenceInterval: 0,
        dueDate: "2026-03-15",
      },
    ]) {
      await seedViaFetch(page, "/tasks", task);
    }
    console.log("  ✅ 4 tasks");

    // Equipment
    for (const eq of [
      {
        eqName: "Hayward Super Pump",
        eqCategory: "pump",
        eqManufacturer: "Hayward",
        eqModel: "SP2607X10",
        eqSerialNumber: "HW-2024-001",
        eqInstallDate: "2024-03-15",
        eqWarrantyExpiry: "2027-03-15",
      },
      {
        eqName: "Pentair Clean & Clear Filter",
        eqCategory: "filter",
        eqManufacturer: "Pentair",
        eqModel: "CC150",
        eqSerialNumber: "PT-2023-042",
        eqInstallDate: "2023-06-01",
        eqWarrantyExpiry: "2026-06-01",
      },
    ]) {
      await seedViaFetch(page, "/equipment", eq);
    }
    console.log("  ✅ 2 equipment items");

    // Chemicals
    for (const chem of [
      {
        chemName: "Liquid Chlorine",
        chemType: "sanitizer",
        chemStockAmount: 10,
        chemStockUnit: "gallons",
        chemAlertThreshold: 3,
      },
      {
        chemName: "Muriatic Acid",
        chemType: "pH adjuster",
        chemStockAmount: 1.5,
        chemStockUnit: "gallons",
        chemAlertThreshold: 2,
      },
      {
        chemName: "Stabilizer (CYA)",
        chemType: "stabilizer",
        chemStockAmount: 0,
        chemStockUnit: "lbs",
        chemAlertThreshold: 5,
      },
    ]) {
      await seedViaFetch(page, "/chemicals", chem);
    }
    console.log("  ✅ 3 chemicals");

    // -----------------------------------------------------------------------
    // Populated states
    // -----------------------------------------------------------------------
    console.log("\n📋 Capturing populated states...");

    await clickTab(page, "Water Chemistry");
    await screenshot(page, "chemistry-data.png");

    await clickTab(page, "Tasks");
    await screenshot(page, "tasks-data.png");

    // Complete the first task (overdue "Backwash filter") by clicking its
    // complete button — the first unchecked circular button in the list
    try {
      await page
        .locator('button[title="Mark complete"]')
        .first()
        .click();
      await page.waitForTimeout(800);
      // Show completed section
      await page.locator('button:has-text("Show completed")').click();
      await page.waitForTimeout(500);
      await screenshot(page, "tasks-completed.png");
    } catch {
      console.warn("  ⚠️  Could not complete a task");
    }

    await clickTab(page, "Equipment");
    await screenshot(page, "equipment-data.png");

    await clickTab(page, "Chemicals");
    await screenshot(page, "chemicals-data.png");

    // -----------------------------------------------------------------------
    // Modals
    // -----------------------------------------------------------------------
    console.log("\n📋 Capturing modals...");
    await captureModal(
      page,
      "Water Chemistry",
      "+ Add Test",
      "chemistry-modal.png",
    );
    await captureModal(page, "Tasks", "+ Add Task", "tasks-modal.png");
    await captureModal(
      page,
      "Equipment",
      "+ Add Equipment",
      "equipment-modal.png",
    );
    await captureModal(
      page,
      "Chemicals",
      "+ Add Chemical",
      "chemicals-modal.png",
    );

    // -----------------------------------------------------------------------
    // Admin page
    // -----------------------------------------------------------------------
    console.log("\n📋 Capturing admin page...");
    await clickTab(page, "Admin");
    await screenshot(page, "admin-users.png");

    await browser.close();
    console.log("\n✅ All screenshots captured!");
    console.log(`   📂 ${SCREENSHOTS_DIR}/`);
  } finally {
    server.kill("SIGTERM");
    // Give the process a moment to exit before cleaning up
    await new Promise((r) => setTimeout(r, 500));
    if (existsSync(DB_PATH)) rmSync(DB_PATH);
  }
}

main().catch((err) => {
  console.error("❌ Error:", err);
  process.exit(1);
});
