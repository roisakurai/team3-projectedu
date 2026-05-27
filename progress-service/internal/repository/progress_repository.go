package repository

import (
	"context"
	"time"

	"progress-service/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProgressRepository interface {
	// StudentProgress
	CreateStudentProgress(ctx context.Context, sp *model.StudentProgress) error
	GetStudentProgress(ctx context.Context, studentID, classID string) (*model.StudentProgress, error)
	GetStudentProgressByStudentID(ctx context.Context, studentID string) ([]*model.StudentProgress, error)
	GetStudentProgressByClassID(ctx context.Context, classID string) ([]*model.StudentProgress, error)
	IncrementStudentProgressField(ctx context.Context, studentID, classID, field string, amount int) error
	UpdateStudentAverageGrade(ctx context.Context, studentID, classID string, avg float64) error
	UpdateStudentLastActivity(ctx context.Context, studentID, classID string) error
	CountDistinctStudents(ctx context.Context) (int64, error)
	CountDistinctClasses(ctx context.Context) (int64, error)
	GetPlatformAverageGrade(ctx context.Context) (float64, error)
	GetTotalSubmissions(ctx context.Context) (int64, error)

	// MaterialProgress
	CreateMaterialProgress(ctx context.Context, mp *model.MaterialProgress) error
	GetMaterialProgress(ctx context.Context, studentID, materialID string) (*model.MaterialProgress, error)
	GetMaterialProgressByStudentClass(ctx context.Context, studentID, classID string) ([]*model.MaterialProgress, error)
	SetMaterialCompleted(ctx context.Context, studentID, materialID string, completedAt time.Time) error

	// AssignmentProgress
	CreateAssignmentProgress(ctx context.Context, ap *model.AssignmentProgress) error
	GetAssignmentProgress(ctx context.Context, studentID, assignmentID string) (*model.AssignmentProgress, error)
	GetAssignmentProgressByStudentClass(ctx context.Context, studentID, classID string) ([]*model.AssignmentProgress, error)
	UpsertAssignmentSubmitted(ctx context.Context, studentID, assignmentID, classID string, submittedAt time.Time) error
	UpdateAssignmentGraded(ctx context.Context, studentID, assignmentID string, grade float64, gradedAt time.Time) error

	// GradeRecords
	InsertGradeRecord(ctx context.Context, gr *model.GradeRecord) error
	GetGradeRecordsByStudent(ctx context.Context, studentID string) ([]*model.GradeRecord, error)
}

type progressRepository struct {
	db *mongo.Database
}

func NewProgressRepository(db *mongo.Database) ProgressRepository {
	return &progressRepository{db: db}
}

func (r *progressRepository) studentProgressColl() *mongo.Collection {
	return r.db.Collection("student_progress")
}

func (r *progressRepository) materialProgressColl() *mongo.Collection {
	return r.db.Collection("material_progress")
}

func (r *progressRepository) assignmentProgressColl() *mongo.Collection {
	return r.db.Collection("assignment_progress")
}

func (r *progressRepository) gradeRecordsColl() *mongo.Collection {
	return r.db.Collection("grade_records")
}

// EnsureIndexes creates all required MongoDB indexes.
func EnsureIndexes(ctx context.Context, db *mongo.Database) error {
	// student_progress unique index
	_, err := db.Collection("student_progress").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "student_id", Value: 1}, {Key: "class_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	// material_progress unique index
	_, err = db.Collection("material_progress").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "student_id", Value: 1}, {Key: "material_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	// assignment_progress unique index
	_, err = db.Collection("assignment_progress").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "student_id", Value: 1}, {Key: "assignment_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	// grade_records index on student_id
	_, err = db.Collection("grade_records").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "student_id", Value: 1}},
	})
	if err != nil {
		return err
	}

	return nil
}

// ─── StudentProgress ──────────────────────────────────────────────────────────

func (r *progressRepository) CreateStudentProgress(ctx context.Context, sp *model.StudentProgress) error {
	sp.ID = primitive.NewObjectID()
	now := time.Now().UTC()
	sp.CreatedAt = now
	sp.UpdatedAt = now
	sp.LastActivityAt = now
	_, err := r.studentProgressColl().InsertOne(ctx, sp)
	return err
}

