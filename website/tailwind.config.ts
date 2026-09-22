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
          DEFAULT: '#0F1B2D',
          light: '#162235',
          border: '#1E3050',
        },
        gold: {
          DEFAULT: '#C9A84C',
          light: '#E8C87A',
        },
        cream: {
          DEFAULT: '#F0EDE8',
          muted: '#8A9BB0',
        },
        // Keep old class names working while pages migrate off paper/wax.
        paper: {
          DEFAULT: '#0F1B2D',
          deep: '#0B1422',
          raised: '#162235',
          night: '#0A1628',
        },
        ink: {
          DEFAULT: '#F0EDE8',
          soft: '#C5C0B8',
          faint: '#8A9BB0',
        },
        wax: {
          DEFAULT: '#C9A84C',
          deep: '#A88B38',
        },
        hairline: 'rgba(201, 168, 76, 0.22)',
      },
      fontFamily: {
        serif: ['var(--font-playfair)', 'Georgia', 'serif'],
        sans: ['var(--font-inter)', 'system-ui', 'sans-serif'],
        mono: ['var(--font-inter)', 'system-ui', 'sans-serif'],
      },
      backgroundImage: {
        'gold-gradient': 'linear-gradient(135deg, #C9A84C, #E8C87A)',
        'navy-gradient': 'linear-gradient(180deg, #0F1B2D 0%, #162235 100%)',
      },
    },
  },
  plugins: [],
}

export default config
