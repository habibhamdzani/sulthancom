import type { User, Category, Brand, Product, Supplier, Customer, ServiceTicket, Receivable, Transaction, StockMovement, StoreSettings, DashboardStats } from '$lib/types';

export const mockUsers: User[] = [
  { id: '1', name: 'Sulthan Admin', email: 'admin@sulthancom.com', role: 'super_admin', isActive: true, createdAt: '2024-01-01' },
  { id: '2', name: 'Budi Inventory', email: 'budi@sulthancom.com', role: 'admin', isActive: true, createdAt: '2024-01-15' },
  { id: '3', name: 'Ani Kasir', email: 'ani@sulthancom.com', role: 'kasir', isActive: true, createdAt: '2024-02-01' },
  { id: '4', name: 'Rudi Teknisi', email: 'rudi@sulthancom.com', role: 'teknisi', isActive: true, createdAt: '2024-02-15' },
];

export const mockCategories: Category[] = [
  { id: '1', name: 'Laptop', icon: 'laptop', color: '#3B82F6', isActive: true, productCount: 25 },
  { id: '2', name: 'PC Rakitan', icon: 'monitor', color: '#8B5CF6', isActive: true, productCount: 15 },
  { id: '3', name: 'Printer', icon: 'printer', color: '#10B981', isActive: true, productCount: 10 },
  { id: '4', name: 'Monitor', icon: 'tv', color: '#F59E0B', isActive: true, productCount: 12 },
  { id: '5', name: 'Sparepart', icon: 'cpu', color: '#EF4444', isActive: true, productCount: 150 },
  { id: '6', name: 'Networking', icon: 'wifi', color: '#06B6D4', isActive: true, productCount: 30 },
  { id: '7', name: 'Aksesoris', icon: 'headphones', color: '#EC4899', isActive: true, productCount: 80 },
  { id: '8', name: 'Service', icon: 'wrench', color: '#6366F1', isActive: true, productCount: 10 },
];

export const mockBrands: Brand[] = [
  { id: '1', name: 'ASUS', isActive: true, productCount: 45 },
  { id: '2', name: 'Lenovo', isActive: true, productCount: 35 },
  { id: '3', name: 'HP', isActive: true, productCount: 30 },
  { id: '4', name: 'Dell', isActive: true, productCount: 25 },
  { id: '5', name: 'Acer', isActive: true, productCount: 20 },
  { id: '6', name: 'MSI', isActive: true, productCount: 18 },
  { id: '7', name: 'Logitech', isActive: true, productCount: 40 },
  { id: '8', name: 'Epson', isActive: true, productCount: 15 },
];

