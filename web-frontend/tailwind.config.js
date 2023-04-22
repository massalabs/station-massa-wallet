/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      backgroundImage: {
        'landing-page': "url('./src/assets/bg-image-landing-page.svg')",
      },
    },
    colors: {
      'bg-primary': '#151A26',
      green: '#1AE19D',
    },
    letterSpacing: {
      massa: '-0.01em',
    },
  },
  plugins: [],
};
