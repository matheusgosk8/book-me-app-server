package handlers

import (
    "encoding/json"
    "net/http"
    "time"

    "github.com/matheusgosk8/book-me-server/ent/session"
    "github.com/matheusgosk8/book-me-server/internal/db"
    "github.com/matheusgosk8/book-me-server/internal/utils"
    log "github.com/sirupsen/logrus"
)

type RefreshRequest struct {
    RefreshToken string `json:"refresh_token"`
}

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
    var req RefreshRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Token inválido", http.StatusBadRequest)
        return
    }

    // 1. Valida o Token JWT e extrai o ID do usuário
    claims, err := utils.ValidateToken(req.RefreshToken)
    if err != nil || claims.Type != "refresh" {
        http.Error(w, "Refresh token inválido ou expirado", http.StatusUnauthorized)
        return
    }

    // 2. Busca a sessão no banco de dados pelo refresh_token
    sess, err := db.Client.Session.
        Query().
        Where(session.RefreshTokenEQ(req.RefreshToken)).
        WithUser(). // Carrega o usuário relacionado
        Only(r.Context())

    if err != nil {
        log.Warnf("Sessão não encontrada para token: %s", req.RefreshToken[:20])
        http.Error(w, "Sessão não encontrada", http.StatusUnauthorized)
        return
    }

    // 3. Verifica se a sessão expirou
    if sess.ExpiresAt.Before(time.Now()) {
        log.Warnf("Sessão expirada para user: %s", sess.Edges.User.Email)
        http.Error(w, "Sessão expirada", http.StatusUnauthorized)
        return
    }

    // 4. Gera novo Access Token
    newAccessToken, _, err := utils.GenerateTokens(sess.Edges.User.ID.String(), sess.Edges.User.UserType)
    if err != nil {
        log.WithError(err).Error("[Refresh] Falha ao regenerar tokens")
        http.Error(w, "Erro ao renovar acesso", http.StatusInternalServerError)
        return
    }

    // 5. Atualiza last_login_at da sessão
    _, err = db.Client.Session.
        UpdateOne(sess).
        SetLastLoginAt(time.Now()).
        Save(r.Context())

    if err != nil {
        log.WithError(err).Error("[Refresh] Falha ao atualizar sessão")
        http.Error(w, "Erro ao atualizar sessão", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "access_token": newAccessToken,
    })
}