export const mockProducts: Product[] = [
  { id: '1', name: 'ASUS VivoBook 14', sku: 'LP-ASU-VB14', barcode: '8901234567890', categoryId: '1', brandId: '1', type: 'physical', costPrice: 7500000, sellingPrice: 8999000, stock: 12, minStock: 5, description: 'Laptop ASUS VivoBook 14 inch, Intel Core i5, RAM 8GB, SSD 512GB', images: [], isActive: true, createdAt: '2024-03-01', updatedAt: '2024-03-15' },
  { id: '2', name: 'Lenovo ThinkPad E14', sku: 'LP-LNV-TP14', barcode: '8901234567891', categoryId: '1', brandId: '2', type: 'physical', costPrice: 9000000, sellingPrice: 10999000, stock: 8, minStock: 3, description: 'Laptop Lenovo ThinkPad E14, Intel Core i7, RAM 16GB, SSD 512GB', images: [], isActive: true, createdAt: '2024-03-01', updatedAt: '2024-03-15' },
  { id: '3', name: 'PC Gaming RTX 4060', sku: 'PC-GAM-4060', categoryId: '2', brandId: '6', type: 'physical', costPrice: 12000000, sellingPrice: 14999000, stock: 5, minStock: 2, description: 'PC Gaming dengan RTX 4060, Ryzen 5 7600X, RAM 16GB, SSD 1TB', images: [], isActive: true, createdAt: '2024-03-05', updatedAt: '2024-03-15' },
  { id: '4', name: 'PC Office Core i3', sku: 'PC-OFF-I3', categoryId: '2', type: 'physical', costPrice: 4500000, sellingPrice: 5499000, stock: 15, minStock: 5, description: 'PC Office dengan Intel Core i3, RAM 8GB, SSD 256GB', images: [], isActive: true, createdAt: '2024-03-05', updatedAt: '2024-03-15' },
  { id: '5', name: 'Epson L3250', sku: 'PR-ESP-L3250', barcode: '8901234567892', categoryId: '3', brandId: '8', type: 'physical', costPrice: 2800000, sellingPrice: 3299000, stock: 10, minStock: 3, description: 'Printer Epson L3250 Inkjet All-in-One', images: [], isActive: true, createdAt: '2024-03-10', updatedAt: '2024-03-15' },
  { id: '6', name: 'Monitor LG 24 inch', sku: 'MN-LG-24', barcode: '8901234567893', categoryId: '4', type: 'physical', costPrice: 1500000, sellingPrice: 1899000, stock: 20, minStock: 5, description: 'Monitor LG 24 inch IPS Full HD', images: [], isActive: true, createdAt: '2024-03-10', updatedAt: '2024-03-15' },
  { id: '7', name: 'Intel Core i5-13400F', sku: 'SP-INT-I5', barcode: '8901234567894', categoryId: '5', type: 'sparepart', costPrice: 2200000, sellingPrice: 2599000, stock: 25, minStock: 10, description: 'Processor Intel Core i5-13400F', images: [], isActive: true, createdAt: '2024-03-12', updatedAt: '2024-03-15' },
  { id: '8', name: 'ASUS ROG Strix B760-F', sku: 'SP-ASU-B760', categoryId: '5', brandId: '1', type: 'sparepart', costPrice: 2800000, sellingPrice: 3299000, stock: 15, minStock: 5, description: 'Motherboard ASUS ROG Strix B760-F', images: [], isActive: true, createdAt: '2024-03-12', updatedAt: '2024-03-15' },
  { id: '9', name: 'RAM DDR5 16GB', sku: 'SP-RAM-D516', categoryId: '5', type: 'sparepart', costPrice: 650000, sellingPrice: 799000, stock: 40, minStock: 15, description: 'RAM DDR5 16GB 4800MHz', images: [], isActive: true, createdAt: '2024-03-12', updatedAt: '2024-03-15' },
  { id: '10', name: 'SSD NVMe 512GB', sku: 'SP-SSD-512', categoryId: '5', type: 'sparepart', costPrice: 500000, sellingPrice: 649000, stock: 50, minStock: 20, description: 'SSD NVMe M.2 512GB', images: [], isActive: true, createdAt: '2024-03-12', updatedAt: '2024-03-15' },
  { id: '11', name: 'Logitech G102', sku: 'AX-LOG-G102', barcode: '8901234567895', categoryId: '7', brandId: '7', type: 'physical', costPrice: 150000, sellingPrice: 199000, stock: 100, minStock: 30, description: 'Mouse Gaming Logitech G102', images: [], isActive: true, createdAt: '2024-03-15', updatedAt: '2024-03-15' },
  { id: '12', name: 'Keyboard Mechanical RK61', sku: 'AX-RK-RK61', categoryId: '7', type: 'physical', costPrice: 350000, sellingPrice: 449000, stock: 45, minStock: 15, description: 'Keyboard Mechanical RK61 RGB', images: [], isActive: true, createdAt: '2024-03-15', updatedAt: '2024-03-15' },
  { id: '13', name: 'Service Install Ulang OS', sku: 'SRV-OS', categoryId: '8', type: 'service', costPrice: 0, sellingPrice: 150000, stock: 999, minStock: 0, description: 'Jasa install ulang Windows/Linux', images: [], isActive: true, createdAt: '2024-03-01', updatedAt: '2024-03-15' },
  { id: '14', name: 'Service Ganti LCD Laptop', sku: 'SRV-LCD', categoryId: '8', type: 'service', costPrice: 0, sellingPrice: 200000, stock: 999, minStock: 0, description: 'Jasa ganti LCD laptop (tidak termasuk sparepart)', images: [], isActive: true, createdAt: '2024-03-01', updatedAt: '2024-03-15' },
  { id: '15', name: 'Headset Gaming Fantech', sku: 'AX-FT-HG', categoryId: '7', type: 'physical', costPrice: 180000, sellingPrice: 249000, stock: 3, minStock: 10, description: 'Headset Gaming Fantech HG11', images: [], isActive: true, createdAt: '2024-03-15', updatedAt: '2024-03-15' },
  { id: '16', name: 'TP-Link Archer C6', sku: 'NT-TPL-C6', categoryId: '6', type: 'physical', costPrice: 350000, sellingPrice: 449000, stock: 18, minStock: 5, description: 'Router WiFi TP-Link Archer C6 AC1200', images: [], isActive: true, createdAt: '2024-03-18', updatedAt: '2024-03-18' },
];

