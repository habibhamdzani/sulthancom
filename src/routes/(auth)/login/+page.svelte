<script lang="ts">
  import { onMount } from 'svelte';
  import { login, isAuthenticated } from '$lib/stores/auth';
  import { Lock, User } from 'lucide-svelte';
  import { goto } from '$app/navigation';
  import { mockStoreSettings } from '$lib/data/mock';

  let email = $state('');
  let password = $state('');
  let error = $state('');
  let loading = $state(false);
  let storeName = $state('Toko');
  let storeAddress = $state('');
  let storeLogo = $state<string | null>(null);

  onMount(() => {
    try {
      const saved = localStorage.getItem('storeSettings');
      if (saved) {
        const s = JSON.parse(saved);
        if (s.name) storeName = s.name;
        if (s.address) storeAddress = s.address;
      }
    } catch {}
    try {
      const logo = localStorage.getItem('storeLogo');
      if (logo) storeLogo = logo;
    } catch {}
  });

  $effect(() => {
    if (!$isAuthenticated) {
      email = '';
      password = '';
      error = '';
    }
  });

  async function handleSubmit() {
    if (!email || !password) {
      error = 'Email dan password harus diisi';
      return;
    }
    loading = true;
    error = '';
    const result = await login(email, password);
    loading = false;
    if (result.success) {
      goto('/dashboard');
    } else {
      error = result.error || 'Login gagal';
    }
  }
</script>

<svelte:head>
  <title>Login - {storeName}</title>
</svelte:head>

<div class="min-h-screen gradient-bg-modern flex items-center justify-center p-4">
  <div class="glass-card w-full max-w-md">
    <div class="p-6">
      <div class="text-center mb-6">
        {#if storeLogo}
          <div class="inline-flex items-center justify-center w-20 h-20 rounded-2xl mb-4 overflow-hidden mx-auto icon-container-3d">
            <img src={storeLogo} alt={storeName} class="w-full h-full object-contain relative z-10" />
          </div>
        {:else}
          <div class="inline-flex items-center justify-center w-16 h-16 icon-container-3d mb-4">
            <span class="text-2xl font-bold text-primary relative z-10">{storeName.charAt(0)}</span>
          </div>
        {/if}
        <h1 class="text-2xl font-bold">{storeName}</h1>
        {#if storeAddress}
          <p class="text-base-content text-sm">{storeAddress}</p>
        {:else}
          <p class="text-base-content text-sm">POS & Service Management System</p>
        {/if}
      </div>

      {#if error}
        <div class="alert alert-error alert-sm mb-4">
          <span>{error}</span>
        </div>
      {/if}

      <form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}>
        <div class="form-control mb-4">
          <label class="label">
            <span class="label-text font-medium">Email</span>
          </label>
          <div class="relative">
            <User class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content icon-3d" />
            <input
              type="email"
              class="input input-bordered w-full pl-10 glass"
              placeholder="admin@sulthancom.com"
              bind:value={email}
            />
          </div>
        </div>

        <div class="form-control mb-6">
          <label class="label">
            <span class="label-text font-medium">Password</span>
          </label>
          <div class="relative">
            <Lock class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content icon-3d" />
            <input
              type="password"
              class="input input-bordered w-full pl-10 glass"
              placeholder="password"
              bind:value={password}
            />
          </div>
        </div>

        <button type="submit" class="btn btn-primary w-full glass-strong" disabled={loading}>
          {#if loading}
            <span class="loading loading-spinner loading-sm"></span>
          {:else}
            Login
          {/if}
        </button>
      </form>
    </div>
  </div>
</div>