func (r *progressRepository) GetStudentProgress(ctx context.Context, studentID, classID string) (*model.StudentProgress, error) {
	var sp model.StudentProgress
	err := r.studentProgressColl().FindOne(ctx, bson.M{
		"student_id": studentID,
		"class_id":   classID,
	}).Decode(&sp)
	if err != nil {
		return nil, err
	}
	return &sp, nil
}

func (r *progressRepository) GetStudentProgressByStudentID(ctx context.Context, studentID string) ([]*model.StudentProgress, error) {
	cursor, err := r.studentProgressColl().Find(ctx, bson.M{"student_id": studentID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*model.StudentProgress
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *progressRepository) GetStudentProgressByClassID(ctx context.Context, classID string) ([]*model.StudentProgress, error) {
	cursor, err := r.studentProgressColl().Find(ctx, bson.M{"class_id": classID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*model.StudentProgress
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *progressRepository) IncrementStudentProgressField(ctx context.Context, studentID, classID, field string, amount int) error {
	_, err := r.studentProgressColl().UpdateOne(ctx,
		bson.M{"student_id": studentID, "class_id": classID},
		bson.M{
			"$inc": bson.M{field: amount},
			"$set": bson.M{"updated_at": time.Now().UTC()},
		},
	)
	return err
}

func (r *progressRepository) UpdateStudentAverageGrade(ctx context.Context, studentID, classID string, avg float64) error {
	_, err := r.studentProgressColl().UpdateOne(ctx,
		bson.M{"student_id": studentID, "class_id": classID},
		bson.M{"$set": bson.M{
			"average_grade": avg,
			"updated_at":    time.Now().UTC(),
		}},
	)
	return err
}

func (r *progressRepository) UpdateStudentLastActivity(ctx context.Context, studentID, classID string) error {
	now := time.Now().UTC()
	_, err := r.studentProgressColl().UpdateOne(ctx,
		bson.M{"student_id": studentID, "class_id": classID},
		bson.M{"$set": bson.M{
			"last_activity_at": now,
			"updated_at":       now,
		}},
	)
	return err
}

func (r *progressRepository) CountDistinctStudents(ctx context.Context) (int64, error) {
	vals, err := r.studentProgressColl().Distinct(ctx, "student_id", bson.M{})
	if err != nil {
		return 0, err
	}
	return int64(len(vals)), nil
}

func (r *progressRepository) CountDistinctClasses(ctx context.Context) (int64, error) {
	vals, err := r.studentProgressColl().Distinct(ctx, "class_id", bson.M{})
	if err != nil {
		return 0, err
	}
	return int64(len(vals)), nil
}

func (r *progressRepository) GetPlatformAverageGrade(ctx context.Context) (float64, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.M{
			"_id":           nil,
			"average_grade": bson.M{"$avg": "$average_grade"},
		}}},
	}
	cursor, err := r.studentProgressColl().Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result []struct {
		AverageGrade float64 `bson:"average_grade"`
	}
	if err := cursor.All(ctx, &result); err != nil {
		return 0, err
	}
	if len(result) == 0 {
		return 0, nil
	}
	return result[0].AverageGrade, nil
}

func (r *progressRepository) GetTotalSubmissions(ctx context.Context) (int64, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$submitted_assignments"},
		}}},
	}
	cursor, err := r.studentProgressColl().Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result []struct {
		Total int64 `bson:"total"`
	}
	if err := cursor.All(ctx, &result); err != nil {
		return 0, err
	}
	if len(result) == 0 {
		return 0, nil
	}
	return result[0].Total, nil
}

// ─── MaterialProgress ─────────────────────────────────────────────────────────

func (r *progressRepository) CreateMaterialProgress(ctx context.Context, mp *model.MaterialProgress) error {
	mp.ID = primitive.NewObjectID()
	mp.CreatedAt = time.Now().UTC()
	_, err := r.materialProgressColl().InsertOne(ctx, mp)
	return err
}

