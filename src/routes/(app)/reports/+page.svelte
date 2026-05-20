<script lang="ts">
  import { salesChartData, categorySalesData, paymentMethodData, mockTransactions, mockServiceTickets, mockProducts } from '$lib/data/mock';
  import { formatCurrency, formatDate } from '$lib/utils/helpers';
  import Chart from '$lib/components/ui/Chart.svelte';
  import { Download, Printer, Calendar } from 'lucide-svelte';

  let reportType = $state('sales');
  let dateRange = $state('monthly');

  const salesChartOptions = {
    chart: { type: 'bar', height: 300, toolbar: { show: false } },
    plotOptions: { bar: { horizontal: false, columnWidth: '55%', borderRadius: 4 } },
    dataLabels: { enabled: false },
    xaxis: { categories: salesChartData.map(d => formatDate(d.date)) },
    yaxis: { labels: { formatter: (val: number) => formatCurrency(val).replace('Rp', '').trim() } },
    tooltip: { y: { formatter: (val: number) => formatCurrency(val) } },
    colors: ['#3B82F6'],
    grid: { borderColor: '#f1f5f9' },
  };

  const salesChartSeries = [{ name: 'Penjualan', data: salesChartData.map(d => d.sales) }];
  const transactionChartSeries = [{ name: 'Transaksi', data: salesChartData.map(d => d.transactions) }];

  const categoryChartOptions = {
    chart: { type: 'donut', height: 280 },
    labels: categorySalesData.map(d => d.name),
    colors: categorySalesData.map(d => d.color),
    legend: { position: 'bottom' },
    plotOptions: { pie: { donut: { size: '65%' } } },
  };

  const categoryChartSeries = categorySalesData.map(d => d.value);

  let lowStockProducts = $derived(mockProducts.filter(p => p.stock <= p.minStock));
  let bestSelling = $derived([...mockProducts].sort((a, b) => b.stock - a.stock).slice(0, 5));
</script>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-2xl font-bold">Laporan & Statistik</h1>
      <p class="text-base-content">Analisis data penjualan dan bisnis</p>
    </div>
    <div class="flex gap-2">
      <button class="btn btn-outline btn-sm gap-2" onclick={() => window.print()}><Download class="w-4 h-4" /> Export PDF</button>
      <button class="btn btn-outline btn-sm gap-2" onclick={() => window.print()}><Printer class="w-4 h-4" /> Print</button>
    </div>
  </div>

  <div class="flex items-center gap-3">
    <div class="tabs tabs-boxed bg-base-200">
      <button class="tab {reportType === 'sales' ? 'tab-active' : ''}" onclick={() => reportType = 'sales'}>Penjualan</button>
      <button class="tab {reportType === 'product' ? 'tab-active' : ''}" onclick={() => reportType = 'product'}>Produk</button>
      <button class="tab {reportType === 'service' ? 'tab-active' : ''}" onclick={() => reportType = 'service'}>Service</button>
    </div>
    <select class="select select-bordered select-sm" bind:value={dateRange}>
      <option value="daily">Harian</option>
      <option value="weekly">Mingguan</option>
      <option value="monthly">Bulanan</option>
      <option value="yearly">Tahunan</option>
    </select>
  </div>

  {#if reportType === 'sales'}
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="glass-card">
        <div class="p-5">
          <h2 class="card-title">Grafik Penjualan</h2>
          <Chart options={salesChartOptions} series={salesChartSeries} type="bar" height={300} />
        </div>
      </div>
      <div class="glass-card">
        <div class="p-5">
          <h2 class="card-title">Grafik Transaksi</h2>
          <Chart options={{...salesChartOptions, colors: ['#10B981']}} series={transactionChartSeries} type="line" height={300} />
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="glass-card">
        <div class="p-5">
          <h2 class="card-title">Penjualan per Kategori</h2>
          <Chart options={categoryChartOptions} series={categoryChartSeries} type="donut" height={280} />
        </div>
      </div>
      <div class="glass-card">
        <div class="p-5">
          <h2 class="card-title">Metode Pembayaran</h2>
          <div class="space-y-3">
            {#each paymentMethodData as method}
              <div class="flex items-center gap-3">
                <div class="w-3 h-3 rounded-full bg-primary"></div>
                <span class="flex-1">{method.name}</span>
                <div class="w-32 bg-base-200 rounded-full h-2">
                  <div class="bg-primary h-2 rounded-full" style="width: {method.value}%"></div>
                </div>
                <span class="text-sm font-medium w-12 text-right">{method.value}%</span>
              </div>
            {/each}
          </div>
        </div>
      </div>
    </div>
  {:else if reportType === 'product'}
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="glass-card">
        <div class="p-5">
          <h2 class="card-title">Stok Menipis</h2>
          <div class="overflow-x-auto">
            <table class="table table-sm">
              <thead>
                <tr>
                  <th>Produk</th>
                  <th>Stok</th>
                  <th>Min. Stok</th>
                </tr>
              </thead>
              <tbody>
                {#each lowStockProducts as product}
                  <tr>
                    <td>{product.name}</td>
                    <td><span class="text-error font-bold">{product.stock}</span></td>
                    <td>{product.minStock}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        </div>
      </div>
      <div class="glass-card">
        <div class="p-5">
          <h2 class="card-title">Produk dengan Stok Terbanyak</h2>
          <div class="overflow-x-auto">
            <table class="table table-sm">
              <thead>
                <tr>
                  <th>Produk</th>
                  <th>Stok</th>
                  <th>Harga</th>
                </tr>
              </thead>
              <tbody>
                {#each bestSelling as product}
                  <tr>
                    <td>{product.name}</td>
                    <td>{product.stock}</td>
                    <td>{formatCurrency(product.sellingPrice)}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  {:else}
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="glass-card">
        <div class="p-5">
          <h2 class="card-title">Status Service</h2>
          <div class="space-y-3">
            {#each ['waiting', 'checked', 'in_progress', 'waiting_parts', 'completed', 'cancelled'] as status}
              {@const count = mockServiceTickets.filter(t => t.status === status).length}
              <div class="flex items-center justify-between">
                <span class="capitalize">{status.replace('_', ' ')}</span>
                <span class="badge badge-sm">{count}</span>
              </div>
            {/each}
          </div>
        </div>
      </div>
      <div class="glass-card">
        <div class="p-5">
          <h2 class="card-title">Riwayat Service</h2>
          <div class="overflow-x-auto">
            <table class="table table-sm">
              <thead>
                <tr>
                  <th>Tiket</th>
                  <th>Pelanggan</th>
                  <th>Status</th>
                </tr>
              </thead>
              <tbody>
                {#each mockServiceTickets as ticket}
                  <tr>
                    <td class="font-mono text-xs">{ticket.ticketNumber}</td>
                    <td>{ticket.customerName}</td>
                    <td class="capitalize">{ticket.status.replace('_', ' ')}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  {/if}
</div>
