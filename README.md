# BootstrapMe

**BootstrapMe** is a universal CLI tool to bootstrap projects in various languages and frameworks. It reads YAML configuration files from `~/.config/bootstrapme/<language>/<framework>.yaml` and guides you through selecting a language, framework, and optionally asking for details like username or project name before running the necessary steps.

## Installation
`sudo bash install.sh`

## Features

- Supports multiple languages and frameworks.
- Automatically runs post-install commands.
- Uses template variables (`{{ variable_name }}`) in commands and files.
- Modern TUI with Lipgloss styling and ASCII art title.
- Error lines highlighted in red for easy debugging.
- Can run non-interactively by supplying defaults to external bootstrap tools (like `create-next-app`).

## Writing Custom Config Files

**Directory:** `~/.config/bootstrapme/<language>/<name>.yaml`

**Example (Go + Gin):**

```yaml
# ~/.config/bootstrapme/go/gin.yaml
name: gin
description: "Bootstrap a Gin-based Go server"
language: go
framework: gin
variables:
  module_name: "github.com/{{ username }}/myginapp"
post_install_commands:
  - "go mod init {{ module_name }}"
  - "go get github.com/gin-gonic/gin"
```

**Fields:**

-   **name**: Name of the preset (shown in TUI).
-   **description**: A short description shown in TUI.
-   **language**: The language category (e.g., `go`, `js`, `python`).
-   **framework**: The framework/tool name (e.g., `gin`, `createreactapp`, `createnextapp`).
-   **variables**: Key-value pairs for template substitution.
-   **files**: An array of files to create (with `path` and `content`).
-   **post_install_commands**: Shell commands run after file creation. Use `{{ variables }}` to substitute values.

### Installers With Interactive Steps
See the below example of using `create-next-app` with a set of default flags
**Directory:** `~/.config/bootstrapme/js/next.yaml`

```yaml 
# ~/.config/bootstrapme/js/createnextapp.yaml
name: Create Next App 
description: "Bootstrap a Next.js project"
language: js
framework: Next 
post_install_commands:
  - "npx create-next-app {{ project_name }} --ts --use-npm --eslint --src-dir --tailwind --app --no-turbopack --no-import-alias"
```

## Contributing
Raise an issue and link your PR to it
