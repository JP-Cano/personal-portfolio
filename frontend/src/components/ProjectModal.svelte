<script lang="ts">
  // Reactive project detail modal. Opens in response to a `project:open`
  // CustomEvent dispatched by the (server-rendered) project cards, replacing the
  // previous window.openProjectModal global + manual DOM population.
  interface ProjectData {
    name: string;
    description: string;
    url?: string;
    startDate: string;
    endDate?: string;
    technologies?: string;
  }

  let open = $state(false);
  let visible = $state(false);
  let project = $state<ProjectData | null>(null);

  const isActive = $derived(project ? !project.endDate : false);
  const techArray = $derived(
    project?.technologies
      ? project.technologies
          .split(",")
          .map((t) => t.trim())
          .filter(Boolean)
      : []
  );

  function formatDate(dateString?: string): string {
    if (!dateString) return "";
    return new Date(dateString).toLocaleDateString("en-US", {
      month: "short",
      year: "numeric",
    });
  }

  function openModal(data: ProjectData) {
    project = data;
    open = true;
    document.body.style.overflow = "hidden";
    requestAnimationFrame(() => {
      visible = true;
    });
  }

  function closeModal() {
    visible = false;
    setTimeout(() => {
      open = false;
      project = null;
      document.body.style.overflow = "";
    }, 250);
  }

  function reset() {
    visible = false;
    open = false;
    project = null;
    document.body.style.overflow = "";
  }

  $effect(() => {
    const onOpen = (e: Event) => {
      const detail = (e as CustomEvent<ProjectData>).detail;
      if (detail) openModal(detail);
    };
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape" && open) closeModal();
    };

    window.addEventListener("project:open", onOpen);
    document.addEventListener("keydown", onKey);
    // Reset on Astro View Transitions navigation.
    document.addEventListener("astro:after-swap", reset);

    return () => {
      window.removeEventListener("project:open", onOpen);
      document.removeEventListener("keydown", onKey);
      document.removeEventListener("astro:after-swap", reset);
    };
  });
</script>

{#if open && project}
  <div
    class={`fixed inset-0 z-[1000] flex items-center justify-center bg-[rgba(0,0,0,0.6)] p-3 backdrop-blur-md transition-opacity duration-[250ms] md:p-4 ${
      visible ? "opacity-100" : "opacity-0"
    }`}
    role="presentation"
    onclick={(e) => {
      if (e.target === e.currentTarget) closeModal();
    }}
  >
    <div
      class={`panel flex max-h-[90vh] w-full max-w-[600px] flex-col overflow-y-auto overflow-x-hidden transition-transform duration-[250ms] ease-[cubic-bezier(0.4,0,0.2,1)] ${
        visible ? "scale-100 translate-y-0" : "scale-[0.98] translate-y-2"
      }`}
      role="dialog"
      aria-modal="true"
      aria-label={project.name}
    >
      <div
        class="flex items-start justify-between gap-4 border-b border-border p-5 md:p-6"
      >
        <div class="flex min-w-0 flex-1 flex-col gap-3">
          <h2
            class="m-0 break-words text-[1.4rem] leading-tight md:text-[1.5rem]"
          >
            {project.name}
          </h2>
          {#if isActive}
            <span class="chip w-fit lowercase text-primary">
              <span
                class="h-1.5 w-1.5 rounded-full bg-current animate-[pulse-dot_1.6s_ease-in-out_infinite]"
              ></span>
              active
            </span>
          {:else}
            <span class="chip w-fit lowercase text-secondary">completed</span>
          {/if}
        </div>
        <button
          type="button"
          class="inline-flex h-9 w-9 shrink-0 items-center justify-center rounded-sm border border-border text-text-light transition-colors hover:border-primary hover:text-primary"
          onclick={closeModal}
          aria-label="Close modal"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="22"
            height="22"
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
      <div class="flex flex-col gap-6 p-5 md:p-6">
        <p class="m-0 text-[0.98rem] leading-[1.7] text-text-light">
          {project.description}
        </p>

        {#if techArray.length > 0}
          <div class="flex flex-wrap gap-2">
            {#each techArray as tech, i (tech + i)}
              <span class="chip">{tech}</span>
            {/each}
          </div>
        {/if}

        <div
          class="flex flex-wrap gap-6 rounded-sm border border-border bg-bg-alt p-4"
        >
          <div class="flex flex-col gap-1">
            <span
              class="font-mono text-[0.7rem] font-medium lowercase tracking-[0.02em] text-primary"
              >Started</span
            >
            <span class="font-mono text-[0.92rem] font-medium text-text"
              >{formatDate(project.startDate)}</span
            >
          </div>
          {#if project.endDate}
            <div class="flex flex-col gap-1">
              <span
                class="font-mono text-[0.7rem] font-medium lowercase tracking-[0.02em] text-primary"
                >Completed</span
              >
              <span class="font-mono text-[0.92rem] font-medium text-text"
                >{formatDate(project.endDate)}</span
              >
            </div>
          {/if}
        </div>

        {#if project.url}
          <a
            href={project.url}
            class="btn-primary group/link"
            target="_blank"
            rel="noopener noreferrer"
          >
            <span>Visit Project</span>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="18"
              height="18"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="transition-transform group-hover/link:translate-x-0.5 group-hover/link:-translate-y-0.5"
            >
              <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"
              ></path>
              <polyline points="15 3 21 3 21 9"></polyline>
              <line x1="10" y1="14" x2="21" y2="3"></line>
            </svg>
          </a>
        {/if}
      </div>
    </div>
  </div>
{/if}
