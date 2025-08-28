package auth

import (
	"encoding/json"
	"net/http"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api/middleware"
	jwtpgk "github.com/happilymarrieddad/order-management-v3/api/internal/jwt"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	"golang.org/x/crypto/bcrypt"
)

// @Summary      User Login
// @Description  Authenticates a user and returns a JWT token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials body      LoginPayload             true  "Login Credentials"
// @Success      200         {object}  LoginResponse            "Successfully authenticated"
// @Failure      400         {object}  middleware.ErrorResponse "Bad Request - Invalid input"
// @Failure      401         {object}  middleware.ErrorResponse "Unauthorized - Invalid credentials"
// @Failure      500         {object}  middleware.ErrorResponse "Internal Server Error"
// @Router       /auth/login [post]
// Login handles user authentication.
func Login(w http.ResponseWriter, r *http.Request) {
	var payload LoginPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := types.Validate(payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "validation failed: "+err.Error())
		return
	}

	repo := middleware.GetRepo(r.Context())

	user, found, err := repo.Users().GetByEmail(r.Context(), payload.Email)
	if err != nil {
		// Log the internal error but return a generic message to the client.
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if !found {
		middleware.WriteError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		middleware.WriteError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	token, err := jwtpgk.GenerateToken(user)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}
