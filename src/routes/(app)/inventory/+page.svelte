<script lang="ts">
  import { mockProducts, mockStockMovements } from '$lib/data/mock';
  import { formatCurrency, formatDate, formatDateTime } from '$lib/utils/helpers';
  import { ArrowDownLeft, ArrowUpRight, Plus, Search } from 'lucide-svelte';

  let activeTab = $state('stock');
  let searchQuery = $state('');
  let showStockModal = $state(false);
  let stockType = $state<'in' | 'out'>('in');
  let stockForm = $state({ productId: '', quantity: 1, notes: '' });

  function openStockModal(type: 'in' | 'out') {
    stockType = type;
    stockForm = { productId: '', quantity: 1, notes: '' };
    showStockModal = true;
  }

  function processStockMovement() {
    showStockModal = false;
  }

  let lowStockProducts = $derived(mockProducts.filter(p => p.stock <= p.minStock));

  function getMovementTypeIcon(type: string) {
    switch (type) {
      case 'in': return ArrowDownLeft;
      case 'out': return ArrowUpRight;
      default: return Plus;
    }
  }

  function getMovementTypeColor(type: string) {
    switch (type) {
      case 'in': return 'text-success';
      case 'out': return 'text-error';
      default: return 'text-warning';
    }
  }

  function getMovementTypeLabel(type: string) {
    const labels: Record<string, string> = { in: 'Stok Masuk', out: 'Stok Keluar', adjustment: 'Penyesuaian', transfer: 'Transfer' };
    return labels[type] || type;
  }
</script>

<div class="space-y-6">
  <div>
    <h1 class="text-2xl font-bold">Inventory Management</h1>
    <p class="text-base-content">Kelola stok dan mutasi barang</p>
  </div>

  <div class="tabs tabs-boxed bg-base-200 w-fit">
    <button class="tab {activeTab === 'stock' ? 'tab-active' : ''}" onclick={() => activeTab = 'stock'}>Stok Produk</button>
    <button class="tab {activeTab === 'low' ? 'tab-active' : ''}" onclick={() => activeTab = 'low'}>Stok Menipis ({lowStockProducts.length})</button>
    <button class="tab {activeTab === 'history' ? 'tab-active' : ''}" onclick={() => activeTab = 'history'}>Riwayat Stok</button>
  </div>

  {#if activeTab === 'stock'}
    <div class="glass-card">
      <div class="p-5">
        <div class="flex items-center gap-3 mb-4">
          <div class="form-control flex-1">
            <div class="relative">
              <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content" />
              <input type="text" class="input input-bordered input-sm w-full pl-10" placeholder="Cari produk..." bind:value={searchQuery} />
            </div>
          </div>
          <button class="btn btn-primary btn-sm" onclick={() => openStockModal('in')}>Stok Masuk</button>
          <button class="btn btn-outline btn-sm" onclick={() => openStockModal('out')}>Stok Keluar</button>
        </div>

        <div class="overflow-x-auto">
          <table class="table table-sm">
            <thead>
              <tr>
                <th>SKU</th>
                <th>Nama Produk</th>
                <th>Stok Saat Ini</th>
                <th>Min. Stok</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              {#each mockProducts.filter(p => !searchQuery || p.name.toLowerCase().includes(searchQuery.toLowerCase())) as product}
                <tr class="hover">
                  <td class="font-mono text-xs">{product.sku}</td>
                  <td>{product.name}</td>
                  <td class="font-medium">{product.stock}</td>
                  <td>{product.minStock}</td>
                  <td>
                    {#if product.stock === 0}
                      <span class="badge badge-sm badge-error">Habis</span>
                    {:else if product.stock <= product.minStock}
                      <span class="badge badge-sm badge-warning">Menipis</span>
                    {:else}
                      <span class="badge badge-sm badge-success">Aman</span>
                    {/if}
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  {:else if activeTab === 'low'}
    <div class="glass-card">
      <div class="p-5">
        {#if lowStockProducts.length === 0}
          <div class="text-center py-8 text-base-content">
            <p>Semua stok aman</p>
          </div>
        {:else}
          <div class="overflow-x-auto">
            <table class="table table-sm">
              <thead>
                <tr>
                  <th>SKU</th>
                  <th>Nama Produk</th>
                  <th>Stok</th>
                  <th>Min. Stok</th>
                  <th>Kekurangan</th>
                </tr>
              </thead>
              <tbody>
                {#each lowStockProducts as product}
                  <tr class="hover">
                    <td class="font-mono text-xs">{product.sku}</td>
                    <td>{product.name}</td>
                    <td><span class="text-error font-bold">{product.stock}</span></td>
                    <td>{product.minStock}</td>
                    <td>{product.minStock - product.stock}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        {/if}
      </div>
    </div>
  {:else}
    <div class="glass-card">
      <div class="p-5">
        <div class="overflow-x-auto">
          <table class="table table-sm">
            <thead>
              <tr>
                <th>Tanggal</th>
                <th>Produk</th>
                <th>Tipe</th>
                <th>Qty</th>
                <th>Stok Sebelumnya</th>
                <th>Stok Baru</th>
                <th>Keterangan</th>
                <th>Oleh</th>
              </tr>
            </thead>
            <tbody>
              {#each mockStockMovements as movement}
                <tr class="hover">
                  <td class="text-xs">{formatDateTime(movement.createdAt)}</td>
                  <td>{movement.productName}</td>
                  <td>
                    <span class="flex items-center gap-1 {getMovementTypeColor(movement.type)}">
                      <svelte:component this={getMovementTypeIcon(movement.type)} class="w-4 h-4" />
                      {getMovementTypeLabel(movement.type)}
                    </span>
                  </td>
                  <td class="font-medium">{movement.quantity > 0 ? '+' : ''}{movement.quantity}</td>
                  <td>{movement.previousStock}</td>
                  <td>{movement.newStock}</td>
                  <td class="text-xs">{movement.notes || '-'}</td>
                  <td class="text-xs">{movement.createdBy}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  {/if}
</div>

{#if showStockModal}
  <dialog class="modal modal-open" onclick={() => showStockModal = false}>
    <div class="modal-box" onclick={(e) => e.stopPropagation()}>
      <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2" onclick={() => showStockModal = false}>✕</button>
      <h3 class="font-bold text-lg mb-4">{stockType === 'in' ? 'Stok Masuk' : 'Stok Keluar'}</h3>
      <div class="space-y-3">
        <div class="form-control"><label class="label"><span class="label-text">Produk</span></label>
          <select class="select select-bordered select-sm w-full" bind:value={stockForm.productId}>
            <option value="">Pilih Produk</option>
            {#each mockProducts as p}
              <option value={p.id}>{p.name} (Stok: {p.stock})</option>
            {/each}
          </select>
        </div>
        <div class="form-control"><label class="label"><span class="label-text">Jumlah</span></label><input type="number" class="input input-bordered input-sm" bind:value={stockForm.quantity} min="1" /></div>
        <div class="form-control"><label class="label"><span class="label-text">Keterangan</span></label><textarea class="textarea textarea-bordered textarea-sm" placeholder="Keterangan" bind:value={stockForm.notes}></textarea></div>
        <div class="flex justify-end gap-2 mt-4">
          <button class="btn btn-ghost btn-sm" onclick={() => showStockModal = false}>Batal</button>
          <button class="btn btn-primary btn-sm" onclick={processStockMovement}>Proses</button>
        </div>
      </div>
    </div>
  </dialog>
{/if}
