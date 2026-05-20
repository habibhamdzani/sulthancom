import { writable, derived } from 'svelte/store';
import type { CartItem } from '$lib/types';

export const cart = writable<CartItem[]>([]);
export const cartDiscount = writable(0);
export const cartTax = writable(0);
export const selectedPaymentMethod = writable('cash');
export const selectedCustomer = writable<string | null>(null);
export const transactionNotes = writable('');

export const cartSubtotal = derived(cart, ($cart) =>
  $cart.reduce((sum, item) => sum + item.subtotal, 0)
);

export const cartTotal = derived(
  [cartSubtotal, cartDiscount, cartTax],
  ([$subtotal, $discount, $tax]) => $subtotal - $discount + $tax
);

export function addToCart(item: Omit<CartItem, 'subtotal'>) {
  cart.update((items) => {
    const existing = items.find((i) => i.productId === item.productId);
    if (existing) {
      return items.map((i) =>
        i.productId === item.productId
          ? { ...i, quantity: i.quantity + item.quantity, subtotal: (i.quantity + item.quantity) * i.price }
          : i
      );
    }
    return [...items, { ...item, subtotal: item.quantity * item.price }];
  });
}

export function updateQuantity(productId: string, quantity: number) {
  if (quantity <= 0) {
    removeFromCart(productId);
    return;
  }
  cart.update((items) =>
    items.map((i) =>
      i.productId === productId ? { ...i, quantity, subtotal: quantity * i.price } : i
    )
  );
}

export function removeFromCart(productId: string) {
  cart.update((items) => items.filter((i) => i.productId !== productId));
}

export function clearCart() {
  cart.set([]);
  cartDiscount.set(0);
  cartTax.set(0);
  selectedPaymentMethod.set('cash');
  selectedCustomer.set(null);
  transactionNotes.set('');
}
