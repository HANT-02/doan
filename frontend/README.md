# Frontend Project

This is the React frontend for the Golang backend, built with Vite + TypeScript + Tailwind CSS.

## Setup

1.  **Install dependencies**:
    ```bash
    cd frontend
    npm install
    ```

2.  **Environment Variables**:
    Copy `.env.example` to `.env` (optional, defaults are configured).
    ```bash
    VITE_API_BASE_URL=http://localhost:9000/v1/auth
    ```
    *Note: The vite config sets up a proxy, so calls to `/v1` are forwarded to `http://localhost:9000` by default.*

3.  **Run Development Server**:
    ```bash
    npm run dev
    ```
    Access at `http://localhost:5173`.

## Features implemented

- **Login**: `/login` (Email/Password) -> Saves JWT to memory/localStorage.
- **Register**: `/register` (Name, Email, Password).
- **Profile**: `/` (Protected) -> Shows user info, Logout.
- **Change Password**: `/change-password` (Protected).
- **Forgot Password**: `/forgot-password` -> Request reset link.
- **Reset Password**: `/reset-password?token=...` -> Set new password.

## Architecture

- **AuthContext**: Manages user state and tokens.
- **Axios Interceptor**: Handles API calls and potential errors.
- **UI Components**: Built with Tailwind CSS (shadcn-like structure).
- **Zod Validation**: Used for all forms.

## Security Note

Passwords are sent in plaintext over HTTPS (simulated locally) as requested (`password_enc` field in DTOs), because the backend handles hashing.
Refresh tokens are stored in `localStorage` to persist sessions, while Access tokens are kept in memory (Context) and re-hydrated on load.
