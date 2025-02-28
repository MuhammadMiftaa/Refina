import { VitePWA } from "vite-plugin-pwa";
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import path from "path";

// https://vitejs.dev/config/
export default defineConfig({
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  plugins: [
    react(),
    tailwindcss(),
    VitePWA({
      devOptions: {
        enabled: true,
      },
      strategies: "injectManifest",
      srcDir: "src",
      filename: "sw.ts",
      registerType: "autoUpdate",
      injectManifest: {
        swDest: "dist/sw.js",
      },
      manifest: {
        name: "Refina",
        short_name: "Refina",
        description: "Rekapan Finansialmu",
        icons: [
          {
            src: "pwa-64x64.png",
            sizes: "64x64",
            type: "image/png",
          },
          {
            src: "pwa-192x192.png",
            sizes: "192x192",
            type: "image/png",
          },
          {
            src: "pwa-512x512.png",
            sizes: "512x512",
            type: "image/png",
          },
          {
            src: "maskable-icon-512x512.png",
            sizes: "512x512",
            type: "image/png",
            purpose: "maskable",
          },
        ],
        theme_color: "#E6E6FA",
        background_color: "#E6E6FA",
        start_url: "/",
        display: "standalone",
        orientation: "portrait",
        // gcm_sender_id: "131099091498"
      },
      workbox: {
        globPatterns: ["**/*.{js,css,html,ico,png,svg}"],
      },
    }),
  ],
});
