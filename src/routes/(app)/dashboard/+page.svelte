<script lang="ts">
  import StatCard from '$lib/components/ui/StatCard.svelte';
  import { mockDashboardStats, salesChartData, categorySalesData, paymentMethodData, mockTransactions, mockServiceTickets } from '$lib/data/mock';
  import { formatCurrency, formatDate, getServiceStatusLabel, getServiceStatusColor } from '$lib/utils/helpers';
  import { TrendingUp, ShoppingCart, Wrench, CreditCard, Package, AlertTriangle, DollarSign, BarChart3 } from 'lucide-svelte';
  import Chart from '$lib/components/ui/Chart.svelte';

  const stats = mockDashboardStats;

  const salesChartOptions = {
    chart: { type: 'area', height: 300, toolbar: { show: false }, zoom: { enabled: false } },
    dataLabels: { enabled: false },
    stroke: { curve: 'smooth', width: 2 },
    xaxis: { categories: salesChartData.map(d => formatDate(d.date)) },
    yaxis: { labels: { formatter: (val: number) => formatCurrency(val).replace('Rp', '').trim() } },
    tooltip: { y: { formatter: (val: number) => formatCurrency(val) } },
    colors: ['#3B82F6'],
    fill: { type: 'gradient', gradient: { shadeIntensity: 1, opacityFrom: 0.4, opacityTo: 0.1, stops: [0, 90, 100] } },
    grid: { borderColor: '#f1f5f9' },
  };

  const salesChartSeries = [{ name: 'Penjualan', data: salesChartData.map(d => d.sales) }];

  const categoryChartOptions = {
    chart: { type: 'donut', height: 280 },
    labels: categorySalesData.map(d => d.name),
    colors: categorySalesData.map(d => d.color),
    legend: { position: 'bottom' },
    tooltip: { y: { formatter: (val: number) => `${val}%` } },
    plotOptions: { pie: { donut: { size: '65%', labels: { show: true, total: { show: true, label: 'Total', formatter: () => '100%' } } } } },
  };

  const categoryChartSeries = categorySalesData.map(d => d.value);

  const paymentChartOptions = {
    chart: { type: 'bar', height: 280, toolbar: { show: false } },
    plotOptions: { bar: { horizontal: false, columnWidth: '55%', borderRadius: 4 } },
    dataLabels: { enabled: false },
    xaxis: { categories: paymentMethodData.map(d => d.name) },
    colors: ['#10B981'],
    grid: { borderColor: '#f1f5f9' },
  };

  const paymentChartSeries = [{ name: 'Persentase', data: paymentMethodData.map(d => d.value) }];
</script>

<div class="space-y-6">
  <div>
    <h1 class="text-2xl font-bold">Dashboard</h1>
    <p class="text-base-content">Ringkasan data toko SulthanCom</p>
  </div>

  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
    <StatCard title="Penjualan Hari Ini" value={formatCurrency(stats.todaySales)} icon={DollarSign} color="success" trend={{ value: 12.5, positive: true }} />
    <StatCard title="Transaksi Hari Ini" value={stats.todayTransactions} icon={ShoppingCart} color="primary" trend={{ value: 8.2, positive: true }} />
    <StatCard title="Service Aktif" value={stats.activeServices} icon={Wrench} color="warning" />
    <StatCard title="Total Piutang" value={formatCurrency(stats.totalReceivables)} icon={CreditCard} color="error" />
  </div>

  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
    <StatCard title="Total Produk" value={stats.totalProducts} icon={Package} color="info" />
    <StatCard title="Stok Menipis" value={stats.lowStockProducts} icon={AlertTriangle} color="error" />
    <StatCard title="Produk Terlaris" value={stats.bestSellingProduct.name} icon={TrendingUp} color="success" />
    <StatCard title="Produk Tidak Laku" value={stats.slowMovingProduct.name} icon={BarChart3} color="neutral" />
  </div>

  <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
    <div class="glass-card lg:col-span-2">
      <div class="p-5">
        <h2 class="text-lg font-bold mb-4">Grafik Penjualan 7 Hari Terakhir</h2>
        <Chart options={salesChartOptions} series={salesChartSeries} type="area" height={300} />
      </div>
    </div>

    <div class="glass-card">
      <div class="p-5">
        <h2 class="text-lg font-bold mb-4">Penjualan per Kategori</h2>
        <Chart options={categoryChartOptions} series={categoryChartSeries} type="donut" height={280} />
      </div>
    </div>
  </div>

  <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
    <div class="glass-card">
      <div class="p-5">
        <h2 class="text-lg font-bold mb-4">Metode Pembayaran</h2>
        <Chart options={paymentChartOptions} series={paymentChartSeries} type="bar" height={280} />
      </div>
    </div>

    <div class="glass-card">
      <div class="p-5">
        <h2 class="text-lg font-bold mb-4">Transaksi Terbaru</h2>
        <div class="overflow-x-auto">
          <table class="table table-sm">
            <thead>
              <tr>
                <th>Invoice</th>
                <th>Total</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              {#each mockTransactions.slice(0, 5) as tx}
                <tr class="hover">
                  <td class="font-mono text-xs">{tx.invoiceNumber}</td>
                  <td>{formatCurrency(tx.total)}</td>
                  <td><span class="badge badge-sm badge-success">Lunas</span></td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>

  <div class="glass-card">
    <div class="p-5">
      <h2 class="text-lg font-bold mb-4">Service Terbaru</h2>
      <div class="overflow-x-auto">
        <table class="table table-sm">
          <thead>
            <tr>
              <th>No. Tiket</th>
              <th>Pelanggan</th>
              <th>Device</th>
              <th>Keluhan</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            {#each mockServiceTickets as ticket}
              <tr class="hover">
                <td class="font-mono text-xs">{ticket.ticketNumber}</td>
                <td>{ticket.customerName}</td>
                <td>{ticket.deviceType}</td>
                <td class="max-w-xs truncate">{ticket.complaint}</td>
                <td><span class="badge badge-sm {getServiceStatusColor(ticket.status)}">{getServiceStatusLabel(ticket.status)}</span></td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>
  </div>
</div>
