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
      /*  colors: {
         'blue': '#1fb6ff',
         'purple': '#7e5bef',
         'pink': '#ff49db',
         'orange': '#ff7849',
         'green': '#13ce66',
         'yellow': '#ffc82c',
         'gray-dark': '#273444',
         'gray': '#8492a6',
         'gray-light': '#d3dce6',
       }, */
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
      borderRadius: {
        '4xl': '2rem',
      },
      listStyleImage: {
        check: "url(/assets/media/icons/circle-check.png)",
      }
    }
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}

