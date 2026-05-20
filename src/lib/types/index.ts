export type Role = 'super_admin' | 'admin' | 'kasir' | 'teknisi';

export interface User {
  id: string;
  name: string;
  email: string;
  role: Role;
  avatar?: string;
  phone?: string;
  isActive: boolean;
  createdAt: string;
}

export interface Category {
  id: string;
  name: string;
  icon?: string;
  color?: string;
  parentId?: string;
  isActive: boolean;
  productCount: number;
}

export interface Brand {
  id: string;
  name: string;
  logo?: string;
  isActive: boolean;
  productCount: number;
}

export type ProductType = 'physical' | 'service' | 'sparepart' | 'bundle';

export interface Product {
  id: string;
  name: string;
  sku: string;
  barcode?: string;
  categoryId: string;
  brandId?: string;
  type: ProductType;
  costPrice: number;
  sellingPrice: number;
  stock: number;
  minStock: number;
  description?: string;
  images: string[];
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface Supplier {
  id: string;
  name: string;
  phone?: string;
  email?: string;
  address?: string;
  isActive: boolean;
  createdAt: string;
}

export type StockMovementType = 'in' | 'out' | 'adjustment' | 'transfer';

export interface StockMovement {
  id: string;
  productId: string;
  productName: string;
  type: StockMovementType;
  quantity: number;
  previousStock: number;
  newStock: number;
  notes?: string;
  createdAt: string;
  createdBy: string;
}

export interface Customer {
  id: string;
  name: string;
  phone?: string;
  email?: string;
  address?: string;
  notes?: string;
  totalTransactions: number;
  totalSpent: number;
  createdAt: string;
}

export type PaymentMethod = 'cash' | 'dana' | 'gopay' | 'shopeepay' | 'bca' | 'bri' | 'bni' | 'mandiri' | 'other';

export type PaymentStatus = 'paid' | 'dp' | 'unpaid' | 'partial';

export interface CartItem {
  productId: string;
  name: string;
  price: number;
  quantity: number;
  discount: number;
  subtotal: number;
}

export interface Transaction {
  id: string;
  invoiceNumber: string;
  customerId?: string;
  customerName?: string;
  customerAddress?: string;
  items: CartItem[];
  subtotal: number;
  discount: number;
  tax: number;
  total: number;
  paymentMethod: PaymentMethod;
  paymentStatus: PaymentStatus;
  paidAmount: number;
  remainingAmount: number;
  dpAmount?: number;
  dueDate?: string;
  notes?: string;
  cashierId: string;
  cashierName: string;
  createdAt: string;
}

export type ServiceStatus = 'waiting' | 'checked' | 'in_progress' | 'waiting_parts' | 'completed' | 'picked_up' | 'cancelled';

export interface ServiceTicket {
  id: string;
  ticketNumber: string;
  customerId?: string;
  customerName: string;
  customerPhone?: string;
  deviceType: string;
  deviceBrand?: string;
  serialNumber?: string;
  complaint: string;
  diagnosis?: string;
  estimatedCost: number;
  actualCost?: number;
  estimatedCompletion?: string;
  status: ServiceStatus;
  sparepartsUsed: {
    productId: string;
    name: string;
    quantity: number;
    price: number;
  }[];
  technicianNotes: {
    id: string;
    note: string;
    createdAt: string;
    createdBy: string;
  }[];
  photos: string[];
  technicianId: string;
  technicianName: string;
  createdAt: string;
  updatedAt: string;
}

export interface Receivable {
  id: string;
  transactionId: string;
  invoiceNumber: string;
  customerId: string;
  customerName: string;
  totalAmount: number;
  paidAmount: number;
  remainingAmount: number;
  dueDate: string;
  status: 'unpaid' | 'partial' | 'paid' | 'overdue';
  payments: {
    id: string;
    amount: number;
    method: PaymentMethod;
    date: string;
    notes?: string;
  }[];
  createdAt: string;
}

export interface StoreSettings {
  name: string;
  address: string;
  phone: string;
  email: string;
  website?: string;
  logo?: string;
  taxRate: number;
  currency: string;
  invoiceFooter?: string;
  theme: string;
}

export interface DashboardStats {
  todaySales: number;
  todayTransactions: number;
  activeServices: number;
  totalReceivables: number;
  totalProducts: number;
  lowStockProducts: number;
  bestSellingProduct: {
    name: string;
    quantity: number;
  };
  slowMovingProduct: {
    name: string;
    quantity: number;
  };
}
