<script lang="ts">
  import Sidebar from '$lib/components/layout/Sidebar.svelte';
  import Header from '$lib/components/layout/Header.svelte';
  import { isAuthenticated } from '$lib/stores/auth';

  let collapsed = $state(false);
</script>

<div class="min-h-screen bg-base-100">
  {#if $isAuthenticated}
    <Sidebar bind:collapsed />
    <div class="main-content transition-all duration-300" class:expanded={collapsed}>
      <Header />
      <main class="p-6">
        <slot />
      </main>
    </div>
  {:else}
    <slot />
  {/if}
</div>
