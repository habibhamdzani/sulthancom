<script lang="ts">
  import { mockTransactions } from '$lib/data/mock';
  import { formatCurrency, formatDate, getPaymentStatusLabel, getPaymentStatusColor } from '$lib/utils/helpers';
  import { Search, Eye, Printer } from 'lucide-svelte';

  let searchQuery = $state('');
  let filteredTransactions = $derived(mockTransactions.filter(tx =>
    !searchQuery || tx.invoiceNumber.toLowerCase().includes(searchQuery.toLowerCase()) || (tx.customerName && tx.customerName.toLowerCase().includes(searchQuery.toLowerCase()))
  ));
</script>

<div class="space-y-6">
  <div>
    <h1 class="text-2xl font-bold">Riwayat Transaksi</h1>
    <p class="text-base-content">Daftar semua transaksi penjualan</p>
  </div>

  <div class="glass-card">
    <div class="p-5">
      <div class="flex items-center gap-3 mb-4">
        <div class="form-control flex-1">
          <div class="relative">
            <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content" />
            <input type="text" class="input input-bordered input-sm w-full pl-10" placeholder="Cari invoice atau pelanggan..." bind:value={searchQuery} />
          </div>
        </div>
      </div>

      <div class="overflow-x-auto">
        <table class="table table-sm">
          <thead>
            <tr>
              <th>Invoice</th>
              <th>Tanggal</th>
              <th>Pelanggan</th>
              <th>Items</th>
              <th>Total</th>
              <th>Pembayaran</th>
              <th>Status</th>
              <th>Kasir</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            {#each filteredTransactions as tx}
              <tr class="hover">
                <td class="font-mono text-xs">{tx.invoiceNumber}</td>
                <td>{formatDate(tx.createdAt)}</td>
                <td>{tx.customerName || 'Umum'}</td>
                <td>{tx.items.length} item</td>
                <td class="font-medium">{formatCurrency(tx.total)}</td>
                <td class="capitalize">{tx.paymentMethod}</td>
                <td><span class="badge badge-sm {getPaymentStatusColor(tx.paymentStatus)}">{getPaymentStatusLabel(tx.paymentStatus)}</span></td>
                <td class="text-xs">{tx.cashierName}</td>
                <td>
                  <div class="flex gap-1">
                    <button class="btn btn-ghost btn-xs" title="Detail" onclick={() => window.open('/invoices', '_self')}><Eye class="w-4 h-4" /></button>
                    <button class="btn btn-ghost btn-xs" title="Cetak"><Printer class="w-4 h-4" /></button>
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
