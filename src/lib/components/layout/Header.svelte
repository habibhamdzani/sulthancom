<script lang="ts">
  import { currentUser, currentTheme, applyTheme } from '$lib/stores/auth';
  import { Bell, Moon, Sun, Search } from 'lucide-svelte';

  let searchQuery = $state('');
  let theme = $state('light');

  $effect(() => {
    theme = $currentTheme;
  });

  function toggleTheme() {
    const next = theme === 'light' ? 'dark' : 'light';
    applyTheme(next);
  }
</script>

<header class="sticky top-0 z-30 h-16 glass border-b border-base-300/50 flex items-center px-6 gap-4">
  <div class="flex-1 max-w-md">
    <div class="relative">
      <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content icon-3d" />
      <input
        type="text"
        placeholder="Cari sesuatu..."
        class="input input-bordered input-sm w-full pl-10 glass"
        bind:value={searchQuery}
      />
    </div>
  </div>

  <div class="flex items-center gap-2 ml-auto">
    <button class="btn btn-ghost btn-sm btn-circle glass" onclick={toggleTheme}>
      {#if theme === 'light'}
        <Moon class="w-5 h-5 icon-3d" />
      {:else}
        <Sun class="w-5 h-5 icon-3d" />
      {/if}
    </button>

    <div class="dropdown dropdown-end">
      <button class="btn btn-ghost btn-sm btn-circle indicator glass">
        <Bell class="w-5 h-5 icon-3d" />
        <span class="badge badge-xs badge-primary indicator-item"></span>
      </button>
      <div class="dropdown-content dropdown-end mt-4 w-80 glass-card rounded-box shadow-lg z-50">
        <div class="p-4 border-b border-base-200/50">
          <h3 class="font-semibold">Notifikasi</h3>
        </div>
        <div class="p-4 text-center text-sm text-base-content">
          Tidak ada notifikasi baru
        </div>
      </div>
    </div>

    {#if $currentUser}
      <div class="dropdown dropdown-end">
        <button class="btn btn-ghost btn-sm gap-2 glass">
          <div class="icon-container-3d w-8 h-8 flex items-center justify-center">
            <span class="text-primary font-bold text-xs relative z-10">{$currentUser.name.charAt(0)}</span>
          </div>
          <span class="hidden md:inline text-sm font-medium">{$currentUser.name}</span>
        </button>
        <div class="dropdown-content dropdown-end mt-4 w-56 glass-card rounded-box shadow-lg z-50">
          <ul class="menu menu-sm p-2">
            <li><a href="/settings">Pengaturan</a></li>
            <li><a href="/users">Manajemen User</a></li>
            <li class="border-t border-base-200/50 mt-1 pt-1"><a href="/login" class="text-error">Logout</a></li>
          </ul>
        </div>
      </div>
    {/if}
  </div>
</header>
