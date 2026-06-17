<script lang="ts">
  // Admin island: manage the clients that belong to a single experience.
  // Reads/writes go through the /api/admin/experiences/:id/clients proxy routes.
  interface ClientView {
    id: number;
    name: string;
    url: string;
    startDate: string;
    endDate: string;
    description: string;
    achievements: string[];
    responsibilities: string[];
    technologies: string[];
  }

  type ArrayField = "achievements" | "responsibilities" | "technologies";

  type RawClient = {
    id: number;
    name?: string;
    url?: string;
    startDate?: string;
    endDate?: string;
    description?: string;
    achievements?: string[];
    responsibilities?: string[];
    technologies?: string[];
  };

  interface Props {
    experienceId: string | number;
    clients?: ClientView[];
  }

  let { experienceId, clients: initialClients = [] }: Props = $props();

  let clients = $state<ClientView[]>(initialClients);
  let formOpen = $state(false);
  let editingId = $state<number | null>(null);
  let saving = $state(false);
  let formError = $state<string | null>(null);
  let form = $state<ClientView>(emptyForm());

  const basePath = $derived(`/api/admin/experiences/${experienceId}/clients`);

  function emptyForm(): ClientView {
    return {
      id: 0,
      name: "",
      url: "",
      startDate: "",
      endDate: "",
      description: "",
      achievements: [],
      responsibilities: [],
      technologies: [],
    };
  }

  function toInputDate(value?: string): string {
    if (!value) return "";
    const d = new Date(value);
    if (isNaN(d.getTime())) return "";
    return `${d.getUTCFullYear()}-${String(d.getUTCMonth() + 1).padStart(2, "0")}`;
  }

  function normalize(raw: RawClient): ClientView {
    return {
      id: raw.id,
      name: raw.name ?? "",
      url: raw.url ?? "",
      startDate: toInputDate(raw.startDate),
      endDate: toInputDate(raw.endDate),
      description: raw.description ?? "",
      achievements: raw.achievements ?? [],
      responsibilities: raw.responsibilities ?? [],
      technologies: raw.technologies ?? [],
    };
  }

  function openAdd() {
    form = emptyForm();
    editingId = null;
    formError = null;
    formOpen = true;
  }

  function openEdit(client: ClientView) {
    form = {
      ...client,
      achievements: [...client.achievements],
      responsibilities: [...client.responsibilities],
      technologies: [...client.technologies],
    };
    editingId = client.id;
    formError = null;
    formOpen = true;
  }

  function closeForm() {
    formOpen = false;
    editingId = null;
    form = emptyForm();
  }

  function addItem(field: ArrayField) {
    form[field] = [...form[field], ""];
  }

  function removeItem(field: ArrayField, index: number) {
    form[field] = form[field].filter((_, i) => i !== index);
  }

  function formatPeriod(client: ClientView): string {
    const fmt = (d: string) =>
      d
        ? new Date(d).toLocaleDateString("en-US", {
            month: "short",
            year: "numeric",
          })
        : "Present";
    return `${fmt(client.startDate)} - ${fmt(client.endDate)}`;
  }

  async function submit(e: SubmitEvent) {
    e.preventDefault();
    saving = true;
    formError = null;

    const payload = {
      name: form.name,
      url: form.url || undefined,
      start_date: form.startDate,
      end_date: form.endDate || undefined,
      description: form.description || undefined,
      achievements: form.achievements.map((s) => s.trim()).filter(Boolean),
      responsibilities: form.responsibilities
        .map((s) => s.trim())
        .filter(Boolean),
      technologies: form.technologies.map((s) => s.trim()).filter(Boolean),
    };

    const url = editingId ? `${basePath}/${editingId}` : basePath;
    const method = editingId ? "PATCH" : "POST";

    try {
      const res = await fetch(url, {
        method,
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
        credentials: "include",
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || "Failed to save client");

      const saved = normalize(data.data as RawClient);
      if (editingId) {
        clients = clients.map((c) => (c.id === editingId ? saved : c));
      } else {
        clients = [...clients, saved];
      }
      closeForm();
    } catch (err) {
      formError = err instanceof Error ? err.message : "An error occurred";
    } finally {
      saving = false;
    }
  }

  async function remove(client: ClientView) {
    if (!confirm(`Delete client "${client.name}"? This cannot be undone.`)) {
      return;
    }
    formError = null;
    try {
      const res = await fetch(`${basePath}/${client.id}`, {
        method: "DELETE",
        credentials: "include",
      });
      if (!res.ok) {
        const data = await res.json();
        throw new Error(data.error || "Failed to delete client");
      }
      clients = clients.filter((c) => c.id !== client.id);
    } catch (err) {
      formError = err instanceof Error ? err.message : "An error occurred";
    }
  }
</script>

{#snippet arrayField(label: string, field: ArrayField)}
  <div class="admin-form-group">
    <label class="admin-label">{label}</label>
    {#each form[field] as _item, i (i)}
      <div class="admin-array-row">
        <input
          type="text"
          class="admin-input"
          bind:value={form[field][i]}
          placeholder={label}
        />
        <button
          type="button"
          class="admin-btn admin-btn-ghost admin-btn-sm"
          onclick={() => removeItem(field, i)}
        >
          Remove
        </button>
      </div>
    {/each}
    <button
      type="button"
      class="admin-btn admin-btn-ghost admin-btn-sm"
      onclick={() => addItem(field)}
    >
      + Add {label.toLowerCase()}
    </button>
  </div>
{/snippet}

<div class="admin-card">
  <div class="admin-card-header">
    <h2 class="admin-card-title">Clients</h2>
    {#if !formOpen}
      <button class="admin-btn admin-btn-primary" onclick={openAdd}>
        Add Client
      </button>
    {/if}
  </div>

  {#if formError}
    <div class="admin-alert admin-alert-error">{formError}</div>
  {/if}

  {#if formOpen}
    <form class="admin-form" onsubmit={submit}>
      <div class="admin-form-group">
        <label class="admin-label admin-label-required" for="client-name">
          Client name
        </label>
        <input
          id="client-name"
          type="text"
          class="admin-input"
          bind:value={form.name}
          placeholder="e.g. Acme Corp"
          required
        />
      </div>

      <div class="admin-form-group">
        <label class="admin-label" for="client-url">Client URL</label>
        <input
          id="client-url"
          type="url"
          class="admin-input"
          bind:value={form.url}
          placeholder="https://example.com"
        />
      </div>

      <div class="admin-form-row">
        <div class="admin-form-group">
          <label class="admin-label admin-label-required" for="client-start">
            Start date
          </label>
          <input
            id="client-start"
            type="month"
            class="admin-input"
            bind:value={form.startDate}
            required
          />
        </div>
        <div class="admin-form-group">
          <label class="admin-label" for="client-end">End date</label>
          <input
            id="client-end"
            type="month"
            class="admin-input"
            bind:value={form.endDate}
          />
        </div>
      </div>

      <div class="admin-form-group">
        <label class="admin-label" for="client-desc">Description</label>
        <textarea
          id="client-desc"
          class="admin-textarea"
          rows="3"
          bind:value={form.description}
          placeholder="Short summary of the engagement..."
        ></textarea>
      </div>

      {@render arrayField("Responsibilities", "responsibilities")}
      {@render arrayField("Achievements", "achievements")}
      {@render arrayField("Technologies", "technologies")}

      <div class="admin-form-actions">
        <button
          type="button"
          class="admin-btn admin-btn-ghost"
          onclick={closeForm}
        >
          Cancel
        </button>
        <button
          type="submit"
          class="admin-btn admin-btn-primary"
          disabled={saving}
        >
          {saving ? "Saving..." : editingId ? "Update Client" : "Create Client"}
        </button>
      </div>
    </form>
  {/if}

  {#if !formOpen}
    {#if clients.length === 0}
      <div class="admin-empty-state">
        <p>No clients yet. Add the first one for this experience.</p>
      </div>
    {:else}
      <div class="admin-table-container">
        <table class="admin-table">
          <thead>
            <tr>
              <th>Client</th>
              <th>Period</th>
              <th>Technologies</th>
              <th style="width: 120px">Actions</th>
            </tr>
          </thead>
          <tbody>
            {#each clients as client (client.id)}
              <tr>
                <td>
                  <strong>{client.name}</strong>
                </td>
                <td class="admin-cell-period">{formatPeriod(client)}</td>
                <td class="admin-cell-tech"
                  >{client.technologies.join(", ") || "-"}</td
                >
                <td>
                  <div class="admin-table-actions">
                    <button
                      type="button"
                      class="admin-btn admin-btn-ghost admin-btn-sm"
                      onclick={() => openEdit(client)}
                    >
                      Edit
                    </button>
                    <button
                      type="button"
                      class="admin-btn admin-btn-ghost admin-btn-sm"
                      style="color: var(--admin-danger)"
                      onclick={() => remove(client)}
                    >
                      Delete
                    </button>
                  </div>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {/if}
  {/if}
</div>
