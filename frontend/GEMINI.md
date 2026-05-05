# Hackathon Judge - Frontend

This is the frontend component of the **Hackathon Judge** project. It is built using modern web technologies and follows a standard Vite-based React architecture.

## Project Overview

- **Purpose:** A web application designed to facilitate the judging process for hackathons.
- **Architecture:** React Single Page Application (SPA).
- **Core Technologies:**
  - **Framework:** [React 19](https://react.dev/)
  - **Build Tool:** [Vite 8](https://vite.dev/)
  - **Language:** [TypeScript](https://www.typescriptlang.org/)
  - **Linting:** [ESLint](https://eslint.org/)
  - **UI/Unit Testing:** [Vitest](https://vitest.dev/) & [React Testing Library](https://testing-library.com/docs/react-testing-library/intro/)
  - **E2E Testing:** [Playwright](https://playwright.dev/)

## Directory Structure

- `src/`: Main source code.
  - `App.tsx`: Main application component.
  - `main.tsx`: Entry point for the React application.
  - `assets/`: Static assets like images and SVGs.
- `tests/e2e/`: Playwright end-to-end tests.
- `public/`: Public assets that are served as-is (e.g., `favicon.svg`, `icons.svg`).
- `index.html`: The HTML template.

## Building and Running

The project uses `npm` for package management.

### Key Commands

- **Start Development Server:**
  ```bash
  npm run dev
  ```
  Starts the Vite development server with Hot Module Replacement (HMR).

- **Build for Production:**
  ```bash
  npm run build
  ```
  Runs TypeScript type checking and builds the production-ready assets in the `dist/` directory.

- **Linting:**
  ```bash
  npm run lint
  ```
  Runs ESLint to check for code quality and style issues.

- **Preview Production Build:**
  ```bash
  npm run preview
  ```
  Serves the locally built production files for testing.

- **Run UI/Unit Tests:**
  ```bash
  npm run test
  ```
  Runs Vitest unit and component tests. Use `npm run test:ui` to open the interactive UI.

- **Run E2E Tests:**
  ```bash
  npm run test:e2e
  ```
  Runs Playwright end-to-end tests.

## Development Conventions

- **Type Safety:** Use TypeScript for all new components and logic.
- **Styling:** The project uses **Tailwind CSS v4**. See [`DESIGN.md`](./DESIGN.md) for the complete design system, brand guidelines, color palette, and typography rules.
- **Hooks:** Leverage React Hooks for state and side-effect management.
- **Component Pattern:** Prefer functional components.
- **Testing:** Add UI/Component tests alongside the component file in `src/` (e.g., `Component.test.tsx`). Add E2E tests in the `tests/e2e/` directory.

## Related Components

This project is part of a larger monorepo which includes:
- `backend/`: (Currently empty) Intended for the server-side logic.
- `agent/`: (Currently empty) Intended for AI-related agents or automation.
