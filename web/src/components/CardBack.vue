<template>
  <svg
    class="card-back-svg"
    viewBox="0 0 100 175"
    preserveAspectRatio="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <defs>
      <linearGradient id="cb-bg" x1="0" y1="0" x2="0" y2="1">
        <stop offset="0%" stop-color="#241828" />
        <stop offset="50%" stop-color="#160c1a" />
        <stop offset="100%" stop-color="#241828" />
      </linearGradient>
      <linearGradient id="cb-gold" x1="0" y1="0" x2="0" y2="1">
        <stop offset="0%" stop-color="#e5d8b0" />
        <stop offset="50%" stop-color="#c2a35f" />
        <stop offset="100%" stop-color="#8a6e2e" />
      </linearGradient>
      <radialGradient id="cb-glow" cx="0.5" cy="0.5" r="0.45">
        <stop offset="0%" stop-color="#c2a35f" stop-opacity="0.2" />
        <stop offset="70%" stop-color="#c2a35f" stop-opacity="0.03" />
        <stop offset="100%" stop-color="#c2a35f" stop-opacity="0" />
      </radialGradient>
    </defs>

    <!-- Background -->
    <rect width="100" height="175" fill="url(#cb-bg)" />
    <rect width="100" height="175" fill="url(#cb-glow)" />

    <!-- Subtle dot pattern -->
    <g fill="url(#cb-gold)" opacity="0.1">
      <circle cx="20" cy="30" r="0.5" />
      <circle cx="80" cy="30" r="0.5" />
      <circle cx="20" cy="145" r="0.5" />
      <circle cx="80" cy="145" r="0.5" />
      <circle cx="15" cy="87.5" r="0.4" />
      <circle cx="85" cy="87.5" r="0.4" />
      <circle cx="35" cy="50" r="0.35" />
      <circle cx="65" cy="50" r="0.35" />
      <circle cx="35" cy="125" r="0.35" />
      <circle cx="65" cy="125" r="0.35" />
    </g>

    <!-- Outer border -->
    <rect x="3.5" y="3.5" width="93" height="168" rx="3" fill="none" stroke="url(#cb-gold)" stroke-width="0.7" />
    <!-- Inner border -->
    <rect x="6.5" y="6.5" width="87" height="162" rx="2" fill="none" stroke="url(#cb-gold)" stroke-width="0.3" opacity="0.6" />

    <!-- Corner ornaments -->
    <g fill="url(#cb-gold)">
      <circle cx="8" cy="8" r="1" />
      <circle cx="92" cy="8" r="1" />
      <circle cx="8" cy="167" r="1" />
      <circle cx="92" cy="167" r="1" />
    </g>
    <g fill="none" stroke="url(#cb-gold)" stroke-width="0.35" opacity="0.7">
      <path d="M 6.5 13 Q 6.5 6.5 13 6.5" />
      <path d="M 93.5 13 Q 93.5 6.5 87 6.5" />
      <path d="M 6.5 162 Q 6.5 168.5 13 168.5" />
      <path d="M 93.5 162 Q 93.5 168.5 87 168.5" />
    </g>

    <!-- Central mandala -->
    <g transform="translate(50, 87.5)">
      <!-- Long rays -->
      <g stroke="url(#cb-gold)" stroke-width="0.4" opacity="0.55">
        <line v-for="r in longRays" :key="'L' + r.i" :x1="r.x1" :y1="r.y1" :x2="r.x2" :y2="r.y2" />
      </g>
      <!-- Short rays -->
      <g stroke="url(#cb-gold)" stroke-width="0.25" opacity="0.35">
        <line v-for="r in shortRays" :key="'S' + r.i" :x1="r.x1" :y1="r.y1" :x2="r.x2" :y2="r.y2" />
      </g>

      <!-- Outer ring -->
      <circle r="28" fill="none" stroke="url(#cb-gold)" stroke-width="0.5" />
      <circle r="25" fill="none" stroke="url(#cb-gold)" stroke-width="0.2" opacity="0.4" />

      <!-- Inner ring -->
      <circle r="16" fill="none" stroke="url(#cb-gold)" stroke-width="0.4" />

      <!-- 6-pointed star (hexagram) -->
      <path :d="starPath" fill="url(#cb-gold)" opacity="0.85" />

      <!-- Center dot -->
      <circle r="1.8" fill="#160c1a" />
      <circle r="1.8" fill="none" stroke="url(#cb-gold)" stroke-width="0.3" />
    </g>

    <!-- Top decoration -->
    <g transform="translate(50, 20)" fill="url(#cb-gold)">
      <path d="M 0,-7 L 2,-2 L 7,0 L 2,2 L 0,7 L -2,2 L -7,0 L -2,-2 Z" />
    </g>
    <g fill="none" stroke="url(#cb-gold)" stroke-width="0.25" opacity="0.4">
      <line x1="34" y1="20" x2="40" y2="20" />
      <line x1="60" y1="20" x2="66" y2="20" />
    </g>

    <!-- Bottom decoration -->
    <g transform="translate(50, 155)" fill="url(#cb-gold)">
      <path d="M 0,-7 L 2,-2 L 7,0 L 2,2 L 0,7 L -2,2 L -7,0 L -2,-2 Z" />
    </g>
    <g fill="none" stroke="url(#cb-gold)" stroke-width="0.25" opacity="0.4">
      <line x1="34" y1="155" x2="40" y2="155" />
      <line x1="60" y1="155" x2="66" y2="155" />
    </g>

    <!-- Side flourishes -->
    <g fill="none" stroke="url(#cb-gold)" stroke-width="0.3" opacity="0.4">
      <path d="M 10 87.5 Q 14 84 18 87.5 Q 14 91 10 87.5 Z" />
      <path d="M 90 87.5 Q 86 84 82 87.5 Q 86 91 90 87.5 Z" />
    </g>
  </svg>
</template>

<script setup>
import { computed } from 'vue';

const longRays = computed(() => {
  const rays = [];
  const innerR = 17;
  const outerR = 30;
  for (let i = 0; i < 8; i++) {
    const angle = (i * 45 - 90) * Math.PI / 180;
    rays.push({
      i,
      x1: innerR * Math.cos(angle),
      y1: innerR * Math.sin(angle),
      x2: outerR * Math.cos(angle),
      y2: outerR * Math.sin(angle),
    });
  }
  return rays;
});

const shortRays = computed(() => {
  const rays = [];
  const innerR = 19;
  const outerR = 27;
  for (let i = 0; i < 8; i++) {
    const angle = (i * 45 + 22.5 - 90) * Math.PI / 180;
    rays.push({
      i,
      x1: innerR * Math.cos(angle),
      y1: innerR * Math.sin(angle),
      x2: outerR * Math.cos(angle),
      y2: outerR * Math.sin(angle),
    });
  }
  return rays;
});

const starPath = computed(() => {
  const outerR = 32;
  const innerR = 18.48; // outerR * cos(30°) ≈ 16 * 0.866
  const points = [];
  for (let i = 0; i < 12; i++) {
    const r = i % 2 === 0 ? outerR : innerR;
    const angle = (i * 30 - 90) * Math.PI / 180;
    points.push(`${(r * Math.cos(angle)).toFixed(2)},${(r * Math.sin(angle)).toFixed(2)}`);
  }
  return `M ${points[0]} L ${points.slice(1).join(' L ')} Z`;
});
</script>

<style scoped>
.card-back-svg {
  width: 100%;
  height: 100%;
  display: block;
  pointer-events: none;
}
</style>
