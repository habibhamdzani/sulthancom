<script lang="ts">
  import { onMount } from 'svelte';
  import { mockTransactions } from '$lib/data/mock';
  import { formatCurrency, formatDate } from '$lib/utils/helpers';
  import { Printer, Download, Eye, Search } from 'lucide-svelte';

  let searchQuery = $state('');
  let selectedTx = $state<any>(null);
  let showModal = $state(false);

  let storeName = $state('SULTHAN COMPUTER');
  let storeAddress = $state('Jl. TGH. Abdul Karim, Kecamatan Kediri, Lombok Barat, Nusa Tenggara Barat');
  let storePhone = $state('081907772549');
  let storeEmail = $state('habibnhamdzani@gmail.com');
  let storeLogo = $state<string | null>(null);
  let stampPreview = $state<string | null>(null);
  let taxRate = $state(11);
  let invoiceSettings = $state({
    headerColor: '#1d4ed8',
    totalColor: '#1d4ed8',
    textColor: '#111827',
    borderColor: '#e5e7eb',
    fontSize: '13',
    fontFamily: 'Segoe UI',
    showStamp: true,
    showLogo: true,
    headerStyle: 'modern' as string,
    footerText: 'TERIMA KASIH ATAS KEPERCAYAAN ANDA',
  });

  onMount(() => {
    try {
      const saved = localStorage.getItem('invoiceDesign');
      if (saved) {
        const s = JSON.parse(saved);
        Object.assign(invoiceSettings, s);
        if (s.storeName) storeName = s.storeName.toUpperCase();
        if (s.storeAddress) storeAddress = s.storeAddress;
        if (s.storePhone) storePhone = s.storePhone;
        if (s.storeEmail) storeEmail = s.storeEmail;
        if (s.taxRate != null) taxRate = Number(s.taxRate);
      }
    } catch {}
    try {
      const logo = localStorage.getItem('storeLogo');
      if (logo) storeLogo = logo;
    } catch {}
    try {
      const stamp = localStorage.getItem('invoiceStamp');
      if (stamp) stampPreview = stamp;
    } catch {}
  });

  let filteredTransactions = $derived(mockTransactions.filter(tx =>
    !searchQuery || tx.invoiceNumber.toLowerCase().includes(searchQuery.toLowerCase()) || (tx.customerName && tx.customerName.toLowerCase().includes(searchQuery.toLowerCase()))
  ));

  function calcTax(tx: any) {
    const taxable = tx.subtotal - (tx.discount || 0);
    return Math.round(taxable * taxRate / 100);
  }

  function calcTotal(tx: any) {
    return tx.subtotal - (tx.discount || 0) + calcTax(tx);
  }

  function viewInvoice(tx: any) {
    selectedTx = tx;
    showModal = true;
  }

  function printInvoice(tx: any) {
    console.log('printInvoice called', tx);
    const win = window.open('', '_blank');
    if (!win) {
      alert('Popup diblokir browser. Izinkan popup untuk situs ini.');
      return;
    }
    const taxAmt = calcTax(tx);
    const finalTotal = calcTotal(tx);
    const items = tx.items.map((i: any) => `
      <div class="row">
        <div class="desc">${i.name}</div>
        <div class="qty">${i.quantity}</div>
        <div class="price">Rp ${i.price.toLocaleString('id-ID')}</div>
        <div class="sum">Rp ${i.subtotal.toLocaleString('id-ID')}</div>
      </div>
    `).join('');
    const emptyRows = Array(Math.max(0, 5 - tx.items.length)).fill('').map(() => '<div class="row empty"><div>&nbsp;</div><div></div><div></div><div></div></div>').join('');
    const headerBg = invoiceSettings.headerStyle === 'minimal' ? `background:transparent; border-bottom:2px solid ${invoiceSettings.headerColor}; color:${invoiceSettings.headerColor}` :
                     invoiceSettings.headerStyle === 'classic' ? `background:linear-gradient(135deg, ${invoiceSettings.headerColor}, ${invoiceSettings.headerColor}dd)` :
                     `background:${invoiceSettings.headerColor}`;
    const totalBg = invoiceSettings.headerStyle === 'minimal' ? `background:transparent; border:2px solid ${invoiceSettings.totalColor}; color:${invoiceSettings.totalColor}` :
                    `background:${invoiceSettings.totalColor}`;
    const html = `<!DOCTYPE html>
<html lang="id">
<head><meta charset="UTF-8"><meta name="viewport" content="width=210mm"><title>Invoice ${tx.invoiceNumber}</title>
<style>
  @page { margin: 0; size: A4 portrait; }
  * { margin: 0; padding: 0; box-sizing: border-box; }
  body { font-family: '${invoiceSettings.fontFamily}', -apple-system, sans-serif; background: #d1d5db; display: flex; justify-content: center; padding: 24px; font-size:${invoiceSettings.fontSize}px; }
  .page { width: 210mm; min-height: 296mm; background: #fff; padding: 40px 48px; box-shadow: 0 8px 30px rgba(0,0,0,.18); }
  .hdr { display: flex; align-items: flex-start; gap: 16px; padding-bottom: 16px; border-bottom: 1px solid ${invoiceSettings.borderColor}; margin-bottom: 24px; }
  .hdr-logo { width: 56px; height: 56px; border-radius: 8px; object-fit: contain; flex-shrink: 0; }
  .hdr-logo-fb { width: 56px; height: 56px; border-radius: 8px; background: ${invoiceSettings.headerColor}; display: flex; align-items: center; justify-content: center; color: #fff; font-size: 22px; font-weight: 800; flex-shrink: 0; }
  .hdr-info h1 { font-size: 22px; font-weight: 800; color: ${invoiceSettings.textColor}; }
  .hdr-info .addr { font-size: 12px; color: #6b7280; line-height: 1.4; margin-top: 2px; }
  .hdr-info .phone { font-size: 12px; font-weight: 700; color: ${invoiceSettings.textColor}; margin-top: 2px; }
  .hdr-info .email { font-size: 12px; color: ${invoiceSettings.headerColor}; margin-top: 1px; }
  .info { display: flex; justify-content: space-between; margin-bottom: 24px; }
  .bill .lbl { font-size: 10px; font-weight: 700; color: #9ca3af; text-transform: uppercase; letter-spacing: 1px; margin-bottom: 4px; }
  .bill .name { font-size: 14px; font-weight: 700; color: ${invoiceSettings.textColor}; }
  .bill .adr { font-size: 12px; color: #6b7280; margin-top: 2px; }
  .meta-r .rw { display: flex; margin-bottom: 2px; }
  .meta-r .rw .lb { font-size: 12px; color: #374151; font-weight: 600; width: 80px; }
  .meta-r .rw .vl { font-size: 12px; font-weight: 700; color: ${invoiceSettings.textColor}; }
  .tbl-wrap { border: 2px solid ${invoiceSettings.borderColor}; border-radius: 12px; overflow: hidden; margin-bottom: 24px; }
  .tbl-hdr { display: flex; ${headerBg}; font-size: 11px; font-weight: 700; text-transform: uppercase; letter-spacing: .5px; }
  .tbl-hdr > div { padding: 12px 16px; }
  .tbl-hdr > div:first-child { flex: 2; }
  .tbl-hdr > div:nth-child(2) { flex: 0.8; text-align: center; }
  .tbl-hdr > div:nth-child(3) { flex: 1; text-align: center; }
  .tbl-hdr > div:last-child { flex: 1; text-align: right; }
  .row { display: flex; border-bottom: 1px solid ${invoiceSettings.borderColor}; }
  .row:last-child { border-bottom: none; }
  .row > div { padding: 14px 16px; }
  .row > div:first-child { flex: 2; }
  .row > div:nth-child(2) { flex: 0.8; text-align: center; color: #6b7280; font-weight: 600; }
  .row > div:nth-child(3) { flex: 1; text-align: center; color: #6b7280; }
  .row > div:last-child { flex: 1; text-align: right; font-weight: 700; color: ${invoiceSettings.textColor}; }
  .row .desc { font-weight: 600; color: ${invoiceSettings.textColor}; }
  .sum { display: flex; justify-content: flex-end; margin-bottom: 32px; }
  .sum-in { width: 240px; }
  .sum .ln { display: flex; justify-content: space-between; padding: 3px 0; font-size: 12px; color: #374151; }
  .sum .ln.disc { color: #dc2626; }
  .sum .ln.ht { padding: 10px 16px; margin-top: 8px; border-radius: 8px; ${totalBg}; }
  .sum .ln.ht .t { color: #fff; font-weight: 700; font-size: 12px; text-transform: uppercase; letter-spacing: .5px; }
  .sum .ln.ht .a { color: #fff; font-weight: 700; font-size: 16px; font-style: italic; }
  .ftr { display: flex; justify-content: space-between; align-items: stretch; margin-top: 40px; }
  .sig { text-align: center; height: 140px; display: flex; flex-direction: column; justify-content: space-between; align-items: center; }
  .sig h4 { font-size: 11px; font-weight: 700; color: ${invoiceSettings.textColor}; text-transform: uppercase; letter-spacing: .5px; }
  .sig .stamp { flex: 1; display: flex; align-items: center; justify-content: center; }
  .sig .ln { width: 180px; border-top: 1px solid ${invoiceSettings.textColor}; }
  .sig .stamp img { width: 80px; height: auto; display: block; }
  .sig .stamp .logo-fb { width: 56px; height: 56px; border-radius: 50%; background: ${invoiceSettings.headerColor}; display: flex; align-items: center; justify-content: center; color: #fff; font-size: 22px; font-weight: 800; }
  .thx { text-align: center; margin-top: 32px; font-size: 10px; color: #cbd5e1; letter-spacing: 2px; text-transform: uppercase; }
  @media print { body { padding: 0; background: #fff; } .page { box-shadow: none; padding: 40px 48px; } .thx { color: #94a3b8; } }
</style></head><body>
<div class="page">
  <div class="hdr">
    ${invoiceSettings.showLogo && storeLogo ? `<img src="${storeLogo}" class="hdr-logo" />` : invoiceSettings.showLogo ? '<div class="hdr-logo-fb">SC</div>' : ''}
    <div class="hdr-info">
      <h1>${storeName}</h1>
      <div class="addr">${storeAddress}</div>
      <div class="phone">${storePhone}</div>
      <div class="email">${storeEmail}</div>
    </div>
  </div>
  <div class="info">
    <div class="bill">
      <div class="lbl">Kepada</div>
      <div class="name">${tx.customerName || 'Pelanggan Umum'}</div>
      <div class="adr">${tx.customerAddress || '-'}</div>
    </div>
    <div class="meta-r">
      <div class="rw"><span class="lb">No. Invoice</span><span class="vl">${tx.invoiceNumber}</span></div>
      <div class="rw"><span class="lb">Tanggal</span><span class="vl">${formatDate(tx.createdAt)}</span></div>
      <div class="rw"><span class="lb">Jatuh Tempo</span><span class="vl">${tx.dueDate ? formatDate(tx.dueDate) : '-'}</span></div>
      <div class="rw"><span class="lb">Pembayaran</span><span class="vl" style="text-transform:capitalize">${tx.paymentMethod}</span></div>
    </div>
  </div>
  <div class="tbl-wrap">
    <div class="tbl-hdr"><div>Deskripsi</div><div>Qty</div><div>Harga</div><div>Jumlah</div></div>
    ${items}${emptyRows}
  </div>
  <div class="sum">
    <div class="sum-in">
      <div class="ln"><span>Subtotal</span><span>Rp ${tx.subtotal.toLocaleString('id-ID')}</span></div>
      ${tx.discount > 0 ? `<div class="ln disc"><span>Diskon</span><span>-Rp ${tx.discount.toLocaleString('id-ID')}</span></div>` : ''}
      ${taxRate > 0 ? `<div class="ln"><span>Pajak PPN ${taxRate}%</span><span>Rp ${taxAmt.toLocaleString('id-ID')}</span></div>` : ''}
      <div class="ln ht"><span class="t">Total</span><span class="a">Rp ${finalTotal.toLocaleString('id-ID')}</span></div>
    </div>
  </div>
  <div class="ftr">
    <div class="sig"><h4>Tanda Terima</h4><div class="stamp"></div><div class="ln"></div></div>
    ${invoiceSettings.showStamp ? `<div class="sig"><h4>Hormat Kami</h4>
      <div class="stamp">
        ${stampPreview ? `<img src="${stampPreview}" alt="Stamp" />` : '<div class="logo-fb">SC</div>'}
      </div>
      <div class="ln"></div>
    </div>` : ''}
  </div>
  <div style="height:32px"></div>
  <div class="thx">${invoiceSettings.footerText}</div>
</div>
<script>window.print();window.close();<\\/script>
</body></html>`;
    console.log('Writing HTML to window');
    win.document.write(html);
    win.document.close();
  }

  function downloadPdf(tx: any) {
    printInvoice(tx);
  }

  function getPaymentBadge(status: string) {
    const map: Record<string, string> = { paid: 'badge-success', dp: 'badge-warning', unpaid: 'badge-error' };
    return `badge-sm ${map[status] || 'badge-ghost'}`;
  }

  function getPaymentLabel(status: string) {
    const map: Record<string, string> = { paid: 'Lunas', dp: 'DP', unpaid: 'Piutang' };
    return map[status] || status;
  }
