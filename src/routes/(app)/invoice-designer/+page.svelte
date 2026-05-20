<script lang="ts">
  import { onMount } from 'svelte';
  import { Save, Upload, Eye, Palette, Layout, Type } from 'lucide-svelte';
  import { mockTransactions } from '$lib/data/mock';

  let saved = $state(false);
  let previewTx = $state(mockTransactions[0]);
  let activeTab = $state('general');

  let settings = $state({
    storeName: 'SULTHAN COMPUTER',
    storeAddress: 'Jl. TGH. Abdul Karim, Kecamatan Kediri, Lombok Barat, Nusa Tenggara Barat',
    storePhone: '081907772549',
    storeEmail: 'habibnhamdzani@gmail.com',
    taxRate: 11,
    footerText: 'TERIMA KASIH ATAS KEPERCAYAAN ANDA',
    headerColor: '#1d4ed8',
    totalColor: '#1d4ed8',
    textColor: '#111827',
    borderColor: '#e5e7eb',
    fontSize: '13',
    fontFamily: 'Segoe UI',
    showStamp: true,
    showLogo: true,
    stampPosition: 'right' as 'left' | 'right' | 'center',
    logoPosition: 'left' as 'left' | 'center' | 'right',
    headerStyle: 'modern' as 'modern' | 'classic' | 'minimal',
  });

  let logoPreview = $state<string | null>(null);
  let stampPreview = $state<string | null>(null);
  let logoInput: HTMLInputElement | undefined = $state();
  let stampInput: HTMLInputElement | undefined = $state();

  onMount(() => {
    try {
      const saved = localStorage.getItem('invoiceDesign');
      if (saved) Object.assign(settings, JSON.parse(saved));
    } catch {}
    try {
      const logo = localStorage.getItem('storeLogo');
      if (logo) logoPreview = logo;
    } catch {}
    try {
      const stamp = localStorage.getItem('invoiceStamp');
      if (stamp) stampPreview = stamp;
    } catch {}
  });

  function handleSave() {
    try {
      localStorage.setItem('invoiceDesign', JSON.stringify(settings));
      localStorage.setItem('storeSettings', JSON.stringify({
        name: settings.storeName,
        address: settings.storeAddress,
        phone: settings.storePhone,
        email: settings.storeEmail,
        taxRate: settings.taxRate,
        invoiceFooter: settings.footerText,
        theme: 'light',
      }));
      if (logoPreview) localStorage.setItem('storeLogo', logoPreview);
      if (stampPreview) localStorage.setItem('invoiceStamp', stampPreview);
    } catch {}
    saved = true;
    setTimeout(() => saved = false, 3000);
  }

  function triggerUpload(type: 'logo' | 'stamp') {
    if (type === 'logo') logoInput?.click();
    else stampInput?.click();
  }

  function handleFileUpload(e: Event, type: 'logo' | 'stamp') {
    const input = e.target as HTMLInputElement;
    const file = input.files?.[0];
    if (!file) return;
    const reader = new FileReader();
    reader.onload = (ev) => {
      const dataUrl = ev.target?.result as string;
      if (type === 'logo') {
        logoPreview = dataUrl;
        try { localStorage.setItem('storeLogo', dataUrl); } catch {}
      } else {
        stampPreview = dataUrl;
        try { localStorage.setItem('invoiceStamp', dataUrl); } catch {}
      }
    };
    reader.readAsDataURL(file);
  }

  function getHeaderBg() {
    return settings.headerStyle === 'modern' ? `background:${settings.headerColor}` : 
           settings.headerStyle === 'classic' ? `background:linear-gradient(135deg, ${settings.headerColor}, ${settings.headerColor}dd)` :
           `background:transparent; border-bottom:2px solid ${settings.headerColor}; color:${settings.headerColor}`;
  }

  function getTotalBg() {
    return settings.headerStyle === 'minimal' ? `background:transparent; border:2px solid ${settings.totalColor}; color:${settings.totalColor}` :
           `background:${settings.totalColor}`;
  }
