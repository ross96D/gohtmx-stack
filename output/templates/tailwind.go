package templates

const TailwindConfig = `/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/*.templ", "./views/**/*.templ"],
  theme: {
    extend: {},
  },
  plugins: [],
}

`

const TailwindInput = `@tailwind base;
@tailwind components;
@tailwind utilities;
`
