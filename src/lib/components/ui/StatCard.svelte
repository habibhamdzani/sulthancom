<script lang="ts">
  import type { ComponentType } from 'svelte';
  import type { LucideIcon } from 'lucide-svelte';

  interface Props {
    title: string;
    value: string | number;
    icon: ComponentType<LucideIcon>;
    trend?: { value: number; positive: boolean };
    color?: string;
    onClick?: () => void;
  }

  let { title, value, icon, trend, color = 'primary', onClick }: Props = $props();

  const colorMap: Record<string, { bg: string; text: string; gradient: string }> = {
    primary: { bg: 'bg-primary/10', text: 'text-primary', gradient: 'from-blue-500/20 to-blue-600/5' },
    success: { bg: 'bg-success/10', text: 'text-success', gradient: 'from-green-500/20 to-green-600/5' },
    warning: { bg: 'bg-warning/10', text: 'text-warning', gradient: 'from-yellow-500/20 to-yellow-600/5' },
    error: { bg: 'bg-error/10', text: 'text-error', gradient: 'from-red-500/20 to-red-600/5' },
    info: { bg: 'bg-info/10', text: 'text-info', gradient: 'from-cyan-500/20 to-cyan-600/5' },
    neutral: { bg: 'bg-neutral/10', text: 'text-neutral', gradient: 'from-gray-500/20 to-gray-600/5' },
    accent: { bg: 'bg-accent/10', text: 'text-accent', gradient: 'from-purple-500/20 to-purple-600/5' },
  };

  let iconBg = $derived(colorMap[color]?.bg ?? 'bg-primary/10');
  let iconColor = $derived(colorMap[color]?.text ?? 'text-primary');
  let iconGradient = $derived(colorMap[color]?.gradient ?? 'from-blue-500/20 to-blue-600/5');
</script>

<div class="glass-card cursor-pointer" onclick={onClick}>
  <div class="p-5">
    <div class="flex items-start justify-between">
      <div>
        <p class="text-sm text-base-content">{title}</p>
        <h3 class="text-2xl font-bold mt-1">{value}</h3>
        {#if trend}
          <p class="text-xs mt-2 {trend.positive ? 'text-success' : 'text-error'}">
            {trend.positive ? '↑' : '↓'} {Math.abs(trend.value)}% dari kemarin
          </p>
        {/if}
      </div>
      <div class="icon-container-3d bg-gradient-to-br {iconGradient}">
        <svelte:component this={icon} class="w-6 h-6 {iconColor} icon-3d relative z-10" />
      </div>
    </div>
  </div>
</div>
