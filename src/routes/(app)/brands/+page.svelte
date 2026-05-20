<script lang="ts">
  import { mockBrands } from '$lib/data/mock';
  import { Plus, Edit } from 'lucide-svelte';

  let brands = $state(mockBrands);
  let showAddModal = $state(false);
  let newBrand = $state({ name: '' });
  let editingBrand = $state<any>(null);
  let showEditModal = $state(false);

  function resetNewBrand() { newBrand = { name: '' }; }

  function saveBrand() {
    if (!newBrand.name) return;
    brands.push({ id: crypto.randomUUID(), name: newBrand.name, productCount: 0, isActive: true });
    showAddModal = false;
    resetNewBrand();
  }

  function editBrand(brand: any) {
    editingBrand = { ...brand };
    showEditModal = true;
  }

  function updateBrand() {
    if (!editingBrand) return;
    const idx = brands.findIndex(b => b.id === editingBrand.id);
    if (idx >= 0) brands[idx] = { ...editingBrand };
    showEditModal = false;
    editingBrand = null;
  }
</script>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-2xl font-bold">Brand Produk</h1>
      <p class="text-base-content">Kelola brand/merek produk</p>
    </div>
    <button class="btn btn-primary btn-sm gap-2" onclick={() => showAddModal = true}>
      <Plus class="w-4 h-4" /> Tambah Brand
    </button>
  </div>

  <div class="glass-card">
    <div class="p-5">
      <div class="overflow-x-auto">
        <table class="table table-sm">
          <thead>
            <tr>
              <th>Nama Brand</th>
              <th>Jumlah Produk</th>
              <th>Status</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            {#each brands as brand}
              <tr class="hover">
                <td class="font-medium">{brand.name}</td>
                <td>{brand.productCount} produk</td>
                <td><span class="badge badge-sm {brand.isActive ? 'badge-success' : 'badge-ghost'}">{brand.isActive ? 'Aktif' : 'Nonaktif'}</span></td>
                <td>
                  <button class="btn btn-ghost btn-xs" onclick={() => editBrand(brand)}><Edit class="w-4 h-4" /></button>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>
  </div>
</div>

{#if showAddModal}
  <dialog class="modal modal-open" onclick={() => showAddModal = false}>
    <div class="modal-box" onclick={(e) => e.stopPropagation()}>
      <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2" onclick={() => { showAddModal = false; resetNewBrand(); }}>✕</button>
      <h3 class="font-bold text-lg mb-4">Tambah Brand</h3>
      <div class="space-y-3">
        <div class="form-control"><label class="label"><span class="label-text">Nama Brand</span></label><input type="text" class="input input-bordered input-sm" placeholder="Nama brand" bind:value={newBrand.name} /></div>
        <div class="flex justify-end gap-2 mt-4">
          <button class="btn btn-ghost btn-sm" onclick={() => { showAddModal = false; resetNewBrand(); }}>Batal</button>
          <button class="btn btn-primary btn-sm" onclick={saveBrand}>Simpan</button>
        </div>
      </div>
    </div>
  </dialog>
{/if}

{#if showEditModal && editingBrand}
  <dialog class="modal modal-open" onclick={() => showEditModal = false}>
    <div class="modal-box" onclick={(e) => e.stopPropagation()}>
      <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2" onclick={() => showEditModal = false}>✕</button>
      <h3 class="font-bold text-lg mb-4">Edit Brand</h3>
      <div class="space-y-3">
        <div class="form-control"><label class="label"><span class="label-text">Nama Brand</span></label><input type="text" class="input input-bordered input-sm" bind:value={editingBrand.name} /></div>
        <div class="flex justify-end gap-2 mt-4">
          <button class="btn btn-ghost btn-sm" onclick={() => showEditModal = false}>Batal</button>
          <button class="btn btn-primary btn-sm" onclick={updateBrand}>Simpan</button>
        </div>
      </div>
    </div>
  </dialog>
{/if}