export const mockSuppliers: Supplier[] = [
  { id: '1', name: 'PT Distributor Komputer', phone: '021-12345678', email: 'sales@distkom.com', address: 'Jakarta', isActive: true, createdAt: '2024-01-01' },
  { id: '2', name: 'CV Teknik Mandiri', phone: '021-87654321', email: 'info@teknikmandiri.com', address: 'Surabaya', isActive: true, createdAt: '2024-01-15' },
  { id: '3', name: 'UD Elektronik Jaya', phone: '031-11122233', email: 'order@elektronikjaya.com', address: 'Surabaya', isActive: true, createdAt: '2024-02-01' },
];

export const mockCustomers: Customer[] = [
  { id: '1', name: 'Ahmad Fauzi', phone: '081234567890', email: 'ahmad@email.com', address: 'Jl. Merdeka No. 10', totalTransactions: 5, totalSpent: 15000000, createdAt: '2024-01-10' },
  { id: '2', name: 'Siti Nurhaliza', phone: '081234567891', email: 'siti@email.com', address: 'Jl. Sudirman No. 25', totalTransactions: 3, totalSpent: 8500000, createdAt: '2024-02-05' },
  { id: '3', name: 'Bambang Setiawan', phone: '081234567892', address: 'Jl. Gatot Subroto No. 5', totalTransactions: 8, totalSpent: 22000000, createdAt: '2024-01-20' },
  { id: '4', name: 'Dewi Lestari', phone: '081234567893', email: 'dewi@email.com', address: 'Jl. Ahmad Yani No. 15', totalTransactions: 2, totalSpent: 4500000, createdAt: '2024-03-01' },
  { id: '5', name: 'Rizky Pratama', phone: '081234567894', address: 'Jl. Diponegoro No. 30', totalTransactions: 1, totalSpent: 150000, createdAt: '2024-03-10' },
];

