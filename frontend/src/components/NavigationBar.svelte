<script lang="ts">
  import Logo from "@/components/Logo.svelte";
  import ThemeToggle from "@/components/ThemeToggle.svelte";

  interface Props {
    cvHref: string;
  }

  let { cvHref }: Props = $props();

  let scrolled = $state(false);
  let menuOpen = $state(false);
  let isHome = $state(true);

  const baseLinks = [
    { hash: "#about", label: "About" },
    { hash: "#experience", label: "Experience" },
    { hash: "#projects", label: "Projects" },
    { hash: "#skills", label: "Skills" },
  ];

  function updateLocation() {
    isHome = window.location.pathname === "/";
  }

  function navClick(e: MouseEvent, href: string) {
    e.preventDefault();
    if (href.startsWith("#") && isHome) {
      try {
        const target = document.querySelector(href);
        if (target) {
          target.scrollIntoView({ behavior: "smooth", block: "start" });
          menuOpen = false;
        }
      } catch {
        // Invalid selector
      }
    } else {
      window.location.href = href;
    }
  }

  $effect(() => {
    updateLocation();
    const onScroll = () => {
      scrolled = window.scrollY > 80;
    };
    onScroll();
    window.addEventListener("scroll", onScroll, { passive: true });
    return () => window.removeEventListener("scroll", onScroll);
  });
</script>

<nav
  class={`fixed inset-x-0 top-0 z-[1000] px-8 border-b border-transparent transition-[padding,background,border-color] duration-300 ${
    scrolled
      ? "py-[0.85rem] bg-[var(--nav-surface)] backdrop-blur-[12px] border-b-border"
      : "py-5"
  }`}
  aria-label="Primary"
>
  <div class="max-w-[1320px] mx-auto flex justify-between items-center gap-6">
    <a
      href={isHome ? "#hero" : "/"}
      class="h-[30px] text-text transition-transform duration-[250ms] hover:-translate-y-px"
      aria-label="Home"
      onclick={(e) => navClick(e, isHome ? "#hero" : "/")}
    >
      <Logo />
    </a>

    <ul
      class={`flex items-center gap-7 m-0 p-0 max-[820px]:fixed max-[820px]:top-16 max-[820px]:flex-col max-[820px]:items-start max-[820px]:gap-5 max-[820px]:p-7 max-[820px]:bg-[var(--nav-surface)] max-[820px]:backdrop-blur-[12px] max-[820px]:border max-[820px]:border-border max-[820px]:rounded max-[820px]:shadow-[0_20px_50px_rgba(0,0,0,0.25)] max-[820px]:transition-[right] max-[820px]:duration-300 ${
        menuOpen ? "max-[820px]:right-6" : "max-[820px]:right-[-110%]"
      }`}
    >
      {#each baseLinks as link (link.hash)}
        <li>
          <a
            href={isHome ? link.hash : `/${link.hash}`}
            class="font-sans text-[1rem] font-medium text-text-light transition-colors duration-[250ms] hover:text-text"
            onclick={(e) => navClick(e, isHome ? link.hash : `/${link.hash}`)}
          >
            {link.label}
          </a>
        </li>
      {/each}
      <li>
        <a
          href={isHome ? "#contact" : "/#contact"}
          class="border-primary/45 text-primary hover:bg-primary/10 inline-flex items-center justify-center rounded-[var(--radius-sm)] border px-[1.05rem] py-2 font-sans text-[1rem] font-medium transition-colors duration-[250ms]"
          onclick={(e) => navClick(e, isHome ? "#contact" : "/#contact")}
        >
          Contact
        </a>
      </li>
    </ul>

    <div class="flex items-center gap-[0.6rem]">
      <ThemeToggle />
      <a
        href={cvHref}
        download="Juan_Pablo_Cano_CV.docx"
        class="btn-ghost"
        title="Download CV"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="16"
          height="16"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          aria-hidden="true"
        >
          <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
          <polyline points="7 10 12 15 17 10"></polyline>
          <line x1="12" y1="15" x2="12" y2="3"></line>
        </svg>
        <span class="max-[820px]:hidden">cv</span>
      </a>
      <button
        class="hidden max-[820px]:flex flex-col gap-[5px] p-1.5 bg-transparent"
        aria-label="Toggle menu"
        aria-expanded={menuOpen}
        onclick={() => (menuOpen = !menuOpen)}
      >
        <span
          class={`block w-[22px] h-0.5 rounded-sm bg-text transition-all duration-300 ${
            menuOpen ? "[transform:rotate(45deg)_translateY(7px)]" : ""
          }`}
        ></span>
        <span
          class={`block w-[22px] h-0.5 rounded-sm bg-text transition-all duration-300 ${
            menuOpen ? "opacity-0" : ""
          }`}
        ></span>
        <span
          class={`block w-[22px] h-0.5 rounded-sm bg-text transition-all duration-300 ${
            menuOpen ? "[transform:rotate(-45deg)_translateY(-7px)]" : ""
          }`}
        ></span>
      </button>
    </div>
  </div>
</nav>
