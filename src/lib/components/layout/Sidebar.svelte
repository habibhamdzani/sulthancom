<script lang="ts">
  import { currentUser, sidebarCollapsed, userRole, logout } from '$lib/stores/auth';
  import {
    LayoutDashboard, ShoppingCart, Package, Warehouse, Wrench, Users, ReceiptText,
    BarChart3, Settings, ChevronLeft, ChevronRight, LogOut, Box, Tags, Building2,
    CreditCard, FileText, Palette
  } from 'lucide-svelte';
  import { hasAccess, type Role } from '$lib/stores/auth';
  import { getRoleLabel } from '$lib/utils/helpers';
  import { page } from '$app/stores';

  let collapsed = $state(false);
  let currentRole = $derived($userRole);

  $effect(() => {
    collapsed = $sidebarCollapsed;
  });

  interface MenuItem {
    label: string;
    icon: typeof LayoutDashboard;
    href: string;
    roles?: Role[];
  }

  interface MenuGroup {
    label: string;
    items: MenuItem[];
  }

  const menuGroups: MenuGroup[] = [
    {
      label: 'Main',
      items: [
        { label: 'Dashboard', icon: LayoutDashboard, href: '/dashboard', roles: ['super_admin', 'admin', 'kasir', 'teknisi'] },
      ],
    },
    {
      label: 'Master Data',
      items: [
        { label: 'Produk', icon: Box, href: '/products', roles: ['super_admin', 'admin'] },
        { label: 'Kategori', icon: Tags, href: '/categories', roles: ['super_admin', 'admin'] },
        { label: 'Brand', icon: Building2, href: '/brands', roles: ['super_admin', 'admin'] },
        { label: 'Supplier', icon: Building2, href: '/suppliers', roles: ['super_admin', 'admin'] },
      ],
    },
    {
      label: 'Inventory',
      items: [
        { label: 'Stok', icon: Warehouse, href: '/inventory', roles: ['super_admin', 'admin'] },
      ],
    },
    {
      label: 'Penjualan',
      items: [
        { label: 'Kasir', icon: ShoppingCart, href: '/pos', roles: ['super_admin', 'kasir'] },
        { label: 'Transaksi', icon: ReceiptText, href: '/transactions', roles: ['super_admin', 'kasir'] },
        { label: 'Invoice', icon: FileText, href: '/invoices', roles: ['super_admin', 'kasir'] },
        { label: 'Invoice Designer', icon: Palette, href: '/invoice-designer', roles: ['super_admin', 'admin'] },
      ],
    },
    {
      label: 'Service',
      items: [
        { label: 'Tiket Service', icon: Wrench, href: '/services', roles: ['super_admin', 'teknisi'] },
      ],
    },
    {
      label: 'Keuangan',
      items: [
        { label: 'Pelanggan', icon: Users, href: '/customers', roles: ['super_admin', 'kasir'] },
        { label: 'Piutang', icon: CreditCard, href: '/receivables', roles: ['super_admin', 'kasir'] },
      ],
    },
    {
      label: 'Laporan',
      items: [
        { label: 'Report', icon: BarChart3, href: '/reports', roles: ['super_admin', 'admin', 'kasir'] },
      ],
    },
    {
      label: 'Pengaturan',
      items: [
        { label: 'Pengaturan Toko', icon: Settings, href: '/settings', roles: ['super_admin'] },
        { label: 'Manajemen User', icon: Users, href: '/users', roles: ['super_admin'] },
      ],
    },
  ];

  function toggleSidebar() {
    collapsed = !collapsed;
    $sidebarCollapsed = collapsed;
  }
</script>

<div class="sidebar fixed top-0 left-0 h-screen glass border-r border-base-300/50 flex flex-col z-40 transition-all duration-300" class:collapsed={collapsed}>
  <div class="flex items-center gap-3 px-4 h-16 border-b border-base-300/50">
    <div class="icon-container-3d w-10 h-10 flex items-center justify-center">
      <span class="text-primary font-bold text-sm relative z-10">SC</span>
    </div>
    {#if !collapsed}
      <div>
        <h1 class="font-bold text-base">SulthanCom</h1>
        <p class="text-xs text-base-content">POS System</p>
      </div>
    {/if}
    <button class="btn btn-ghost btn-xs btn-circle ml-auto" onclick={toggleSidebar}>
      {#if collapsed}
        <ChevronRight class="w-4 h-4 icon-3d" />
      {:else}
        <ChevronLeft class="w-4 h-4 icon-3d" />
      {/if}
    </button>
  </div>

  <nav class="flex-1 overflow-y-auto py-4">
    {#each menuGroups as group}
      {#if group.items.some(item => hasAccess(item.roles || [], currentRole))}
        {#if !collapsed}
          <div class="px-4 py-2 text-xs font-semibold text-base-content uppercase">{group.label}</div>
        {:else}
          <div class="px-2 py-2 flex justify-center">
            <div class="w-8 h-px bg-base-300/50"></div>
          </div>
        {/if}
        <ul class="menu menu-sm w-full gap-1">
          {#each group.items as item}
            {#if hasAccess(item.roles || [], currentRole)}
              <li>
                <a href={item.href} class="flex items-center gap-3 rounded-xl transition-all duration-200 hover:bg-base-200/50" class:active={$page.url.pathname.startsWith(item.href)}>
                  <div class="icon-container-3d p-2">
                    <svelte:component this={item.icon} class="w-5 h-5 flex-shrink-0 icon-3d" />
                  </div>
                  {#if !collapsed}
                    <span class="font-medium">{item.label}</span>
                  {/if}
                </a>
              </li>
            {/if}
          {/each}
        </ul>
      {/if}
    {/each}
  </nav>

  <div class="border-t border-base-300/50 p-4">
    {#if $currentUser}
      <div class="flex items-center gap-3">
        <div class="icon-container-3d w-10 h-10 flex items-center justify-center">
          <span class="text-primary font-bold text-sm relative z-10">{$currentUser.name.charAt(0)}</span>
        </div>
        {#if !collapsed}
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium truncate">{$currentUser.name}</p>
            <p class="text-xs text-base-content">{getRoleLabel($currentUser.role)}</p>
          </div>
          <button class="btn btn-ghost btn-xs btn-circle" onclick={logout}>
            <LogOut class="w-4 h-4 icon-3d" />
          </button>
        {/if}
      </div>
    {/if}
  </div>
</div>
