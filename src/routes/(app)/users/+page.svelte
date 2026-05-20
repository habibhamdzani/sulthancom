<script lang="ts">
  import { mockUsers } from '$lib/data/mock';
  import { getRoleLabel } from '$lib/utils/helpers';
  import { Plus, Edit, Trash2, KeyRound, UserPlus } from 'lucide-svelte';

  let users = $state(mockUsers);
  let showAddModal = $state(false);
  let showEditModal = $state(false);
  let editingUser = $state<any>(null);
  let newUser = $state({ name: '', email: '', role: 'kasir', password: '' });

  function resetNewUser() { newUser = { name: '', email: '', role: 'kasir', password: '' }; }

  function saveUser() {
    if (!newUser.name || !newUser.email) return;
    users.push({ id: crypto.randomUUID(), name: newUser.name, email: newUser.email, role: newUser.role, isActive: true, createdAt: new Date().toISOString().split('T')[0] });
    showAddModal = false;
    resetNewUser();
  }

  function updateUser() {
    if (!editingUser) return;
    const idx = users.findIndex(u => u.id === editingUser.id);
    if (idx >= 0) users[idx] = { ...editingUser };
    showEditModal = false;
    editingUser = null;
  }

  function deleteUser(id: string) {
    if (confirm('Hapus user ini?')) {
      users = users.filter(u => u.id !== id);
    }
  }

  let showResetPasswordModal = $state(false);
  let resetPasswordUser = $state<any>(null);
  let resetPasswordForm = $state({ newPassword: '', confirmPassword: '' });

  function openResetPassword(user: any) {
    resetPasswordUser = user;
    resetPasswordForm = { newPassword: '', confirmPassword: '' };
    showResetPasswordModal = true;
  }

  function handleResetPassword() {
    if (!resetPasswordForm.newPassword || resetPasswordForm.newPassword.length < 6) {
      alert('Password minimal 6 karakter');
      return;
    }
    if (resetPasswordForm.newPassword !== resetPasswordForm.confirmPassword) {
      alert('Konfirmasi password tidak cocok');
      return;
    }
    alert(`Password untuk ${resetPasswordUser.name} berhasil direset`);
    showResetPasswordModal = false;
    resetPasswordUser = null;
  }

  function getRoleBadge(role: string) {
    const colors: Record<string, string> = {
      super_admin: 'badge-error',
      admin: 'badge-primary',
      kasir: 'badge-success',
      teknisi: 'badge-warning',
    };
    return colors[role] || 'badge-ghost';
  }