func (r *progressRepository) GetMaterialProgress(ctx context.Context, studentID, materialID string) (*model.MaterialProgress, error) {
	var mp model.MaterialProgress
	err := r.materialProgressColl().FindOne(ctx, bson.M{
		"student_id":  studentID,
		"material_id": materialID,
	}).Decode(&mp)
	if err != nil {
		return nil, err
	}
	return &mp, nil
}

func (r *progressRepository) GetMaterialProgressByStudentClass(ctx context.Context, studentID, classID string) ([]*model.MaterialProgress, error) {
	cursor, err := r.materialProgressColl().Find(ctx, bson.M{
		"student_id": studentID,
		"class_id":   classID,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*model.MaterialProgress
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *progressRepository) SetMaterialCompleted(ctx context.Context, studentID, materialID string, completedAt time.Time) error {
	_, err := r.materialProgressColl().UpdateOne(ctx,
		bson.M{"student_id": studentID, "material_id": materialID},
		bson.M{"$set": bson.M{
			"is_completed": true,
			"completed_at": completedAt,
		}},
		options.Update().SetUpsert(true),
	)
	return err
}

// ─── AssignmentProgress ───────────────────────────────────────────────────────

func (r *progressRepository) CreateAssignmentProgress(ctx context.Context, ap *model.AssignmentProgress) error {
	ap.ID = primitive.NewObjectID()
	now := time.Now().UTC()
	ap.CreatedAt = now
	ap.UpdatedAt = now
	_, err := r.assignmentProgressColl().InsertOne(ctx, ap)
	return err
}

func (r *progressRepository) GetAssignmentProgress(ctx context.Context, studentID, assignmentID string) (*model.AssignmentProgress, error) {
	var ap model.AssignmentProgress
	err := r.assignmentProgressColl().FindOne(ctx, bson.M{
		"student_id":    studentID,
		"assignment_id": assignmentID,
	}).Decode(&ap)
	if err != nil {
		return nil, err
	}
	return &ap, nil
}

func (r *progressRepository) GetAssignmentProgressByStudentClass(ctx context.Context, studentID, classID string) ([]*model.AssignmentProgress, error) {
	cursor, err := r.assignmentProgressColl().Find(ctx, bson.M{
		"student_id": studentID,
		"class_id":   classID,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*model.AssignmentProgress
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *progressRepository) UpsertAssignmentSubmitted(ctx context.Context, studentID, assignmentID, classID string, submittedAt time.Time) error {
	now := time.Now().UTC()
	_, err := r.assignmentProgressColl().UpdateOne(ctx,
		bson.M{"student_id": studentID, "assignment_id": assignmentID},
		bson.M{"$set": bson.M{
			"status":       model.AssignmentStatusSubmitted,
			"submitted_at": submittedAt,
			"class_id":     classID,
			"updated_at":   now,
		},
			"$setOnInsert": bson.M{
				"_id":        primitive.NewObjectID(),
				"created_at": now,
			},
		},
		options.Update().SetUpsert(true),
	)
	return err
}

func (r *progressRepository) UpdateAssignmentGraded(ctx context.Context, studentID, assignmentID string, grade float64, gradedAt time.Time) error {
	_, err := r.assignmentProgressColl().UpdateOne(ctx,
		bson.M{"student_id": studentID, "assignment_id": assignmentID},
		bson.M{"$set": bson.M{
			"status":     model.AssignmentStatusGraded,
			"grade":      grade,
			"graded_at":  gradedAt,
			"updated_at": time.Now().UTC(),
		}},
	)
	return err
}

// ─── GradeRecords ─────────────────────────────────────────────────────────────

func (r *progressRepository) InsertGradeRecord(ctx context.Context, gr *model.GradeRecord) error {
	gr.ID = primitive.NewObjectID()
	_, err := r.gradeRecordsColl().InsertOne(ctx, gr)
	return err
}

func (r *progressRepository) GetGradeRecordsByStudent(ctx context.Context, studentID string) ([]*model.GradeRecord, error) {
	cursor, err := r.gradeRecordsColl().Find(ctx, bson.M{"student_id": studentID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*model.GradeRecord
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
