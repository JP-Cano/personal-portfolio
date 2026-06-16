<script lang="ts">
  // Admin island: CRUD for experiences. Replaces the previous vanilla
  // window.openModal / global-script approach with reactive Svelte state.
  interface ExperienceRow {
    id: number;
    title: string;
    company: string;
    url: string;
    location: string;
    type: string;
    startDate: string;
    endDate: string;
    description: string;
  }

  interface FormState {
    title: string;
    company: string;
    url: string;
    location: string;
    type: string;
    start_date: string;
    end_date: string;
    description: string;
  }

  interface Props {
    experiences?: ExperienceRow[];
  }

  let { experiences: initial = [] }: Props = $props();

  let experiences = $state<ExperienceRow[]>(initial);
  let modalOpen = $state(false);
  let editingId = $state<number | null>(null);
  let saving = $state(false);
  let formError = $state<string | null>(null);
  let deleteTarget = $state<ExperienceRow | null>(null);
  let deleting = $state(false);

  let form = $state<FormState>(emptyForm());

  function emptyForm(): FormState {
    return {
      title: "",
      company: "",
      url: "",
      location: "",
      type: "Remote",
      start_date: "",
      end_date: "",
      description: "",
    };
  }

  function toInputDate(value: string): string {
    if (!value) return "";
    const date = new Date(value);
    return isNaN(date.getTime()) ? "" : date.toISOString().split("T")[0];
  }

  function formatDate(value: string): string {
    if (!value) return "Present";
    const date = new Date(value);
    return isNaN(date.getTime())
      ? value
      : date.toLocaleDateString("en-US", { year: "numeric", month: "short" });
  }

  function openAdd() {
    form = emptyForm();
    editingId = null;
    formError = null;
    modalOpen = true;
  }

  function openEdit(row: ExperienceRow) {
    form = {
      title: row.title,
      company: row.company,
      url: row.url,
      location: row.location,
      type: row.type,
      start_date: toInputDate(row.startDate),
      end_date: toInputDate(row.endDate),
      description: row.description,
    };
    editingId = row.id;
    formError = null;
    modalOpen = true;
  }

  function closeModal() {
    modalOpen = false;
    editingId = null;
    form = emptyForm();
  }

  function rowFromForm(id: number): ExperienceRow {
    return {
      id,
      title: form.title,
      company: form.company,
      url: form.url,
      location: form.location,
      type: form.type,
      startDate: form.start_date,
      endDate: form.end_date,
      description: form.description,
    };
  }

  async function save(e: SubmitEvent) {
    e.preventDefault();
    saving = true;
    formError = null;

    const payload = {
      title: form.title,
      company: form.company,
      url: form.url || undefined,
      location: form.location,
      type: form.type,
      start_date: form.start_date,
      end_date: form.end_date || undefined,
      description: form.description,
    };

    const url = editingId
      ? `/api/admin/experiences/${editingId}`
      : "/api/admin/experiences";
    const method = editingId ? "PATCH" : "POST";

    try {
      const res = await fetch(url, {
        method,
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
        credentials: "include",
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || "Failed to save experience");

      if (editingId) {
        const updated = rowFromForm(editingId);
        experiences = experiences.map((r) =>
          r.id === editingId ? updated : r
        );
      } else {
        const newId = (data.data?.id as number | undefined) ?? Date.now();
        experiences = [...experiences, rowFromForm(newId)];
      }
      closeModal();
    } catch (err) {
      formError = err instanceof Error ? err.message : "An error occurred";
    } finally {
      saving = false;
    }
  }

  async function confirmDelete() {
    if (!deleteTarget) return;
    deleting = true;
    formError = null;
    try {
      const res = await fetch(`/api/admin/experiences/${deleteTarget.id}`, {
        method: "DELETE",
        credentials: "include",
      });
      if (!res.ok) {
        const data = await res.json();
        throw new Error(data.error || "Failed to delete");
      }
      experiences = experiences.filter((r) => r.id !== deleteTarget!.id);
      deleteTarget = null;
    } catch (err) {
      formError = err instanceof Error ? err.message : "An error occurred";
    } finally {
      deleting = false;
    }
  }

  $effect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        modalOpen = false;
        deleteTarget = null;
      }
    };
    document.addEventListener("keydown", onKey);
    return () => document.removeEventListener("keydown", onKey);
  });
</script>

