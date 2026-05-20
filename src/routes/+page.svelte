<script lang="ts">
  import { goto } from '$app/navigation';
  import { isAuthenticated } from '$lib/stores/auth';

  $effect(() => {
    if ($isAuthenticated) {
      goto('/dashboard');
    } else {
      goto('/login');
    }
  });
</script>

<div class="min-h-screen flex items-center justify-center">
  <span class="loading loading-spinner loading-lg"></span>
</div>
