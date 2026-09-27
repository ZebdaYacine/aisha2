import type { NextConfig } from "next";
const publicMinioEndpoint = process.env.MINIO_PUBLIC_ENDPOINT ?? "localhost:9000";
const [publicMinioHost, publicMinioPort] = publicMinioEndpoint.split(":");
const nextConfig: NextConfig = {
  output: "standalone",
  reactStrictMode: true,
  images: {
    // Product media uses short-lived MinIO URLs that are reachable by the
    // browser, but not by the frontend container through localhost.
    unoptimized: true,
    remotePatterns: [
      {
        protocol: process.env.MINIO_USE_SSL === "true" ? "https" : "http",
        hostname: publicMinioHost,
        ...(publicMinioPort ? { port: publicMinioPort } : {}),
      },
    ],
  },
};
export default nextConfig;