<div class="admin-card">
  <div class="admin-card-header">
    <h2 class="admin-card-title">All Experiences</h2>
    <button class="admin-btn admin-btn-primary" onclick={openAdd}>
      Add Experience
    </button>
  </div>

  {#if formError && !modalOpen && !deleteTarget}
    <div class="admin-alert admin-alert-error">{formError}</div>
  {/if}

  {#if experiences.length === 0}
    <div class="admin-empty-state">
      <p>No experiences yet. Add your first experience!</p>
    </div>
  {:else}
    <div class="admin-table-container">
      <table class="admin-table">
        <thead>
          <tr>
            <th>Title</th>
            <th>Location</th>
            <th>Period</th>
            <th style="width: 220px">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each experiences as exp (exp.id)}
            <tr>
              <td>
                <strong>{exp.title}</strong>
                <div class="admin-cell-sub">{exp.company}</div>
              </td>
              <td>
                <span>{exp.location || "-"}</span>
                <div class="admin-cell-sub">{exp.type}</div>
              </td>
              <td class="admin-cell-period">
                {formatDate(exp.startDate)} - {formatDate(exp.endDate)}
              </td>
              <td>
                <div class="admin-table-actions">
                  <a
                    href={`/admin/experiences/${exp.id}`}
                    class="admin-btn admin-btn-ghost admin-btn-sm"
                    title="Manage clients"
                  >
                    Clients
                  </a>
                  <button
                    type="button"
                    class="admin-btn admin-btn-ghost admin-btn-sm"
                    onclick={() => openEdit(exp)}
                  >
                    Edit
                  </button>
                  <button
                    type="button"
                    class="admin-btn admin-btn-ghost admin-btn-sm"
                    style="color: var(--admin-danger)"
                    onclick={() => (deleteTarget = exp)}
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
</div>

{#if modalOpen}
  <div
    class="admin-modal-overlay"
    style="display: flex"
    role="presentation"
    onclick={(e) => {
      if (e.target === e.currentTarget) closeModal();
    }}
  >
    <div class="admin-modal" role="dialog" aria-modal="true">
      <div class="admin-modal-header">
        <h2 class="admin-modal-title">
          {editingId ? "Edit Experience" : "Add Experience"}
        </h2>
        <button
          type="button"
          class="admin-modal-close"
          aria-label="Close modal"
          onclick={closeModal}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="20"
            height="20"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>
      <div class="admin-modal-body">
        {#if formError}
          <div class="admin-alert admin-alert-error">{formError}</div>
        {/if}
        <form class="admin-form" onsubmit={save}>
          <div class="admin-form-group">
            <label class="admin-label admin-label-required" for="exp-title">
              Job Title
            </label>
            <input
              id="exp-title"
              type="text"
              class="admin-input"
              bind:value={form.title}
              placeholder="e.g. Senior Software Engineer"
              required
            />
          </div>

          <div class="admin-form-group">
            <label class="admin-label admin-label-required" for="exp-company">
              Company
            </label>
            <input
              id="exp-company"
              type="text"
              class="admin-input"
              bind:value={form.company}
              placeholder="e.g. Google"
              required
            />
          </div>

          <div class="admin-form-group">
            <label class="admin-label" for="exp-url">Company URL</label>
            <input
              id="exp-url"
              type="url"
              class="admin-input"
              bind:value={form.url}
              placeholder="https://example.com"
            />
          </div>

          <div class="admin-form-row">
            <div class="admin-form-group">
              <label class="admin-label" for="exp-location">Location</label>
              <input
                id="exp-location"
                type="text"
                class="admin-input"
                bind:value={form.location}
                placeholder="e.g. San Francisco, CA"
              />
            </div>
            <div class="admin-form-group">
              <label class="admin-label admin-label-required" for="exp-type">
                Work Type
              </label>
              <select
                id="exp-type"
                class="admin-select"
                bind:value={form.type}
                required
              >
                <option value="Remote">Remote</option>
                <option value="On Site">On Site</option>
                <option value="Hybrid">Hybrid</option>
              </select>
            </div>
          </div>

          <div class="admin-form-row">
            <div class="admin-form-group">
              <label class="admin-label admin-label-required" for="exp-start">
                Start Date
              </label>
              <input
                id="exp-start"
                type="date"
                class="admin-input"
                bind:value={form.start_date}
                required
              />
            </div>
            <div class="admin-form-group">
              <label class="admin-label" for="exp-end">End Date</label>
              <input
                id="exp-end"
                type="date"
                class="admin-input"
                bind:value={form.end_date}
              />
            </div>
          </div>

          <div class="admin-form-group">
            <label class="admin-label" for="exp-desc">Description</label>
            <textarea
              id="exp-desc"
              class="admin-textarea"
              rows="4"
              bind:value={form.description}
              placeholder="Describe your role and achievements..."
            ></textarea>
          </div>

          <div class="admin-form-actions">
            <button
              type="button"
              class="admin-btn admin-btn-ghost"
              onclick={closeModal}
            >
              Cancel
            </button>
            <button
              type="submit"
              class="admin-btn admin-btn-primary"
              disabled={saving}
            >
              {saving
                ? "Saving..."
                : editingId
                  ? "Update Experience"
                  : "Create Experience"}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
{/if}

{#if deleteTarget}
  <div
    class="admin-modal-overlay"
    style="display: flex"
    role="presentation"
    onclick={(e) => {
      if (e.target === e.currentTarget) deleteTarget = null;
    }}
  >
    <div class="admin-modal" role="dialog" aria-modal="true">
      <div class="admin-modal-header">
        <h2 class="admin-modal-title">Delete Experience</h2>
        <button
          type="button"
          class="admin-modal-close"
          aria-label="Close modal"
          onclick={() => (deleteTarget = null)}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="20"
            height="20"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>
      <div class="admin-modal-body">
        {#if formError}
          <div class="admin-alert admin-alert-error">{formError}</div>
        {/if}
        <p class="admin-confirm-text">
          Are you sure you want to delete <strong>{deleteTarget.title}</strong>?
        </p>
        <p class="admin-confirm-warning">This action cannot be undone.</p>
        <div class="admin-form-actions">
          <button
            type="button"
            class="admin-btn admin-btn-ghost"
            onclick={() => (deleteTarget = null)}
          >
            Cancel
          </button>
          <button
            type="button"
            class="admin-btn admin-btn-danger"
            disabled={deleting}
            onclick={confirmDelete}
          >
            {deleting ? "Deleting..." : "Delete"}
          </button>
        </div>
      </div>
    </div>
  </div>
{/if}
