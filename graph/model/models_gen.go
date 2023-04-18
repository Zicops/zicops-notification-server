// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type ClassRoomFlags struct {
	ID                     *string   `json:"id"`
	IsClassroomStarted     *bool     `json:"is_classroom_started"`
	IsParticipantsPresent  *bool     `json:"is_participants_present"`
	IsAdDisplayed          *bool     `json:"is_ad_displayed"`
	IsBreak                *bool     `json:"is_break"`
	IsModeratorJoined      *bool     `json:"is_moderator_joined"`
	IsTrainerJoined        *bool     `json:"is_trainer_joined"`
	AdVideoURL             *string   `json:"ad_video_url"`
	IsMicrophoneEnabled    *bool     `json:"is_microphone_enabled"`
	IsVideoSharingEnabled  *bool     `json:"is_video_sharing_enabled"`
	IsScreenSharingEnabled *bool     `json:"is_screen_sharing_enabled"`
	IsChatEnabled          *bool     `json:"is_chat_enabled"`
	IsQaEnabled            *bool     `json:"is_qa_enabled"`
	Quiz                   []*string `json:"quiz"`
}

type ClassRoomFlagsInput struct {
	ID                     *string   `json:"id"`
	IsClassroomStarted     *bool     `json:"is_classroom_started"`
	IsParticipantsPresent  *bool     `json:"is_participants_present"`
	IsAdDisplayed          *bool     `json:"is__ad_displayed"`
	IsBreak                *bool     `json:"is_break"`
	IsModeratorJoined      *bool     `json:"is_moderator_joined"`
	IsTrainerJoined        *bool     `json:"is_trainer_joined"`
	AdVideoURL             *string   `json:"ad_video_url"`
	IsMicrophoneEnabled    *bool     `json:"is_microphone_enabled"`
	IsVideoSharingEnabled  *bool     `json:"is_video_sharing_enabled"`
	IsScreenSharingEnabled *bool     `json:"is_screen_sharing_enabled"`
	IsChatEnabled          *bool     `json:"is_chat_enabled"`
	IsQaEnabled            *bool     `json:"is_qa_enabled"`
	Quiz                   []*string `json:"quiz"`
}

type FirestoreData struct {
	Title     string  `json:"title"`
	Body      string  `json:"body"`
	CreatedAt int     `json:"created_at"`
	UserID    string  `json:"user_id"`
	IsRead    bool    `json:"is_read"`
	MessageID string  `json:"message_id"`
	Link      *string `json:"link"`
	LspID     string  `json:"lsp_id"`
}

type FirestoreDataInput struct {
	Title     string `json:"title"`
	Body      string `json:"body"`
	IsRead    bool   `json:"is_read"`
	MessageID string `json:"message_id"`
}

type FirestoreMessage struct {
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedAt int    `json:"created_at"`
	UserID    string `json:"user_id"`
	MessageID string `json:"message_id"`
	IsRead    bool   `json:"is_read"`
	Link      string `json:"link"`
	LspID     string `json:"lsp_id"`
}

type Messages struct {
	Body      *string `json:"body"`
	MeetingID *string `json:"meeting_id"`
	UserID    *string `json:"user_id"`
	Time      *int    `json:"time"`
	ChatType  *string `json:"chat_type"`
}

type Notification struct {
	Statuscode string  `json:"statuscode"`
	Error      *string `json:"error"`
	UserID     *string `json:"user_id"`
}

type NotificationInput struct {
	Title  string    `json:"title"`
	Body   string    `json:"body"`
	UserID []*string `json:"user_id"`
}

type PaginatedNotifications struct {
	Messages         []*FirestoreMessage `json:"messages"`
	NextPageSnapShot *string             `json:"nextPageSnapShot"`
}

type PaginatedTagsData struct {
	Data             []*TagsData `json:"data"`
	PrevPageSnapShot *string     `json:"prevPageSnapShot"`
}

type PollResponse struct {
	ID       *string `json:"id"`
	PollID   *string `json:"poll_id"`
	Response *string `json:"response"`
	UserIds  *string `json:"user_ids"`
}

type PollResponseInput struct {
	ID       *string `json:"id"`
	PollID   *string `json:"poll_id"`
	Response *string `json:"response"`
	UserIds  *string `json:"user_ids"`
}

type PollResults struct {
	PollID        *string         `json:"poll_id"`
	Question      *string         `json:"question"`
	PollResponses []*PollResponse `json:"poll_responses"`
}

type Polls struct {
	ID        *string   `json:"id"`
	MeetingID *string   `json:"meeting_id"`
	CourseID  *string   `json:"course_id"`
	TopicID   *string   `json:"topic_id"`
	Question  *string   `json:"question"`
	Options   []*string `json:"options"`
	PollIds   []*string `json:"poll_ids"`
	Status    *string   `json:"status"`
}

type PollsInput struct {
	ID        *string   `json:"id"`
	MeetingID *string   `json:"meeting_id"`
	CourseID  *string   `json:"course_id"`
	TopicID   *string   `json:"topic_id"`
	Question  *string   `json:"question"`
	Options   []*string `json:"options"`
	PollIds   []*string `json:"poll_ids"`
	Status    *string   `json:"status"`
}

type PublishedQuiz struct {
	ID     *string `json:"id"`
	QuizID *string `json:"quizId"`
}

type TagsData struct {
	UserLspID *string   `json:"user_lsp_id"`
	UserID    *string   `json:"user_id"`
	Tags      []*string `json:"tags"`
	LspID     *string   `json:"lsp_id"`
}

type UserDetails struct {
	UserID    *string `json:"user_id"`
	UserLspID *string `json:"user_lsp_id"`
}
