<script lang="ts">
  import { mockServiceTickets, mockCustomers } from '$lib/data/mock';
  import { formatCurrency, formatDate, getServiceStatusLabel, getServiceStatusColor, generateServiceNumber } from '$lib/utils/helpers';
  import { Plus, Search, Eye, Edit, MessageSquare, Wrench } from 'lucide-svelte';
  import { currentUser } from '$lib/stores/auth';

  let tickets = $state(mockServiceTickets);
  let searchQuery = $state('');
  let selectedStatus = $state('all');
  let showAddModal = $state(false);
  let selectedTicket = $state<any>(null);
  let showDetailModal = $state(false);
  let newNote = $state('');

  let filteredTickets = $derived(tickets.filter(t => {
    const matchSearch = !searchQuery || t.customerName.toLowerCase().includes(searchQuery.toLowerCase()) || t.ticketNumber.toLowerCase().includes(searchQuery.toLowerCase());
    const matchStatus = selectedStatus === 'all' || t.status === selectedStatus;
    return matchSearch && matchStatus;
  }));

  const statusOptions = [
    { id: 'all', label: 'Semua' },
    { id: 'waiting', label: 'Menunggu' },
    { id: 'checked', label: 'Dicek' },
    { id: 'in_progress', label: 'Dikerjakan' },
    { id: 'waiting_parts', label: 'Tunggu Part' },
    { id: 'completed', label: 'Selesai' },
    { id: 'picked_up', label: 'Diambil' },
    { id: 'cancelled', label: 'Batal' },
  ];

  function addNote() {
    if (!newNote || !selectedTicket) return;
    selectedTicket.technicianNotes.push({
      id: Date.now().toString(),
      note: newNote,
      createdAt: new Date().toISOString(),
      createdBy: $currentUser?.name || 'Teknisi',
    });
    newNote = '';
  }
