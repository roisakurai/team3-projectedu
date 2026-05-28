package model

// AssignmentCreatedEvent is published on "lms:assignment:created"
type AssignmentCreatedEvent struct {
	Event        string `json:"event"`
	AssignmentID string `json:"assignment_id"`
	ClassID      string `json:"class_id"`
	Title        string `json:"title"`
	Deadline     string `json:"deadline"`
	TeacherID    string `json:"teacher_id"`
	Timestamp    string `json:"timestamp"`
}

// AssignmentSubmittedEvent is published on "lms:assignment:submitted"
type AssignmentSubmittedEvent struct {
	Event        string `json:"event"`
	SubmissionID string `json:"submission_id"`
	AssignmentID string `json:"assignment_id"`
	StudentID    string `json:"student_id"`
	ClassID      string `json:"class_id"`
	Timestamp    string `json:"timestamp"`
}

// AssignmentGradedEvent is published on "lms:assignment:graded"
type AssignmentGradedEvent struct {
	Event        string  `json:"event"`
	SubmissionID string  `json:"submission_id"`
	AssignmentID string  `json:"assignment_id"`
	StudentID    string  `json:"student_id"`
	ClassID      string  `json:"class_id"`
	Grade        float64 `json:"grade"`
	Timestamp    string  `json:"timestamp"`
}

// MaterialCreatedEvent is published on "lms:material:created"
type MaterialCreatedEvent struct {
	Event      string `json:"event"`
	MaterialID string `json:"material_id"`
	ClassID    string `json:"class_id"`
	Title      string `json:"title"`
	TeacherID  string `json:"teacher_id"`
	Timestamp  string `json:"timestamp"`
}

// MaterialCompletedEvent is published on "lms:material:completed"
type MaterialCompletedEvent struct {
	Event      string `json:"event"`
	StudentID  string `json:"student_id"`
	MaterialID string `json:"material_id"`
	ClassID    string `json:"class_id"`
	Timestamp  string `json:"timestamp"`
}

// StudentJoinedClassEvent is published on "lms:class:student_joined"
type StudentJoinedClassEvent struct {
	Event     string `json:"event"`
	StudentID string `json:"student_id"`
	ClassID   string `json:"class_id"`
	Timestamp string `json:"timestamp"`
}
