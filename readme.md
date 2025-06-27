# Go-Web

A Golang web app for fun.

Users can login using Discord OAuth2. You can create, view, delete, and edit tasks.

---

## Prerequisites

- Go 1.20+  
- Node.js & pnpm  
- PostgreSQL (running locally or via Docker)

---

## Setup

### Install frontend dependencies

```bash
pnpm install
```

---

## Build steps

### Build and bundle TypeScript (with esbuild)

```bash
pnpm run ts:build
```

Or for watch mode (rebuild on changes):

```bash
pnpm run ts:watch
```

This outputs bundled JS to `web/static/js/bundle.js`.

---

### Build CSS (Tailwind)

```bash
pnpm run css:build
```

Or for watch mode (rebuild on changes):

```bash
pnpm run css:watch
```

This compiles Tailwind CSS from `web/static/css/global.css` to `web/static/css/tailwind.css`.

---

### Build full frontend (using Vite)

```bash
pnpm run build
```

Or for watch mode:

```bash
pnpm run dev
```

> Note: Vite builds and watches the frontend assets. You may prefer to use esbuild or Vite depending on your workflow.

---

## Build the Go backend

From the project root or `cmd/go-web`:

```bash
go build -o cmd/go-web/tmp/main.exe ./cmd/go-web
```

Or simply:

```bash
go build .
```

---

## Run the application

Run the compiled binary:

```bash
./cmd/go-web/tmp/main.exe
```

---

## Available Scripts

| Command           | Description                           |
|-------------------|-------------------------------------|
| `pnpm install`    | Install project dependencies         |
| `pnpm run build`  | Build frontend assets with Vite      |
| `pnpm run dev`    | Build and watch frontend assets with Vite |
| `pnpm run css:build` | Compile Tailwind CSS              |
| `pnpm run css:watch` | Compile and watch Tailwind CSS    |
| `pnpm run ts:build` | Bundle TypeScript with esbuild     |
| `pnpm run ts:watch` | Watch and bundle TypeScript with esbuild |
| `go build .`      | Build Go backend executable          |

---

## License

This project is licensed under the MIT License. See `LICENSE` for details.
