package utils

// MapLibsToPackages converts user-friendly library names returned from the
// survey prompt to the actual npm package names.
func MapLibsToPackages(libs []string) []string {
    mapping := map[string]string{
        "ESLint":          "eslint",
        "Prettier":        "prettier",
        "TailwindCSS":     "tailwindcss",
        "React Router":    "react-router-dom",
        "Zustand":         "zustand",
        "Axios":           "axios",
        "Jest":            "jest",
        "Vue Router":      "vue-router",
        "Pinia":           "pinia",
        "Vitest":          "vitest",
    }

    var packages []string
    for _, l := range libs {
        if p, ok := mapping[l]; ok {
            packages = append(packages, p)
        }
    }
    return packages
}
