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
        paper: {
          DEFAULT: '#ece6d8',
          deep: '#e2d8c6',
          raised: '#f4efe6',
          night: '#1c1713',
        },
        ink: {
          DEFAULT: '#1d1814',
          soft: '#5a5148',
          faint: '#85786c',
        },
        wax: {
          DEFAULT: '#7a2c28',
          deep: '#5c211e',
        },
        hairline: 'rgba(29, 24, 20, 0.14)',
      },
      fontFamily: {
        serif: ['var(--font-newsreader)', 'Georgia', 'Times New Roman', 'serif'],
        mono: ['var(--font-plex-mono)', 'ui-monospace', 'monospace'],
      },
    },
  },
  plugins: [],
}

export default config
