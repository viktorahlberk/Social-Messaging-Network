package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/models"
	"social-network/pkg/utils"
	ws "social-network/pkg/wsServer"
	"strings"
)

func (handler *Handler) NewEvent(wsServer *ws.Server, w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	if r.Method != "POST" {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	/* ---------------------------- read incoming data --------------------------- */
	// Try to decode the JSON request to Event
	var event models.Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	event.ID = utils.UniqueId()
	event.AuthorID = r.Context().Value(utils.UserKey).(string)
	/* -------------------- check if user is a meber of group ------------------- */
	var isMember = false
	isAdmin, err := handler.repos.GroupRepo.IsAdmin(event.GroupID, event.AuthorID)
	if err != nil {
		utils.RespondWithError(w, "Error on reading role", 200)
		return
	}
	if !isAdmin {
		isMember, err = handler.repos.GroupRepo.IsMember(event.GroupID, event.AuthorID)
		if err != nil {
			utils.RespondWithError(w, "Error on checking if is group member", 200)
			return
		}
	}
	if !isMember && !isAdmin {
		utils.RespondWithError(w, "Not a member", 200)
		return
	}
	/* ------------------------- save event in database ------------------------- */
	if err = handler.repos.EventRepo.Save(event); err != nil {
		utils.RespondWithError(w, "Internal server error", 200)
		return
	}
	/* ----------------- if user going also save as participant ----------------- */
	if strings.ToUpper(event.Going) == "YES" {
		if err = handler.repos.EventRepo.AddParticipant(event.ID, event.AuthorID); err != nil {
			utils.RespondWithError(w, "Internal server error", 200)
			return
		}
	}
	/* -------------------- save new notification about event ------------------- */
	// get all group members
	members, err := handler.repos.GroupRepo.GetMembers(event.GroupID)
	if err != nil {
		utils.RespondWithError(w, "Internal server error", 200)
		return
	}

	// for each member create notification
	for i := 0; i < len(members); i++ {
		newNotif := models.Notification{
			ID:       utils.UniqueId(),
			TargetID: members[i].ID,
			Type:     "EVENT",
			Content:  event.ID,
			Sender:   event.AuthorID,
		}
		// save notification in database
		if members[i].ID != event.AuthorID {
			err = handler.repos.NotifRepo.Save(newNotif)
			if err != nil {
				utils.RespondWithError(w, "Internal server error", 200)
				return
			}
		}
		
		// NOTIFY  GROUP MEMBER ABOUT THE NEW EVENT IF ONLINE
		for client := range wsServer.Clients {
			if client.ID == members[i].ID && client.ID != event.AuthorID{
				client.SendNotification(newNotif)
			}
		}

	}
	utils.RespondWithEvents(w, []models.Event{event}, 200)
}

// requestId: notification.id,
// eventId: notification.event.id,
// response: reqResponse,
// Handles clients reaction to participation in event
// waits for POST req with eventID as "id" and user status "going" with response YES or NO
func (handler *Handler) Participate(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	if r.Method != "POST" {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	/* --------------------------- read incoming data --------------------------- */
	type Response struct {
		EventID   string `json:"eventId"`
		RequestID string `json:"requestId"` //notif id
		Response  string `json:"response"` //YES || NO
	}
	var response Response
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	// get current user
	userId := r.Context().Value(utils.UserKey).(string)

	/* ---------------- check that event id and response provided --------------- */
	if len(response.EventID) == 0 || len(response.Response) == 0 {
		utils.RespondWithError(w, "Provided incomplete data", 200)
		return
	}

	/* ------------------- check if response alredy registerd ------------------- */
	isParticipating, err := handler.repos.EventRepo.IsParticipating(response.EventID, userId)
	if err != nil {
		utils.RespondWithError(w, "Internal server error", 200)
		return
	}
	/* ----------------------------- handle response ---------------------------- */
	if response.Response == "YES" && !isParticipating {
		if err = handler.repos.EventRepo.AddParticipant(response.EventID, userId); err != nil {
			utils.RespondWithError(w, "Internal server error", 200)
			return
		}
	} else if strings.ToUpper(response.Response) == "NO" && isParticipating {
		if err = handler.repos.EventRepo.RemoveParticipant(response.EventID, userId); err != nil {
			utils.RespondWithError(w, "Internal server error", 200)
			return
		}
	}
	/* --------------------------- remove notificaton -------------------------- */
	if len(response.RequestID) !=0 {//participation activated form notification
		if err = handler.repos.NotifRepo.Delete(response.RequestID); err != nil {
			utils.RespondWithError(w, "Internal server error", 200)
			return
		}
	}else{ //participation activated without noification
		// delete notification if exists
		notif := models.Notification{Type: "EVENT", TargetID: userId, Content:response.EventID}
		if err = handler.repos.NotifRepo.DeleteByType(notif); err!=nil{
			utils.RespondWithError(w, "Internal server error", 200)
			return
		}

	}
	utils.RespondWithSuccess(w, "Data saved successfully", 200)
}
