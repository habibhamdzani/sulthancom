<script lang="ts">
  import { mockReceivables } from '$lib/data/mock';
  import { formatCurrency, formatDate, getPaymentStatusLabel, getPaymentStatusColor } from '$lib/utils/helpers';
  import { Search, Eye, DollarSign } from 'lucide-svelte';

  let receivables = $state(mockReceivables);
  let searchQuery = $state('');
  let selectedStatus = $state('all');
  let showPaymentModal = $state(false);
  let selectedReceivable = $state<any>(null);
  let paymentAmount = $state('');
  let paymentMethod = $state('cash');

  let filteredReceivables = $derived(receivables.filter(r => {
    const matchSearch = !searchQuery || r.customerName.toLowerCase().includes(searchQuery.toLowerCase()) || r.invoiceNumber.toLowerCase().includes(searchQuery.toLowerCase());
    const matchStatus = selectedStatus === 'all' || r.status === selectedStatus;
    return matchSearch && matchStatus;
  }));

  const statusOptions = [
    { id: 'all', label: 'Semua' },
    { id: 'unpaid', label: 'Belum Dibayar' },
    { id: 'partial', label: 'Sebagian' },
    { id: 'paid', label: 'Lunas' },
    { id: 'overdue', label: 'Lewat Tempo' },
  ];

  function getStatusBadge(status: string) {
    const colors: Record<string, string> = { unpaid: 'badge-error', partial: 'badge-warning', paid: 'badge-success', overdue: 'badge-error' };
    const labels: Record<string, string> = { unpaid: 'Belum Dibayar', partial: 'Sebagian', paid: 'Lunas', overdue: 'Lewat Tempo' };
    return { color: colors[status] || 'badge-ghost', label: labels[status] || status };
  }

  function makePayment() {
    if (!selectedReceivable || !paymentAmount) return;
    selectedReceivable.payments.push({
      id: Date.now().toString(),
      amount: parseInt(paymentAmount),
      method: paymentMethod,
      date: new Date().toISOString(),
    });
    selectedReceivable.paidAmount += parseInt(paymentAmount);
    selectedReceivable.remainingAmount -= parseInt(paymentAmount);
    if (selectedReceivable.remainingAmount <= 0) {
      selectedReceivable.status = 'paid';
    } else {
      selectedReceivable.status = 'partial';
    }
    showPaymentModal = false;
    paymentAmount = '';
  }
</script>

<div class="space-y-6">
  <div>
    <h1 class="text-2xl font-bold">Manajemen Piutang</h1>
    <p class="text-base-content">Kelola piutang dan pembayaran</p>
  </div>

  <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
    <div class="glass-card">
      <div class="p-5 p-5">
        <p class="text-sm text-base-content">Total Piutang</p>
        <p class="text-2xl font-bold">{formatCurrency(receivables.reduce((sum, r) => sum + r.remainingAmount, 0))}</p>
      </div>
    </div>
    <div class="glass-card">
      <div class="p-5 p-5">
        <p class="text-sm text-base-content">Belum Dibayar</p>
        <p class="text-2xl font-bold text-error">{formatCurrency(receivables.filter(r => r.status === 'unpaid').reduce((sum, r) => sum + r.remainingAmount, 0))}</p>
      </div>
    </div>
    <div class="glass-card">
      <div class="p-5 p-5">
        <p class="text-sm text-base-content">Sebagian Dibayar</p>
        <p class="text-2xl font-bold text-warning">{formatCurrency(receivables.filter(r => r.status === 'partial').reduce((sum, r) => sum + r.remainingAmount, 0))}</p>
      </div>
    </div>
  </div>

  <div class="flex items-center gap-3">
    <div class="form-control flex-1">
      <div class="relative">
        <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content" />
        <input type="text" class="input input-bordered input-sm w-full pl-10" placeholder="Cari piutang..." bind:value={searchQuery} />
      </div>
    </div>
    <div class="flex gap-1">
      {#each statusOptions as status}
        <button class="btn btn-sm {selectedStatus === status.id ? 'btn-primary' : 'btn-outline'}" onclick={() => selectedStatus = status.id}>{status.label}</button>
      {/each}
    </div>
  </div>

  <div class="glass-card">
    <div class="p-5">
      <div class="overflow-x-auto">
        <table class="table table-sm">
          <thead>
            <tr>
              <th>Invoice</th>
              <th>Pelanggan</th>
              <th>Total</th>
              <th>Sudah Dibayar</th>
              <th>Sisa</th>
              <th>Jatuh Tempo</th>
              <th>Status</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            {#each filteredReceivables as r}
              {@const badge = getStatusBadge(r.status)}
              <tr class="hover">
                <td class="font-mono text-xs">{r.invoiceNumber}</td>
                <td>{r.customerName}</td>
                <td>{formatCurrency(r.totalAmount)}</td>
                <td>{formatCurrency(r.paidAmount)}</td>
                <td class="font-medium text-error">{formatCurrency(r.remainingAmount)}</td>
                <td>{formatDate(r.dueDate)}</td>
                <td><span class="badge badge-sm {badge.color}">{badge.label}</span></td>
                <td>
                  <div class="flex gap-1">
                    <button class="btn btn-ghost btn-xs" disabled={r.status === 'paid'} onclick={() => { selectedReceivable = r; showPaymentModal = true; }}>
                      <DollarSign class="w-4 h-4" />
                    </button>
                    <button class="btn btn-ghost btn-xs"><Eye class="w-4 h-4" /></button>
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

{#if showPaymentModal && selectedReceivable}
  <div class="modal modal-open">
    <div class="modal-box">
      <h3 class="text-lg font-bold mb-4">Bayar Piutang</h3>
      <div class="bg-base-200 rounded-lg p-4 mb-4">
        <p class="text-sm"><span class="text-base-content">Invoice:</span> {selectedReceivable.invoiceNumber}</p>
        <p class="text-sm"><span class="text-base-content">Pelanggan:</span> {selectedReceivable.customerName}</p>
        <p class="text-sm"><span class="text-base-content">Sisa:</span> <span class="text-error font-bold">{formatCurrency(selectedReceivable.remainingAmount)}</span></p>
      </div>
      <div class="space-y-4">
        <div class="form-control">
          <label class="label"><span class="label-text font-medium">Nominal Pembayaran</span></label>
          <input type="number" class="input input-bordered w-full" bind:value={paymentAmount} placeholder="Masukkan nominal" max={selectedReceivable.remainingAmount} />
        </div>
        <div class="form-control">
          <label class="label"><span class="label-text font-medium">Metode Pembayaran</span></label>
          <select class="select select-bordered w-full" bind:value={paymentMethod}>
            <option value="cash">Tunai</option>
            <option value="dana">DANA</option>
            <option value="gopay">GoPay</option>
            <option value="bca">BCA</option>
            <option value="bri">BRI</option>
            <option value="bni">BNI</option>
            <option value="mandiri">Mandiri</option>
          </select>
        </div>
      </div>
      <div class="modal-action">
        <button class="btn btn-ghost" onclick={() => showPaymentModal = false}>Batal</button>
        <button class="btn btn-primary" onclick={makePayment}>Proses Pembayaran</button>
      </div>
    </div>
  </div>
{/if}
