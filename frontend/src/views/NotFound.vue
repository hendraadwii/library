<template>
  <div class="min-h-screen bg-gray-50 flex flex-col justify-center items-center px-4 sm:px-6 lg:px-8">
    <div class="text-center">
      <h1 class="text-9xl font-bold text-indigo-600">404</h1>
      <h2 class="mt-4 text-3xl font-extrabold text-gray-900 tracking-tight sm:text-4xl">Page not found</h2>
      <p class="mt-6 text-base text-gray-500">Sorry, we couldn't find the page you're looking for.</p>
      <div class="mt-10">
        <router-link to="/" class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
          Back to home
        </router-link>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'NotFound'
}
</script> 