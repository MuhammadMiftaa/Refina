/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      fontFamily: {
        inter: ["Inter"],
        poppins: ["Poppins"],
      },
      colors: {
        primary: "#060610",
        secondary: "#1d1e22",
      },
    },
  },
  plugins: [],
};
