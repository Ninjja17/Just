/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        admin: {
          50: '#f5f7fa',
          100: '#e4e7eb',
          500: '#3f51b5',
          600: '#3949a1',
          700: '#2f3d8a',
        },
      },
    },
  },
  plugins: [],
}
