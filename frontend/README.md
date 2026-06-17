# Enterprise Order Management Frontend

React + Vite demo client for the Golang backend API.

## Setup

```powershell
npm install
copy .env.example .env
npm run dev
```

Default environment:

```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

## Auth Flow

- `POST /auth/register`
- `POST /auth/login`
- `POST /auth/refresh-token`
- `POST /auth/logout`
- `GET /auth/me`

The demo stores `access_token`, `refresh_token`, and `user` in `localStorage`.

For production, the refresh token should be stored in an `httpOnly`, `secure`, `sameSite` cookie instead of `localStorage`. This reduces the risk of token theft from XSS attacks.

## Scripts

```powershell
npm run dev
npm run build
npm run preview
```
