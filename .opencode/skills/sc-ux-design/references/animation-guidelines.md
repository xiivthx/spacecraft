# Animation Guidelines

Quality-control reference for UI transitions and motion design. Apply these rules when implementing animated interfaces.

## Duration Standards

- **Micro-interactions** (button press, hover, toggle): 150–300ms
- **Complex transitions** (page enter/exit, modal open/close): 200–400ms
- **Maximum**: 400ms for any UI transition. Longer durations feel sluggish.
- **Exception**: Ambient/background animations (loading indicators, decorative) may run longer if non-blocking.

## Easing Rules

- **Enter/Appear**: `ease-out` (fast start, gentle settle). Makes elements feel responsive.
  - CSS: `ease-out`, `cubic-bezier(0, 0, 0.2, 1)`
  - Preferred: `cubic-bezier(0.16, 1, 0.3, 1)` (ease-out-expo)
- **Exit/Disappear**: `ease-in` (gentle start, fast finish). Elements leaving should not linger.
  - CSS: `ease-in`, `cubic-bezier(0.4, 0, 1, 1)`
- **Continuous/Indeterminate**: `linear` only when motion has no start or end (spinners, marquees).
- **No linear for UI**: Linear easing on discrete UI transitions feels mechanical and unfinished.

## Reduced Motion

- **Must**: Respect `prefers-reduced-motion: reduce` media query.
- **Must**: Disable all non-essential animations when reduced motion is requested.
- **Must**: Keep essential functional animations (spinners, progress bars) but reduce their duration.
- **Implementation**:

  ```css
  @media (prefers-reduced-motion: reduce) {
    *, *::before, *::after {
      animation-duration: 0.01ms !important;
      animation-iteration-count: 1 !important;
      transition-duration: 0.01ms !important;
    }
  }
  ```

## Anti-Patterns

### 1. No bounce or elastic easing

Bounce (`cubic-bezier` with overshoot) and elastic easing on interface elements feels dated and tacky. Reserve spring physics for things that are actually physical (drag-to-dismiss, pull-to-refresh). Interface elements that appear/disappear should use smooth ease-out curves.

**Bad**: `cubic-bezier(0.68, -0.55, 0.27, 1.55)` — dialog "springs" in with overshoot
**Good**: `cubic-bezier(0.16, 1, 0.3, 1)` — dialog slides in smoothly

### 2. No width/height animation

Animating `width`, `height`, `padding`, or `margin` triggers layout recalculations on every frame, causing jank and poor performance. Use `transform` and `opacity` instead.

**Bad**: `transition: width 300ms` — forces layout thrash
**Good**: `transition: transform 300ms` with `scaleX()` or `grid-template-rows`

### 3. No decorative-only animation

Motion must convey meaning. An animation that serves no functional purpose (bouncing logo, animated background gradient, pulsing decorative elements) distracts and drains attention.

**Bad**: Hero section background subtly shifting gradient colors on loop
**Good**: Card scales up 2% on hover to indicate interactivity

### 4. Animations must be interruptible

Animations that must complete before accepting new input feel unresponsive. If the user triggers a new animation mid-flight, the old one should be interruptible and the new one should start immediately.

**Bad**: CSS `animation` that plays to completion before responding to pointer events
**Good**: CSS `transition` (inherently interruptible) or Web Animation API with cancel-on-new-input

### 5. No image hover transform (AI tell)

Scaling or rotating images on hover is a recurring AI-generated signature. Let imagery sit still, or use a subtler, purposeful interaction (e.g., slight brightness shift).

## Motion Philosophy

- **Motion conveys meaning**: Animate to show relationship (parent→child), state change (open→closed), or hierarchy (important→detail).
- **Motion guides attention**: The most important change should move first or move most.
- **Stagger for lists**: When multiple items appear, stagger by 30–60ms per item (total under 400ms for the group).
- **Exit before enter**: When replacing content, animate the old content out before animating the new content in. Avoid simultaneous crossfade.
- **One thing at a time**: Avoid animating more than 2–3 elements simultaneously. The eye can only track so much.

## References

- [impeccable.style/slop — Motion](https://impeccable.style/slop#motion) — 3 motion detection rules
- [Material Design — Motion](https://m3.material.io/styles/motion) — duration, easing, and choreography patterns
- [prefers-reduced-motion](https://developer.mozilla.org/en-US/docs/Web/CSS/@media/prefers-reduced-motion) — MDN reference
