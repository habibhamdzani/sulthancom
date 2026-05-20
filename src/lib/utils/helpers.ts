export function formatCurrency(amount: number): string {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(amount);
}

export function formatDate(date: string | Date): string {
  const d = typeof date === 'string' ? new Date(date) : date;
  return new Intl.DateTimeFormat('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  }).format(d);
}

export function formatDateTime(date: string | Date): string {
  const d = typeof date === 'string' ? new Date(date) : date;
  return new Intl.DateTimeFormat('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(d);
}

export function generateId(): string {
  return Date.now().toString(36) + Math.random().toString(36).substr(2);
}

export function generateInvoiceNumber(): string {
  const date = new Date();
  const dateStr = date.toISOString().slice(0, 10).replace(/-/g, '');
  const random = Math.floor(Math.random() * 999).toString().padStart(3, '0');
  return `INV-${dateStr}-${random}`;
}

export function generateServiceNumber(): string {
  const date = new Date();
  const dateStr = date.toISOString().slice(0, 10).replace(/-/g, '');
  const random = Math.floor(Math.random() * 999).toString().padStart(3, '0');
  return `SRV-${dateStr}-${random}`;
}

export function getRoleLabel(role: string): string {
  const labels: Record<string, string> = {
    super_admin: 'Super Admin',
    admin: 'Admin',
    kasir: 'Kasir',
    teknisi: 'Teknisi',
  };
  return labels[role] || role;
}

export function getServiceStatusLabel(status: string): string {
  const labels: Record<string, string> = {
    waiting: 'Menunggu',
    checked: 'Dicek',
    in_progress: 'Dalam Pengerjaan',
    waiting_parts: 'Menunggu Sparepart',
    completed: 'Selesai',
    picked_up: 'Sudah Diambil',
    cancelled: 'Dibatalkan',
  };
  return labels[status] || status;
}

export function getServiceStatusColor(status: string): string {
  const colors: Record<string, string> = {
    waiting: 'badge-ghost',
    checked: 'badge-info',
    in_progress: 'badge-warning',
    waiting_parts: 'badge-error',
    completed: 'badge-success',
    picked_up: 'badge-primary',
    cancelled: 'badge-neutral',
  };
  return colors[status] || 'badge-ghost';
}

export function getPaymentStatusLabel(status: string): string {
  const labels: Record<string, string> = {
    paid: 'Lunas',
    dp: 'DP',
    unpaid: 'Belum Lunas',
    partial: 'Sebagian',
  };
  return labels[status] || status;
}

export function getPaymentStatusColor(status: string): string {
  const colors: Record<string, string> = {
    paid: 'badge-success',
    dp: 'badge-warning',
    unpaid: 'badge-error',
    partial: 'badge-info',
  };
  return colors[status] || 'badge-ghost';
}
