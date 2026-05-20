<script lang="ts">
  import type { Snippet } from 'svelte';
  import { isAuthenticated, loadUser } from '$lib/stores/auth';
  import { goto } from '$app/navigation';
  import Sidebar from '$lib/components/layout/Sidebar.svelte';
  import Header from '$lib/components/layout/Header.svelte';

  interface Props {
    children: Snippet;
  }

  let { children }: Props = $props();
  let loaded = $state(false);

  $effect(() => {
    if (!loaded) {
      loadUser().then(() => {
        loaded = true;
        if (!$isAuthenticated) {
          goto('/login');
        }
      });
    }
  });

  $effect(() => {
    if (loaded && !$isAuthenticated) {
      goto('/login');
    }
  });
</script>

{#if loaded && $isAuthenticated}
  <div class="min-h-screen bg-base-100">
    <Sidebar />
    <div class="main-content">
      <Header />
      <main class="p-6">
        {@render children()}
      </main>
    </div>
  </div>
{/if}
