package nuxtjs

const (
	NuxtjsPackageJson = `{
  "name": "{{.ProjectName}}",
  "private": true,
  "scripts": {
    "build": "nuxt build",
    "dev": "nuxt dev",
    "generate": "nuxt generate",
    "preview": "nuxt preview",
    "postinstall": "nuxt prepare"
  },
  "devDependencies": {
    "@nuxt/devtools": "latest",
    "nuxt": "^3.8.0"
  },
  "dependencies": {
    "@nuxtjs/tailwindcss": "^6.8.4",
    "vue": "^3.3.0"
  }
}`

	NuxtjsAppVue = `<template>
  <div>
    <NuxtWelcome />
  </div>
</template>
`

	NuxtjsIndexPage = `<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-4xl font-bold text-center mb-8">
      Welcome to {{.ProjectName}}
    </h1>
    <p class="text-lg text-center text-gray-600">
      Get started by editing <code class="bg-gray-100 px-2 py-1 rounded">pages/index.vue</code>
    </p>
    <div class="flex justify-center mt-8">
      <NuxtLink 
        to="/about" 
        class="bg-green-500 hover:bg-green-600 text-white px-6 py-3 rounded-lg transition-colors"
      >
        Learn More
      </NuxtLink>
    </div>
  </div>
</template>

<script setup>
// This is a Nuxt 3 page with auto-imports
useHead({
  title: '{{.ProjectName}} - Home'
})
</script>
`

	NuxtjsAboutPage = `<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-4xl font-bold text-center mb-8">
      About {{.ProjectName}}
    </h1>
    <p class="text-lg text-center text-gray-600 mb-8">
      This is a Nuxt.js application with TypeScript and Tailwind CSS.
    </p>
    <div class="flex justify-center">
      <NuxtLink 
        to="/" 
        class="bg-blue-500 hover:bg-blue-600 text-white px-6 py-3 rounded-lg transition-colors"
      >
        Go Home
      </NuxtLink>
    </div>
  </div>
</template>

<script setup>
useHead({
  title: '{{.ProjectName}} - About'
})
</script>
`

	NuxtjsLayoutDefault = `<template>
  <div>
    <header class="bg-white shadow-sm border-b">
      <nav class="container mx-auto px-4 py-4">
        <div class="flex justify-between items-center">
          <NuxtLink to="/" class="text-xl font-bold text-gray-800">
            {{.ProjectName}}
          </NuxtLink>
          <div class="space-x-4">
            <NuxtLink 
              to="/" 
              class="text-gray-600 hover:text-gray-800 transition-colors"
            >
              Home
            </NuxtLink>
            <NuxtLink 
              to="/about" 
              class="text-gray-600 hover:text-gray-800 transition-colors"
            >
              About
            </NuxtLink>
          </div>
        </div>
      </nav>
    </header>
    
    <main>
      <slot />
    </main>
    
    <footer class="bg-gray-50 border-t mt-16">
      <div class="container mx-auto px-4 py-8 text-center text-gray-600">
        <p>&copy; 2024 {{.ProjectName}}. Built with Nuxt.js</p>
      </div>
    </footer>
  </div>
</template>
`
)
