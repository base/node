export async function handleRequest(req: any, res: any): Promise<void> {
  try {
    // TODO: keep your existing request handling logic here
    // Example:
    // const result = await doSomething(req);
    // return res.status(200).json(result);

    return res.status(200).json({ ok: true });
  } catch (err) {
    const error = err as Error;

    // Extended error logging with request context (safe defaults)
    console.error("API request failed", {
      method: req?.method,
      url: req?.originalUrl ?? req?.url,
      params: req?.params,
      query: req?.query,
      body: req?.body,
      error: error?.message,
      stack: error?.stack,
    });

    res.status(500).json({ message: "Internal server error" });
  }
}
