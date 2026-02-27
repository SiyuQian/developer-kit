# PM Skill Design — Product Discovery & Market Research

## Overview

A general-purpose product manager skill that helps any project discover market needs through parallel deep research. Triggered via `/pm`, it runs 3 specialized agents concurrently to analyze competitors, user pain points, and industry trends, then synthesizes findings into 2-3 prioritized feature recommendations for collaborative decision-making.

## Trigger

- User invokes `/pm` or says "help me do market research" / "帮我做市场调研"
- Accepts a product direction as argument, e.g., `/pm "a task management tool for developers"`

## Process Flow

### Phase 1: Scope Clarification
- Ask 1-2 questions to clarify product positioning
- Determine target users, core scenarios, and competitive landscape
- Establish search keywords for the agents

### Phase 2: Parallel Deep Research (3 agents)

**Agent 1 — Competitor Analyst**
- Search for alternative products, feature comparisons, pricing models
- Search patterns: `"{type} alternatives 2026"`, `"{competitor} features pricing"`, `"{a} vs {b}"`
- Output: competitor table (name, features, pricing, strengths, weaknesses, popularity)

**Agent 2 — User Pain Analyst**
- Search Reddit, forums, review sites for real user complaints and wishlists
- Search patterns: `"{type} pain points reddit"`, `"{type} missing features forum"`, `"{type} frustrating"`
- Output: pain point list (description, frequency, source, whether competitors solved it)

**Agent 3 — Trend Analyst**
- Search for industry trends, emerging technologies, market growth data
- Search patterns: `"{type} trends 2026"`, `"{type} emerging technology"`, `"{tech} market growth"`
- Output: trend list (direction, maturity, related tech, potential opportunity)

### Phase 3: Synthesis
- Cross-validate findings across all 3 agents
- Scoring criteria:
  - Market demand intensity (how many users mention it)
  - Competitive gap (where competitors are weak)
  - Trend alignment (matches industry direction)
  - Implementation feasibility (rough estimate based on project's tech stack)
- Rank and select top 2-3 feature recommendations

### Phase 4: Collaborative Decision
- Present each feature with: description, market evidence, competitor reference, priority recommendation
- Discuss with user, answer follow-up questions
- Reach consensus on what to build

## File Structure

```
skills/pm/
├── SKILL.md
└── references/
    └── search-patterns.md
```

## Key Design Decisions

1. **WebSearch as primary data source** — no external API dependencies
2. **Output stays in conversation** — no auto-write to Trello or files; user decides next steps
3. **User-invoked only** — triggered via `/pm`, not automatically
4. **3 parallel agents** — maximizes coverage while staying focused
5. **2-3 feature output** — actionable, not overwhelming
