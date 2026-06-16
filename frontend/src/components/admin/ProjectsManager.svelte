<script lang="ts">
  // Admin island: CRUD for projects. Technologies are edited as a dynamic list
  // but persisted as the comma-separated string the Project model expects.
  interface ProjectRow {
    id: number;
    name: string;
    description: string;
    url: string;
    startDate: string;
    endDate: string;
    technologies: string;
  }

  interface FormState {
    name: string;
    description: string;
    url: string;
    start_date: string;
    end_date: string;
    technologies: string[];
  }

  interface Props {
    projects?: ProjectRow[];
  }

  let { projects: initial = [] }: Props = $props();

  let projects = $state<ProjectRow[]>(initial);
  let modalOpen = $state(false);
  let editingId = $state<number | null>(null);
  let saving = $state(false);
  let formError = $state<string | null>(null);
  let deleteTarget = $state<ProjectRow | null>(null);
  let deleting = $state(false);

  let form = $state<FormState>(emptyForm());

  function emptyForm(): FormState {
    return {
      name: "",
      description: "",
      url: "",
      start_date: "",
      end_date: "",
      technologies: [""],
    };
  }

  function splitTech(value: string): string[] {
    const list = value
      .split(",")
      .map((t) => t.trim())
      .filter(Boolean);
    return list.length > 0 ? list : [""];
  }

  function toInputDate(value: string): string {
    if (!value) return "";
    const date = new Date(value);
    return isNaN(date.getTime()) ? "" : date.toISOString().split("T")[0];
  }

  function formatDate(value: string): string {
    if (!value) return "Ongoing";
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

  function openEdit(row: ProjectRow) {
    form = {
      name: row.name,
      description: row.description,
      url: row.url,
      start_date: toInputDate(row.startDate),
      end_date: toInputDate(row.endDate),
      technologies: splitTech(row.technologies),
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

  function addTech() {
    form.technologies = [...form.technologies, ""];
  }

  function removeTech(index: number) {
    form.technologies = form.technologies.filter((_, i) => i !== index);
  }

  function rowFromForm(id: number, techString: string): ProjectRow {
    return {
      id,
      name: form.name,
      description: form.description,
      url: form.url,
      startDate: form.start_date,
      endDate: form.end_date,
      technologies: techString,
    };
  }

  async function save(e: SubmitEvent) {
    e.preventDefault();
    saving = true;
    formError = null;

    const techString = form.technologies
      .map((t) => t.trim())
      .filter(Boolean)
      .join(",");

    const payload = {
      name: form.name,
      description: form.description,
      url: form.url || undefined,
      start_date: form.start_date,
      end_date: form.end_date || undefined,
      technologies: techString || undefined,
    };

    const url = editingId
      ? `/api/admin/projects/${editingId}`
      : "/api/admin/projects";
    const method = editingId ? "PATCH" : "POST";

    try {
      const res = await fetch(url, {
        method,
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
        credentials: "include",
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || "Failed to save project");

      if (editingId) {
        const updated = rowFromForm(editingId, techString);
        projects = projects.map((r) => (r.id === editingId ? updated : r));
      } else {
        const newId = (data.data?.id as number | undefined) ?? Date.now();
        projects = [...projects, rowFromForm(newId, techString)];
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
      const res = await fetch(`/api/admin/projects/${deleteTarget.id}`, {
        method: "DELETE",
        credentials: "include",
      });
      if (!res.ok) {
        const data = await res.json();
        throw new Error(data.error || "Failed to delete");
      }
      projects = projects.filter((r) => r.id !== deleteTarget!.id);
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
    <h2 class="admin-card-title">All Projects</h2>
    <button class="admin-btn admin-btn-primary" onclick={openAdd}>
      Add Project
    </button>
  </div>

  {#if formError && !modalOpen && !deleteTarget}
    <div class="admin-alert admin-alert-error">{formError}</div>
  {/if}

  {#if projects.length === 0}
    <div class="admin-empty-state">
      <p>No projects yet. Add your first project!</p>
    </div>
  {:else}
    <div class="admin-table-container">
      <table class="admin-table">
        <thead>
          <tr>
            <th>Project</th>
            <th>URL</th>
            <th>Period</th>
            <th style="width: 160px">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each projects as project (project.id)}
            <tr>
              <td>
                <strong>{project.name}</strong>
                <div class="admin-cell-sub admin-cell-ellipsis">
                  {project.description}
                </div>
              </td>
              <td>
                {#if project.url}
                  <a
                    href={project.url}
                    target="_blank"
                    rel="noopener noreferrer"
                    class="admin-link"
                  >
                    View Project
                  </a>
                {:else}
                  <span class="admin-cell-sub">-</span>
                {/if}
              </td>
              <td class="admin-cell-period">
                {formatDate(project.startDate)} - {formatDate(project.endDate)}
              </td>
              <td>
                <div class="admin-table-actions">
                  <button
                    type="button"
                    class="admin-btn admin-btn-ghost admin-btn-sm"
                    onclick={() => openEdit(project)}
                  >
                    Edit
                  </button>
                  <button
                    type="button"
                    class="admin-btn admin-btn-ghost admin-btn-sm"
                    style="color: var(--admin-danger)"
                    onclick={() => (deleteTarget = project)}
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
          {editingId ? "Edit Project" : "Add Project"}
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
            <label class="admin-label admin-label-required" for="proj-name">
              Project Name
            </label>
            <input
              id="proj-name"
              type="text"
              class="admin-input"
              bind:value={form.name}
              placeholder="e.g. E-commerce Platform"
              required
            />
          </div>

          <div class="admin-form-group">
            <label class="admin-label admin-label-required" for="proj-desc">
              Description
            </label>
            <textarea
              id="proj-desc"
              class="admin-textarea"
              rows="4"
              bind:value={form.description}
              placeholder="Describe your project..."
              required
            ></textarea>
          </div>

          <div class="admin-form-group">
            <label class="admin-label" for="proj-url">Project URL</label>
            <input
              id="proj-url"
              type="url"
              class="admin-input"
              bind:value={form.url}
              placeholder="https://github.com/username/project"
            />
          </div>

          <div class="admin-form-row">
            <div class="admin-form-group">
              <label class="admin-label admin-label-required" for="proj-start">
                Start Date
              </label>
              <input
                id="proj-start"
                type="date"
                class="admin-input"
                bind:value={form.start_date}
                required
              />
            </div>
            <div class="admin-form-group">
              <label class="admin-label" for="proj-end">End Date</label>
              <input
                id="proj-end"
                type="date"
                class="admin-input"
                bind:value={form.end_date}
              />
            </div>
          </div>

          <div class="admin-form-group">
            <label class="admin-label">Technologies</label>
            {#each form.technologies as _tech, i (i)}
              <div class="admin-array-row">
                <input
                  type="text"
                  class="admin-input"
                  bind:value={form.technologies[i]}
                  placeholder="e.g. React, Node.js"
                />
                <button
                  type="button"
                  class="admin-btn admin-btn-ghost admin-btn-sm"
                  style="color: var(--admin-danger)"
                  onclick={() => removeTech(i)}
                >
                  Remove
                </button>
              </div>
            {/each}
            <button
              type="button"
              class="admin-btn admin-btn-ghost admin-btn-sm"
              onclick={addTech}
            >
              + Add technology
            </button>
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
                  ? "Update Project"
                  : "Create Project"}
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
        <h2 class="admin-modal-title">Delete Project</h2>
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
          Are you sure you want to delete <strong>{deleteTarget.name}</strong>?
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
