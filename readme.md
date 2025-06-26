# A Golang web app for fun

Wanted to try to make a go web app. It has OAuth2. Users can login using discord.
You can create, view, delete and edit tasks.

---

## Build & Run Instructions

Follow these steps to build and run the project:

### 📦 Install dependencies

```bash
pnpm install
```

---

## Build steps

Run the following commands to build the project:

### Build TypeScript / JavaScript

```bash
pnpm build
```

Compiles and bundles your TypeScript/JavaScript files into `web/static/js/`.

---

### Build CSS

```bash
pnpm css-build
```

---

### Build Go binary

```bash
go build .
```

---

## Run the application

After building the project, you can run the generated binary:

```bash
./cmd/go-web/go-web.exe
```

---

## Available Scripts

| Command          | Description                              |
|:-----------------|:-----------------------------------------|
| `pnpm install`   | Install project dependencies             |
| `pnpm build`     | Compile and bundle TypeScript/JavaScript |
| `pnpm css-build` | Compile Tailwind or other CSS files      |
| `go build .`     | Build Go executable                      |

---

## License

This project is licensed under the MIT License. See `LICENSE` for details.