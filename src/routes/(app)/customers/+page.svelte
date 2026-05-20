<script lang="ts">
  import { mockCustomers } from '$lib/data/mock';
  import { formatCurrency, formatDate } from '$lib/utils/helpers';
  import { Plus, Search, Edit, Eye, Phone, Mail, MapPin } from 'lucide-svelte';

  let customers = $state(mockCustomers);
  let searchQuery = $state('');
  let showAddModal = $state(false);
  let showDetailModal = $state(false);
  let showEditModal = $state(false);
  let viewingCustomer = $state<any>(null);
  let editingCustomer = $state<any>(null);
  let newCustomer = $state({ name: '', phone: '', email: '', address: '', notes: '' });

  function resetNewCustomer() { newCustomer = { name: '', phone: '', email: '', address: '', notes: '' }; }

  function saveCustomer() {
    if (!newCustomer.name) return;
    customers.push({ id: crypto.randomUUID(), ...newCustomer, totalTransactions: 0, totalSpent: 0, createdAt: new Date().toISOString(), lastTransaction: null });
    showAddModal = false;
    resetNewCustomer();
  }

  function updateCustomer() {
    if (!editingCustomer) return;
    const idx = customers.findIndex(c => c.id === editingCustomer.id);
    if (idx >= 0) customers[idx] = { ...editingCustomer };
    showEditModal = false;
    editingCustomer = null;
  }

  let filteredCustomers = $derived(customers.filter(c =>
    !searchQuery || c.name.toLowerCase().includes(searchQuery.toLowerCase()) || (c.phone && c.phone.includes(searchQuery))
  ));
