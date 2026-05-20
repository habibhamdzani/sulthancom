<script lang="ts">
  import { mockProducts, mockCategories, mockCustomers } from '$lib/data/mock';
  import { formatCurrency, generateInvoiceNumber } from '$lib/utils/helpers';
  import { cart, cartSubtotal, cartTotal, cartDiscount, cartTax, addToCart, updateQuantity, removeFromCart, clearCart, selectedPaymentMethod, selectedCustomer, transactionNotes } from '$lib/stores/cart';
  import { Search, Minus, Plus, Trash2, ShoppingCart, User, Receipt, X } from 'lucide-svelte';
  import { currentUser } from '$lib/stores/auth';

  let searchQuery = $state('');
  let selectedCategory = $state('all');
  let paymentStatus = $state('paid');
  let dpAmount = $state('');
  let dueDate = $state('');
  let showPaymentModal = $state(false);
  let showSuccessModal = $state(false);
  let lastTransaction = $state<any>(null);

  let filteredProducts = $derived.by(() => {
    let result = mockProducts.filter(p => p.type !== 'service' && p.isActive);
    if (searchQuery) {
      const q = searchQuery.toLowerCase();
      result = result.filter(p => p.name.toLowerCase().includes(q) || p.sku.toLowerCase().includes(q));
    }
    if (selectedCategory !== 'all') {
      result = result.filter(p => p.categoryId === selectedCategory);
    }
    return result;
  });

  function getCategoryName(id: string) {
    return mockCategories.find(c => c.id === id)?.name || '';
  }

  function handleAddToCart(product: any) {
    addToCart({
      productId: product.id,
      name: product.name,
      price: product.sellingPrice,
      quantity: 1,
      discount: 0,
    });
  }

  function handleCheckout() {
    if (cart.length === 0) return;
    showPaymentModal = true;
  }

  function processPayment() {
    const total = $cartTotal;
    let paidAmount = total;
    let remainingAmount = 0;

    if (paymentStatus === 'dp') {
      paidAmount = parseInt(dpAmount) || 0;
      remainingAmount = total - paidAmount;
    } else if (paymentStatus === 'unpaid') {
      paidAmount = 0;
      remainingAmount = total;
    }

    lastTransaction = {
      invoiceNumber: generateInvoiceNumber(),
      items: [...$cart],
      subtotal: $cartSubtotal,
      discount: $cartDiscount,
      tax: $cartTax,
      total: total,
      paymentMethod: $selectedPaymentMethod,
      paymentStatus: paymentStatus,
      paidAmount: paidAmount,
      remainingAmount: remainingAmount,
      customer: $selectedCustomer,
      cashier: $currentUser?.name,
      date: new Date().toISOString(),
    };

    showPaymentModal = false;
    showSuccessModal = true;
    clearCart();
  }

  const paymentMethods = [
    { id: 'cash', label: 'Tunai', icon: '💵' },
    { id: 'dana', label: 'DANA', icon: '💙' },
    { id: 'gopay', label: 'GoPay', icon: '💚' },
    { id: 'shopeepay', label: 'ShopeePay', icon: '🧡' },
    { id: 'bca', label: 'BCA', icon: '🏦' },
    { id: 'bri', label: 'BRI', icon: '🏦' },
    { id: 'bni', label: 'BNI', icon: '🏦' },
    { id: 'mandiri', label: 'Mandiri', icon: '🏦' },
  ];
</script>

