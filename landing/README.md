# Astro Starter Kit: Basics

```sh
npm create astro@latest -- --template basics
```

> 🧑‍🚀 **Seasoned astronaut?** Delete this file. Have fun!

## 🚀 Project Structure

Inside of your Astro project, you'll see the following folders and files:

```text
/
├── public/
│   └── favicon.svg
├── src
│   ├── assets
│   │   └── astro.svg
│   ├── components
│   │   └── Welcome.astro
│   ├── layouts
│   │   └── Layout.astro
│   └── pages
│       └── index.astro
└── package.json
```

To learn more about the folder structure of an Astro project, refer to [our guide on project structure](https://docs.astro.build/en/basics/project-structure/).

## 🧞 Commands

All commands are run from the root of the project, from a terminal:

| Command                   | Action                                           |
| :------------------------ | :----------------------------------------------- |
| `npm install`             | Installs dependencies                            |
| `npm run dev`             | Starts local dev server at `localhost:4321`      |
| `npm run build`           | Build your production site to `./dist/`          |
| `npm run preview`         | Preview your build locally, before deploying     |
| `npm run astro ...`       | Run CLI commands like `astro add`, `astro check` |
| `npm run astro -- --help` | Get help using the Astro CLI                     |

## 👀 Want to learn more?

Feel free to check [our documentation](https://docs.astro.build) or jump into our [Discord server](https://astro.build/chat).

## Tailwind Setup

This project is configured to use Tailwind CSS with a local build instead of the CDN. The following files were added:

- `tailwind.config.cjs` — Tailwind config with theme colors and font extension
- `postcss.config.cjs` — PostCSS setup to run Tailwind and Autoprefixer
- `src/assets/styles.css` — Global stylesheet with Tailwind directives and project variables

Install the required developer dependencies and then run the dev server:

```bash
# Install required dev dependencies
npm install -D tailwindcss postcss autoprefixer @tailwindcss/postcss @tailwindcss/forms
# then run your usual dev server
npm run dev
```

Notes:
- `src/layouts/Layout.astro` now imports `src/assets/styles.css` and no longer uses the CDN script.
- If you already have an existing Tailwind setup, back up `tailwind.config.cjs` before running `npx tailwindcss init -p`.