</script>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-2xl font-bold">Manajemen Pelanggan</h1>
      <p class="text-base-content">Kelola data pelanggan</p>
    </div>
    <button class="btn btn-primary btn-sm gap-2" onclick={() => showAddModal = true}>
      <Plus class="w-4 h-4" /> Tambah Pelanggan
    </button>
  </div>

  <div class="glass-card">
    <div class="p-5">
      <div class="flex items-center gap-3 mb-4">
        <div class="form-control flex-1">
          <div class="relative">
            <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content" />
            <input type="text" class="input input-bordered input-sm w-full pl-10" placeholder="Cari pelanggan..." bind:value={searchQuery} />
          </div>
        </div>
      </div>

      <div class="overflow-x-auto">
        <table class="table table-sm">
          <thead>
            <tr>
              <th>Nama</th>
              <th>Telepon</th>
              <th>Email</th>
              <th>Alamat</th>
              <th>Transaksi</th>
              <th>Total Belanja</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            {#each filteredCustomers as customer}
              <tr class="hover">
                <td>
                  <div class="flex items-center gap-3">
                    <div class="avatar placeholder">
                      <div class="bg-neutral text-neutral-content rounded-full w-8">
                        <span class="text-xs">{customer.name.charAt(0)}</span>
                      </div>
                    </div>
                    <span class="font-medium">{customer.name}</span>
                  </div>
                </td>
                <td>{customer.phone || '-'}</td>
                <td>{customer.email || '-'}</td>
                <td class="max-w-xs truncate">{customer.address || '-'}</td>
                <td>{customer.totalTransactions}</td>
                <td class="font-medium">{formatCurrency(customer.totalSpent)}</td>
                <td>
                  <div class="flex gap-1">
                    <button class="btn btn-ghost btn-xs" onclick={() => { viewingCustomer = customer; showDetailModal = true; }}><Eye class="w-4 h-4" /></button>
                    <button class="btn btn-ghost btn-xs" onclick={() => { editingCustomer = customer; showEditModal = true; }}><Edit class="w-4 h-4" /></button>
                  </div>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>
  </div>
</div>

{#if showDetailModal && viewingCustomer}
  <dialog class="modal modal-open" onclick={() => showDetailModal = false}>
    <div class="modal-box" onclick={(e) => e.stopPropagation()}>
      <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2" onclick={() => showDetailModal = false}>✕</button>
      <h3 class="font-bold text-lg mb-4">{viewingCustomer.name}</h3>
      <dl class="space-y-2 text-sm">
        <div class="flex justify-between"><dt class="text-base-content">Telepon</dt><dd>{viewingCustomer.phone || '-'}</dd></div>
        <div class="flex justify-between"><dt class="text-base-content">Email</dt><dd>{viewingCustomer.email || '-'}</dd></div>
        <div class="flex justify-between"><dt class="text-base-content">Alamat</dt><dd>{viewingCustomer.address || '-'}</dd></div>
        <div class="flex justify-between"><dt class="text-base-content">Total Transaksi</dt><dd>{viewingCustomer.totalTransactions}</dd></div>
        <div class="flex justify-between"><dt class="text-base-content">Total Belanja</dt><dd>{formatCurrency(viewingCustomer.totalSpent)}</dd></div>
      </dl>
    </div>
  </dialog>
{/if}

{#if showEditModal && editingCustomer}
  <dialog class="modal modal-open" onclick={() => showEditModal = false}>
    <div class="modal-box" onclick={(e) => e.stopPropagation()}>
      <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2" onclick={() => showEditModal = false}>✕</button>
      <h3 class="font-bold text-lg mb-4">Edit Pelanggan</h3>
      <div class="space-y-3">
        <div class="form-control"><label class="label"><span class="label-text">Nama</span></label><input type="text" class="input input-bordered input-sm" bind:value={editingCustomer.name} /></div>
        <div class="grid grid-cols-2 gap-3">
          <div class="form-control"><label class="label"><span class="label-text">Telepon</span></label><input type="text" class="input input-bordered input-sm" bind:value={editingCustomer.phone} /></div>
          <div class="form-control"><label class="label"><span class="label-text">Email</span></label><input type="email" class="input input-bordered input-sm" bind:value={editingCustomer.email} /></div>
        </div>
        <div class="form-control"><label class="label"><span class="label-text">Alamat</span></label><textarea class="textarea textarea-bordered textarea-sm" bind:value={editingCustomer.address}></textarea></div>
        <div class="flex justify-end gap-2 mt-4">
          <button class="btn btn-ghost btn-sm" onclick={() => showEditModal = false}>Batal</button>
          <button class="btn btn-primary btn-sm" onclick={updateCustomer}>Simpan</button>
        </div>
      </div>
    </div>
  </dialog>
{/if}

{#if showAddModal}
  <dialog class="modal modal-open" onclick={() => showAddModal = false}>
    <div class="modal-box" onclick={(e) => e.stopPropagation()}>
      <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2" onclick={() => { showAddModal = false; resetNewCustomer(); }}>✕</button>
      <h3 class="text-lg font-bold mb-4">Tambah Pelanggan Baru</h3>
      <div class="space-y-3">
        <div class="form-control"><label class="label"><span class="label-text">Nama</span></label><input type="text" class="input input-bordered input-sm w-full" placeholder="Nama lengkap" bind:value={newCustomer.name} /></div>
        <div class="grid grid-cols-2 gap-3">
          <div class="form-control"><label class="label"><span class="label-text">Telepon</span></label><input type="text" class="input input-bordered input-sm" placeholder="08xxxxxxxxxx" bind:value={newCustomer.phone} /></div>
          <div class="form-control"><label class="label"><span class="label-text">Email</span></label><input type="email" class="input input-bordered input-sm" placeholder="email@contoh.com" bind:value={newCustomer.email} /></div>
        </div>
        <div class="form-control"><label class="label"><span class="label-text">Alamat</span></label><textarea class="textarea textarea-bordered textarea-sm" rows="2" placeholder="Alamat lengkap" bind:value={newCustomer.address}></textarea></div>
        <div class="form-control"><label class="label"><span class="label-text">Catatan</span></label><textarea class="textarea textarea-bordered textarea-sm" rows="2" placeholder="Catatan tambahan" bind:value={newCustomer.notes}></textarea></div>
        <div class="flex justify-end gap-2 mt-4">
          <button class="btn btn-ghost btn-sm" onclick={() => { showAddModal = false; resetNewCustomer(); }}>Batal</button>
          <button class="btn btn-primary btn-sm" onclick={saveCustomer}>Simpan</button>
        </div>
      </div>
    </div>
  </dialog>
{/if}
