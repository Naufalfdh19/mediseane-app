/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    colors: {
      white: "#FFFFFF",
      blue: "#0054a3",
      "primary-green": "#008081",
      "secondary-green": "#E7FFEF",
      "primary-white": "#FAFAFA",
      "primary-gray": "#545454",
      "secondary-gray": "#D4D4D4",
      "primary-black": "#232321",
      "primary-blue": "#5384B1",
      "primary-red": "#AE2E2E",
      "secondary-red": "#FFD3D3",
      "primary-orange": "#ee4d2d",
      "primary-purple": "#9D00FF",
    },
    fontFamily: {
      sans: ["Poppins"],
    },
    extend: {},
  },
  plugins: [],
};
