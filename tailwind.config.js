/** @type {import('tailwindcss').Config} */
module.exports = {
  content: {
    relative: true,
    files: [
      "./handlers/*.{html,js}",
      "./ui/comps/**/*.{html,js,vugu}",
      "./internal/handlers/page-handler.go",
      "./static/media/svg/*.svg",

    ],
  },
  darkMode: "class",
  theme: {
    backgroundImage: {
      go: "url(/assets/media/patterns/gopattern.svg)",
    },


    fontFamily: {
      sans: ['"Open Sans"', 'sans-serif'],
      serif: ['Merriweather', 'serif'],
      logo: ['"Courier Prime"', 'serif'],
      'body': ['"Open Sans"',],
    },
    extend: {
      colors: {
        primary: {
          DEFAULT: "#111827",
          foreground: "#f9fafb",
        },
        secondary: {
          DEFAULT: "#f1f5f9",
          foreground: "#1f2937",
        },
        destructive: {
          DEFAULT: "#ef4444",
          foreground: "#f9fafb",
        },
        muted: {
          DEFAULT: "#f1f5f9",
          foreground: "#6b7280",
        },
        accent: {
          DEFAULT: "#f1f5f9",
          foreground: "#1f2937",
        },
        card: {
          DEFAULT: "#ffffff",
          foreground: "#1f2937",
        },
      },
      borderRadius: {
        lg: "0.5rem",
        md: "calc(0.5rem - 2px)",
        sm: "calc(0.5rem - 4px)",
        '4xl': '2rem',
      },
      backgroundColor: {
        'default': '#333',
        'landing-modal':'#FFF3DE',
        'landing-cta':'#A9694F',

      },
      backgroundImage: {
        'landing-hero':"url('/assets/media/svg/artboard.svg')",
        'landing-hero-png':"url('/assets/media/illustrations/artboard.png')",
      },
      spacing: {
        '8xl': '96rem',
        '9xl': '128rem',
      },
      listStyleImage: {
        check: "url(/assets/media/icons/circle-check.png)",
      },
      boxShadow: {
        'sm': '0 1px 2px 0 rgb(0 0 0 / 0.05)',
        DEFAULT: '0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1)',
        'md': '0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1)',
        'lg': '0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1)',
        'xl': '0 20px 25px -5px rgb(0 0 0 / 0.1), 0 8px 10px -6px rgb(0 0 0 / 0.1)',
        '2xl': '0 25px 50px -12px rgb(0 0 0 / 0.25)',
        'inner': 'inset 0 2px 4px 0 rgb(0 0 0 / 0.05)',
        'none': 'none'
      },
    }
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}

