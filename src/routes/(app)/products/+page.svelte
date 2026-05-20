<script lang="ts">
  import { mockProducts, mockCategories, mockBrands } from '$lib/data/mock';
  import { formatCurrency } from '$lib/utils/helpers';
  import { Plus, Search, Filter, Edit, Trash2, Archive, Eye } from 'lucide-svelte';

  let products = $state(mockProducts);
  let searchQuery = $state('');
  let selectedCategory = $state('all');
  let selectedBrand = $state('all');
  let showAddModal = $state(false);
  let showEditModal = $state(false);
  let editingProduct = $state<any>(null);
  let showDetailModal = $state(false);
  let viewingProduct = $state<any>(null);
  let currentPage = $state(1);
  const perPage = 10;

  let newProduct = $state({ name: '', sku: '', barcode: '', categoryId: '', brandId: '', type: 'physical', costPrice: 0, sellingPrice: 0, stock: 0, minStock: 0 });

  function resetNewProduct() {
    newProduct = { name: '', sku: '', barcode: '', categoryId: '', brandId: '', type: 'physical', costPrice: 0, sellingPrice: 0, stock: 0, minStock: 0 };
  }

  function saveProduct() {
    const item = { id: crypto.randomUUID(), ...newProduct, isActive: true, createdAt: new Date().toISOString() };
    products.push(item);
    showAddModal = false;
    resetNewProduct();
  }

  function updateProduct() {
    if (!editingProduct) return;
    const idx = products.findIndex(p => p.id === editingProduct.id);
    if (idx >= 0) products[idx] = { ...editingProduct };
    showEditModal = false;
    editingProduct = null;
  }

  function archiveProduct(id: string) {
    const idx = products.findIndex(p => p.id === id);
    if (idx >= 0) { products[idx] = { ...products[idx], isActive: !products[idx].isActive }; }
  }

  let filteredProducts = $derived.by(() => {
    let result = products;
    if (searchQuery) {
      const q = searchQuery.toLowerCase();
      result = result.filter(p => p.name.toLowerCase().includes(q) || p.sku.toLowerCase().includes(q));
    }
    if (selectedCategory !== 'all') {
      result = result.filter(p => p.categoryId === selectedCategory);
    }
    if (selectedBrand !== 'all') {
      result = result.filter(p => p.brandId === selectedBrand);
    }
    return result;
  });

  let totalPages = $derived(Math.ceil(filteredProducts.length / perPage));
  let paginatedProducts = $derived(filteredProducts.slice((currentPage - 1) * perPage, currentPage * perPage));

  function getCategoryName(id: string) {
    return mockCategories.find(c => c.id === id)?.name || '-';
  }

  function getBrandName(id?: string) {
    if (!id) return '-';
    return mockBrands.find(b => b.id === id)?.name || '-';
  }

  function getProductType(type: string) {
    const types: Record<string, string> = { physical: 'Barang', service: 'Jasa', sparepart: 'Sparepart', bundle: 'Bundle' };
    return types[type] || type;
  }
