<script lang="ts">
  import { mockSuppliers } from '$lib/data/mock';
  import { Plus, Edit } from 'lucide-svelte';

  let suppliers = $state(mockSuppliers);
  let showAddModal = $state(false);
  let showEditModal = $state(false);
  let editingSupplier = $state<any>(null);
  let newSupplier = $state({ name: '', phone: '', email: '', address: '' });

  function resetNewSupplier() { newSupplier = { name: '', phone: '', email: '', address: '' }; }

  function saveSupplier() {
    if (!newSupplier.name) return;
    suppliers.push({ id: crypto.randomUUID(), ...newSupplier, isActive: true });
    showAddModal = false;
    resetNewSupplier();
  }

  function editSupplier(supplier: any) {
    editingSupplier = { ...supplier };
    showEditModal = true;
  }

  function updateSupplier() {
    if (!editingSupplier) return;
    const idx = suppliers.findIndex(s => s.id === editingSupplier.id);
    if (idx >= 0) suppliers[idx] = { ...editingSupplier };
    showEditModal = false;
    editingSupplier = null;
  }
</script>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-2xl font-bold">Supplier</h1>
      <p class="text-base-content">Kelola data supplier</p>
    </div>
    <button class="btn btn-primary btn-sm gap-2" onclick={() => showAddModal = true}>
      <Plus class="w-4 h-4" /> Tambah Supplier
    </button>
  </div>

  <div class="glass-card">
    <div class="p-5">
      <div class="overflow-x-auto">
        <table class="table table-sm">
          <thead>
            <tr>
              <th>Nama Supplier</th>
              <th>Telepon</th>
              <th>Email</th>
              <th>Alamat</th>
              <th>Status</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            {#each suppliers as supplier}
              <tr class="hover">
                <td class="font-medium">{supplier.name}</td>
                <td>{supplier.phone || '-'}</td>
                <td>{supplier.email || '-'}</td>
                <td>{supplier.address || '-'}</td>
                <td><span class="badge badge-sm {supplier.isActive ? 'badge-success' : 'badge-ghost'}">{supplier.isActive ? 'Aktif' : 'Nonaktif'}</span></td>
                <td><button class="btn btn-ghost btn-xs" onclick={() => editSupplier(supplier)}><Edit class="w-4 h-4" /></button></td>
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
      <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2" onclick={() => { showAddModal = false; resetNewSupplier(); }}>✕</button>
      <h3 class="font-bold text-lg mb-4">Tambah Supplier</h3>
      <div class="space-y-3">
        <div class="form-control"><label class="label"><span class="label-text">Nama Supplier</span></label><input type="text" class="input input-bordered input-sm" placeholder="Nama" bind:value={newSupplier.name} /></div>
        <div class="grid grid-cols-2 gap-3">
          <div class="form-control"><label class="label"><span class="label-text">Telepon</span></label><input type="text" class="input input-bordered input-sm" placeholder="Telepon" bind:value={newSupplier.phone} /></div>
          <div class="form-control"><label class="label"><span class="label-text">Email</span></label><input type="email" class="input input-bordered input-sm" placeholder="Email" bind:value={newSupplier.email} /></div>
        </div>
        <div class="form-control"><label class="label"><span class="label-text">Alamat</span></label><textarea class="textarea textarea-bordered textarea-sm" placeholder="Alamat" bind:value={newSupplier.address}></textarea></div>
        <div class="flex justify-end gap-2 mt-4">
          <button class="btn btn-ghost btn-sm" onclick={() => { showAddModal = false; resetNewSupplier(); }}>Batal</button>
          <button class="btn btn-primary btn-sm" onclick={saveSupplier}>Simpan</button>
        </div>
      </div>
    </div>
  </dialog>
{/if}

{#if showEditModal && editingSupplier}
  <dialog class="modal modal-open" onclick={() => showEditModal = false}>
    <div class="modal-box" onclick={(e) => e.stopPropagation()}>
      <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2" onclick={() => showEditModal = false}>✕</button>
      <h3 class="font-bold text-lg mb-4">Edit Supplier</h3>
      <div class="space-y-3">
        <div class="form-control"><label class="label"><span class="label-text">Nama Supplier</span></label><input type="text" class="input input-bordered input-sm" bind:value={editingSupplier.name} /></div>
        <div class="grid grid-cols-2 gap-3">
          <div class="form-control"><label class="label"><span class="label-text">Telepon</span></label><input type="text" class="input input-bordered input-sm" bind:value={editingSupplier.phone} /></div>
          <div class="form-control"><label class="label"><span class="label-text">Email</span></label><input type="email" class="input input-bordered input-sm" bind:value={editingSupplier.email} /></div>
        </div>
        <div class="form-control"><label class="label"><span class="label-text">Alamat</span></label><textarea class="textarea textarea-bordered textarea-sm" bind:value={editingSupplier.address}></textarea></div>
        <div class="flex justify-end gap-2 mt-4">
          <button class="btn btn-ghost btn-sm" onclick={() => showEditModal = false}>Batal</button>
          <button class="btn btn-primary btn-sm" onclick={updateSupplier}>Simpan</button>
        </div>
      </div>
    </div>
  </dialog>
{/if}
