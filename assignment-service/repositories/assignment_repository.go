package repositories

import (
	"context"
	"errors"
	"time"

	"assignment-service/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	ErrNotFound      = errors.New("document not found")
	ErrAlreadyExists = errors.New("document already exists")
)

// AssignmentRepository defines the persistence interface.
type AssignmentRepository interface {
	// Assignments
	CreateAssignment(ctx context.Context, a *models.Assignment) error
	GetAssignmentByID(ctx context.Context, id string) (*models.Assignment, error)
	ListAssignmentsByClassID(ctx context.Context, classID string) ([]*models.Assignment, error)
	UpdateAssignment(ctx context.Context, id string, update bson.M) error
	DeleteAssignment(ctx context.Context, id string) error

	// Submissions
	CreateSubmission(ctx context.Context, s *models.Submission) error
	GetSubmissionByID(ctx context.Context, id string) (*models.Submission, error)
	GetSubmissionByAssignmentAndStudent(ctx context.Context, assignmentID, studentID string) (*models.Submission, error)
	ListSubmissionsByAssignment(ctx context.Context, assignmentID string) ([]*models.Submission, error)
	UpdateSubmission(ctx context.Context, id string, update bson.M) error
}

type mongoRepository struct {
	assignments *mongo.Collection
	submissions *mongo.Collection
}

// NewMongoRepository creates a new repository and ensures indexes.
func NewMongoRepository(db *mongo.Database) (AssignmentRepository, error) {
	repo := &mongoRepository{
		assignments: db.Collection("assignments"),
		submissions: db.Collection("submissions"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := repo.ensureIndexes(ctx); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *mongoRepository) ensureIndexes(ctx context.Context) error {
	assignmentIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "class_id", Value: 1}},
			Options: options.Index().SetName("idx_class_id"),
		},
		{
			Keys:    bson.D{{Key: "teacher_id", Value: 1}},
			Options: options.Index().SetName("idx_teacher_id"),
		},
	}

	if _, err := r.assignments.Indexes().CreateMany(ctx, assignmentIndexes); err != nil {
		return err
	}

	submissionIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "assignment_id", Value: 1}},
			Options: options.Index().SetName("idx_sub_assignment_id"),
		},
		{
			Keys:    bson.D{{Key: "student_id", Value: 1}},
			Options: options.Index().SetName("idx_sub_student_id"),
		},
		{
			// Unique constraint: one submission per student per assignment
			Keys: bson.D{
				{Key: "assignment_id", Value: 1},
				{Key: "student_id", Value: 1},
			},
			Options: options.Index().
				SetName("idx_sub_assignment_student_unique").
				SetUnique(true),
		},
	}

	if _, err := r.submissions.Indexes().CreateMany(ctx, submissionIndexes); err != nil {
		return err
	}

	return nil
}

// ── Assignment methods ────────────────────────────────────────────────────────

func (r *mongoRepository) CreateAssignment(ctx context.Context, a *models.Assignment) error {
	a.ID = bson.NewObjectID()
	now := time.Now().UTC()
	a.CreatedAt = now
	a.UpdatedAt = now

	_, err := r.assignments.InsertOne(ctx, a)
	return err
}

func (r *mongoRepository) GetAssignmentByID(ctx context.Context, id string) (*models.Assignment, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrNotFound
	}

	var a models.Assignment
	err = r.assignments.FindOne(ctx, bson.M{"_id": oid}).Decode(&a)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotFound
	}
	return &a, err
}

func (r *mongoRepository) ListAssignmentsByClassID(ctx context.Context, classID string) ([]*models.Assignment, error) {
	cursor, err := r.assignments.Find(ctx, bson.M{"class_id": classID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var assignments []*models.Assignment
	if err := cursor.All(ctx, &assignments); err != nil {
		return nil, err
	}
	return assignments, nil
}

func (r *mongoRepository) UpdateAssignment(ctx context.Context, id string, update bson.M) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return ErrNotFound
	}

	update["updated_at"] = time.Now().UTC()
	result, err := r.assignments.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": update})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *mongoRepository) DeleteAssignment(ctx context.Context, id string) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return ErrNotFound
	}

	result, err := r.assignments.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}

// ── Submission methods ────────────────────────────────────────────────────────

func (r *mongoRepository) CreateSubmission(ctx context.Context, s *models.Submission) error {
	s.ID = bson.NewObjectID()
	now := time.Now().UTC()
	s.CreatedAt = now
	s.UpdatedAt = now

	_, err := r.submissions.InsertOne(ctx, s)
	if mongo.IsDuplicateKeyError(err) {
		return ErrAlreadyExists
	}
	return err
}

func (r *mongoRepository) GetSubmissionByID(ctx context.Context, id string) (*models.Submission, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrNotFound
	}

	var s models.Submission
	err = r.submissions.FindOne(ctx, bson.M{"_id": oid}).Decode(&s)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotFound
	}
	return &s, err
}

func (r *mongoRepository) GetSubmissionByAssignmentAndStudent(ctx context.Context, assignmentID, studentID string) (*models.Submission, error) {
	var s models.Submission
	err := r.submissions.FindOne(ctx, bson.M{
		"assignment_id": assignmentID,
		"student_id":    studentID,
	}).Decode(&s)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotFound
	}
	return &s, err
}

func (r *mongoRepository) ListSubmissionsByAssignment(ctx context.Context, assignmentID string) ([]*models.Submission, error) {
	cursor, err := r.submissions.Find(ctx, bson.M{"assignment_id": assignmentID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var submissions []*models.Submission
	if err := cursor.All(ctx, &submissions); err != nil {
		return nil, err
	}
	return submissions, nil
}

func (r *mongoRepository) UpdateSubmission(ctx context.Context, id string, update bson.M) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return ErrNotFound
	}

	update["updated_at"] = time.Now().UTC()
	result, err := r.submissions.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": update})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}
