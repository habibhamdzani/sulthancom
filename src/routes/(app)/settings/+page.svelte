<script lang="ts">
  import { onMount } from 'svelte';
  import { mockStoreSettings } from '$lib/data/mock';
  import { applyTheme, currentTheme } from '$lib/stores/auth';
  import { Save, Upload, Image } from 'lucide-svelte';

  let settings = $state({ ...mockStoreSettings });
  let saved = $state(false);
  let logoPreview = $state<string | null>(null);
  let fileInput: HTMLInputElement | undefined = $state();

  onMount(() => {
    try {
      const stored = localStorage.getItem('storeSettings');
      if (stored) Object.assign(settings, JSON.parse(stored));
    } catch {}
    try {
      const logo = localStorage.getItem('storeLogo');
      if (logo) logoPreview = logo;
    } catch {}
    settings.theme = $currentTheme;
  });

  function triggerUpload() {
    fileInput?.click();
  }

  function handleFileUpload(e: Event) {
    const input = e.target as HTMLInputElement;
    const file = input.files?.[0];
    if (!file) return;
    const reader = new FileReader();
    reader.onload = (ev) => {
      const dataUrl = ev.target?.result as string;
      logoPreview = dataUrl;
      try { localStorage.setItem('storeLogo', dataUrl); } catch {}
    };
    reader.readAsDataURL(file);
  }

  function handleSave() {
    try {
      localStorage.setItem('storeSettings', JSON.stringify(settings));
    } catch {}
    applyTheme(settings.theme);
    saved = true;
    setTimeout(() => saved = false, 3000);
  }
</script>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-2xl font-bold">Pengaturan Toko</h1>
      <p class="text-base-content">Konfigurasi pengaturan sistem</p>
    </div>
    {#if saved}
      <div class="badge badge-success gap-2">Tersimpan!</div>
    {/if}
  </div>

  <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
    <div class="glass-card">
      <div class="p-5">
        <h2 class="card-title">Informasi Toko</h2>

        <div class="form-control mb-4">
          <label class="label"><span class="label-text font-medium">Logo Toko</span></label>
            <div class="flex items-center gap-4">
              <div class="w-20 h-20 bg-base-200 rounded-lg flex items-center justify-center overflow-hidden">
                {#if logoPreview}
                  <img src={logoPreview} alt="Logo" class="w-full h-full object-contain" />
                {:else}
                  <span class="text-2xl font-bold text-primary">SC</span>
                {/if}
              </div>
              <input type="file" accept="image/*" class="hidden" bind:this={fileInput} onchange={handleFileUpload} />
              <button class="btn btn-outline btn-sm gap-2" onclick={triggerUpload}><Upload class="w-4 h-4" /> Upload Logo</button>
            </div>
        </div>

        <div class="form-control mb-4">
          <label class="label"><span class="label-text font-medium">Nama Toko</span></label>
          <input type="text" class="input input-bordered w-full" bind:value={settings.name} />
        </div>

        <div class="form-control mb-4">
          <label class="label"><span class="label-text font-medium">Alamat</span></label>
          <textarea class="textarea textarea-bordered w-full" bind:value={settings.address} rows="2"></textarea>
        </div>

        <div class="grid grid-cols-2 gap-4 mb-4">
          <div class="form-control">
            <label class="label"><span class="label-text font-medium">Telepon</span></label>
            <input type="text" class="input input-bordered w-full" bind:value={settings.phone} />
          </div>
          <div class="form-control">
            <label class="label"><span class="label-text font-medium">Email</span></label>
            <input type="email" class="input input-bordered w-full" bind:value={settings.email} />
          </div>
        </div>

        <div class="form-control mb-4">
          <label class="label"><span class="label-text font-medium">Website</span></label>
          <input type="text" class="input input-bordered w-full" bind:value={settings.website} />
        </div>

        <button class="btn btn-primary" onclick={handleSave}>
          <Save class="w-4 h-4" /> Simpan Pengaturan
        </button>
      </div>
    </div>

    <div class="space-y-6">
      <div class="glass-card">
        <div class="p-5">
          <h2 class="card-title">Pengaturan Transaksi</h2>

          <div class="form-control mb-4">
            <label class="label"><span class="label-text font-medium">Pajak (%)</span></label>
            <input type="number" class="input input-bordered w-full" bind:value={settings.taxRate} min="0" max="100" />
          </div>

          <div class="form-control mb-4">
            <label class="label"><span class="label-text font-medium">Mata Uang</span></label>
            <select class="select select-bordered w-full" bind:value={settings.currency}>
              <option value="IDR">IDR - Rupiah</option>
              <option value="USD">USD - Dollar</option>
            </select>
          </div>

          <div class="form-control">
            <label class="label"><span class="label-text font-medium">Footer Invoice</span></label>
            <textarea class="textarea textarea-bordered w-full" bind:value={settings.invoiceFooter} rows="2"></textarea>
          </div>
        </div>
      </div>

      <div class="glass-card">
        <div class="p-5">
          <h2 class="card-title">Tema</h2>
          <div class="grid grid-cols-3 gap-3">
            {#each ['light', 'dark', 'cupcake', 'corporate', 'retro'] as theme}
              <button class="btn btn-sm btn-outline capitalize" class:btn-primary={settings.theme === theme} onclick={() => { settings.theme = theme; applyTheme(theme); }}>
                {theme}
              </button>
            {/each}
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
