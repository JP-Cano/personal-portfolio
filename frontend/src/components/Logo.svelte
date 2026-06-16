<script lang="ts">
  /**
   * Wordmark logo: "JPC" set in the Clash Display face (--font-logo), filled
   * with the luminous indigo -> violet -> blue gradient that slowly drifts
   * (shimmers faster on hover). Theme-aware via the gradient tokens. Motion
   * stills under prefers-reduced-motion. Sized by the consumer via font-size.
   */
  interface Props {
    class?: string;
  }

  let { class: className = "" }: Props = $props();
</script>

<span class="logo {className}" role="img" aria-label="Juan Pablo Cano">JPC</span
>

<style>
  .logo {
    display: inline-block;
    font-family: var(--font-logo), sans-serif;
    font-weight: 700;
    font-size: 1.6rem;
    letter-spacing: -0.05em;
    line-height: 1;
    user-select: none;
    background-image: linear-gradient(
      100deg,
      var(--grad-1) 0%,
      var(--grad-2) 45%,
      var(--grad-3) 90%
    );
    background-size: 220% 100%;
    -webkit-background-clip: text;
    background-clip: text;
    color: transparent;
    animation: logo-shimmer 8s linear infinite;
    transition: filter 0.3s ease;
  }

  .logo:hover {
    animation-duration: 2.4s;
    filter: brightness(1.18);
  }

  /* Self-contained so the shimmer also runs in /admin (no global.css there). */
  @keyframes logo-shimmer {
    from {
      background-position: 0% 50%;
    }
    to {
      background-position: 200% 50%;
    }
  }

  @media (prefers-reduced-motion: reduce) {
    .logo {
      animation: none;
    }
  }
</style>