</script>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-2xl font-bold">Invoice</h1>
      <p class="text-base-content">Kelola dan cetak invoice</p>
    </div>
  </div>

  <div class="glass-card">
    <div class="p-5">
      <div class="flex items-center gap-3 mb-4">
        <div class="form-control flex-1">
          <div class="relative">
            <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content" />
            <input type="text" class="input input-bordered input-sm w-full pl-10" placeholder="Cari invoice..." bind:value={searchQuery} />
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
              <th>Status</th>
              <th class="text-center">Aksi</th>
            </tr>
          </thead>
          <tbody>
            {#each filteredTransactions as tx}
              <tr class="hover">
                <td class="font-mono text-xs font-medium">{tx.invoiceNumber}</td>
                <td class="text-sm">{formatDate(tx.createdAt)}</td>
                <td><span class="font-medium">{tx.customerName || 'Umum'}</span></td>
                <td>{tx.items.length} item</td>
                <td class="font-semibold">{formatCurrency(calcTotal(tx))}</td>
                <td><span class="badge {getPaymentBadge(tx.paymentStatus)}">{getPaymentLabel(tx.paymentStatus)}</span></td>
                <td>
                  <div class="flex gap-1 justify-center">
                    <button class="btn btn-ghost btn-xs btn-square" onclick={() => viewInvoice(tx)} title="Lihat"><Eye class="w-4 h-4" /></button>
                    <button class="btn btn-ghost btn-xs btn-square" onclick={() => printInvoice(tx)} title="Cetak"><Printer class="w-4 h-4" /></button>
                    <button class="btn btn-ghost btn-xs btn-square" onclick={() => downloadPdf(tx)} title="PDF"><Download class="w-4 h-4" /></button>
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

{#if showModal && selectedTx}
  <dialog class="modal modal-open" onclick={() => showModal = false}>
    <div class="modal-box max-w-[210mm] p-0 bg-transparent shadow-none" onclick={(e) => e.stopPropagation()}>
      <div class="sticky top-0 z-10 flex justify-end p-3">
        <button class="btn btn-sm btn-circle btn-ghost bg-base-100 shadow-md" onclick={() => showModal = false}>✕</button>
      </div>
      <div class="bg-white rounded-xl shadow-xl p-10 mx-4" style="font-family:{invoiceSettings.fontFamily}, sans-serif; font-size:{invoiceSettings.fontSize}px;">
        <div class="flex items-start gap-4 pb-4 border-b mb-6" style="border-color:{invoiceSettings.borderColor}">
          {#if invoiceSettings.showLogo && storeLogo}
            <img src={storeLogo} alt="Logo" class="w-14 h-14 rounded-lg object-contain flex-shrink-0" />
          {:else if invoiceSettings.showLogo}
            <div class="w-14 h-14 rounded-lg flex items-center justify-center text-white text-xl font-extrabold flex-shrink-0" style="background:{invoiceSettings.headerColor}">SC</div>
          {/if}
          <div>
            <h1 class="text-xl font-extrabold" style="color:{invoiceSettings.textColor}">{storeName}</h1>
            <p class="text-xs mt-0.5 leading-snug" style="color:#6b7280">{storeAddress}</p>
            <p class="text-xs font-bold mt-0.5" style="color:{invoiceSettings.textColor}">{storePhone}</p>
            <p class="text-xs mt-0.5" style="color:{invoiceSettings.headerColor}">{storeEmail}</p>
          </div>
        </div>

          <div class="flex justify-between mb-6">
            <div>
              <p class="text-[10px] font-bold text-slate-400 uppercase tracking-wide mb-1">Kepada</p>
              <p class="text-sm font-bold" style="color:{invoiceSettings.textColor}">{selectedTx.customerName || 'Pelanggan Umum'}</p>
              <p class="text-xs mt-0.5" style="color:#6b7280">{selectedTx.customerAddress || '-'}</p>
            </div>
            <div class="space-y-0.5">
              <div class="flex gap-4 items-center"><span class="text-xs text-slate-600 font-semibold w-20">No. Invoice</span><span class="text-xs font-bold" style="color:{invoiceSettings.textColor}">{selectedTx.invoiceNumber}</span></div>
              <div class="flex gap-4 items-center"><span class="text-xs text-slate-600 font-semibold w-20">Tanggal</span><span class="text-xs font-bold" style="color:{invoiceSettings.textColor}">{formatDate(selectedTx.createdAt)}</span></div>
              <div class="flex gap-4 items-center"><span class="text-xs text-slate-600 font-semibold w-20">Jatuh Tempo</span><span class="text-xs font-bold" style="color:{invoiceSettings.textColor}">{selectedTx.dueDate ? formatDate(selectedTx.dueDate) : '-'}</span></div>
              <div class="flex gap-4 items-center"><span class="text-xs text-slate-600 font-semibold w-20">Pembayaran</span><span class="text-xs font-bold capitalize" style="color:{invoiceSettings.textColor}">{selectedTx.paymentMethod}</span></div>
            </div>
          </div>

        <div class="border-2 rounded-xl overflow-hidden mb-6" style="border-color:{invoiceSettings.borderColor}">
          <div class="flex text-white text-[11px] font-bold uppercase tracking-wide" style={invoiceSettings.headerStyle === 'minimal' ? `background:transparent; border-bottom:2px solid ${invoiceSettings.headerColor}; color:${invoiceSettings.headerColor}` : invoiceSettings.headerStyle === 'classic' ? `background:linear-gradient(135deg, ${invoiceSettings.headerColor}, ${invoiceSettings.headerColor}dd)` : `background:${invoiceSettings.headerColor}`}>
            <div class="flex-[2] px-4 py-3">Deskripsi</div>
            <div class="flex-[0.8] px-4 py-3 text-center">Qty</div>
            <div class="flex-1 px-4 py-3 text-center">Harga</div>
            <div class="flex-1 px-4 py-3 text-right">Jumlah</div>
          </div>
          {#each selectedTx.items as item}
            <div class="flex border-b" style="border-color:{invoiceSettings.borderColor}">
              <div class="flex-[2] px-4 py-3.5"><span class="font-semibold" style="color:{invoiceSettings.textColor}">{item.name}</span></div>
              <div class="flex-[0.8] px-4 py-3.5 text-center font-semibold text-slate-600">{item.quantity}</div>
              <div class="flex-1 px-4 py-3.5 text-center text-slate-600">Rp {item.price.toLocaleString('id-ID')}</div>
              <div class="flex-1 px-4 py-3.5 text-right font-bold" style="color:{invoiceSettings.textColor}">Rp {item.subtotal.toLocaleString('id-ID')}</div>
            </div>
          {/each}
          {#each Array(Math.max(0, 5 - selectedTx.items.length)) as _}
            <div class="flex border-b" style="border-color:{invoiceSettings.borderColor}">
              <div class="flex-[2] px-4 py-3.5">&nbsp;</div>
              <div class="flex-[0.8] px-4 py-3.5"></div>
              <div class="flex-1 px-4 py-3.5"></div>
              <div class="flex-1 px-4 py-3.5"></div>
            </div>
          {/each}
        </div>

        <div class="flex justify-end mb-8">
          <div style="width:240px">
            <div class="flex justify-between py-0.5 text-xs text-slate-700"><span>Subtotal</span><span>Rp {selectedTx.subtotal.toLocaleString('id-ID')}</span></div>
            {#if selectedTx.discount > 0}
              <div class="flex justify-between py-0.5 text-xs text-red-600"><span>Diskon</span><span>-Rp {selectedTx.discount.toLocaleString('id-ID')}</span></div>
            {/if}
            {#if taxRate > 0}
              <div class="flex justify-between py-0.5 text-xs text-slate-700"><span>Pajak PPN {taxRate}%</span><span>Rp {calcTax(selectedTx).toLocaleString('id-ID')}</span></div>
            {/if}
            <div class="flex justify-between items-center mt-2 px-4 py-2.5 rounded-lg text-white" style={invoiceSettings.headerStyle === 'minimal' ? `background:transparent; border:2px solid ${invoiceSettings.totalColor}; color:${invoiceSettings.totalColor}` : `background:${invoiceSettings.totalColor}`}>
              <span class="font-bold text-xs uppercase tracking-wide">Total</span>
              <span class="font-bold text-lg italic">Rp {calcTotal(selectedTx).toLocaleString('id-ID')}</span>
            </div>
          </div>
        </div>

        <div class="flex justify-between mt-16">
          <div class="flex flex-col items-center h-[140px] justify-between">
            <p class="text-[11px] font-bold uppercase tracking-wide" style="color:{invoiceSettings.textColor}">Tanda Terima</p>
            <div class="w-44" style="border-top:1px solid {invoiceSettings.textColor}"></div>
          </div>
          {#if invoiceSettings.showStamp}
            <div class="flex flex-col items-center h-[140px] justify-between">
              <p class="text-[11px] font-bold uppercase tracking-wide" style="color:{invoiceSettings.textColor}">Hormat Kami</p>
              <div class="flex-1 flex items-center justify-center">
                {#if stampPreview}
                  <img src={stampPreview} alt="Stamp" class="w-20 h-auto" />
                {:else}
                  <div class="w-14 h-14 rounded-full flex items-center justify-center text-white text-xl font-extrabold" style="background:{invoiceSettings.headerColor}">SC</div>
                {/if}
              </div>
              <div class="w-44" style="border-top:1px solid {invoiceSettings.textColor}"></div>
            </div>
          {/if}
        </div>
        <div class="h-8"></div>
        <p class="text-center mt-8 text-[10px] tracking-[2px] uppercase" style="color:#d1d5db">{invoiceSettings.footerText}</p>
      </div>
      <div class="flex justify-center gap-3 py-4">
        <button class="btn btn-primary btn-sm gap-2" onclick={() => { printInvoice(selectedTx); }}><Printer class="w-4 h-4" /> Cetak</button>
        <button class="btn btn-outline btn-sm gap-2" onclick={() => { downloadPdf(selectedTx); }}><Download class="w-4 h-4" /> Download PDF</button>
      </div>
    </div>
  </dialog>
{/if}