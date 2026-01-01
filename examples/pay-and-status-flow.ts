/**
 * High-level example of how you might structure a "pay + status" flow
 * with Base Account SDK.
 *
 * This file is intentionally framework-agnostic and does not import
 * any concrete SDK functions. It is a skeleton you can adapt to your app.
 *
 * Steps this example describes:
 * 1. Initialize an SDK client.
 * 2. Create a payment with metadata.
 * 3. Poll or fetch payment status.
 */

// Replace this type with the real client type from the SDK.
type BaseAccountClient = {
  // TODO: replace with real method from the SDK.
  // pay(options: { amount: string; to: string; memo?: string }): Promise<{ id: string }>;
  // getPaymentStatus(id: string): Promise<"pending" | "succeeded" | "failed">;
};

async function createClient(): Promise<BaseAccountClient> {
  /**
   * TODO:
   *  - Import the real factory from the SDK (for example `createAccountClient`)
   *  - Pass in network, apiKey, wallet, etc.
   *
   * Example (pseudo-code):
   *
   * import { createAccountClient } from "@base-org/account-sdk";
   *
   * return createAccountClient({
   *   network: "base-sepolia",
   *   apiKey: process.env.BASE_API_KEY!,
   *   // ...
   * });
   */
  throw new Error("createClient is a placeholder. Replace it with real SDK initialization.");
}

async function examplePayAndCheckStatus() {
  console.log("Starting pay + status example…");

  const client = await createClient();

  // 1. Create payment (pseudo-code)
  const payment = await (async () => {
    // TODO: replace with real `client.pay` or equivalent method.
    // return client.pay({
    //   amount: "1.00", // USD or token units depending on your setup
    //   to: "0xRecipientAddressHere",
    //   memo: "Order #1234",
    // });

    throw new Error("Replace this block with a real pay call.");
  })();

  console.log("Created payment:", payment);

  // 2. Poll payment status (pseudo-code)
  const paymentId = (payment as any).id;

  // Simple polling loop illustration.
  // In a real app use better retry / backoff strategy.
  for (let attempt = 0; attempt < 10; attempt++) {
    console.log(`Checking status (attempt ${attempt + 1})…`);

    // TODO: replace with real `client.getPaymentStatus(paymentId)` call.
    // const status = await client.getPaymentStatus(paymentId);
    const status = "pending"; // placeholder

    console.log(`Current status: ${status}`);

    if (status === "succeeded" || status === "failed") {
      console.log("Final status reached, stopping polling.");
      break;
    }

    await new Promise((resolve) => setTimeout(resolve, 3_000));
  }

  console.log("Pay + status example finished.");
}

// Allow running via `ts-node` or compiled JS.
if (require.main === module) {
  examplePayAndCheckStatus().catch((error) => {
    console.error("Error in pay + status example:", error);
    process.exitCode = 1;
  });
}
