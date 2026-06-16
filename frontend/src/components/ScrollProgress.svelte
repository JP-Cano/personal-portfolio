<script lang="ts">
  let width = $state(0);

  function update() {
    const windowHeight = window.innerHeight;
    const documentHeight = document.documentElement.scrollHeight;
    const scrollTop = window.scrollY || document.documentElement.scrollTop;
    const denominator = documentHeight - windowHeight;
    const pct = denominator > 0 ? (scrollTop / denominator) * 100 : 0;
    width = Math.min(Math.max(pct, 0), 100);
  }

  $effect(() => {
    update();
    window.addEventListener("scroll", update, { passive: true });
    window.addEventListener("resize", update);
    return () => {
      window.removeEventListener("scroll", update);
      window.removeEventListener("resize", update);
    };
  });
</script>

<div class="fixed top-0 left-0 w-full h-[3px] bg-transparent z-[9999]">
  <div
    class="h-full w-0 bg-[linear-gradient(90deg,var(--grad-1),var(--grad-2),var(--grad-3))] transition-[width] duration-100 ease-out"
    style={`width: ${width}%`}
  ></div>
</div>