</script>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-2xl font-bold">Invoice Designer</h1>
      <p class="text-base-content">Desain dan kustomisasi tampilan invoice</p>
    </div>
    {#if saved}
      <div class="badge badge-success gap-2">Tersimpan!</div>
    {/if}
  </div>

  <div class="tabs tabs-boxed bg-base-200 w-fit">
    <button class="tab {activeTab === 'general' ? 'tab-active' : ''}" onclick={() => activeTab = 'general'}><Layout class="w-4 h-4" /> Umum</button>
    <button class="tab {activeTab === 'style' ? 'tab-active' : ''}" onclick={() => activeTab = 'style'}><Palette class="w-4 h-4" /> Warna & Style</button>
    <button class="tab {activeTab === 'images' ? 'tab-active' : ''}" onclick={() => activeTab = 'images'}><Upload class="w-4 h-4" /> Gambar</button>
    <button class="tab {activeTab === 'preview' ? 'tab-active' : ''}" onclick={() => activeTab = 'preview'}><Eye class="w-4 h-4" /> Preview</button>
  </div>

  {#if activeTab === 'general'}
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="glass-card">
        <div class="p-5">
          <h2 class="card-title text-base">Informasi Toko</h2>
          <div class="space-y-3 mt-2">
            <div class="form-control"><label class="label"><span class="label-text">Nama Toko</span></label><input type="text" class="input input-bordered input-sm" bind:value={settings.storeName} /></div>
            <div class="form-control"><label class="label"><span class="label-text">Alamat</span></label><textarea class="textarea textarea-bordered textarea-sm" rows="2" bind:value={settings.storeAddress}></textarea></div>
            <div class="grid grid-cols-2 gap-3">
              <div class="form-control"><label class="label"><span class="label-text">Telepon</span></label><input type="text" class="input input-bordered input-sm" bind:value={settings.storePhone} /></div>
              <div class="form-control"><label class="label"><span class="label-text">Email</span></label><input type="email" class="input input-bordered input-sm" bind:value={settings.storeEmail} /></div>
            </div>
            <div class="grid grid-cols-2 gap-3">
              <div class="form-control"><label class="label"><span class="label-text">Pajak PPN (%)</span></label><input type="number" class="input input-bordered input-sm" bind:value={settings.taxRate} min="0" max="100" /></div>
              <div class="form-control"><label class="label"><span class="label-text">Ukuran Font</span></label>
                <select class="select select-bordered select-sm" bind:value={settings.fontSize}>
                  <option value="11">Kecil (11px)</option><option value="13">Normal (13px)</option><option value="15">Besar (15px)</option>
                </select>
              </div>
            </div>
            <div class="form-control"><label class="label"><span class="label-text">Footer Text</span></label><input type="text" class="input input-bordered input-sm" bind:value={settings.footerText} /></div>
          </div>
        </div>
      </div>

      <div class="glass-card">
        <div class="p-5">
          <h2 class="card-title text-base">Layout Invoice</h2>
          <div class="space-y-3 mt-2">
            <div class="form-control"><label class="label"><span class="label-text">Style Header</span></label>
              <select class="select select-bordered select-sm" bind:value={settings.headerStyle}>
                <option value="modern">Modern (Solid)</option><option value="classic">Classic (Gradient)</option><option value="minimal">Minimal (Border)</option>
              </select>
            </div>
            <div class="form-control"><label class="label"><span class="label-text">Posisi Logo</span></label>
              <div class="flex gap-2">
                {#each ['left', 'center', 'right'] as pos}
                  <button class="btn btn-sm {settings.logoPosition === pos ? 'btn-primary' : 'btn-outline'}" onclick={() => settings.logoPosition = pos as any}>{pos === 'left' ? 'Kiri' : pos === 'center' ? 'Tengah' : 'Kanan'}</button>
                {/each}
              </div>
            </div>
            <div class="form-control"><label class="label"><span class="label-text">Posisi Stamp</span></label>
              <div class="flex gap-2">
                {#each ['left', 'right', 'center'] as pos}
                  <button class="btn btn-sm {settings.stampPosition === pos ? 'btn-primary' : 'btn-outline'}" onclick={() => settings.stampPosition = pos as any}>{pos === 'left' ? 'Kiri' : pos === 'center' ? 'Tengah' : 'Kanan'}</button>
                {/each}
              </div>
            </div>
            <label class="label cursor-pointer justify-start gap-3">
              <input type="checkbox" class="toggle toggle-primary" checked={settings.showLogo} onchange={() => settings.showLogo = !settings.showLogo} />
              <span class="label-text">Tampilkan Logo</span>
            </label>
            <label class="label cursor-pointer justify-start gap-3">
              <input type="checkbox" class="toggle toggle-primary" checked={settings.showStamp} onchange={() => settings.showStamp = !settings.showStamp} />
              <span class="label-text">Tampilkan Stamp</span>
            </label>
          </div>
        </div>
      </div>
    </div>
  {:else if activeTab === 'style'}
    <div class="glass-card">
      <div class="p-5">
        <h2 class="card-title text-base">Warna Invoice</h2>
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mt-2">
          <div class="form-control">
            <label class="label"><span class="label-text">Header Tabel</span></label>
            <div class="flex items-center gap-2">
              <input type="color" class="w-12 h-10 rounded cursor-pointer" bind:value={settings.headerColor} />
              <input type="text" class="input input-bordered input-sm flex-1 font-mono" bind:value={settings.headerColor} />
            </div>
          </div>
          <div class="form-control">
            <label class="label"><span class="label-text">Kotak Total</span></label>
            <div class="flex items-center gap-2">
              <input type="color" class="w-12 h-10 rounded cursor-pointer" bind:value={settings.totalColor} />
              <input type="text" class="input input-bordered input-sm flex-1 font-mono" bind:value={settings.totalColor} />
            </div>
          </div>
          <div class="form-control">
            <label class="label"><span class="label-text">Teks Utama</span></label>
            <div class="flex items-center gap-2">
              <input type="color" class="w-12 h-10 rounded cursor-pointer" bind:value={settings.textColor} />
              <input type="text" class="input input-bordered input-sm flex-1 font-mono" bind:value={settings.textColor} />
            </div>
          </div>
          <div class="form-control">
            <label class="label"><span class="label-text">Border</span></label>
            <div class="flex items-center gap-2">
              <input type="color" class="w-12 h-10 rounded cursor-pointer" bind:value={settings.borderColor} />
              <input type="text" class="input input-bordered input-sm flex-1 font-mono" bind:value={settings.borderColor} />
            </div>
          </div>
        </div>
        <div class="mt-4">
          <label class="label"><span class="label-text">Font Family</span></label>
          <div class="flex gap-2 flex-wrap">
            {#each ['Segoe UI', 'Arial', 'Georgia', 'Courier New', 'Verdana'] as font}
              <button class="btn btn-sm {settings.fontFamily === font ? 'btn-primary' : 'btn-outline'}" style="font-family:{font}" onclick={() => settings.fontFamily = font}>{font}</button>
            {/each}
          </div>
        </div>
      </div>
    </div>
  {:else if activeTab === 'images'}
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <div class="glass-card">
        <div class="p-5">
          <h2 class="card-title text-base">Logo Toko</h2>
          <div class="flex items-center gap-4 mt-2">
            <div class="w-24 h-24 bg-base-200 rounded-lg flex items-center justify-center overflow-hidden">
              {#if logoPreview}
                <img src={logoPreview} alt="Logo" class="w-full h-full object-contain" />
              {:else}
                <span class="text-2xl font-bold text-primary">SC</span>
              {/if}
            </div>
            <div class="flex-1">
              <input type="file" accept="image/*" class="hidden" bind:this={logoInput} onchange={(e) => handleFileUpload(e, 'logo')} />
              <button class="btn btn-outline btn-sm gap-2 w-full" onclick={() => triggerUpload('logo')}><Upload class="w-4 h-4" /> Upload Logo</button>
              {#if logoPreview}
                <button class="btn btn-ghost btn-xs text-error mt-2 w-full" onclick={() => logoPreview = null}>Hapus Logo</button>
              {/if}
            </div>
          </div>
          <p class="text-xs text-base-content mt-2">Format: PNG, JPG, SVG. Maks 2MB.</p>
        </div>
      </div>

      <div class="glass-card">
        <div class="p-5">
          <h2 class="card-title text-base">Stamp / Cap Toko</h2>
          <div class="flex items-center gap-4 mt-2">
            <div class="w-24 h-24 bg-base-200 rounded-lg flex items-center justify-center overflow-hidden">
              {#if stampPreview}
                <img src={stampPreview} alt="Stamp" class="w-full h-full object-contain" />
              {:else}
                <span class="text-xs text-base-content text-center">Belum ada stamp</span>
              {/if}
            </div>
            <div class="flex-1">
              <input type="file" accept="image/*" class="hidden" bind:this={stampInput} onchange={(e) => handleFileUpload(e, 'stamp')} />
              <button class="btn btn-outline btn-sm gap-2 w-full" onclick={() => triggerUpload('stamp')}><Upload class="w-4 h-4" /> Upload Stamp</button>
              {#if stampPreview}
                <button class="btn btn-ghost btn-xs text-error mt-2 w-full" onclick={() => stampPreview = null}>Hapus Stamp</button>
              {/if}
            </div>
          </div>
          <p class="text-xs text-base-content mt-2">Stamp tampil di bagian "Hormat Kami" pada invoice.</p>
        </div>
      </div>
    </div>
  {:else}
    <div class="glass-card">
      <div class="p-5">
        <div class="flex justify-between items-center mb-4">
          <h2 class="card-title text-base">Preview Invoice</h2>
          <select class="select select-bordered select-sm" bind:value={previewTx}>
            {#each mockTransactions as tx}
              <option value={tx}>{tx.invoiceNumber} - {tx.customerName || 'Umum'}</option>
            {/each}
          </select>
        </div>
        <div class="bg-white rounded-xl shadow-lg p-8 mx-auto max-w-2xl" style="font-family:{settings.fontFamily}, sans-serif; font-size:{settings.fontSize}px;">
          <div class="flex items-start gap-4 pb-4 border-b mb-6" style="border-color:{settings.borderColor}">
            {#if settings.showLogo && logoPreview}
              <img src={logoPreview} alt="Logo" class="w-14 h-14 rounded-lg object-contain flex-shrink-0" />
            {:else if settings.showLogo}
              <div class="w-14 h-14 rounded-lg flex items-center justify-center text-white text-xl font-extrabold flex-shrink-0" style="background:{settings.headerColor}">SC</div>
            {/if}
            <div class="flex-1">
              <h1 class="text-xl font-extrabold" style="color:{settings.textColor}">{settings.storeName}</h1>
              <p class="text-xs mt-0.5 leading-snug" style="color:#6b7280">{settings.storeAddress}</p>
              <p class="text-xs font-bold mt-0.5" style="color:{settings.textColor}">{settings.storePhone}</p>
              <p class="text-xs mt-0.5" style="color:{settings.headerColor}">{settings.storeEmail}</p>
            </div>
          </div>

          <div class="flex justify-between mb-6">
            <div>
              <p class="text-[10px] font-bold text-slate-400 uppercase tracking-wide mb-1">Kepada</p>
              <p class="text-sm font-bold" style="color:{settings.textColor}">{previewTx.customerName || 'Pelanggan Umum'}</p>
              <p class="text-xs mt-0.5" style="color:#6b7280">{previewTx.customerAddress || '-'}</p>
            </div>
            <div class="space-y-0.5">
              <div class="flex gap-4 items-center"><span class="text-xs text-slate-600 font-semibold w-20">No. Invoice</span><span class="text-xs font-bold" style="color:{settings.textColor}">{previewTx.invoiceNumber}</span></div>
              <div class="flex gap-4 items-center"><span class="text-xs text-slate-600 font-semibold w-20">Tanggal</span><span class="text-xs font-bold" style="color:{settings.textColor}">{new Date(previewTx.createdAt).toLocaleDateString('id-ID')}</span></div>
              <div class="flex gap-4 items-center"><span class="text-xs text-slate-600 font-semibold w-20">Jatuh Tempo</span><span class="text-xs font-bold" style="color:{settings.textColor}">{previewTx.dueDate ? new Date(previewTx.dueDate).toLocaleDateString('id-ID') : '-'}</span></div>
              <div class="flex gap-4 items-center"><span class="text-xs text-slate-600 font-semibold w-20">Pembayaran</span><span class="text-xs font-bold capitalize" style="color:{settings.textColor}">{previewTx.paymentMethod}</span></div>
            </div>
          </div>

          <div class="border rounded-lg overflow-hidden mb-6" style="border-color:{settings.borderColor}">
            <div class="flex text-white text-[11px] font-bold uppercase tracking-wide" style={getHeaderBg()}>
              <div class="flex-[2] px-3.5 py-2.5">Deskripsi</div>
              <div class="flex-[0.8] px-3.5 py-2.5 text-center">Qty</div>
              <div class="flex-1 px-3.5 py-2.5 text-center">Harga</div>
              <div class="flex-1 px-3.5 py-2.5 text-right">Jumlah</div>
            </div>
            {#each previewTx.items as item}
              <div class="flex border-b" style="border-color:{settings.borderColor}">
                <div class="flex-[2] px-3.5 py-2.5"><span class="font-semibold" style="color:{settings.textColor}">{item.name}</span></div>
                <div class="flex-[0.8] px-3.5 py-2.5 text-center font-semibold" style="color:#6b7280">{item.quantity}</div>
                <div class="flex-1 px-3.5 py-2.5 text-center" style="color:#6b7280">Rp {item.price.toLocaleString('id-ID')}</div>
                <div class="flex-1 px-3.5 py-2.5 text-right font-bold" style="color:{settings.textColor}">Rp {item.subtotal.toLocaleString('id-ID')}</div>
              </div>
            {/each}
          </div>

          <div class="flex justify-end mb-8">
            <div style="width:240px">
              <div class="flex justify-between py-0.5 text-xs" style="color:#374151"><span>Subtotal</span><span>Rp {previewTx.subtotal.toLocaleString('id-ID')}</span></div>
              {#if previewTx.discount > 0}
                <div class="flex justify-between py-0.5 text-xs text-red-600"><span>Diskon</span><span>-Rp {previewTx.discount.toLocaleString('id-ID')}</span></div>
              {/if}
              {#if settings.taxRate > 0}
                <div class="flex justify-between py-0.5 text-xs" style="color:#374151"><span>Pajak PPN {settings.taxRate}%</span><span>Rp {Math.round((previewTx.subtotal - (previewTx.discount || 0)) * settings.taxRate / 100).toLocaleString('id-ID')}</span></div>
              {/if}
              <div class="flex justify-between items-center mt-2 px-4 py-2.5 rounded-lg text-white" style={getTotalBg()}>
                <span class="font-bold text-xs uppercase tracking-wide">Total</span>
                <span class="font-bold text-lg italic">Rp {(previewTx.subtotal - (previewTx.discount || 0) + Math.round((previewTx.subtotal - (previewTx.discount || 0)) * settings.taxRate / 100)).toLocaleString('id-ID')}</span>
              </div>
            </div>
          </div>

          <div class="flex justify-between mt-16">
            <div class="flex flex-col items-center">
              <p class="text-[11px] font-bold uppercase tracking-wide mb-16" style="color:{settings.textColor}">Tanda Terima</p>
              <div class="w-44" style="border-top:1px solid {settings.textColor}"></div>
            </div>
            {#if settings.showStamp}
              <div class="flex flex-col items-center">
                <p class="text-[11px] font-bold uppercase tracking-wide mb-2" style="color:{settings.textColor}">Hormat Kami</p>
                {#if stampPreview}
                  <img src={stampPreview} alt="Stamp" class="w-20 h-auto mb-1" />
                {:else}
                  <div class="w-14 h-14 rounded-full mb-1 flex items-center justify-center text-white text-xl font-extrabold" style="background:{settings.headerColor}">SC</div>
                {/if}
                <div class="w-44" style="border-top:1px solid {settings.textColor}"></div>
              </div>
            {/if}
          </div>

          <!-- Spacer to push lines to bottom -->
          <div class="h-8"></div>
          <p class="text-center mt-8 text-[10px] tracking-[2px] uppercase" style="color:#d1d5db">{settings.footerText}</p>
        </div>
      </div>
    </div>
  {/if}

  <div class="flex justify-end">
    <button class="btn btn-primary gap-2" onclick={handleSave}><Save class="w-4 h-4" /> Simpan Desain</button>
  </div>
</div>