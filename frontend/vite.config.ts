import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

const BACKEND = "http://127.0.0.1:8080";

export default defineConfig({
  plugins: [
    react(),
    {
      name: "redirect-proxy",
      configureServer(server) {
        server.middlewares.use(async (req, res, next) => {
          // Skip API (handled by proxy ниже), Vite internal, и не-GET
          if (
            req.url!.startsWith("/api") ||
            req.url!.startsWith("/@") ||
            req.url!.startsWith("/node_modules") ||
            req.method !== "GET"
          ) {
            return next();
          }

          try {
            const resp = await fetch(`${BACKEND}${req.url}`);
            if (resp.status !== 404) {
              for (const [key, val] of resp.headers) {
                // nginx/strip заголовки кусочные
                if (key.toLowerCase() === "transfer-encoding") continue;
                res.setHeader(key, val);
              }
              res.writeHead(resp.status);
              const body = await resp.text();
              res.end(body);
              return;
            }
          } catch {
            // backend недоступен — отдать SPA
          }

          next();
        });
      },
    },
  ],
  server: {
    port: 5173,
    proxy: {
      "/api": {
        target: BACKEND,
        changeOrigin: true,
      },
    },
  },
});
