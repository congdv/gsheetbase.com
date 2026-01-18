Act as a Senior Full-Stack Engineer and Architect. I want to build a "Sheets-to-API" platform. 

# Tech Stack:
- Backend: Go (Golang) using a "Light BFF" pattern.
- Frontend: React with Vite (TypeScript), Tailwind CSS.
- Infrastructure: Monorepo structured for Railway deployment (separate /ui and /api folders).
- Database: PostgreSQL (for user metadata) and Redis (for API caching).

# Project Goal:
Allow users to connect their Google Sheets via OAuth and generate a public/private REST API endpoint that returns their sheet data as clean JSON.

# Suggestions for Copilot
- Provide the simplest correct solution first.
- Explain trade-offs briefly if multiple approaches exist.
- Do not rewrite unrelated code.
- Do not need to run the app