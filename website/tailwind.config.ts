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
          light: '#1C3066',
          border: '#2B3C70',
        },
        gold: {
          DEFAULT: '#D1AB4C',
          light: '#E2C46E',
        },
        cream: {
          DEFAULT: '#FFFFFF',
          muted: '#8E98B8',
        },
        // Keep old class names working while pages migrate off paper/wax.
        paper: {
          DEFAULT: '#12214C',
          deep: '#0E1A3E',
          raised: '#1C3066',
          night: '#0C1638',
        },
        ink: {
          DEFAULT: '#FFFFFF',
          soft: '#C5CBE0',
          faint: '#8E98B8',
        },
        wax: {
          DEFAULT: '#D1AB4C',
          deep: '#B08E3A',
        },
        hairline: 'rgba(209, 171, 76, 0.22)',
      },
      fontFamily: {
        serif: ['var(--font-playfair)', 'Georgia', 'serif'],
        sans: ['var(--font-inter)', 'system-ui', 'sans-serif'],
        mono: ['var(--font-inter)', 'system-ui', 'sans-serif'],
      },
      backgroundImage: {
        'gold-gradient': 'linear-gradient(135deg, #D1AB4C, #E2C46E)',
        'navy-gradient': 'linear-gradient(180deg, #12214C 0%, #162756 50%, #1A2D61 100%)',
      },
    },
  },
  plugins: [],
}

export default config