export const mockTransactions: Transaction[] = [
  { id: '1', invoiceNumber: 'INV-20240315-001', customerId: '1', customerName: 'Ahmad Fauzi', customerAddress: 'Jl. Merdeka No. 45, Mataram', items: [{ productId: '1', name: 'ASUS VivoBook 14', price: 8999000, quantity: 1, discount: 0, subtotal: 8999000 }, { productId: '11', name: 'Logitech G102', price: 199000, quantity: 1, discount: 0, subtotal: 199000 }], subtotal: 9198000, discount: 100000, tax: 0, total: 9098000, paymentMethod: 'bca', paymentStatus: 'paid', paidAmount: 9098000, remainingAmount: 0, cashierId: '3', cashierName: 'Ani Kasir', createdAt: '2024-03-15T10:30:00' },
  { id: '2', invoiceNumber: 'INV-20240315-002', items: [{ productId: '13', name: 'Service Install Ulang OS', price: 150000, quantity: 1, discount: 0, subtotal: 150000 }], subtotal: 150000, discount: 0, tax: 0, total: 150000, paymentMethod: 'cash', paymentStatus: 'paid', paidAmount: 150000, remainingAmount: 0, cashierId: '3', cashierName: 'Ani Kasir', createdAt: '2024-03-15T11:00:00' },
  { id: '3', invoiceNumber: 'INV-20240315-003', customerId: '2', customerName: 'Siti Nurhaliza', customerAddress: 'Perumahan Green Hills Blok A3, Lombok Barat', items: [{ productId: '3', name: 'PC Gaming RTX 4060', price: 14999000, quantity: 1, discount: 500000, subtotal: 14499000 }], subtotal: 14499000, discount: 0, tax: 0, total: 14499000, paymentMethod: 'dana', paymentStatus: 'dp', paidAmount: 5000000, remainingAmount: 9499000, dpAmount: 5000000, dueDate: '2024-04-15', cashierId: '3', cashierName: 'Ani Kasir', createdAt: '2024-03-15T14:00:00' },
  { id: '4', invoiceNumber: 'INV-20240316-001', customerId: '3', customerName: 'Bambang Setiawan', items: [{ productId: '5', name: 'Epson L3250', price: 3299000, quantity: 2, discount: 0, subtotal: 6598000 }], subtotal: 6598000, discount: 200000, tax: 0, total: 6398000, paymentMethod: 'cash', paymentStatus: 'unpaid', paidAmount: 0, remainingAmount: 6398000, dueDate: '2024-04-16', notes: 'Bayar tempo 30 hari', cashierId: '3', cashierName: 'Ani Kasir', createdAt: '2024-03-16T09:00:00' },
];

export const mockStockMovements: StockMovement[] = [
  { id: '1', productId: '1', productName: 'ASUS VivoBook 14', type: 'in', quantity: 10, previousStock: 2, newStock: 12, notes: 'Restock dari distributor', createdAt: '2024-03-15T08:00:00', createdBy: 'Budi Inventory' },
  { id: '2', productId: '11', productName: 'Logitech G102', type: 'out', quantity: 5, previousStock: 105, newStock: 100, notes: 'Terjual', createdAt: '2024-03-15T10:30:00', createdBy: 'System' },
  { id: '3', productId: '15', productName: 'Headset Gaming Fantech', type: 'adjustment', quantity: -2, previousStock: 5, newStock: 3, notes: 'Barang rusak, dikembalikan ke supplier', createdAt: '2024-03-15T16:00:00', createdBy: 'Budi Inventory' },
];

export const mockServiceTickets: ServiceTicket[] = [
  { id: '1', ticketNumber: 'SRV-20240310-001', customerId: '1', customerName: 'Ahmad Fauzi', customerPhone: '081234567890', deviceType: 'Laptop', deviceBrand: 'ASUS', serialNumber: 'ASN123456', complaint: 'Laptop tidak bisa nyala, layar hitam', diagnosis: 'Kemungkinan RAM longgar atau motherboard bermasalah', estimatedCost: 350000, actualCost: 250000, estimatedCompletion: '2024-03-12', status: 'completed', sparepartsUsed: [{ productId: '9', name: 'RAM DDR5 16GB', quantity: 1, price: 799000 }], technicianNotes: [{ id: '1', note: 'Cek RAM, ternyata longgar. Pasang ulang dan test.', createdAt: '2024-03-11T10:00:00', createdBy: 'Rudi Teknisi' }, { id: '2', note: 'Sudah selesai, laptop normal kembali.', createdAt: '2024-03-12T14:00:00', createdBy: 'Rudi Teknisi' }], photos: [], technicianId: '4', technicianName: 'Rudi Teknisi', createdAt: '2024-03-10T09:00:00', updatedAt: '2024-03-12T14:00:00' },
  { id: '2', ticketNumber: 'SRV-20240314-002', customerId: '4', customerName: 'Dewi Lestari', customerPhone: '081234567893', deviceType: 'Printer', deviceBrand: 'Epson', complaint: 'Printer tidak bisa mencetak, lampu berkedip', diagnosis: '', estimatedCost: 150000, estimatedCompletion: '2024-03-16', status: 'in_progress', sparepartsUsed: [], technicianNotes: [{ id: '1', note: 'Sedang dicek, kemungkinan head printer kotor.', createdAt: '2024-03-15T09:00:00', createdBy: 'Rudi Teknisi' }], photos: [], technicianId: '4', technicianName: 'Rudi Teknisi', createdAt: '2024-03-14T10:00:00', updatedAt: '2024-03-15T09:00:00' },
  { id: '3', ticketNumber: 'SRV-20240316-003', customerName: 'Rizky Pratama', customerPhone: '081234567894', deviceType: 'PC', complaint: 'PC sering restart sendiri', diagnosis: 'Belum dicek', estimatedCost: 200000, status: 'waiting', sparepartsUsed: [], technicianNotes: [], photos: [], technicianId: '4', technicianName: 'Rudi Teknisi', createdAt: '2024-03-16T08:00:00', updatedAt: '2024-03-16T08:00:00' },
];