</script>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-2xl font-bold">Service Management</h1>
      <p class="text-base-content">Kelola tiket service dan perbaikan</p>
    </div>
    <button class="btn btn-primary btn-sm gap-2" onclick={() => showAddModal = true}>
      <Plus class="w-4 h-4" /> Tiket Baru
    </button>
  </div>

  <div class="flex items-center gap-3">
    <div class="form-control flex-1">
      <div class="relative">
        <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content" />
        <input type="text" class="input input-bordered input-sm w-full pl-10" placeholder="Cari tiket atau pelanggan..." bind:value={searchQuery} />
      </div>
    </div>
    <div class="flex gap-1 overflow-x-auto">
      {#each statusOptions as status}
        <button class="btn btn-sm {selectedStatus === status.id ? 'btn-primary' : 'btn-outline'}" onclick={() => selectedStatus = status.id}>{status.label}</button>
      {/each}
    </div>
  </div>

  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
    {#each filteredTickets as ticket}
      <div class="glass-card cursor-pointer hover:shadow-md transition-shadow" onclick={() => { selectedTicket = ticket; showDetailModal = true; }}>
        <div class="p-5 p-4">
          <div class="flex items-start justify-between mb-2">
            <span class="font-mono text-xs">{ticket.ticketNumber}</span>
            <span class="badge badge-sm {getServiceStatusColor(ticket.status)}">{getServiceStatusLabel(ticket.status)}</span>
          </div>
          <h3 class="font-semibold">{ticket.customerName}</h3>
          <p class="text-sm text-base-content">{ticket.deviceType} {ticket.deviceBrand || ''}</p>
          <p class="text-xs text-base-content line-clamp-2">{ticket.complaint}</p>
          <div class="flex items-center justify-between mt-2 pt-2 border-t border-base-200">
            <span class="text-xs text-base-content">{formatDate(ticket.createdAt)}</span>
            {#if ticket.estimatedCost}
              <span class="text-sm font-medium">{formatCurrency(ticket.estimatedCost)}</span>
            {/if}
          </div>
        </div>
      </div>
    {/each}
  </div>
</div>

{#if showDetailModal && selectedTicket}
  <div class="modal modal-open">
    <div class="modal-box w-11/12 max-w-3xl">
      <div class="flex items-start justify-between mb-4">
        <div>
          <h3 class="text-lg font-bold">{selectedTicket.ticketNumber}</h3>
          <p class="text-sm text-base-content">{selectedTicket.customerName} - {selectedTicket.customerPhone}</p>
        </div>
        <button class="btn btn-ghost btn-sm btn-circle" onclick={() => showDetailModal = false}>✕</button>
      </div>

      <div class="grid grid-cols-2 gap-4 mb-4">
        <div>
          <label class="label"><span class="label-text font-medium">Device</span></label>
          <p>{selectedTicket.deviceType} {selectedTicket.deviceBrand || ''}</p>
        </div>
        <div>
          <label class="label"><span class="label-text font-medium">Serial Number</span></label>
          <p>{selectedTicket.serialNumber || '-'}</p>
        </div>
        <div>
          <label class="label"><span class="label-text font-medium">Keluhan</span></label>
          <p>{selectedTicket.complaint}</p>
        </div>
        <div>
          <label class="label"><span class="label-text font-medium">Diagnosa</span></label>
          <p>{selectedTicket.diagnosis || 'Belum didiagnosa'}</p>
        </div>
        <div>
          <label class="label"><span class="label-text font-medium">Estimasi Biaya</span></label>
          <p class="font-bold">{formatCurrency(selectedTicket.estimatedCost)}</p>
        </div>
        <div>
          <label class="label"><span class="label-text font-medium">Status</span></label>
          <select class="select select-bordered select-sm w-full" bind:value={selectedTicket.status}>
            {#each statusOptions.filter(s => s.id !== 'all') as status}
              <option value={status.id}>{status.label}</option>
            {/each}
          </select>
        </div>
      </div>

      <div class="divider">Timeline</div>

      <div class="space-y-3 max-h-48 overflow-y-auto mb-4">
        {#each selectedTicket.technicianNotes as note}
          <div class="bg-base-200 rounded-lg p-3">
            <p class="text-sm">{note.note}</p>
            <p class="text-xs text-base-content mt-1">{note.createdBy} - {formatDate(note.createdAt)}</p>
          </div>
        {/each}
      </div>

      <div class="flex gap-2">
        <input type="text" class="input input-bordered input-sm flex-1" placeholder="Tambah catatan..." bind:value={newNote} onkeydown={(e) => { if (e.key === 'Enter') addNote(); }} />
        <button class="btn btn-primary btn-sm" onclick={addNote}>Tambah</button>
      </div>

      <div class="modal-action">
        <button class="btn btn-ghost" onclick={() => showDetailModal = false}>Tutup</button>
        <button class="btn btn-outline btn-sm gap-2" onclick={() => window.print()}><Wrench class="w-4 h-4" /> Cetak Nota</button>
      </div>
    </div>
  </div>
{/if}

{#if showAddModal}
  <div class="modal modal-open">
    <div class="modal-box w-11/12 max-w-2xl">
      <h3 class="text-lg font-bold mb-4">Tiket Service Baru</h3>
      <div class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
          <div class="form-control">
            <label class="label"><span class="label-text font-medium">Pelanggan</span></label>
            <select class="select select-bordered w-full">
              <option value="">Pilih pelanggan</option>
              {#each mockCustomers as c}
                <option value={c.id}>{c.name}</option>
              {/each}
            </select>
          </div>
          <div class="form-control">
            <label class="label"><span class="label-text font-medium">No. HP</span></label>
            <input type="text" class="input input-bordered w-full" placeholder="08xxxxxxxxxx" />
          </div>
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div class="form-control">
            <label class="label"><span class="label-text font-medium">Jenis Device</span></label>
            <select class="select select-bordered w-full">
              <option>Laptop</option>
              <option>PC</option>
              <option>Printer</option>
              <option>Monitor</option>
              <option>Lainnya</option>
            </select>
          </div>
          <div class="form-control">
            <label class="label"><span class="label-text font-medium">Merek</span></label>
            <input type="text" class="input input-bordered w-full" placeholder="ASUS, Lenovo, dll" />
          </div>
        </div>
        <div class="form-control">
          <label class="label"><span class="label-text font-medium">Keluhan</span></label>
          <textarea class="textarea textarea-bordered w-full" rows="3" placeholder="Deskripsikan keluhan pelanggan"></textarea>
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div class="form-control">
            <label class="label"><span class="label-text font-medium">Estimasi Biaya</span></label>
            <input type="number" class="input input-bordered w-full" placeholder="0" />
          </div>
          <div class="form-control">
            <label class="label"><span class="label-text font-medium">Estimasi Selesai</span></label>
            <input type="date" class="input input-bordered w-full" />
          </div>
        </div>
      </div>
      <div class="modal-action">
        <button class="btn btn-ghost" onclick={() => showAddModal = false}>Batal</button>
        <button class="btn btn-primary" onclick={() => showAddModal = false}>Buat Tiket</button>
      </div>
    </div>
  </div>
{/if}
