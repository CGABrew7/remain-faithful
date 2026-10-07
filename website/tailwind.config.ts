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
          mid: '#1C3066',
          deep: '#0C1D4C',
          border: 'rgba(255,255,255,0.10)',
        },
        gold: {
          DEFAULT: '#D1AB4C',
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
        'navy-gradient': 'linear-gradient(180deg, #12214C 0%, #1C3066 100%)',
      },
      boxShadow: {
        card: 'inset 0 1px 0 rgba(255,255,255,0.06), 0 18px 40px -28px rgba(0,0,0,0.75)',
      },
    },
  },
  plugins: [],
}

export default config
