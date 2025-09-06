/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // RealWorld brand colors
        'conduit-green': '#5CB85C',
        'conduit-green-dark': '#449D44',
        'conduit-gray': '#687077',
        'conduit-gray-light': '#ECEEEF',
        'conduit-red': '#B85C5C',
      },
      fontFamily: {
        'sans': ['Inter', 'system-ui', 'sans-serif'],
        'serif': ['Merriweather', 'Georgia', 'serif'],
      },
      container: {
        center: true,
        padding: '1rem',
        screens: {
          sm: '640px',
          md: '768px',
          lg: '1024px',
          xl: '1140px',
        },
      },
    },
  },
  plugins: [],
}