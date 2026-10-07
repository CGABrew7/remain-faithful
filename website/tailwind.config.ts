import type { Config } from 'tailwindcss'

const config: Config = {
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        navy: {
          DEFAULT: '#12214C',
          light: '#1A2D61',
          mid: '#162756',
          deep: '#0C1D4C',
          border: 'rgba(255,255,255,0.10)',
        },
        gold: {
          DEFAULT: '#D1AB4C',
          light: '#E2C46E',
        },
        cream: {
          DEFAULT: '#FFFFFF',
          muted: '#B4BCD4',
        },
        paper: {
          DEFAULT: '#12214C',
          deep: '#0C1D4C',
          raised: '#1A2D61',
          night: '#0C1D4C',
        },
        ink: {
          DEFAULT: '#FFFFFF',
          soft: '#C5CBE0',
          faint: '#B4BCD4',
        },
        wax: {
          DEFAULT: '#D1AB4C',
          deep: '#B08E3A',
        },
        hairline: 'rgba(255, 255, 255, 0.10)',
      },
      fontFamily: {
        serif: ['var(--font-playfair)', 'Georgia', 'serif'],
        sans: ['var(--font-source)', 'system-ui', 'sans-serif'],
        mono: ['var(--font-source)', 'system-ui', 'sans-serif'],
      },
      backgroundImage: {
        'gold-gradient': 'linear-gradient(180deg, #E2C46E, #D1AB4C)',
        'navy-gradient': 'linear-gradient(180deg, #12214C 0%, #162756 48%, #1A2D61 100%)',
      },
      boxShadow: {
        card: 'inset 0 1px 0 rgba(255,255,255,0.06), 0 18px 40px -28px rgba(0,0,0,0.75)',
      },
    },
  },
  plugins: [],
}

export default config
