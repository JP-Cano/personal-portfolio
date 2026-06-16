<script lang="ts">
  // Interactive island. The actual theme resolution, persistence and no-FOUC
  // application live in ThemeProvider.astro's inline head script; this button
  // only requests a toggle and mirrors the resolved theme in its icon.
  let isDark = $state(false);

  function readTheme() {
    isDark = document.documentElement.getAttribute("data-theme") === "dark";
  }

  function toggle() {
    window.dispatchEvent(new Event("theme:toggle"));
  }

  $effect(() => {
    readTheme();

    const onChanged = () => readTheme();
    const onStorage = (e: StorageEvent) => {
      if (e.key === "theme") readTheme();
    };

    window.addEventListener("theme:changed", onChanged);
    window.addEventListener("storage", onStorage);

    return () => {
      window.removeEventListener("theme:changed", onChanged);
      window.removeEventListener("storage", onStorage);
    };
  });
</script>

<button
  class="btn-ghost w-10 h-10 p-0 hover:translate-y-0"
  data-mode={isDark ? "dark" : "light"}
  aria-pressed={isDark}
  aria-label="Toggle dark mode"
  title="Toggle theme"
  onclick={toggle}
>
  {#if isDark}
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="17"
      height="17"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
      aria-hidden="true"
    >
      <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path>
    </svg>
  {:else}
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="17"
      height="17"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
      aria-hidden="true"
    >
      <circle cx="12" cy="12" r="4"></circle>
      <path
        d="M12 2v2M12 20v2M4.93 4.93l1.41 1.41M17.66 17.66l1.41 1.41M2 12h2M20 12h2M6.34 17.66l-1.41 1.41M19.07 4.93l-1.41 1.41"
      ></path>
    </svg>
  {/if}
</button>