</script>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-2xl font-bold">Manajemen Produk</h1>
      <p class="text-base-content">Kelola produk dan jasa</p>
    </div>
    <button class="btn btn-primary btn-sm gap-2" onclick={() => showAddModal = true}>
      <Plus class="w-4 h-4" /> Tambah Produk
    </button>
  </div>

  <div class="glass-card">
    <div class="p-5">
      <div class="flex flex-wrap gap-3 mb-4">
        <div class="form-control flex-1 min-w-[200px]">
          <div class="relative">
            <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content" />
            <input type="text" class="input input-bordered input-sm w-full pl-10" placeholder="Cari produk..." bind:value={searchQuery} />
          </div>
        </div>
        <select class="select select-bordered select-sm" bind:value={selectedCategory}>
          <option value="all">Semua Kategori</option>
          {#each mockCategories as cat}
            <option value={cat.id}>{cat.name}</option>
          {/each}
        </select>
        <select class="select select-bordered select-sm" bind:value={selectedBrand}>
          <option value="all">Semua Brand</option>
          {#each mockBrands as brand}
            <option value={brand.id}>{brand.name}</option>
          {/each}
        </select>
      </div>

      <div class="overflow-x-auto">
        <table class="table table-sm">
          <thead>
            <tr>
              <th>SKU</th>
              <th>Nama Produk</th>
              <th>Kategori</th>
              <th>Brand</th>
              <th>Tipe</th>
              <th>Harga Modal</th>
              <th>Harga Jual</th>
              <th>Stok</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            {#each paginatedProducts as product}
              <tr class="hover">
                <td class="font-mono text-xs">{product.sku}</td>
                <td>
                  <div class="font-medium">{product.name}</div>
                  {#if product.barcode}
                    <div class="text-xs text-base-content">{product.barcode}</div>
                  {/if}
                </td>
                <td>{getCategoryName(product.categoryId)}</td>
                <td>{getBrandName(product.brandId)}</td>
                <td><span class="badge badge-sm badge-ghost">{getProductType(product.type)}</span></td>
                <td>{formatCurrency(product.costPrice)}</td>
                <td class="font-medium">{formatCurrency(product.sellingPrice)}</td>
                <td>
                  {#if product.stock <= product.minStock}
                    <span class="badge badge-sm badge-error">{product.stock}</span>
                  {:else}
                    <span class="badge badge-sm badge-success">{product.stock}</span>
                  {/if}
                </td>
                <td>
                  <div class="flex gap-1">
                    <button class="btn btn-ghost btn-xs" title="Lihat" onclick={() => { viewingProduct = product; showDetailModal = true; }}><Eye class="w-4 h-4" /></button>
                    <button class="btn btn-ghost btn-xs" title="Edit" onclick={() => { editingProduct = { ...product }; showEditModal = true; }}><Edit class="w-4 h-4" /></button>
                    <button class="btn btn-ghost btn-xs text-error" title="Arsip" onclick={() => archiveProduct(product.id)}><Archive class="w-4 h-4" /></button>
                  </div>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>

      {#if totalPages > 1}
        <div class="flex items-center justify-between mt-4">
          <p class="text-sm text-base-content">Menampilkan {(currentPage - 1) * perPage + 1}-{Math.min(currentPage * perPage, filteredProducts.length)} dari {filteredProducts.length} produk</p>
          <div class="join">
            {#each Array(totalPages) as _, i}
              <button class="join-item btn btn-sm {currentPage === i + 1 ? 'btn-primary' : ''}" onclick={() => currentPage = i + 1}>{i + 1}</button>
            {/each}
          </div>
        </div>
      {/if}
    </div>
  </div>
</div>

{#if showDetailModal && viewingProduct}
  <dialog class="modal modal-open" onclick={() => showDetailModal = false}>
    <div class="modal-box" onclick={(e) => e.stopPropagation()}>
      <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2" onclick={() => showDetailModal = false}>✕</button>
      <h3 class="font-bold text-lg mb-4">{viewingProduct.name}</h3>
      <dl class="space-y-2 text-sm">
        <div class="flex justify-between"><dt class="text-base-content">SKU</dt><dd class="font-mono">{viewingProduct.sku}</dd></div>
        <div class="flex justify-between"><dt class="text-base-content">Barcode</dt><dd>{viewingProduct.barcode || '-'}</dd></div>
        <div class="flex justify-between"><dt class="text-base-content">Kategori</dt><dd>{getCategoryName(viewingProduct.categoryId)}</dd></div>
        <div class="flex justify-between"><dt class="text-base-content">Brand</dt><dd>{getBrandName(viewingProduct.brandId)}</dd></div>
        <div class="flex justify-between"><dt class="text-base-content">Tipe</dt><dd>{getProductType(viewingProduct.type)}</dd></div>
        <div class="flex justify-between"><dt class="text-base-content">Harga Modal</dt><dd>{formatCurrency(viewingProduct.costPrice)}</dd></div>
        <div class="flex justify-between"><dt class="text-base-content">Harga Jual</dt><dd>{formatCurrency(viewingProduct.sellingPrice)}</dd></div>
        <div class="flex justify-between"><dt class="text-base-content">Stok</dt><dd>{viewingProduct.stock}</dd></div>
        <div class="flex justify-between"><dt class="text-base-content">Min. Stok</dt><dd>{viewingProduct.minStock}</dd></div>
        <div class="flex justify-between"><dt class="text-base-content">Status</dt><dd><span class="badge badge-sm {viewingProduct.isActive ? 'badge-success' : 'badge-ghost'}">{viewingProduct.isActive ? 'Aktif' : 'Nonaktif'}</span></dd></div>
      </dl>
    </div>
  </dialog>
{/if}

{#if showEditModal && editingProduct}
  <dialog class="modal modal-open" onclick={() => showEditModal = false}>
    <div class="modal-box" onclick={(e) => e.stopPropagation()}>
      <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2" onclick={() => showEditModal = false}>✕</button>
      <h3 class="font-bold text-lg mb-4">Edit Produk</h3>
      <div class="space-y-3">
        <div class="form-control"><label class="label"><span class="label-text">Nama Produk</span></label><input type="text" class="input input-bordered input-sm" bind:value={editingProduct.name} /></div>
        <div class="grid grid-cols-2 gap-3">
          <div class="form-control"><label class="label"><span class="label-text">SKU</span></label><input type="text" class="input input-bordered input-sm" bind:value={editingProduct.sku} /></div>
          <div class="form-control"><label class="label"><span class="label-text">Barcode</span></label><input type="text" class="input input-bordered input-sm" bind:value={editingProduct.barcode} /></div>
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div class="form-control"><label class="label"><span class="label-text">Kategori</span></label>
            <select class="select select-bordered select-sm" bind:value={editingProduct.categoryId}>
              <option value="">Pilih</option>
              {#each mockCategories as c}<option value={c.id}>{c.name}</option>{/each}
            </select>
          </div>
          <div class="form-control"><label class="label"><span class="label-text">Brand</span></label>
            <select class="select select-bordered select-sm" bind:value={editingProduct.brandId}>
              <option value="">Pilih</option>
              {#each mockBrands as b}<option value={b.id}>{b.name}</option>{/each}
            </select>
          </div>
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div class="form-control"><label class="label"><span class="label-text">Tipe</span></label>
            <select class="select select-bordered select-sm" bind:value={editingProduct.type}>
              <option value="physical">Barang</option><option value="service">Jasa</option><option value="sparepart">Sparepart</option><option value="bundle">Bundle</option>
            </select>
          </div>
          <div></div>
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div class="form-control"><label class="label"><span class="label-text">Harga Modal</span></label><input type="number" class="input input-bordered input-sm" bind:value={editingProduct.costPrice} /></div>
          <div class="form-control"><label class="label"><span class="label-text">Harga Jual</span></label><input type="number" class="input input-bordered input-sm" bind:value={editingProduct.sellingPrice} /></div>
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div class="form-control"><label class="label"><span class="label-text">Stok</span></label><input type="number" class="input input-bordered input-sm" bind:value={editingProduct.stock} /></div>
          <div class="form-control"><label class="label"><span class="label-text">Min. Stok</span></label><input type="number" class="input input-bordered input-sm" bind:value={editingProduct.minStock} /></div>
        </div>
        <div class="flex justify-end gap-2 mt-4">
          <button class="btn btn-ghost btn-sm" onclick={() => showEditModal = false}>Batal</button>
          <button class="btn btn-primary btn-sm" onclick={updateProduct}>Simpan</button>
        </div>
      </div>
    </div>
  </dialog>
{/if}

{#if showAddModal}
  <dialog class="modal modal-open" onclick={() => showAddModal = false}>
    <div class="modal-box" onclick={(e) => e.stopPropagation()}>
      <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2" onclick={() => { showAddModal = false; resetNewProduct(); }}>✕</button>
      <h3 class="font-bold text-lg mb-4">Tambah Produk Baru</h3>
      <div class="space-y-3">
        <div class="form-control"><label class="label"><span class="label-text">Nama Produk</span></label><input type="text" class="input input-bordered input-sm" placeholder="Nama produk" bind:value={newProduct.name} /></div>
        <div class="grid grid-cols-2 gap-3">
          <div class="form-control"><label class="label"><span class="label-text">SKU</span></label><input type="text" class="input input-bordered input-sm" placeholder="SKU" bind:value={newProduct.sku} /></div>
          <div class="form-control"><label class="label"><span class="label-text">Barcode</span></label><input type="text" class="input input-bordered input-sm" placeholder="Barcode" bind:value={newProduct.barcode} /></div>
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div class="form-control"><label class="label"><span class="label-text">Kategori</span></label>
            <select class="select select-bordered select-sm" bind:value={newProduct.categoryId}>
              <option value="">Pilih Kategori</option>
              {#each mockCategories as c}<option value={c.id}>{c.name}</option>{/each}
            </select>
          </div>
          <div class="form-control"><label class="label"><span class="label-text">Brand</span></label>
            <select class="select select-bordered select-sm" bind:value={newProduct.brandId}>
              <option value="">Pilih Brand</option>
              {#each mockBrands as b}<option value={b.id}>{b.name}</option>{/each}
            </select>
          </div>
        </div>
        <div class="form-control"><label class="label"><span class="label-text">Tipe</span></label>
          <select class="select select-bordered select-sm" bind:value={newProduct.type}>
            <option value="physical">Barang</option><option value="service">Jasa</option><option value="sparepart">Sparepart</option><option value="bundle">Bundle</option>
          </select>
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div class="form-control"><label class="label"><span class="label-text">Harga Modal</span></label><input type="number" class="input input-bordered input-sm" bind:value={newProduct.costPrice} /></div>
          <div class="form-control"><label class="label"><span class="label-text">Harga Jual</span></label><input type="number" class="input input-bordered input-sm" bind:value={newProduct.sellingPrice} /></div>
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div class="form-control"><label class="label"><span class="label-text">Stok Awal</span></label><input type="number" class="input input-bordered input-sm" bind:value={newProduct.stock} /></div>
          <div class="form-control"><label class="label"><span class="label-text">Min. Stok</span></label><input type="number" class="input input-bordered input-sm" bind:value={newProduct.minStock} /></div>
        </div>
        <div class="flex justify-end gap-2 mt-4">
          <button class="btn btn-ghost btn-sm" onclick={() => { showAddModal = false; resetNewProduct(); }}>Batal</button>
          <button class="btn btn-primary btn-sm" onclick={saveProduct}>Simpan</button>
        </div>
      </div>
    </div>
  </dialog>
{/if}
