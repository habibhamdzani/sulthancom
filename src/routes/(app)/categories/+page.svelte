<script lang="ts">
  import { mockCategories } from '$lib/data/mock';
  import { Plus, Edit, Trash2, MoreVertical } from 'lucide-svelte';

  let categories = $state(mockCategories);
  let showAddModal = $state(false);
  let newCategory = $state({ name: '', color: '#3B82F6', icon: '📦' });

  function resetNewCategory() { newCategory = { name: '', color: '#3B82F6', icon: '📦' }; }

  function saveCategory() {
    if (!newCategory.name) return;
    categories.push({ id: crypto.randomUUID(), name: newCategory.name, color: newCategory.color, icon: newCategory.icon, productCount: 0, isActive: true });
    showAddModal = false;
    resetNewCategory();
  }

  function toggleCategory(id: string) {
    const idx = categories.findIndex(c => c.id === id);
    if (idx >= 0) categories[idx] = { ...categories[idx], isActive: !categories[idx].isActive };
  }

  const colorOptions = ['#3B82F6', '#EF4444', '#10B981', '#F59E0B', '#8B5CF6', '#EC4899', '#06B6D4', '#F97316'];
</script>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-2xl font-bold">Kategori Produk</h1>
      <p class="text-base-content">Kelola kategori produk</p>
    </div>
    <button class="btn btn-primary btn-sm gap-2" onclick={() => showAddModal = true}>
      <Plus class="w-4 h-4" /> Tambah Kategori
    </button>
  </div>

  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
    {#each categories as cat}
      <div class="glass-card">
        <div class="p-5 p-5">
          <div class="flex items-start justify-between">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-lg flex items-center justify-center" style="background-color: {cat.color}20; color: {cat.color}">
                <span class="text-lg">{cat.icon?.charAt(0).toUpperCase()}</span>
              </div>
              <div>
                <h3 class="font-semibold">{cat.name}</h3>
                <p class="text-xs text-base-content">{cat.productCount} produk</p>
              </div>
            </div>
            <div class="dropdown dropdown-end">
              <button class="btn btn-ghost btn-xs"><MoreVertical class="w-4 h-4" /></button>
              <div class="dropdown-content dropdown-end mt-4 w-32 bg-base-100 rounded-box shadow-lg border border-base-200 z-50">
                <ul class="menu menu-sm p-2">
                  <li><a>Edit</a></li>
                  <li><a class="text-error">Hapus</a></li>
                </ul>
              </div>
            </div>
          </div>
          <div class="mt-3">
            <label class="label cursor-pointer justify-start gap-2">
              <input type="checkbox" class="toggle toggle-success toggle-sm" checked={cat.isActive} onchange={() => toggleCategory(cat.id)} />
              <span class="label-text text-xs">{cat.isActive ? 'Aktif' : 'Nonaktif'}</span>
            </label>
          </div>
        </div>
      </div>
    {/each}
  </div>
</div>

{#if showAddModal}
  <dialog class="modal modal-open" onclick={() => showAddModal = false}>
    <div class="modal-box" onclick={(e) => e.stopPropagation()}>
      <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2" onclick={() => { showAddModal = false; resetNewCategory(); }}>✕</button>
      <h3 class="font-bold text-lg mb-4">Tambah Kategori</h3>
      <div class="space-y-3">
        <div class="form-control"><label class="label"><span class="label-text">Nama Kategori</span></label><input type="text" class="input input-bordered input-sm" placeholder="Nama kategori" bind:value={newCategory.name} /></div>
        <div class="form-control"><label class="label"><span class="label-text">Warna</span></label>
          <div class="flex gap-2">
            {#each colorOptions as color}
              <button class="w-8 h-8 rounded-full border-2 {newCategory.color === color ? 'border-base-content' : 'border-transparent'}" style="background-color: {color}" onclick={() => newCategory.color = color}></button>
            {/each}
          </div>
        </div>
        <div class="form-control"><label class="label"><span class="label-text">Ikon</span></label><input type="text" class="input input-bordered input-sm" placeholder="📦" bind:value={newCategory.icon} /></div>
        <div class="flex justify-end gap-2 mt-4">
          <button class="btn btn-ghost btn-sm" onclick={() => { showAddModal = false; resetNewCategory(); }}>Batal</button>
          <button class="btn btn-primary btn-sm" onclick={saveCategory}>Simpan</button>
        </div>
      </div>
    </div>
  </dialog>
{/if}
