package httpapi

import (
	"strconv"
	"strings"

	"github.com/aisha-platform/aisha/apps/api/internal/features/notification"
	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/valyala/fasthttp"
)

type NotificationHandler struct {
	service *notification.Service
	hub     *notification.Hub
	origins map[string]struct{}
}

func NewNotificationHandler(service *notification.Service, hub *notification.Hub, origins []string) *NotificationHandler {
	allowed := make(map[string]struct{}, len(origins))
	for _, origin := range origins {
		allowed[origin] = struct{}{}
	}
	return &NotificationHandler{service: service, hub: hub, origins: allowed}
}

func (h *NotificationHandler) List(c fiber.Ctx) error {
	principal, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "20"))
	items, total, pageSize, err := h.service.List(c.Context(), principal, page, pageSize)
	if err != nil {
		return WrapAPIError(err, CodeInternalError, "Unable to load notifications.")
	}
	unread, err := h.service.UnreadCount(c.Context(), principal)
	if err != nil {
		return WrapAPIError(err, CodeInternalError, "Unable to load notification status.")
	}
	if page < 1 {
		page = 1
	}
	return c.JSON(map[string]any{"items": items, "page": page, "pageSize": pageSize, "total": total, "unreadCount": unread})
}

func (h *NotificationHandler) MarkRead(c fiber.Ctx) error {
	principal, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	if err := h.service.MarkRead(c.Context(), principal, c.Params("id")); err != nil {
		return WrapAPIError(err, CodeResourceNotFound, "Notification was not found.")
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *NotificationHandler) MarkAllRead(c fiber.Ctx) error {
	principal, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	if err := h.service.MarkAllRead(c.Context(), principal); err != nil {
		return WrapAPIError(err, CodeInternalError, "Unable to update notifications.")
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *NotificationHandler) SocketTicket(c fiber.Ctx) error {
	principal, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	ticket := h.service.IssueSocketTicket(principal)
	if ticket == "" {
		return NewAPIError(CodeInternalError, "Unable to open notification stream.", nil)
	}
	return c.JSON(map[string]string{"ticket": ticket})
}

func (h *NotificationHandler) WebSocket(c fiber.Ctx) error {
	ticket := strings.TrimSpace(c.Query("ticket"))
	userID, err := h.service.ConsumeSocketTicket(ticket)
	if err != nil {
		return NewAPIError(CodeAuthenticationRequired, "Authentication is required.", nil)
	}
	if !c.IsWebSocket() {
		return fiber.ErrUpgradeRequired
	}
	upgrader := websocket.FastHTTPUpgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 2048,
		CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
			origin := string(ctx.Request.Header.Peek("Origin"))
			if origin == "" {
				return true
			}
			_, ok := h.origins[origin]
			return ok
		},
	}
	return upgrader.Upgrade(c.RequestCtx(), func(conn *websocket.Conn) {
		client := h.hub.Add(userID, conn)
		defer h.hub.Remove(userID, client)
		defer client.Close()
		_ = client.WriteJSON(map[string]string{"type": "ready"})
		conn.SetPongHandler(func(string) error { return nil })
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	})
}
