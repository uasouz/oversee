import type { Config } from "tailwindcss";

const config: Config = {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: "#D32F2F", // Security/Alert red
          light: "#FF6659",
          dark: "#9A0007",
        },
        background: "#121212",     // Light gray background
        card: "#FFFFFF",           // White card/container
        text: {
          DEFAULT: "#333333",      // Primary text color
          subtle: "#666666",       // Secondary text
        },
        border: "#DDDDDD",         // Border for cards and tables
      },
    }
    ,
  },
  plugins: [],
};

export default config;