export const mockReceivables: Receivable[] = [
  { id: '1', transactionId: '3', invoiceNumber: 'INV-20240315-003', customerId: '2', customerName: 'Siti Nurhaliza', totalAmount: 14499000, paidAmount: 5000000, remainingAmount: 9499000, dueDate: '2024-04-15', status: 'partial', payments: [{ id: '1', amount: 5000000, method: 'dana', date: '2024-03-15', notes: 'DP' }], createdAt: '2024-03-15' },
  { id: '2', transactionId: '4', invoiceNumber: 'INV-20240316-001', customerId: '3', customerName: 'Bambang Setiawan', totalAmount: 6398000, paidAmount: 0, remainingAmount: 6398000, dueDate: '2024-04-16', status: 'unpaid', payments: [], createdAt: '2024-03-16' },
];

export const mockStoreSettings: StoreSettings = {
  name: 'SulthanCom',
  address: 'Jl. Raya Komputer No. 123, Surabaya',
  phone: '031-12345678',
  email: 'info@sulthancom.com',
  website: 'sulthancom.com',
  taxRate: 0,
  currency: 'IDR',
  invoiceFooter: 'Terima kasih atas kunjungan Anda. Garansi berlaku sesuai ketentuan.',
  theme: 'light',
};

export const mockDashboardStats: DashboardStats = {
  todaySales: 25746000,
  todayTransactions: 4,
  activeServices: 2,
  totalReceivables: 15897000,
  totalProducts: 327,
  lowStockProducts: 3,
  bestSellingProduct: { name: 'Logitech G102', quantity: 45 },
  slowMovingProduct: { name: 'UPS APC 650VA', quantity: 0 },
};

export const salesChartData = [
  { date: '2024-03-10', sales: 5200000, transactions: 3 },
  { date: '2024-03-11', sales: 8500000, transactions: 5 },
  { date: '2024-03-12', sales: 12000000, transactions: 7 },
  { date: '2024-03-13', sales: 7800000, transactions: 4 },
  { date: '2024-03-14', sales: 15500000, transactions: 8 },
  { date: '2024-03-15', sales: 25746000, transactions: 4 },
  { date: '2024-03-16', sales: 6398000, transactions: 1 },
];

export const categorySalesData = [
  { name: 'Laptop', value: 45, color: '#3B82F6' },
  { name: 'PC Rakitan', value: 20, color: '#8B5CF6' },
  { name: 'Sparepart', value: 15, color: '#EF4444' },
  { name: 'Aksesoris', value: 12, color: '#EC4899' },
  { name: 'Printer', value: 5, color: '#10B981' },
  { name: 'Service', value: 3, color: '#6366F1' },
];

export const paymentMethodData = [
  { name: 'Cash', value: 35 },
  { name: 'Transfer Bank', value: 40 },
  { name: 'DANA', value: 15 },
  { name: 'GoPay', value: 5 },
  { name: 'ShopeePay', value: 5 },
];