</script>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-2xl font-bold">Manajemen User</h1>
      <p class="text-base-content">Kelola pengguna dan hak akses</p>
    </div>
    <button class="btn btn-primary btn-sm gap-2" onclick={() => showAddModal = true}>
      <UserPlus class="w-4 h-4" /> Tambah User
    </button>
  </div>

  <div class="glass-card">
    <div class="p-5">
      <div class="overflow-x-auto">
        <table class="table table-sm">
          <thead>
            <tr>
              <th>User</th>
              <th>Email</th>
              <th>Role</th>
              <th>Status</th>
              <th>Bergabung</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            {#each users as user}
              <tr class="hover">
                <td>
                  <div class="flex items-center gap-3">
                    <div class="avatar placeholder">
                      <div class="bg-primary text-primary-content rounded-full w-10">
                        <span>{user.name.charAt(0)}</span>
                      </div>
                    </div>
                    <div>
                      <div class="font-medium">{user.name}</div>
                    </div>
                  </div>
                </td>
                <td>{user.email}</td>
                <td><span class="badge badge-sm {getRoleBadge(user.role)}">{getRoleLabel(user.role)}</span></td>
                <td><span class="badge badge-sm {user.isActive ? 'badge-success' : 'badge-ghost'}">{user.isActive ? 'Aktif' : 'Nonaktif'}</span></td>
                <td class="text-sm">{user.createdAt}</td>
                <td>
                  <div class="flex gap-1">
                    <button class="btn btn-ghost btn-xs" onclick={() => openResetPassword(user)}><KeyRound class="w-4 h-4" /></button>
                    <button class="btn btn-ghost btn-xs" onclick={() => { editingUser = user; showEditModal = true; }}><Edit class="w-4 h-4" /></button>
                    <button class="btn btn-ghost btn-xs text-error" onclick={() => deleteUser(user.id)}><Trash2 class="w-4 h-4" /></button>
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

{#if showEditModal && editingUser}
  <dialog class="modal modal-open" onclick={() => showEditModal = false}>
    <div class="modal-box" onclick={(e) => e.stopPropagation()}>
      <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2" onclick={() => showEditModal = false}>✕</button>
      <h3 class="text-lg font-bold mb-4">Edit User</h3>
      <div class="space-y-3">
        <div class="form-control"><label class="label"><span class="label-text">Nama</span></label><input type="text" class="input input-bordered input-sm w-full" bind:value={editingUser.name} /></div>
        <div class="form-control"><label class="label"><span class="label-text">Email</span></label><input type="email" class="input input-bordered input-sm w-full" bind:value={editingUser.email} /></div>
        <div class="form-control"><label class="label"><span class="label-text">Role</span></label>
          <select class="select select-bordered select-sm w-full" bind:value={editingUser.role}>
            <option value="super_admin">Super Admin</option><option value="admin">Admin</option><option value="kasir">Kasir</option><option value="teknisi">Teknisi</option>
          </select>
        </div>
        <div class="flex justify-end gap-2 mt-4">
          <button class="btn btn-ghost btn-sm" onclick={() => showEditModal = false}>Batal</button>
          <button class="btn btn-primary btn-sm" onclick={updateUser}>Simpan</button>
        </div>
      </div>
    </div>
  </dialog>
{/if}

{#if showResetPasswordModal && resetPasswordUser}
  <dialog class="modal modal-open" onclick={() => showResetPasswordModal = false}>
    <div class="modal-box" onclick={(e) => e.stopPropagation()}>
      <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2" onclick={() => showResetPasswordModal = false}>✕</button>
      <div class="flex items-center gap-3 mb-4">
        <KeyRound class="w-6 h-6 text-primary" />
        <div>
          <h3 class="text-lg font-bold">Reset Password</h3>
          <p class="text-sm text-base-content">{resetPasswordUser.name} — {resetPasswordUser.email}</p>
        </div>
      </div>
      <div class="space-y-3">
        <div class="form-control"><label class="label"><span class="label-text">Password Baru</span></label><input type="password" class="input input-bordered input-sm w-full" placeholder="Min. 6 karakter" bind:value={resetPasswordForm.newPassword} /></div>
        <div class="form-control"><label class="label"><span class="label-text">Konfirmasi Password</span></label><input type="password" class="input input-bordered input-sm w-full" placeholder="Ulangi password" bind:value={resetPasswordForm.confirmPassword} /></div>
        <div class="flex justify-end gap-2 mt-4">
          <button class="btn btn-ghost btn-sm" onclick={() => showResetPasswordModal = false}>Batal</button>
          <button class="btn btn-primary btn-sm" onclick={handleResetPassword}>Reset Password</button>
        </div>
      </div>
    </div>
  </dialog>
{/if}

{#if showAddModal}
  <dialog class="modal modal-open" onclick={() => showAddModal = false}>
    <div class="modal-box" onclick={(e) => e.stopPropagation()}>
      <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2" onclick={() => { showAddModal = false; resetNewUser(); }}>✕</button>
      <h3 class="text-lg font-bold mb-4">Tambah User Baru</h3>
      <div class="space-y-3">
        <div class="form-control"><label class="label"><span class="label-text">Nama</span></label><input type="text" class="input input-bordered input-sm w-full" placeholder="Nama lengkap" bind:value={newUser.name} /></div>
        <div class="form-control"><label class="label"><span class="label-text">Email</span></label><input type="email" class="input input-bordered input-sm w-full" placeholder="email@sulthancom.com" bind:value={newUser.email} /></div>
        <div class="form-control"><label class="label"><span class="label-text">Role</span></label>
          <select class="select select-bordered select-sm w-full" bind:value={newUser.role}>
            <option value="super_admin">Super Admin</option><option value="admin">Admin</option><option value="kasir">Kasir</option><option value="teknisi">Teknisi</option>
          </select>
        </div>
        <div class="form-control"><label class="label"><span class="label-text">Password</span></label><input type="password" class="input input-bordered input-sm w-full" placeholder="Min. 8 karakter" bind:value={newUser.password} /></div>
        <div class="flex justify-end gap-2 mt-4">
          <button class="btn btn-ghost btn-sm" onclick={() => { showAddModal = false; resetNewUser(); }}>Batal</button>
          <button class="btn btn-primary btn-sm" onclick={saveUser}>Simpan</button>
        </div>
      </div>
    </div>
  </dialog>
{/if}
