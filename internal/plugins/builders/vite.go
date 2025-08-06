package builders

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	tilocontext "github.com/ti-lo/tilokit/internal/core/context"
	"github.com/ti-lo/tilokit/internal/utils"
)

// VitePlugin implements Vite build tool integration
type VitePlugin struct{}

// NewVitePlugin creates a new Vite plugin instance
func NewVitePlugin() *VitePlugin {
	return &VitePlugin{}
}

// Name returns the name of the Vite plugin.
func (p *VitePlugin) Name() string {
	return "vite-builder"
}

// Version returns the version of the Vite plugin.
func (p *VitePlugin) Version() string {
	return "1.0.0"
}

// Description returns the description of the Vite plugin.
func (p *VitePlugin) Description() string {
	return "Vite build tool integration with modern configuration"
}

// SupportedFrameworks returns the list of frameworks supported by Vite.
func (p *VitePlugin) SupportedFrameworks() []string {
	return []string{"react", "vue", "svelte", "vanilla"}
}

// SupportedBuildTools returns the list of build tools supported by Vite.
func (p *VitePlugin) SupportedBuildTools() []string {
	return []string{"vite"}
}

// PreGenerate validates Vite compatibility before project generation.
func (p *VitePlugin) PreGenerate(ctx *tilocontext.ExecutionContext) error {
	// Validate Vite compatibility
	framework := ctx.Config.Framework
	if !utils.Contains(p.SupportedFrameworks(), framework) {
		return errors.Errorf("framework %s is not supported by Vite plugin", framework)
	}

	ctx.SetMetadata("vite_config_generated", false)
	return nil
}

func (p *VitePlugin) Generate(ctx *tilocontext.ExecutionContext) error {
	// Generate vite.config.js based on framework
	viteConfig, err := p.generateViteConfig(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to generate Vite config")
	}

	configPath := filepath.Join(ctx.ProjectPath, "vite.config.js")
	if err := utils.WriteFile(configPath, viteConfig); err != nil {
		return errors.Wrap(err, "failed to write Vite config")
	}

	// Generate package.json scripts
	if err := p.updatePackageJsonScripts(ctx); err != nil {
		return errors.Wrap(err, "failed to update package.json scripts")
	}

	ctx.SetMetadata("vite_config_generated", true)
	return nil
}

func (p *VitePlugin) PostGenerate(ctx *tilocontext.ExecutionContext) error {
	// Add Vite-specific development instructions
	ctx.SetMetadata("dev_command", "npm run dev")
	ctx.SetMetadata("build_command", "npm run build")
	ctx.SetMetadata("preview_command", "npm run preview")
	return nil
}

func (p *VitePlugin) generateViteConfig(ctx *tilocontext.ExecutionContext) (string, error) {
	framework := ctx.Config.Framework

	var config string
	switch framework {
	case "react":
		config = `import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    open: true,
    host: true
  },
  build: {
    outDir: 'dist',
    sourcemap: true,
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['react', 'react-dom']
        }
      }
    }
  },
  resolve: {
    alias: {
      '@': '/src'
    }
  }
})`
	case "vue":
		config = `import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    port: 3000,
    open: true,
    host: true
  },
  build: {
    outDir: 'dist',
    sourcemap: true,
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['vue']
        }
      }
    }
  },
  resolve: {
    alias: {
      '@': '/src'
    }
  }
})`
	case "svelte":
		config = `import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()],
  server: {
    port: 3000,
    open: true,
    host: true
  },
  build: {
    outDir: 'dist',
    sourcemap: true
  },
  resolve: {
    alias: {
      '@': '/src'
    }
  }
})`
	default:
		config = `import { defineConfig } from 'vite'

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    port: 3000,
    open: true,
    host: true
  },
  build: {
    outDir: 'dist',
    sourcemap: true
  },
  resolve: {
    alias: {
      '@': '/src'
    }
  }
})`
	}

	return config, nil
}

func (p *VitePlugin) updatePackageJsonScripts(ctx *tilocontext.ExecutionContext) error {
	packageJsonPath := filepath.Join(ctx.ProjectPath, "package.json")

	// Read existing package.json if it exists
	var packageJson map[string]interface{}
	if utils.FileExists(packageJsonPath) {
		dataStr, err := utils.ReadFile(packageJsonPath)
		if err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(dataStr), &packageJson); err != nil {
			return err
		}
	} else {
		packageJson = make(map[string]interface{})
	}

	// Ensure scripts section exists
	if packageJson["scripts"] == nil {
		packageJson["scripts"] = make(map[string]interface{})
	}

	scripts := packageJson["scripts"].(map[string]interface{})

	// Add Vite scripts
	scripts["dev"] = "vite"
	scripts["build"] = "vite build"
	scripts["preview"] = "vite preview"
	scripts["lint"] = "eslint . --ext js,jsx,ts,tsx,vue,svelte --report-unused-disable-directives --max-warnings 0"

	// Write updated package.json
	data, err := json.MarshalIndent(packageJson, "", "  ")
	if err != nil {
		return err
	}

	// Use more restrictive permissions (0600 instead of 0644)
	return os.WriteFile(packageJsonPath, data, 0600)
}
