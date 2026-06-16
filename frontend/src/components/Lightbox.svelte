<script lang="ts">
  /**
   * Self-contained certificate gallery + accessible lightbox.
   *
   * Renders the systems-styled cert cards itself and, on activation, opens a
   * full-screen overlay with the selected image, its title/issuer metadata,
   * previous/next controls and a close button. Handles keyboard (Esc closes,
   * ArrowLeft/Right navigate), a focus trap, backdrop-click close and body
   * scroll lock while open. Replaces the previous GLightbox dependency.
   */
  interface CertImage {
    url: string;
    title: string;
    issuer: string;
    date: string;
  }

  interface Props {
    images?: CertImage[];
  }

  let { images = [] }: Props = $props();

  let open = $state(false);
  let visible = $state(false);
  let index = $state(0);

  let overlay = $state<HTMLDivElement | null>(null);
  let closeBtn = $state<HTMLButtonElement | null>(null);
  let lastFocused: HTMLElement | null = null;

  const current = $derived(images[index] ?? null);
  const hasMany = $derived(images.length > 1);

  function openAt(i: number) {
    index = i;
    open = true;
    lastFocused = document.activeElement as HTMLElement | null;
    document.body.style.overflow = "hidden";
    requestAnimationFrame(() => {
      visible = true;
      closeBtn?.focus();
    });
  }

  function close() {
    visible = false;
    const restore = lastFocused;
    setTimeout(() => {
      open = false;
      document.body.style.overflow = "";
      restore?.focus?.();
    }, 200);
  }

  function reset() {
    visible = false;
    open = false;
    document.body.style.overflow = "";
  }

  function next() {
    if (images.length === 0) return;
    index = (index + 1) % images.length;
  }

  function prev() {
    if (images.length === 0) return;
    index = (index - 1 + images.length) % images.length;
  }

  function trapFocus(e: KeyboardEvent) {
    if (!overlay) return;
    const focusable = overlay.querySelectorAll<HTMLElement>(
      'button, [href], [tabindex]:not([tabindex="-1"])'
    );
    if (focusable.length === 0) return;
    const first = focusable[0];
    const last = focusable[focusable.length - 1];
    if (e.shiftKey) {
      if (document.activeElement === first) {
        e.preventDefault();
        last.focus();
      }
    } else if (document.activeElement === last) {
      e.preventDefault();
      first.focus();
    }
  }

  function onKeydown(e: KeyboardEvent) {
    if (!open) return;
    switch (e.key) {
      case "Escape":
        e.preventDefault();
        close();
        break;
      case "ArrowRight":
        if (hasMany) {
          e.preventDefault();
          next();
        }
        break;
      case "ArrowLeft":
        if (hasMany) {
          e.preventDefault();
          prev();
        }
        break;
      case "Tab":
        trapFocus(e);
        break;
    }
  }

  function onCardKeydown(e: KeyboardEvent, i: number) {
    if (e.key === "Enter" || e.key === " ") {
      e.preventDefault();
      openAt(i);
    }
  }

  $effect(() => {
    document.addEventListener("keydown", onKeydown);
    document.addEventListener("astro:after-swap", reset);
    return () => {
      document.removeEventListener("keydown", onKeydown);
      document.removeEventListener("astro:after-swap", reset);
      document.body.style.overflow = "";
    };
  });
</script>

<div
  class="grid grid-cols-[repeat(auto-fill,minmax(min(300px,100%),1fr))] gap-8"
