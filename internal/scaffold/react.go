package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
)

func GenerateReact(projectName string) error {
	fmt.Println("🚧 Create project React:", projectName)

	if err := os.MkdirAll(projectName, os.ModePerm); err != nil {
		return fmt.Errorf("internal error creating project directory: %w", err)
	}

	packageJSON := `{
  "name": "` + projectName + `",
  "version": "0.1.0",
  "private": true,
  "scripts": {
    "dev": "vite",
    "build": "vite build",
    "preview": "vite preview"
  },
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0"
  },
  "devDependencies": {
    "vite": "^5.0.0"
  }
}`

	err := os.WriteFile(filepath.Join(projectName, "package.json"), []byte(packageJSON), 0644)
	if err != nil {
		return fmt.Errorf("không thể tạo package.json: %w", err)
	}

	indexHTML := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <title>` + projectName + `</title>
</head>
<body>
  <div id="root"></div>
  <script type="module" src="/main.jsx"></script>
</body>
</html>`

	err = os.WriteFile(filepath.Join(projectName, "index.html"), []byte(indexHTML), 0644)
	if err != nil {
		return fmt.Errorf("không thể tạo index.html: %w", err)
	}

	mainJSX := `import React from "react";
import ReactDOM from "react-dom/client";

const App = () => <h1>Hello from ` + projectName + `!</h1>;

ReactDOM.createRoot(document.getElementById("root")).render(<App />);`

	err = os.WriteFile(filepath.Join(projectName, "main.jsx"), []byte(mainJSX), 0644)
	if err != nil {
		return fmt.Errorf("không thể tạo main.jsx: %w", err)
	}

	fmt.Println("✅ React project done:", projectName)
	fmt.Println("👉 You are run:")
	fmt.Println("   cd", projectName)
	fmt.Println("   npm install && npm run dev")
	return nil
}
