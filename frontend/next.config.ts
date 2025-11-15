import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here */
  async rewrites() {
      return [
        { source: "/api/go/:path*", destination: "http://localhost:8000/api/go/:path*"}
      ];
  },
};

export default nextConfig;