>
  {#each images as cert, i (cert.url + i)}
    <button
      type="button"
      class="card card-hover group flex h-full cursor-pointer flex-col overflow-hidden p-0 text-left"
      onclick={() => openAt(i)}
      onkeydown={(e) => onCardKeydown(e, i)}
      aria-label={`View certificate: ${cert.title}`}
    >
      <div
        class="relative w-full overflow-hidden border-b border-border bg-bg-alt pt-[70.7%]"
      >
        <img
          src={cert.url}
          alt={cert.title}
          class="absolute inset-0 h-full w-full object-cover"
          loading="lazy"
        />
        <span
          class="absolute inset-0 flex items-end justify-start bg-[linear-gradient(to_top,rgba(0,0,0,0.55)_0%,rgba(0,0,0,0)_55%)] p-3 opacity-0 transition-opacity duration-300 group-hover:opacity-100 group-focus-visible:opacity-100"
        >
          <span class="chip text-primary">
            <svg
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              aria-hidden="true"
            >
              <path d="M15 3h6v6M9 21H3v-6M21 3l-7 7M3 21l7-7" />
            </svg>
            view
          </span>
        </span>
      </div>
      <div class="flex flex-1 flex-col gap-3 p-5">
        <h3 class="m-0 line-clamp-2 text-[1.05rem] leading-snug">
          {cert.title}
        </h3>
        <div class="mt-auto flex flex-col gap-1.5">
          <span
            class="flex items-center gap-2 font-mono text-[0.78rem] text-text-light"
          >
            <svg
              width="14"
              height="14"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              aria-hidden="true"
              class="shrink-0 text-secondary"
            >
              <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
              <circle cx="12" cy="7" r="4" />
            </svg>
            <span class="truncate">{cert.issuer}</span>
          </span>
          <span
            class="flex items-center gap-2 font-mono text-[0.78rem] text-text-light"
          >
            <svg
              width="14"
              height="14"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              aria-hidden="true"
              class="shrink-0 text-secondary"
            >
              <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
              <line x1="16" y1="2" x2="16" y2="6" />
              <line x1="8" y1="2" x2="8" y2="6" />
              <line x1="3" y1="10" x2="21" y2="10" />
            </svg>
            <span class="truncate">{cert.date}</span>
          </span>
        </div>
      </div>
    </button>
  {/each}
</div>

{#if open && current}
  <div
    class={`fixed inset-0 z-[1000] flex items-center justify-center bg-[rgba(0,0,0,0.78)] p-3 backdrop-blur-md transition-opacity duration-200 md:p-6 ${
      visible ? "opacity-100" : "opacity-0"
    }`}
    bind:this={overlay}
    role="presentation"
    onclick={(e) => {
      if (e.target === e.currentTarget) close();
    }}
  >
    <div
      class={`panel flex max-h-[92vh] w-full max-w-[920px] flex-col overflow-hidden transition-transform duration-200 ease-[cubic-bezier(0.4,0,0.2,1)] ${
        visible ? "scale-100 translate-y-0" : "scale-[0.99] translate-y-2"
      }`}
      role="dialog"
      aria-modal="true"
      aria-label={current.title}
    >
      <header
        class="flex items-center justify-between gap-4 border-b border-border px-4 py-2.5"
      >
        <span class="eyebrow">certificate</span>
        <button
          type="button"
          class="inline-flex h-9 w-9 items-center justify-center rounded-sm border border-border text-text-light transition-colors hover:border-primary hover:text-primary"
          bind:this={closeBtn}
          onclick={close}
          aria-label="Close certificate viewer"
        >
          <svg
            width="22"
            height="22"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            aria-hidden="true"
          >
            <line x1="18" y1="6" x2="6" y2="18" />
            <line x1="6" y1="6" x2="18" y2="18" />
          </svg>
        </button>
      </header>

      <div
        class="relative flex min-h-0 flex-1 items-center justify-center bg-bg-alt"
      >
        {#if hasMany}
          <button
            type="button"
            class="absolute left-2 top-1/2 inline-flex h-9 w-9 -translate-y-1/2 items-center justify-center rounded-sm border border-border bg-surface text-text transition-colors hover:border-primary hover:text-primary sm:left-3.5 sm:h-[42px] sm:w-[42px]"
            onclick={prev}
            aria-label="Previous certificate"
          >
            <svg
              width="22"
              height="22"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
              aria-hidden="true"
            >
              <polyline points="15 18 9 12 15 6" />
            </svg>
          </button>
        {/if}

        <img
          class="block max-h-[calc(92vh-9rem)] max-w-full object-contain"
          src={current.url}
          alt={current.title}
        />

        {#if hasMany}
          <button
            type="button"
            class="absolute right-2 top-1/2 inline-flex h-9 w-9 -translate-y-1/2 items-center justify-center rounded-sm border border-border bg-surface text-text transition-colors hover:border-primary hover:text-primary sm:right-3.5 sm:h-[42px] sm:w-[42px]"
            onclick={next}
            aria-label="Next certificate"
          >
            <svg
              width="22"
              height="22"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
              aria-hidden="true"
            >
              <polyline points="9 18 15 12 9 6" />
            </svg>
          </button>
        {/if}
      </div>

      <footer class="border-t border-border px-5 py-4">
        <h2 class="m-0 text-[1.05rem]">{current.title}</h2>
        <p class="mt-1.5 font-mono text-[0.8rem] text-text-light">
          {current.issuer}<span class="text-secondary"> · </span>{current.date}
        </p>
        {#if hasMany}
          <span
            class="mt-2 inline-block font-mono text-[0.72rem] text-text-light"
          >
            {index + 1} / {images.length}
          </span>
        {/if}
      </footer>
    </div>
  </div>
{/if}