<div class="h-[calc(100vh-7rem)] flex gap-4">
  <div class="flex-1 flex flex-col min-w-0">
    <div class="flex items-center gap-3 mb-4">
      <div class="form-control flex-1">
        <div class="relative">
          <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content" />
          <input type="text" class="input input-bordered input-sm w-full pl-10" placeholder="Scan barcode atau cari produk..." bind:value={searchQuery} />
        </div>
      </div>
      <select class="select select-bordered select-sm" bind:value={selectedCategory}>
        <option value="all">Semua</option>
        {#each mockCategories as cat}
          <option value={cat.id}>{cat.name}</option>
        {/each}
      </select>
    </div>

    <div class="flex-1 overflow-y-auto">
      <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-3">
        {#each filteredProducts as product}
          <div class="glass-card pos-product-card" onclick={() => handleAddToCart(product)}>
            <div class="p-4">
              <div class="w-full h-24 bg-base-200/50 rounded-lg flex items-center justify-center mb-2">
                <span class="text-3xl icon-3d">📦</span>
              </div>
              <h3 class="font-medium text-sm truncate">{product.name}</h3>
              <p class="text-xs text-base-content">{getCategoryName(product.categoryId)}</p>
              <p class="font-bold text-primary text-sm">{formatCurrency(product.sellingPrice)}</p>
              <p class="text-xs {product.stock <= product.minStock ? 'text-error' : 'text-success'}">Stok: {product.stock}</p>
            </div>
          </div>
        {/each}
      </div>
    </div>
  </div>

  <div class="w-96 glass-card rounded-xl flex flex-col">
    <div class="p-4 border-b border-base-200/50">
      <div class="flex items-center justify-between">
        <h2 class="font-bold flex items-center gap-2">
          <ShoppingCart class="w-5 h-5 icon-3d" /> Keranjang
        </h2>
        {#if $cart.length > 0}
          <button class="btn btn-ghost btn-xs text-error" onclick={clearCart}>
            <X class="w-4 h-4 icon-3d" /> Kosongkan
          </button>
        {/if}
      </div>
    </div>

    <div class="flex-1 overflow-y-auto">
      {#if $cart.length === 0}
        <div class="flex flex-col items-center justify-center h-full text-base-content">
          <ShoppingCart class="w-16 h-16 mb-4" />
          <p class="text-sm">Keranjang kosong</p>
          <p class="text-xs">Pilih produk untuk memulai transaksi</p>
        </div>
      {:else}
        {#each $cart as item}
          <div class="cart-item">
            <div class="flex-1 min-w-0">
              <p class="font-medium text-sm truncate">{item.name}</p>
              <p class="text-xs text-base-content">{formatCurrency(item.price)}</p>
            </div>
            <div class="flex items-center gap-2">
              <button class="btn btn-ghost btn-xs btn-circle" onclick={() => updateQuantity(item.productId, item.quantity - 1)}>
                <Minus class="w-3 h-3" />
              </button>
              <span class="text-sm font-medium w-8 text-center">{item.quantity}</span>
              <button class="btn btn-ghost btn-xs btn-circle" onclick={() => updateQuantity(item.productId, item.quantity + 1)}>
                <Plus class="w-3 h-3" />
              </button>
            </div>
            <p class="font-medium text-sm w-24 text-right">{formatCurrency(item.subtotal)}</p>
            <button class="btn btn-ghost btn-xs text-error" onclick={() => removeFromCart(item.productId)}>
              <Trash2 class="w-3 h-3" />
            </button>
          </div>
        {/each}
      {/if}
    </div>

    {#if $cart.length > 0}
      <div class="p-4 border-t border-base-200 space-y-3">
        <div class="flex justify-between text-sm">
          <span class="text-base-content">Subtotal</span>
          <span>{formatCurrency($cartSubtotal)}</span>
        </div>
        <div class="flex justify-between text-sm">
          <span class="text-base-content">Diskon</span>
          <span>- {formatCurrency($cartDiscount)}</span>
        </div>
        <div class="flex justify-between text-lg font-bold border-t border-base-200 pt-2">
          <span>Total</span>
          <span class="text-primary">{formatCurrency($cartTotal)}</span>
        </div>

        <button class="btn btn-primary w-full" onclick={handleCheckout}>
          <Receipt class="w-4 h-4" /> Bayar {formatCurrency($cartTotal)}
        </button>
      </div>
    {/if}
  </div>
</div>

{#if showPaymentModal}
  <div class="modal modal-open">
    <div class="modal-box w-11/12 max-w-2xl">
      <h3 class="text-lg font-bold mb-4">Pembayaran</h3>

      <div class="space-y-4">
        <div>
          <label class="label"><span class="label-text font-medium">Metode Pembayaran</span></label>
          <div class="grid grid-cols-4 gap-2">
            {#each paymentMethods as method}
                <button class="btn btn-sm {$selectedPaymentMethod === method.id ? 'btn-primary' : 'btn-outline'}" onclick={() => $selectedPaymentMethod = method.id}>
                {method.label}
              </button>
            {/each}
          </div>
        </div>

        <div>
          <label class="label"><span class="label-text font-medium">Pelanggan (opsional)</span></label>
            <select class="select select-bordered w-full" bind:value={$selectedCustomer}>
            <option value="">Umum</option>
            {#each mockCustomers as customer}
              <option value={customer.id}>{customer.name}</option>
            {/each}
          </select>
        </div>

        <div>
          <label class="label"><span class="label-text font-medium">Status Pembayaran</span></label>
          <div class="flex gap-2">
            <button class="btn btn-sm {paymentStatus === 'paid' ? 'btn-primary' : 'btn-outline'}" onclick={() => paymentStatus = 'paid'}>Lunas</button>
            <button class="btn btn-sm {paymentStatus === 'dp' ? 'btn-primary' : 'btn-outline'}" onclick={() => paymentStatus = 'dp'}>DP</button>
            <button class="btn btn-sm {paymentStatus === 'unpaid' ? 'btn-primary' : 'btn-outline'}" onclick={() => paymentStatus = 'unpaid'}>Piutang</button>
          </div>
        </div>

        {#if paymentStatus === 'dp'}
          <div>
            <label class="label"><span class="label-text font-medium">Nominal DP</span></label>
            <input type="number" class="input input-bordered w-full" bind:value={dpAmount} placeholder="Masukkan nominal DP" />
            <p class="text-xs text-base-content mt-1">Sisa: {formatCurrency($cartTotal - (parseInt(dpAmount) || 0))}</p>
          </div>
        {/if}

        {#if paymentStatus === 'dp' || paymentStatus === 'unpaid'}
          <div>
            <label class="label"><span class="label-text font-medium">Jatuh Tempo</span></label>
            <input type="date" class="input input-bordered w-full" bind:value={dueDate} />
          </div>
        {/if}

        <div>
          <label class="label"><span class="label-text font-medium">Catatan</span></label>
          <textarea class="textarea textarea-bordered w-full" bind:value={$transactionNotes} rows="2" placeholder="Catatan transaksi (opsional)"></textarea>
        </div>
      </div>

      <div class="modal-action">
        <button class="btn btn-ghost" onclick={() => showPaymentModal = false}>Batal</button>
        <button class="btn btn-primary" onclick={processPayment}>Proses Pembayaran</button>
      </div>
    </div>
  </div>
{/if}

{#if showSuccessModal && lastTransaction}
  <div class="modal modal-open">
    <div class="modal-box text-center">
      <div class="w-16 h-16 bg-success/10 rounded-full flex items-center justify-center mx-auto mb-4">
        <svg class="w-8 h-8 text-success" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
        </svg>
      </div>
      <h3 class="text-lg font-bold">Transaksi Berhasil!</h3>
      <p class="text-base-content mb-4">{lastTransaction.invoiceNumber}</p>

      <div class="text-left bg-base-200 rounded-lg p-4 mb-4">
        <div class="space-y-2 text-sm">
          <div class="flex justify-between"><span>Total</span><span class="font-bold">{formatCurrency(lastTransaction.total)}</span></div>
          <div class="flex justify-between"><span>Metode</span><span>{lastTransaction.paymentMethod}</span></div>
          <div class="flex justify-between"><span>Status</span><span>{lastTransaction.paymentStatus}</span></div>
          {#if lastTransaction.paidAmount < lastTransaction.total}
            <div class="flex justify-between"><span>Sisa</span><span class="text-error">{formatCurrency(lastTransaction.remainingAmount)}</span></div>
          {/if}
        </div>
      </div>

      <div class="flex gap-2 justify-center">
        <button class="btn btn-outline btn-sm" onclick={() => showSuccessModal = false}>Tutup</button>
        <button class="btn btn-primary btn-sm" onclick={() => showSuccessModal = false}>Cetak Struk</button>
      </div>
    </div>
  </div>
{/if}
