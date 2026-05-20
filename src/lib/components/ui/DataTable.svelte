<script lang="ts">
  interface Column<T> {
    key: string;
    label: string;
    render?: (item: T) => string;
  }

  interface Props<T> {
    columns: Column<T>[];
    data: T[];
    loading?: boolean;
    onRowClick?: (item: T) => void;
  }

  let { columns, data, loading = false, onRowClick }: Props<any> = $props();
</script>

<div class="overflow-x-auto">
  <table class="table table-sm">
    <thead>
      <tr>
        {#each columns as col}
          <th>{col.label}</th>
        {/each}
      </tr>
    </thead>
    <tbody>
      {#if loading}
        <tr>
          <td colspan={columns.length} class="text-center py-8">
            <span class="loading loading-spinner loading-md"></span>
          </td>
        </tr>
      {:else if data.length === 0}
        <tr>
          <td colspan={columns.length} class="text-center py-8 text-base-content">
            Tidak ada data
          </td>
        </tr>
      {:else}
        {#each data as item}
          <tr class="hover" onclick={() => onRowClick?.(item)}>
            {#each columns as col}
              <td>
                {#if col.render}
                  {@html col.render(item)}
                {:else}
                  {item[col.key]}
                {/if}
              </td>
            {/each}
          </tr>
        {/each}
      {/if}
    </tbody>
  </table>
</div>
